# 扩展 Release 构建架构设计

日期：2026-11-12

## 目标

在现有 GitHub Release 工作流中增加 Windows ARM64、Linux ARM64 和 Linux RISC-V 64 位构建产物，同时保留现有 Windows AMD64、macOS Universal 和 Linux AMD64 产物。

## 构建方式

采用现有 Runner 交叉编译：

- `windows-latest` 构建 `windows/amd64` 和 `windows/arm64`。
- `macos-latest` 构建并签名 `darwin/universal`，生成 DMG。
- `ubuntu-latest` 构建 `linux/amd64`、`linux/arm64` 和 `linux/riscv64`。

MyGo 默认禁用 cgo，支持从任意宿主编译这些目标，因此不引入 ARM 原生 Runner、QEMU 或额外编译工具链。Test 工作流保持三个操作系统各运行一次，不增加架构维度。

## 工作流与打包

Release 矩阵共有 6 个目标。每个目标继续使用 `platform` 指定 MyGo 目标，使用 `output` 统一表示构建目录、压缩包架构后缀和上传产物名称。

Windows 每个架构分别生成一个 ZIP；macOS 生成一个 Universal DMG；Linux 每个架构分别生成 `.tar.gz` 和 `.deb`。上传产物名称必须包含 `matrix.output`，避免相同操作系统的不同架构互相覆盖。

`publish` 任务仍等待全部矩阵任务成功后执行。发布前必须确认包数量准确：

- Windows ZIP：2 个。
- macOS DMG：1 个。
- Linux `.tar.gz`：3 个。
- Linux `.deb`：3 个。

数量校验失败后不得发布。校验通过时，为全部安装包生成统一的 `SHA256SUMS`，再发布 GitHub Release。

## 文档与兼容性

README 的发布产物说明同步列出全部架构。现有发布触发方式、版本号、游戏代码、素材和 Test 工作流不变。本次变更不创建或推送发布标签，也不触发实际 Release。

## 验证

1. 运行工作流 YAML 结构校验器。
2. 本地交叉构建 Windows ARM64、Linux ARM64 和 Linux RISC-V 64 位目标，确认 MyGo 能生成对应目录和安装包。
3. 运行 `go test -count=1 ./...`、`go vet ./...` 和 `git diff --check`。
4. 检查矩阵产物名称与发布数量校验一致。
5. 推送后核验 GitHub Actions 的 Test 工作流；只有用户创建 `v*` 标签时才验证完整 Release 构建。
