# SRE Workshop

[English](README.md) · **繁體中文**

涵蓋 SRE（Site Reliability Engineering）核心技能的實戰工作坊系列，包含 Docker 容器化、GitHub Actions CI/CD 自動化、Prometheus 監控，最後用一個綜合練習把三者串起來。

## 工作坊列表

| 工作坊 | 時長 | 說明 |
|--------|------|------|
| [Docker](Docker/) | 90 分鐘 | Container、Docker Compose、Dockerfile |
| [CI/CD](CI-CD/) | 100 分鐘 | GitHub Actions、CI pipeline、self-hosted runner 部署 |
| [Prometheus](Prometheus/) | 60 分鐘 | Metrics 收集、告警設定、Discord 通知 |
| [最後練習](final-exercise/) | 25 分鐘 | 透過 CI/CD 部署服務、Prometheus 抓 metric、告警到 Discord |

## 課前檢核清單

工作坊前請跑過一遍。每個項目下面都有安裝方式與驗證指令。

- [ ] Git
- [ ] Go 1.24+
- [ ] Docker（內含 Docker Compose）
- [ ] make
- [ ] VS Code（含 YAML 擴充套件）
- [ ] GitHub 帳號
- [ ] Docker Hub 帳號
- [ ] Discord 帳號

**Windows 使用者**：所有操作都在 **WSL 2** 裡面做。先裝好 WSL，之後每個工具都照「Windows (WSL)」那條路線走。
用系統管理員身份打開 PowerShell 以安裝 WSL 2：

```powershell
wsl --install
```
之後重新開啟 PowerShell 時，要再次進入 WSL，可以輸入下面這段指令：
```powershell
wsl
```

### 安裝與驗證

<details>
<summary><b>Git</b></summary>

- **macOS**：`brew install git`（或第一次執行 `git` 時讓 Xcode CLI tools 自動裝）
- **Windows (WSL)**：`sudo apt update && sudo apt install -y git`

驗證：
```bash
git --version
# git version 2.x
```

</details>

<details>
<summary><b>Go 1.24+</b></summary>

- **macOS**：`brew install go`
- **Windows (WSL)**：`sudo snap install go --classic` — 或從 [go.dev/dl](https://go.dev/dl/) 下載安裝

驗證：
```bash
go version
# go version go1.24.x darwin/arm64
```

</details>

<details>
<summary><b>Docker（內含 Docker Compose）</b></summary>

- **macOS**：`brew install --cask orbstack`，然後啟動 OrbStack
- **Windows (WSL)**：安裝 [Docker Desktop](https://www.docker.com/products/docker-desktop/)，並在 *Settings → Resources → WSL Integration* 啟用 WSL 2 整合。`docker` 指令就可以在 Ubuntu WSL 裡直接用。

驗證：
```bash
docker --version
# Docker version 27.x, build ...

docker compose version
# Docker Compose version v2.x
```

</details>

<details>
<summary><b>make</b></summary>

- **macOS**：Xcode CLI tools 內建，如果沒有請執行 `xcode-select --install`
- **Windows (WSL)**：`sudo apt install -y make`

驗證：
```bash
make --version
# GNU Make 4.x
```

</details>

<details>
<summary><b>VS Code</b></summary>

- **macOS**：`brew install --cask visual-studio-code` — 或從 [code.visualstudio.com](https://code.visualstudio.com/) 下載
- **Windows**：從 [code.visualstudio.com](https://code.visualstudio.com/) 下載，並安裝 [WSL 擴充套件](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-wsl)，這樣才能直接編輯 WSL 裡的檔案

裝好之後從 VS Code 的 Extensions 面板再裝一個 **YAML** 擴充套件。

</details>

<details>
<summary><b>GitHub / Docker Hub / Discord 帳號</b></summary>

註冊即可（不需安裝）：

- **GitHub** — [github.com/signup](https://github.com/signup)
- **Docker Hub** — [hub.docker.com/signup](https://hub.docker.com/signup)
- **Discord** — [discord.com/register](https://discord.com/register)

</details>

## 專案結構

```
sre-workshop/
├── README.md               # 英文版 README
├── README.zh-TW.md         # 本檔案
├── Docker/                 # Docker 工作坊
├── CI-CD/                  # CI/CD 工作坊
├── Prometheus/             # Prometheus 工作坊
└── final-exercise/         # 整天課程結束後的綜合練習
```

## 授權

本教材僅供 SDC 工作坊教學使用。
