# Prometheus Monitoring Workshop

## 課程概覽

**1 小時**的 Prometheus Workshop，從監控基礎概念開始，逐步涵蓋 Prometheus 架構、部署、Alertmanager 告警系統，到使用 Grafana 建立完整的監控系統。完成後你將能從零開始建立一套監控系統。


## 前置需求

- [ ] **Docker** 已安裝 — `docker --version`
- [ ] **Docker Compose** 已安裝 — `docker compose version`
- [ ] **Git** 已安裝 — `git --version`
- [ ] **Discord** 帳號 — 用於接收告警通知

## 課程大綱

| 章節 | 主題 | 時長 |
|------|------|------|
| 01 | Monitoring 概念介紹 | 15 分鐘 |
| 02 | Prometheus 核心概念 | 20 分鐘 |
| 03 | 動手做：部署 Prometheus | 25 分鐘 |
| 04 | Alertmanager 告警系統 | 25 分鐘 |
| 05 | 實戰：完整監控系統 | 30 分鐘 |
| -- | Q&A | 5 分鐘 |

> **合計：2 小時（120 分鐘）**

## 專案結構

```
Prometheus/
├── README.md                          # English version
├── README.zh-TW.md                    # 本檔案
├── 01-monitoring-intro.md             # Monitoring 概念介紹
├── 02-prometheus-fundamentals.md      # Prometheus 架構與 PromQL
├── 03-first-prometheus.md             # 動手做：部署 Prometheus
├── 04-alerting.md                     # Alertmanager 告警系統
├── 05-full-monitoring-stack.md        # 實戰：完整監控系統
├── prometheus-guide.md                # 完整參考指南（PDF 轉檔）
├── prometheus_schedule.md             # 原始課程規劃（中文）
├── Prometheus_Guide_*.pdf             # 原始 PDF 指南
└── examples/                         # 範例設定檔
    └── monitoring-stack/              # 完整監控 Stack
        ├── docker-compose.yml         # 所有服務定義
        └── config/
            ├── prometheus.yml         # Prometheus 設定
            ├── alertmanager.yml       # Alertmanager 設定
            ├── blackbox.yml           # Blackbox Exporter 設定
            ├── rules/
            │   └── alerts.yml         # Alert Rules
            └── grafana/
                └── provisioning/
                    └── datasources/
                        └── prometheus.yml  # Grafana 自動設定
```

## 最終成果

完成所有章節後，你將擁有一套完整的監控系統：

```
你的 Monitoring Stack
├── Prometheus       — 蒐集 metrics、評估告警規則
├── Node Exporter    — 回報主機硬體指標（CPU、記憶體、硬碟）
├── Blackbox Exporter — 從外部探測 HTTP 服務是否正常
├── Alertmanager     — 接收告警、路由通知到 Discord
└── Grafana          — 把 metrics 視覺化成即時 Dashboard
```
