# Extra Release Architectures Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 Windows ARM64、Linux ARM64 和 Linux RISC-V 64 位增加发布产物，并在验证后发布包含六个构建目标的 `v0.1.0`。

**Architecture:** 保留现有三个原生 Runner，在 Windows 和 Ubuntu 上利用 MyGo 的无 cgo 交叉编译能力扩展矩阵。每个矩阵条目以 `output` 作为构建目录、包名和 artifact 名称的统一标识，发布任务按准确数量验收安装包后生成校验和。

**Tech Stack:** GitHub Actions、YAML、Python/PyYAML、Go 1.27.3、MyGo 0.8.3。

## Global Constraints

- 构建目标必须恰好为 `windows/amd64`、`windows/arm64`、`darwin/universal`、`linux/amd64`、`linux/arm64`、`linux/riscv64`。
- Test 工作流保持三个操作系统各运行一次，不增加架构维度。
- 不引入 ARM 原生 Runner、QEMU、额外编译工具链或新的项目依赖。
- 不修改游戏代码、素材和 `mygo.json` 中已经等于 `0.1.0` 的版本号。
- 发布包数量必须为 2 个 Windows ZIP、1 个 macOS DMG、3 个 Linux `.tar.gz`、3 个 Linux `.deb`。
- 用户已授权在验证和推送后创建并推送 `v0.1.0` 标签，触发实际 GitHub Release。

---

### Task 1: 扩展架构校验和 Release 矩阵

**Files:**
- Modify: `scripts/check_workflows.py`
- Modify: `.github/workflows/release.yml`

**Interfaces:**
- Consumes: 现有 `load_workflow(name)`、`validate_release(workflow)` 和 `run_steps(job)`。
- Produces: 六目标矩阵；每个条目包含 `name`、`os`、`platform`、`output`、`kind`，且 `output == platform.replace("/", "-")`。

- [ ] **Step 1: 先加强校验，使旧工作流失败**

在 `validate_release` 中，把平台断言扩展为：

```python
    assert {item["platform"] for item in matrix} == {
        "windows/amd64", "windows/arm64", "darwin/universal",
        "linux/amd64", "linux/arm64", "linux/riscv64",
    }
    assert {item["kind"] for item in matrix} == {"windows", "macos", "linux"}
    for item in matrix:
        assert item["output"] == item["platform"].replace("/", "-")
```

在上传 artifact 的断言后增加：

```python
    assert uploads[0]["with"]["name"] == "release-${{ matrix.output }}"
```

在 `publish` 权限断言后增加发布数量与通配符校验：

```python
    verify = next(
        step["run"] for step in publish["steps"]
        if "sha256sum" in step.get("run", "")
    )
    assert "windows=(package/*windows-*.zip)" in verify
    assert 'test "${#windows[@]}" -eq 2' in verify
    assert 'test "${#macos[@]}" -eq 1' in verify
    assert 'test "${#linux_tar[@]}" -eq 3' in verify
    assert 'test "${#linux_deb[@]}" -eq 3' in verify
```

- [ ] **Step 2: 运行校验并确认旧配置失败**

Run:

```bash
python scripts/check_workflows.py
```

Expected: `test.yml: OK` 后出现 `AssertionError`，因为旧矩阵缺少三个新增平台。

- [ ] **Step 3: 添加三个矩阵条目**

在 Windows AMD64 条目后添加：

```yaml
          - name: Windows arm64
            os: windows-latest
            platform: windows/arm64
            output: windows-arm64
            kind: windows
```

在 Linux AMD64 条目后添加：

```yaml
          - name: Linux arm64
            os: ubuntu-latest
            platform: linux/arm64
            output: linux-arm64
            kind: linux
          - name: Linux riscv64
            os: ubuntu-latest
            platform: linux/riscv64
            output: linux-riscv64
            kind: linux
```

- [ ] **Step 4: 让打包和上传使用矩阵输出目录**

Windows 打包步骤改为：

```yaml
      - name: Package Windows
        if: matrix.kind == 'windows'
        shell: pwsh
        run: |
          New-Item -ItemType Directory -Path package -Force | Out-Null
          Compress-Archive -Path "build/${{ matrix.output }}/*" -DestinationPath "package/flappy-gopher-${env:GITHUB_REF_NAME}-${{ matrix.output }}.zip" -Force
```

Linux 打包步骤改为：

```yaml
      - name: Package Linux
        if: matrix.kind == 'linux'
        shell: bash
        run: |
          set -euo pipefail
          mkdir -p package
          cp build/${{ matrix.output }}/*.tar.gz build/${{ matrix.output }}/*.deb package/
```

上传名称改为：

```yaml
          name: release-${{ matrix.output }}
```

- [ ] **Step 5: 更新发布前数量校验**

发布校验的相关行改为：

```bash
          windows=(package/*windows-*.zip)
          macos=(package/*.dmg)
          linux_tar=(package/*.tar.gz)
          linux_deb=(package/*.deb)
          test "${#windows[@]}" -eq 2
          test "${#macos[@]}" -eq 1
          test "${#linux_tar[@]}" -eq 3
          test "${#linux_deb[@]}" -eq 3
```

- [ ] **Step 6: 运行校验并确认通过**

Run:

```bash
python scripts/check_workflows.py
```

Expected:

```text
test.yml: OK
release.yml: OK
```

- [ ] **Step 7: 本地交叉构建全部新增目标**

Run:

```bash
"/c/Program Files/Go/bin/go.exe" tool mygo build -skip-installer -platform windows/arm64,linux/arm64,linux/riscv64
```

Expected: 命令退出码为 0；生成 `build/windows-arm64/`、`build/linux-arm64/` 和 `build/linux-riscv64/`，两个 Linux 目录均包含 `.tar.gz` 和 `.deb`。

- [ ] **Step 8: 提交工作流和校验器**

```bash
git add scripts/check_workflows.py .github/workflows/release.yml
git commit -m "ci: release additional Windows and Linux architectures"
```

### Task 2: 更新发布文档并完成本地验证

**Files:**
- Modify: `README.md`

**Interfaces:**
- Consumes: Task 1 的六个发布目标。
- Produces: 与 Release 实际产物一致的用户文档。

- [ ] **Step 1: 更新 README 的发布产物说明**

将发布说明段落改为：

```markdown
Release 工作流会在三个原生 Runner 上重新测试并构建，全部成功后发布 Windows amd64/arm64 ZIP、macOS Universal DMG、Linux amd64/arm64/riscv64 的 `.tar.gz` 和 `.deb`，以及统一的 `SHA256SUMS`。新增架构通过 MyGo 在现有 Runner 上交叉编译。macOS 使用 ad hoc 签名，首次打开时可能需要按系统提示确认。
```

- [ ] **Step 2: 运行完整本地验证**

Run:

```bash
python scripts/check_workflows.py && "/c/Program Files/Go/bin/go.exe" test -count=1 ./... && "/c/Program Files/Go/bin/go.exe" vet ./... && git diff --check
```

Expected: 两个 YAML 均为 `OK`，Go 测试输出 `ok flappybird`，其余命令退出码为 0。

- [ ] **Step 3: 提交文档**

```bash
git add README.md
git commit -m "docs: list all release architectures"
```

### Task 3: 推送、发布并核验 v0.1.0

**Files:**
- Modify: 无新的源文件变更；本任务操作 Git 提交、标签和 GitHub Release。

**Interfaces:**
- Consumes: Task 1 和 Task 2 的本地提交，以及 `mygo.json` 中的 `version: 0.1.0`。
- Produces: 远程 `main`、标签 `v0.1.0` 和包含全部 10 个文件（9 个安装包加 `SHA256SUMS`）的 GitHub Release。

- [ ] **Step 1: 检查工作区与版本**

Run:

```bash
git status --short --branch
python -c "import json; assert json.load(open('mygo.json', encoding='utf-8'))['version'] == '0.1.0'; print('version: 0.1.0')"
```

Expected: 工作区没有未提交变更，版本输出为 `0.1.0`。

- [ ] **Step 2: 推送 main**

Run:

```bash
git push origin main
```

Expected: 推送成功，远程 `main` 指向本次最新提交。

- [ ] **Step 3: 核验 Test 工作流**

读取：

```text
https://api.github.com/repos/mengzhuo/flappy-gopher/actions/workflows/test.yml/runs?event=push&per_page=10
```

找到最新提交对应的 run，再读取其 `jobs_url` 并附加 `?per_page=100`。

Expected: Windows、macOS、Ubuntu 三个 job 均为 `completed / success`，Test 和 Vet 步骤均成功。若仍在运行，等待 30–50 秒后查询；若失败，停止发布并修复。

- [ ] **Step 4: 确认标签不存在**

Run:

```bash
git tag --list v0.1.0
git ls-remote --tags origin refs/tags/v0.1.0
```

Expected: 两条命令都没有标签输出。若标签已存在，停止并检查现有 Release，不覆盖标签。

- [ ] **Step 5: 创建并推送发布标签**

Run:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Expected: 标签指向已通过验证的最新 `main` 提交，远程推送成功。

- [ ] **Step 6: 核验六个构建任务与发布任务**

读取：

```text
https://api.github.com/repos/mengzhuo/flappy-gopher/actions/workflows/release.yml/runs?event=push&per_page=10
```

找到 `head_branch` 为 `v0.1.0` 且 `head_sha` 等于标签提交的 run，再读取其 `jobs_url?per_page=100`。

Expected: 6 个 build job 和 1 个 `Publish GitHub Release` job 均为 `completed / success`。若仍在运行，等待 30–50 秒后查询；若失败，读取失败步骤日志并停止声称发布成功。

- [ ] **Step 7: 核验 Release 资产**

读取：

```text
https://api.github.com/repos/mengzhuo/flappy-gopher/releases/tags/v0.1.0
```

Expected: Release 不是 draft，资产包含 2 个 Windows ZIP、1 个 macOS DMG、3 个 Linux `.tar.gz`、3 个 Linux `.deb` 和 1 个 `SHA256SUMS`，总计 10 个文件。
