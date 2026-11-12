# 素材来源与许可证

本项目使用的双翼飞机角色素材采用 CC BY 3.0，应用图标和 Kenney 音效采用 CC0 1.0 Universal。各素材的来源、作者、许可和修改方式如下；CC0 全文见 `LICENSE-CC0.txt`，音频包原始许可说明另存于 `audio/LICENSE.txt`。

## 双翼飞机 Gopher 与应用图标

- 角色原画：Go Gopher，作者 Renee French。
- 素材：`Go_gopher_biplane.jpg`，来自 Wikimedia Commons。
- 来源：https://commons.wikimedia.org/wiki/File:Go_gopher_biplane.jpg
- 原始项目链接：https://golang.org/doc/gopher/
- 许可证：Creative Commons Attribution 3.0（CC BY 3.0）。
- 许可证全文：https://creativecommons.org/licenses/by/3.0
- 使用文件：`sprites/gopher.png`。
- 修改：去除与画布边缘连通的白色纸张背景，按白色封闭区域添加浅 Go 蓝和暖色，保留深色线稿，裁掉右上方独立速度线并输出为 880 × 624 透明 PNG。处理脚本为 `../scripts/process_biplane.py`。
- 动画派生：`gopher-body.png` 与 `gopher-propeller.png` 由 `../scripts/split_propeller.py` 从当前 `gopher.png` 拆分，均为 880 × 624 透明 PNG。主体保留轮毂，桨叶保留两片桨叶和运动弧线；原始 `gopher.png` 不覆盖，作者与 CC BY 3.0 许可证不变。
- 应用图标：`icon.png` 仍使用上一版 Egon Elbre 的 CC0 飞行 Gopher 派生图标，尺寸 256 × 256。
- 图标原始素材：https://github.com/egonelbre/gophers/blob/master/sketch/superhero/flying.png
- 图标原文件 SHA-256：`bf7ee9969140484912285259606645559b43b069ded88110bff73fd59e962a82`。
- 图标许可证：CC0 1.0 Universal；全文见 `LICENSE-CC0.txt`。

## 背景、GC 时间门和地面

运行时终端背景、漂浮关键字、程序化群山、GC 时间门、堆内存块及地址均由本项目 Go 绘制代码生成，不使用第三方图片。毫秒数值和内存地址只是装饰，不来自真实运行时。

## 音效

- 作品：Interface Sounds 1.0
- 作者：Kenney（https://kenney.nl/）
- 官方页面：https://kenney.nl/assets/interface-sounds
- 许可证：CC0
- WAV 转换分发：https://github.com/Calinou/kenney-interface-sounds
- 使用文件：
  - `audio/flap.wav` ← `pluck_001.wav`
  - `audio/score.wav` ← `confirmation_001.wav`
  - `audio/hit.wav` ← `error_008.wav`
- 修改：仅重命名。
