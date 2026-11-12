# Flappy Gopher 三平台 CI 与 Release 设计

**状态：** 已批准，待实施  
**日期：** 2026-11-12  
**范围：** GitHub Actions 测试、Windows/macOS/Linux 构建与 GitHub Releases 发布

## 1. 目标

为 Flappy Gopher 增加可重复的自动化测试与发布流程：

- 推送到 `main` 或向 `main` 提交 Pull Request 时，在 Windows、macOS、Linux 上运行测试和静态检查。
- 推送 `v*` 标签时，在三个原生 Runner 上并行测试和构建，全部成功后发布到同一个 GitHub Release。
- 发布 Windows amd64、macOS Universal、Linux amd64 产物及 SHA-256 校验文件。

## 2. 工作流划分

### 2.1 `.github/workflows/test.yml`

- 触发条件：`push` 到 `main`、`pull_request` 目标为 `main`。
- 权限：`contents: read`。
- 矩阵：`windows-latest`、`macos-latest`、`ubuntu-latest`。
- 每个 Runner 使用 `go.mod` 中的 Go 版本，执行：
  - `go test -count=1 ./...`
  - `go vet ./...`
- 任一平台失败则整个测试工作流失败；矩阵使用 `fail-fast: false`，保留其他平台结果便于诊断。

### 2.2 `.github/workflows/release.yml`

- 触发条件：推送匹配 `v*` 的标签，例如 `v0.1.0`。
- 构建任务权限：`contents: read`；发布任务权限：`contents: write`。
- 构建矩阵使用 `fail-fast: false`，每个平台先执行与测试工作流相同的测试和 vet，再执行 MyGo 构建。
- 每个构建任务只上传整理后的发布文件为 Actions artifact。
- 独立 `publish` 任务依赖全部构建任务成功，下载并合并产物，生成统一的 `SHA256SUMS`，然后创建一个 GitHub Release。
- Release 使用标签名和自动生成的发布说明，不创建草稿，不自动修改版本号或推送标签。

## 3. 平台与产物

| 平台 | Runner | MyGo 目标 | Release 文件 |
| --- | --- | --- | --- |
| Windows amd64 | `windows-latest` | `windows/amd64` | `flappy-gopher-<tag>-windows-amd64.zip` |
| macOS Universal | `macos-latest` | `darwin/universal` | MyGo 生成的 Universal DMG |
| Linux amd64 | `ubuntu-latest` | `linux/amd64` | MyGo 生成的 `.tar.gz` 与 `.deb` |

- Windows 使用 `-skip-installer`，不生成 NSIS 安装器；将 `build/windows-amd64` 的完整内容压缩为 ZIP。
- macOS 的 Universal 构建同时覆盖 Apple Silicon 与 Intel Mac，使用 MyGo 默认 ad hoc 签名并生成 DMG，不配置 Developer ID 或公证。
- Linux 发布 MyGo 默认生成的用户安装归档和 Debian 包，不发布 PKGBUILD。
- 每个构建任务把发布文件放入独立的 `package/` 目录；最终 Release 只上传这些文件与 `SHA256SUMS`，不上传中间 `.app`、单独 exe、缓存或测试截图。

## 4. Actions 与依赖

采用当前 Node 24 兼容版本系列：

- `actions/checkout@v7`
- `actions/setup-go@v7`，使用 `go-version-file: go.mod` 并启用模块缓存
- `actions/upload-artifact@v7`
- `actions/download-artifact@v8`
- `softprops/action-gh-release@v3`

发布只使用 GitHub 自动提供的 `GITHUB_TOKEN`，不增加私有密钥、PAT、签名证书或第三方服务。

## 5. 错误处理与一致性

- 测试、vet、构建、打包或校验任一命令失败都会停止对应任务，不能发布不完整结果。
- 任何一个平台构建失败，`publish` 不运行。
- `publish` 在生成校验文件前检查预期平台文件是否存在；缺文件即失败。
- 同一标签重复运行时，允许更新该标签对应的 Release 资产，但不改变游戏版本配置。
- `mygo.json` 的 `version` 仍是应用包版本来源；标签用于触发和命名 Release，两者由发布者保持一致。

## 6. 非目标

- 不发布 Windows arm64、Linux arm64 或独立 macOS amd64/arm64 包。
- 不生成 Windows 安装器，不做 macOS Developer ID 签名或公证。
- 不修改游戏代码、素材、物理、难度、窗口行为或用户数据。
- 不增加自动版本 bump、定时发布、多仓库镜像或更新服务器。

## 7. 验证

1. 本地解析两个 YAML 文件，检查 GitHub Actions 表达式、矩阵和权限。
2. 本地继续通过 `go test -count=1 ./...` 与 `go vet ./...`。
3. 提交并推送后，确认 `test.yml` 在三平台实际运行结果；如本机 Git 网络不可用，明确区分本地验证与 GitHub 未验证部分。
4. 不为了测试发布流程而擅自推送版本标签。Release 的实际触发由用户决定；可提供 `git tag v0.1.0 && git push origin v0.1.0` 作为后续操作。
