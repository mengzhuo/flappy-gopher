# Flappy Gopher（MyGo）

使用 Go 和 MyGo Native UI 编写的桌面单键街机游戏。控制双翼飞机 Gopher 穿过 `STW MAX` / `STW MIN` 之间的 `GC BUDGET`，避免触发 `STOP THE WORLD`。纯 Go、无 WebView；角色原画为 Renee French 的 CC BY 3.0 双翼飞机素材，图标与 Kenney 音效为 CC0，详见 `resources/ASSETS.md`。

## 操作

- 空格 / 方向上键 / 鼠标左键：`go func()` 推进，结束后重新启动 goroutine
- P / Esc：暂停或继续
- R：返回准备界面
- M：静音或恢复声音
- 窗口失焦：自动暂停，返回后按 P 或 Esc 继续

每穿过一道 GC 时间门得 1 分。每第 5 次会出现 `128 MB freed` 趣味文案，但不额外加分；画面中的毫秒数、内存量和地址都不是真实运行时指标。

门帽毫秒数由门缝整体高度决定：STW MAX 为 06～14ms，STW MIN 为 01～04ms；门缝越靠下数值越大，首个居中的门显示 12ms／03ms。数字仅用于装饰，不影响碰撞或计分。

每局第一个 GC 门在首次 flap 后的 4.0 秒到达飞机，后续门每隔 1.6 秒到达。难度由相邻门的高度修正量逐级提升：前 5 次门间变化固定为 1 个决策段，第 6～10 次随机 1～2 段，第 11～20 次随机 1～3 段，之后随机 1～4 段。决策段按当前 flap 上升顶点及同窗口自由落高度校准，只描述关卡高度差，不限制玩家真实按键次数；门速、门缝和飞行物理保持不变。

## 窗口与画面

启动时，窗口内容宽、高分别取主显示器工作区的 60%。窗口可以自由调整宽高比；逻辑高度和垂直物理保持不变，较宽窗口只会显示更多前方场景，连续 GC 门的到达间隔仍为 1.6 秒。地面地址是装饰值，从 `0x00A0` 开始，每经过一个 32 单位内存块增加 `0x20`，重开新一局时重置。

背景中的远、中、近三层群山完全由 Go 代码按确定性伪随机山脊生成，随地面滚动形成视差，不引入图片素材或新碰撞体。飞机以水平姿态为主，随升降轻微俯仰。所有文字、路径和飞机直接按窗口原生分辨率绘制，不再通过固定高度画布放大整幅画面。

飞机主体与两片桨叶使用同一原始画布的透明派生素材。飞行时桨叶绕原轮毂每秒旋转 2 圈，并跟随飞机整体俯仰；准备界面保持静止，暂停和结束时冻结，重开时相位归零。

画面采用明快的 Go 品牌风格：浅蓝白天空、Go Blue 门体描边、Deep Aqua 文字与门帽、Aqua 通过反馈及 Fuchsia 危险提示。暂停与结束状态使用白色卡片，不再显示深色整屏遮罩。

## 环境

- Go 1.27.3 或兼容的更新版本
- Windows 10 / 11；MyGo 也支持 macOS 和 Linux

本机 Go 安装在 `C:\Program Files\Go\bin\go.exe`。如果终端尚未把它加入 `PATH`，可在 Git Bash 中将下列命令里的 `go` 替换为 `"/c/Program Files/Go/bin/go.exe"`。

## 开发运行

```sh
go tool mygo dev
```

## 测试与静态检查

```sh
go test -count=1 ./...
go vet ./...
```

测试覆盖空气阻力、碰撞、管道回收、计分、反馈生命周期、主题文案、状态切换、键鼠输入、失焦暂停、最高分持久化、素材解码、音效事件、动画时钟、窗口尺寸、世界缩放、门时序、递增地址、飞机俯仰、山脊连续性、视差、原生坐标映射、静音行为、决策段校准、首门四秒到达、难度级别、长期序列可达性、状态冻结、浅色对比度、高度时间映射、桨叶素材拆分、轮毂坐标和旋转时钟。

## 构建 Windows 版

```sh
go tool mygo build -skip-installer -platform windows/amd64
```

构建产物位于 `build/windows-amd64/Flappy Gopher.exe`。图片与 WAV 已嵌入程序，无需另附素材目录。应用标识为 `io.mzh.flappygopher`，最高分保存为该标识对应 MyGo 用户数据目录中的 `high-score.json`。本次改名不继承旧 `dev.mygo.flappybird` 的最高分与窗口位置，也不会删除旧数据；存储或音频设备不可用时不影响游玩。

## GitHub Actions 与发布

推送到 `main` 或提交 Pull Request 时，Test 工作流会在 Windows、macOS 和 Linux 上运行完整测试与 `go vet ./...`。

发布前先更新并提交 `mygo.json` 中的 `version`，再创建同版本的 `v*` 标签：

```sh
git tag v0.1.0
git push origin v0.1.0
```

Release 工作流会在三个原生 Runner 上重新测试并构建，全部成功后发布 Windows amd64 ZIP、macOS Universal DMG、Linux amd64 `.tar.gz` 和 `.deb`，以及统一的 `SHA256SUMS`。macOS 使用 ad hoc 签名，首次打开时可能需要按系统提示确认。

## 素材与许可证

双翼飞机 Gopher 原画为 Renee French 的 CC BY 3.0 素材；应用图标和 Kenney 音效为 CC0。具体作者、原始文件、下载来源和修改方式见 `resources/ASSETS.md`。CC0 全文位于 `resources/LICENSE-CC0.txt`，音频包原始许可说明位于 `resources/audio/LICENSE.txt`。
