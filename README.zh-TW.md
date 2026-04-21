# SRE Workshop

涵蓋 SRE（Site Reliability Engineering）核心技能的實戰工作坊系列，包含 Docker 容器化、GitHub Actions CI/CD 自動化，以及 Prometheus 監控。

## 工作坊列表

| 工作坊 | 時長 | 說明 |
|--------|------|------|
| [Docker](Docker/) | 2 小時 | Docker 基礎、Dockerfile 撰寫、Docker Compose 編排 |
| [CI/CD](CI-CD/) | 2 小時 | GitHub Actions、CI pipeline、Release 與部署自動化 |
| [Prometheus](Prometheus/) | 1 小時 | Metrics 收集、告警設定、Grafana 儀表板 |

> **建議上課順序**：Docker → CI/CD → Prometheus，後續課程預設學員已具備前面課程的知識。

## 課前準備

- [ ] **Git** — `git --version`
- [ ] **Go 1.22+** — `go version`
- [ ] **Docker** — `docker --version`
- [ ] **Docker Compose** — `docker compose version`
- [ ] **GitHub 帳號** — [github.com](https://github.com)
- [ ] **Discord 帳號** — 用於接收 Prometheus 告警通知（選用）
- [ ] **程式編輯器** — 推薦使用 [VS Code](https://code.visualstudio.com/)（搭配 YAML 擴充套件）

## 專案結構

```
sre-workshop/
├── README.md              # 英文版 README
├── README.zh-TW.md        # 本檔案
├── Docker/                # Docker 工作坊
│   ├── README.md          # 工作坊總覽（英文）
│   ├── README.zh-TW.md    # 工作坊總覽（中文）
│   └── docker-workshop.md # 完整教材
├── CI-CD/                 # CI/CD 工作坊
│   ├── README.md          # 工作坊總覽
│   ├── README.zh-TW.md    # 工作坊總覽（中文）
│   ├── 01-cicd-intro.md   # CI/CD 概念介紹
│   ├── 02-github-actions-basics.md
│   ├── 03-go-ci-pipeline.md
│   ├── 04-deployment.md
│   ├── examples/          # 範例 Go 應用程式
│   └── exercises/         # 實作練習
└── Prometheus/            # Prometheus 工作坊
    ├── README.md          # 工作坊總覽
    ├── 01-monitoring-intro.md
    ├── 02-prometheus-fundamentals.md
    ├── 03-first-prometheus.md
    ├── 04-alerting.md
    ├── 05-full-monitoring-stack.md
    └── examples/          # 監控系統設定檔
```

## 授權

本教材僅供 SDC 工作坊教學使用。
