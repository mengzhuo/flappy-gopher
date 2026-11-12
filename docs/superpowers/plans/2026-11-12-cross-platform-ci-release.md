# Flappy Gopher 三平台 CI 与 Release 实施计划

> **执行方式：** 使用 `superpowers:executing-plans` 在当前会话逐项实施；不得派发子代理。步骤使用复选框跟踪，只有观察到预期结果后才能勾选。

**目标：** 添加三平台 GitHub Actions 测试与发布流程，推送 `v*` 标签后发布 Windows ZIP、macOS Universal DMG、Linux tar.gz/deb 和 SHA-256 校验文件。

**架构：** `test.yml` 在三种原生 Runner 上执行同一组 Go 检查；`release.yml` 用三平台矩阵生成独立 artifact，再由一个 `publish` 任务汇总发布。一个标准库断言加 PyYAML 解析的本地脚本检查工作流结构，避免 YAML 缩进、触发器、权限或矩阵配置回退。

**技术栈：** GitHub Actions、Go 1.27.3、MyGo v0.8.3、Python 3 + PyYAML、Windows/macOS/Linux Runner。

**已批准规格：** `docs/superpowers/specs/2026-11-12-cross-platform-ci-release-design.md`

## 全局约束

- 测试触发器：推送到 `main`、向 `main` 提交 Pull Request。
- Release 触发器：推送 `v*` 标签；不得为了验证而擅自创建或推送标签。
- 测试矩阵：`windows-latest`、`macos-latest`、`ubuntu-latest`，均运行 `go test -count=1 ./...` 和 `go vet ./...`。
- Release 矩阵：Windows amd64、macOS Universal、Linux amd64，全部成功后才能发布。
- Actions 固定使用：`actions/checkout@v7`、`actions/setup-go@v7`、`actions/upload-artifact@v7`、`actions/download-artifact@v8`、`softprops/action-gh-release@v3`。
- Windows 使用 `-skip-installer`；macOS 使用 ad hoc 签名；不配置私有密钥或公证。
- `mygo.json` 版本号不自动修改，游戏源码与素材不变。
- 不派发子代理，不创建发布标签，不删除用户数据或旧构建产物。
- 本地 Go 命令使用 `"/c/Program Files/Go/bin/go.exe"`。
- 每个任务完成验证后再提交；推送失败时保留本地提交并如实报告，不能声称远程已更新。

---

### 任务 1：本地工作流校验器与三平台测试工作流

**文件：**
- 创建：`scripts/check_workflows.py`
- 创建：`.github/workflows/test.yml`

**接口与行为：**
- `load_workflow(path)` 使用 PyYAML `safe_load` 读取工作流。
- `validate_test(workflow)` 检查测试触发器、权限、矩阵和 Go 命令。
- YAML 的 `"on"` 必须加引号，避免 PyYAML 按 YAML 1.1 将其解析为布尔值。

- [ ] **步骤 1：先创建只覆盖 Test 工作流的失败校验器**

创建 `scripts/check_workflows.py`：

```python
from pathlib import Path
import yaml

ROOT = Path(__file__).resolve().parents[1]
WORKFLOWS = ROOT / ".github" / "workflows"
GO_COMMANDS = {"go test -count=1 ./...", "go vet ./..."}


def load_workflow(name):
    with (WORKFLOWS / name).open(encoding="utf-8") as stream:
        return yaml.safe_load(stream)


def run_steps(job):
    return {step.get("run") for step in job["steps"] if "run" in step}


def validate_test(workflow):
    assert workflow["name"] == "Test"
    assert workflow["on"]["push"]["branches"] == ["main"]
    assert workflow["on"]["pull_request"]["branches"] == ["main"]
    assert workflow["permissions"]["contents"] == "read"
    job = workflow["jobs"]["test"]
    assert job["strategy"]["fail-fast"] is False
    assert set(job["strategy"]["matrix"]["os"]) == {
        "windows-latest", "macos-latest", "ubuntu-latest"
    }
    assert job["runs-on"] == "${{ matrix.os }}"
    assert GO_COMMANDS <= run_steps(job)
    uses = [step["uses"] for step in job["steps"] if "uses" in step]
    assert uses == ["actions/checkout@v7", "actions/setup-go@v7"]


if __name__ == "__main__":
    validate_test(load_workflow("test.yml"))
    print("test.yml: OK")
```

- [ ] **步骤 2：运行校验器并确认文件缺失导致失败**

运行：

```bash
python scripts/check_workflows.py
```

预期：退出码非 0，出现 `FileNotFoundError`，因为 `.github/workflows/test.yml` 尚不存在。

- [ ] **步骤 3：创建三平台 Test 工作流**

创建 `.github/workflows/test.yml`：

```yaml
name: Test

"on":
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: read

jobs:
  test:
    name: ${{ matrix.os }}
    runs-on: ${{ matrix.os }}
    strategy:
      fail-fast: false
      matrix:
        os: [windows-latest, macos-latest, ubuntu-latest]
    steps:
      - name: Checkout
        uses: actions/checkout@v7
      - name: Set up Go
        uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
          cache: true
          cache-dependency-path: go.sum
      - name: Test
        run: go test -count=1 ./...
      - name: Vet
        run: go vet ./...
```

- [ ] **步骤 4：运行结构校验和完整 Go 检查**

运行：

```bash
python scripts/check_workflows.py
"/c/Program Files/Go/bin/go.exe" test -count=1 ./...
"/c/Program Files/Go/bin/go.exe" vet ./...
```

预期：输出 `test.yml: OK`，Go 测试通过，vet 无诊断。

- [ ] **步骤 5：提交任务 1**

```bash
git add .github/workflows/test.yml scripts/check_workflows.py
git commit -m "ci: add cross-platform test workflow"
```

---

### 任务 2：三平台 Release 构建、汇总与发布

**文件：**
- 修改：`scripts/check_workflows.py`
- 创建：`.github/workflows/release.yml`

**接口与行为：**
- Release 构建任务上传名为 `release-windows`、`release-macos`、`release-linux` 的 artifact。
- `publish` 使用 `needs: build`，合并下载到 `package/`。
- `validate_release(workflow)` 检查三个目标、发布依赖、权限和产物断言。

- [ ] **步骤 1：先扩展校验器，使 Release 缺失时失败**

把 `scripts/check_workflows.py` 完整替换为：

```python
from pathlib import Path
import yaml

ROOT = Path(__file__).resolve().parents[1]
WORKFLOWS = ROOT / ".github" / "workflows"
GO_COMMANDS = {"go test -count=1 ./...", "go vet ./..."}


def load_workflow(name):
    with (WORKFLOWS / name).open(encoding="utf-8") as stream:
        return yaml.safe_load(stream)


def run_steps(job):
    return {step.get("run") for step in job["steps"] if "run" in step}


def validate_test(workflow):
    assert workflow["name"] == "Test"
    assert workflow["on"]["push"]["branches"] == ["main"]
    assert workflow["on"]["pull_request"]["branches"] == ["main"]
    assert workflow["permissions"]["contents"] == "read"
    job = workflow["jobs"]["test"]
    assert job["strategy"]["fail-fast"] is False
    assert set(job["strategy"]["matrix"]["os"]) == {
        "windows-latest", "macos-latest", "ubuntu-latest"
    }
    assert job["runs-on"] == "${{ matrix.os }}"
    assert GO_COMMANDS <= run_steps(job)
    uses = [step["uses"] for step in job["steps"] if "uses" in step]
    assert uses == ["actions/checkout@v7", "actions/setup-go@v7"]


def validate_release(workflow):
    assert workflow["name"] == "Release"
    assert workflow["on"]["push"]["tags"] == ["v*"]
    build = workflow["jobs"]["build"]
    assert build["permissions"]["contents"] == "read"
    assert build["strategy"]["fail-fast"] is False
    matrix = build["strategy"]["matrix"]["include"]
    assert {item["platform"] for item in matrix} == {
        "windows/amd64", "darwin/universal", "linux/amd64"
    }
    assert {item["kind"] for item in matrix} == {"windows", "macos", "linux"}
    assert GO_COMMANDS <= run_steps(build)
    uploads = [
        step for step in build["steps"]
        if step.get("uses") == "actions/upload-artifact@v7"
    ]
    assert len(uploads) == 1
    assert uploads[0]["with"]["if-no-files-found"] == "error"
    publish = workflow["jobs"]["publish"]
    assert publish["needs"] == "build"
    assert publish["permissions"]["contents"] == "write"
    uses = [step["uses"] for step in publish["steps"] if "uses" in step]
    assert uses == ["actions/download-artifact@v8", "softprops/action-gh-release@v3"]
    release = publish["steps"][-1]["with"]
    assert release["files"] == "package/*"
    assert release["generate_release_notes"] is True


if __name__ == "__main__":
    validate_test(load_workflow("test.yml"))
    print("test.yml: OK")
    validate_release(load_workflow("release.yml"))
    print("release.yml: OK")
```

- [ ] **步骤 2：运行校验器并确认 Release 文件缺失**

运行：

```bash
python scripts/check_workflows.py
```

预期：先输出 `test.yml: OK`，随后因 `release.yml` 不存在而出现 `FileNotFoundError`。

- [ ] **步骤 3：创建三平台 Release 工作流**

创建 `.github/workflows/release.yml`：

```yaml
name: Release

"on":
  push:
    tags: ["v*"]

jobs:
  build:
    name: ${{ matrix.name }}
    runs-on: ${{ matrix.os }}
    permissions:
      contents: read
    strategy:
      fail-fast: false
      matrix:
        include:
          - name: Windows amd64
            os: windows-latest
            platform: windows/amd64
            output: windows-amd64
            kind: windows
          - name: macOS Universal
            os: macos-latest
            platform: darwin/universal
            output: darwin-universal
            kind: macos
          - name: Linux amd64
            os: ubuntu-latest
            platform: linux/amd64
            output: linux-amd64
            kind: linux
    steps:
      - name: Checkout
        uses: actions/checkout@v7
      - name: Set up Go
        uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
          cache: true
          cache-dependency-path: go.sum
      - name: Test
        run: go test -count=1 ./...
      - name: Vet
        run: go vet ./...
      - name: Build
        run: go tool mygo build -skip-installer -platform "${{ matrix.platform }}"
      - name: Package Windows
        if: matrix.kind == 'windows'
        shell: pwsh
        run: |
          New-Item -ItemType Directory -Path package -Force | Out-Null
          Compress-Archive -Path "build/windows-amd64/*" -DestinationPath "package/flappy-gopher-${env:GITHUB_REF_NAME}-windows-amd64.zip" -Force
      - name: Package macOS
        if: matrix.kind == 'macos'
        shell: bash
        run: |
          set -euo pipefail
          mkdir -p package
          cp build/darwin-universal/*.dmg package/
      - name: Package Linux
        if: matrix.kind == 'linux'
        shell: bash
        run: |
          set -euo pipefail
          mkdir -p package
          cp build/linux-amd64/*.tar.gz build/linux-amd64/*.deb package/
      - name: Upload package
        uses: actions/upload-artifact@v7
        with:
          name: release-${{ matrix.kind }}
          path: package/
          if-no-files-found: error
          retention-days: 1

  publish:
    name: Publish GitHub Release
    needs: build
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Download packages
        uses: actions/download-artifact@v8
        with:
          pattern: release-*
          path: package
          merge-multiple: true
      - name: Verify packages and generate checksums
        shell: bash
        run: |
          set -euo pipefail
          shopt -s nullglob
          windows=(package/*windows-amd64.zip)
          macos=(package/*.dmg)
          linux_tar=(package/*.tar.gz)
          linux_deb=(package/*.deb)
          test "${#windows[@]}" -eq 1
          test "${#macos[@]}" -eq 1
          test "${#linux_tar[@]}" -eq 1
          test "${#linux_deb[@]}" -eq 1
          cd package
          sha256sum -- *.zip *.dmg *.tar.gz *.deb > SHA256SUMS
      - name: Publish release
        uses: softprops/action-gh-release@v3
        with:
          files: package/*
          generate_release_notes: true
          draft: false
```

`matrix.output` 是构建目录的声明性元数据；打包步骤显式使用相同目录，避免把 shell 拼接隐藏在表达式中。

- [ ] **步骤 4：运行结构和回归检查**

运行：

```bash
python scripts/check_workflows.py
"/c/Program Files/Go/bin/go.exe" test -count=1 ./...
"/c/Program Files/Go/bin/go.exe" vet ./...
```

预期：输出 `test.yml: OK` 和 `release.yml: OK`，Go 检查通过。

- [ ] **步骤 5：提交任务 2**

```bash
git add .github/workflows/release.yml scripts/check_workflows.py
git commit -m "ci: add cross-platform release workflow"
```

---

### 任务 3：发布说明、推送与真实 Actions 核验

**文件：**
- 修改：`README.md`
- 纳入版本控制：`docs/superpowers/plans/2026-11-12-cross-platform-ci-release.md`

**接口与行为：**
- README 说明普通推送只触发测试，版本标签触发发布。
- 实际 Release 不由本次实施触发。

- [ ] **步骤 1：在 README 的构建段落后添加发布说明**

在“构建 Windows 版”段落之后、“素材与许可证”之前插入：

````markdown
## GitHub Actions 与发布

推送到 `main` 或提交 Pull Request 时，Test 工作流会在 Windows、macOS 和 Linux 上运行完整测试与 `go vet ./...`。

发布前先更新并提交 `mygo.json` 中的 `version`，再创建同版本的 `v*` 标签：

```sh
git tag v0.1.0
git push origin v0.1.0
```

Release 工作流会在三个原生 Runner 上重新测试并构建，全部成功后发布 Windows amd64 ZIP、macOS Universal DMG、Linux amd64 `.tar.gz` 和 `.deb`，以及统一的 `SHA256SUMS`。macOS 使用 ad hoc 签名，首次打开时可能需要按系统提示确认。
````

- [ ] **步骤 2：执行提交前完整验证**

运行：

```bash
python scripts/check_workflows.py
"/c/Program Files/Go/bin/go.exe" test -count=1 ./...
"/c/Program Files/Go/bin/go.exe" vet ./...
git diff --check
git status --short
```

预期：两个工作流均为 OK，Go 检查通过，`git diff --check` 无输出；状态只包含预期 README、计划和前序已提交文件。

- [ ] **步骤 3：提交文档与计划**

```bash
git add README.md
git add -f docs/superpowers/plans/2026-11-12-cross-platform-ci-release.md
git commit -m "docs: explain automated releases"
```

- [ ] **步骤 4：推送到用户仓库**

运行：

```bash
git push origin main
```

预期：GitHub `origin/main` 包含设计、两个工作流、校验器和 README 提交。若网络失败，保留全部本地提交，报告错误并停止远程核验，不能声称推送成功。

- [ ] **步骤 5：核验远程提交和 Test Actions**

运行：

```bash
git fetch origin
git status --short --branch
git log --oneline --decorate -4
```

通过公开 GitHub Actions API 或 Actions 页面检查本次 `main` 推送触发的 Test 运行：

`https://api.github.com/repos/mengzhuo/flappy-gopher/actions/runs?event=push&per_page=10`

必须如实记录：
- Windows、macOS、Linux 各 job 的状态与结论。
- 若仍在排队或运行，报告当前状态，不把“已触发”写成“已通过”。
- 若失败，读取对应日志并修复工作流后重新验证。
- 不创建或推送 `v*` 标签；Release 实际发布保持未触发。

- [ ] **步骤 6：最终交付**

最终中文回复列出：
- 两个工作流文件及触发条件。
- 三平台 Release 产物。
- 本地 YAML、Go 测试和 vet 结果。
- 推送提交号与真实 Actions 状态。
- 后续发布命令，不声称已经创建 Release。
