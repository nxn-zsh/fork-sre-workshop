# 04 — Alertmanager 告警系統

> **時間：25 分鐘**

---

## 目錄

- [學習目標](#學習目標)
- [Alert 的完整流程](#alert-的完整流程)
- [Step 1：建立 Alert Rules](#step-1建立-alert-rules)
  - [Alert Rule 結構解說](#alert-rule-結構解說)
  - [Severity 等級設計](#severity-等級設計)
- [Step 2：設定 Alertmanager](#step-2設定-alertmanager)
  - [Alertmanager 設定檔解說](#alertmanager-設定檔解說)
  - [Route 路由設計](#route-路由設計)
- [Step 3：建立 Discord Webhook](#step-3建立-discord-webhook)
- [Step 4：更新 Docker Compose](#step-4更新-docker-compose)
- [Step 5：更新 Prometheus 設定](#step-5更新-prometheus-設定)
- [Step 6：啟動並驗證](#step-6啟動並驗證)
- [Alert 的三種狀態](#alert-的三種狀態)
- [常見問題排解](#常見問題排解)
- [小結與練習題](#小結與練習題)

---

## 學習目標

完成本章節後，你將能夠：

- 撰寫 Prometheus Alert Rules（告警規則）
- 部署和設定 Alertmanager
- 設定 Alertmanager 的 routing 規則
- 建立 Discord Webhook 來接收告警通知
- 理解 Alert 從觸發到通知的完整流程

---

## Alert 的完整流程

在開始實作之前，先理解 alert 從產生到通知的完整流程：

```
Step 1: Prometheus 每 15 秒評估一次 alert rules
  │
  ▼
Step 2: 如果 PromQL 條件在 `for` 設定的時間內持續成立
  │     alert 狀態從 Inactive → Pending → Firing
  ▼
Step 3: Prometheus 把 Firing 的 alert 傳給 Alertmanager
  │
  ▼
Step 4: Alertmanager 依照 labels（severity、service 等）把 alerts 分組
  │
  ▼
Step 5: 等過了 group_wait 時間，Alertmanager 送出通知
  │
  ▼
Step 6: 你在 Discord 收到告警訊息 🚨
```

> 💡 **講師提示：** 先在白板上畫出這個流程，讓學生有全局概念。後面的步驟就是一步一步把這個流程建起來。

---

## Step 1：建立 Alert Rules

Alert rules 告訴 Prometheus「什麼情況下要發出警報」。

在 `Prometheus/examples/monitoring-stack/config/` 下建立 alert rules 檔案：

```bash
mkdir -p Prometheus/examples/monitoring-stack/config/rules
```

建立 `config/rules/alerts.yml`：

```yaml
groups:
  - name: node-exporter-alerts
    rules:
      # Alert 1：Target 掛掉
      - alert: TargetDown
        expr: up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "{{ $labels.instance }} 已無法連線"
          description: "Target {{ $labels.instance }} (job: {{ $labels.job }}) 已經超過 1 分鐘沒有回應 scrape 請求。"

      # Alert 2：記憶體使用率過高
      - alert: HighMemoryUsage
        expr: (1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100 > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "記憶體使用率超過 80%"
          description: "Instance {{ $labels.instance }} 的記憶體使用率已達 {{ $value | printf \"%.1f\" }}%，持續超過 5 分鐘。"

      # Alert 3：硬碟空間不足
      - alert: DiskSpaceLow
        expr: (1 - node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100 > 85
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "硬碟使用率超過 85%"
          description: "Instance {{ $labels.instance }} 的根目錄硬碟使用率已達 {{ $value | printf \"%.1f\" }}%。"

      # Alert 4：CPU 使用率過高
      - alert: HighCPUUsage
        expr: (1 - avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m]))) * 100 > 80
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "CPU 使用率超過 80%"
          description: "Instance {{ $labels.instance }} 的 CPU 使用率已達 {{ $value | printf \"%.1f\" }}%，持續超過 10 分鐘。"
```

### Alert Rule 結構解說

以 `TargetDown` 為例，逐欄位解說：

```yaml
- alert: TargetDown                    # Alert 名稱（PascalCase）
  expr: up == 0                        # PromQL 條件式
  for: 1m                             # 條件必須持續成立多久才真的觸發
  labels:                             # 附加到 alert 上的 labels
    severity: critical                # 決定告警的嚴重等級和路由
  annotations:                        # 給人看的說明文字
    summary: "{{ $labels.instance }} 已無法連線"
    description: "Target {{ $labels.instance }} (job: {{ $labels.job }}) ..."
```

| 欄位 | 說明 |
|------|------|
| `alert` | Alert 的名稱，用 PascalCase，描述 **症狀** 而非原因 |
| `expr` | PromQL 表達式。條件為 true 時 alert 就會被觸發 |
| `for` | 條件必須 **持續成立** 多久，alert 才會真的 fire。防止瞬間的抖動造成誤報 |
| `labels` | 附加到 alert 上的 labels，可以用來做 routing、filtering |
| `annotations` | 給人看的文字說明。支援 Go template 語法，`{{ $labels.xxx }}` 會被替換成實際的 label 值 |

### Severity 等級設計

| 等級 | 代表意義 | 範例 | 處理方式 |
|------|---------|------|---------|
| `critical` | 服務已經受到影響，需要 **立即處理** | 服務掛了、資料庫無法連線 | 立即通知值班人員 |
| `warning` | 有潛在風險或效能下降，**儘快處理** | 記憶體使用率 > 80%、硬碟快滿了 | 排進當天或隔天的工作 |
| `info` | 狀態變化或預防性事件，**知悉即可** | SSL 憑證快到期、服務重啟了 | 有空再看 |

> 💡 **講師提示：** 強調 `for` 的重要性——如果沒有 `for`，metrics 瞬間的波動就會觸發 alert，造成大量的假警報（false positive）。`for: 1m` 表示條件必須持續 1 分鐘才觸發，這大大減少了雜訊。

---

## Step 2：設定 Alertmanager

Alertmanager 負責接收 Prometheus 發來的 alerts，進行 **分組（grouping）、去重複（deduplication）、路由（routing）**，最後送到指定的通知管道。

建立 `config/alertmanager.yml`：

```yaml
global:
  resolve_timeout: 5m    # 如果 alert 停止觸發超過 5 分鐘，自動標記為 resolved

# 通知模板（選用）
# templates:
#   - '/etc/alertmanager/templates/*.tmpl'

# 路由設定
route:
  receiver: "discord-default"        # 預設的通知接收者
  group_by: ["alertname", "severity"] # 用哪些 labels 來分組 alerts
  group_wait: 30s                    # 收到新 alert 後等待多久才發送（等待更多同組的 alerts 一起發）
  group_interval: 5m                 # 同一組 alert 重複通知的最短間隔
  repeat_interval: 4h                # 已發送的 alert 重複提醒的間隔

  # 子路由：根據 severity 送到不同頻道
  routes:
    - match:
        severity: critical
      receiver: "discord-critical"
      repeat_interval: 1h            # critical alerts 每小時重複提醒

# 通知接收者定義
receivers:
  - name: "discord-default"
    discord_configs:
      - webhook_url: "YOUR_DISCORD_WEBHOOK_URL_HERE"
        title: '{{ if eq .Status "firing" }}🔔{{ else }}✅{{ end }} [{{ .Status | toUpper }}] {{ .GroupLabels.alertname }}'
        message: |
          {{ range .Alerts }}
          **{{ .Labels.alertname }}** ({{ .Labels.severity }})
          {{ .Annotations.summary }}
          {{ .Annotations.description }}
          {{ end }}

  - name: "discord-critical"
    discord_configs:
      - webhook_url: "YOUR_DISCORD_CRITICAL_WEBHOOK_URL_HERE"
        title: '{{ if eq .Status "firing" }}🚨{{ else }}✅{{ end }} [{{ .Status | toUpper }}] {{ .GroupLabels.alertname }}'
        message: |
          {{ range .Alerts }}
          **{{ .Labels.alertname }}** ({{ .Labels.severity }})
          {{ .Annotations.summary }}
          {{ .Annotations.description }}
          {{ end }}
```

### Alertmanager 設定檔解說

#### 全域設定

| 設定 | 說明 |
|------|------|
| `resolve_timeout: 5m` | 如果 Prometheus 在 5 分鐘內沒有再發送同一個 alert，Alertmanager 會自動認為這個 alert 已經 resolved |

#### Route 路由設計

```
所有 alerts
  │
  ├── severity == critical → discord-critical（每 1 小時重複提醒）
  │
  └── 其他（default）→ discord-default（每 4 小時重複提醒）
```

| 設定 | 說明 |
|------|------|
| `group_by` | 用哪些 labels 來分組。相同組的 alerts 會被合併成一封通知 |
| `group_wait` | 收到第一個 alert 後等多久才發送，目的是等待同組的其他 alerts 一起送出 |
| `group_interval` | 同一組如果有新的 alert 加入，至少間隔多久才發送更新通知 |
| `repeat_interval` | 已經發送過的 alert 還在 firing，多久再提醒一次 |

> 💡 **講師提示：** 用「郵局信件」的比喻——`group_wait` 就像郵差等了 30 秒，看有沒有更多信要送到同一個地址，再一起送出去。`repeat_interval` 就像如果信件一直沒被領取，每隔 4 小時再提醒一次。

---

## Step 3：建立 Discord Webhook

> **如果你不使用 Discord**，可以跳過這步，改用 Alertmanager 的 Web UI 來觀察 alerts。設定檔中的 `webhook_url` 留著即可。

### 在 Discord 建立 Webhook

1. 打開 Discord，進入你的伺服器
2. 選擇一個頻道（或建立一個新的 `#monitoring-alerts` 頻道）
3. 點擊頻道名稱旁的 ⚙️（頻道設定）
4. 選擇 **Integrations** → **Webhooks** → **New Webhook**
5. 取個名稱（例如 `Prometheus Alerts`）
6. 點擊 **Copy Webhook URL**

### 更新 Alertmanager 設定

把複製的 Webhook URL 貼到 `alertmanager.yml` 中：

```yaml
receivers:
  - name: "discord-default"
    discord_configs:
      - webhook_url: "https://discord.com/api/webhooks/xxxxx/yyyyy"   # ← 貼在這裡
```

> **注意**：Webhook URL 是機密資訊，不要推到公開的 Git repository。可以用環境變數或 `.env` 檔案來管理。

---

## Step 4：更新 Docker Compose

編輯 `docker-compose.yml`，加入 Alertmanager 服務：

```yaml
services:
  prometheus:
    image: prom/prometheus:v3.2.1
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./config/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - ./config/rules:/etc/prometheus/rules:ro
      - prometheus-data:/prometheus
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"
      - "--storage.tsdb.path=/prometheus"
      - "--storage.tsdb.retention.time=7d"
      - "--web.console.libraries=/etc/prometheus/console_libraries"
      - "--web.console.templates=/etc/prometheus/consoles"
    restart: unless-stopped

  node-exporter:
    image: prom/node-exporter:v1.9.0
    container_name: node-exporter
    ports:
      - "9100:9100"
    restart: unless-stopped

  alertmanager:
    image: prom/alertmanager:v0.28.1
    container_name: alertmanager
    ports:
      - "9093:9093"
    volumes:
      - ./config/alertmanager.yml:/etc/alertmanager/alertmanager.yml:ro
    command:
      - "--config.file=/etc/alertmanager/alertmanager.yml"
    restart: unless-stopped

volumes:
  prometheus-data:
```

注意 Prometheus 的 volumes 多了一行 `./config/rules:/etc/prometheus/rules:ro`，把 alert rules 掛載進去。

---

## Step 5：更新 Prometheus 設定

編輯 `config/prometheus.yml`，告訴 Prometheus 去哪裡找 alert rules，以及把 alerts 送去哪裡：

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

# 告訴 Prometheus 要把 alerts 送到 Alertmanager
alerting:
  alertmanagers:
    - static_configs:
        - targets: ["alertmanager:9093"]

# 告訴 Prometheus 去哪裡讀取 alert rules
rule_files:
  - "/etc/prometheus/rules/*.yml"

scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets: ["localhost:9090"]

  - job_name: "node-exporter"
    static_configs:
      - targets: ["node-exporter:9100"]

  # 也監測 Alertmanager 自己的 metrics
  - job_name: "alertmanager"
    static_configs:
      - targets: ["alertmanager:9093"]
```

新增的部分：
- `alerting` — 定義 Alertmanager 的位置
- `rule_files` — 定義 alert rules 檔案的路徑
- 新增了 `alertmanager` 的 scrape job

---

## Step 6：啟動並驗證

### 重啟所有服務

```bash
cd Prometheus/examples/monitoring-stack
docker compose up -d
```

### 驗證各元件

```bash
docker compose ps
```

你應該看到三個 container 都是 **Up** 狀態：

```
NAME            IMAGE                        STATUS         PORTS
alertmanager    prom/alertmanager:v0.28.1     Up             0.0.0.0:9093->9093/tcp
node-exporter   prom/node-exporter:v1.9.0    Up             0.0.0.0:9100->9100/tcp
prometheus      prom/prometheus:v3.2.1       Up             0.0.0.0:9090->9090/tcp
```

### 確認 Alert Rules 已載入

1. 打開 http://localhost:9090/alerts
2. 你應該看到四條 alert rules，狀態都是 **Inactive**（綠色）——代表目前沒有問題
3. 每條 rule 旁邊會顯示 `expr`、`for`、`labels`、`annotations`

### 確認 Alertmanager 正常運作

1. 打開 http://localhost:9093
2. 你會看到 Alertmanager 的 Web UI
3. 目前應該沒有任何 alerts（因為一切正常）

### 確認 Targets

1. 打開 http://localhost:9090/targets
2. 你應該看到三個 targets 都是 **UP**：`prometheus`、`node-exporter`、`alertmanager`

> 💡 **講師提示：** 讓學生確認所有三個 UI 都可以正常打開。如果 Alertmanager 無法啟動，最常見的原因是 YAML 格式錯誤——請他們檢查 `alertmanager.yml` 的縮排。

---

## Alert 的三種狀態

在 Prometheus 的 Alerts 頁面中，每條 alert rule 有三種狀態：

```
Inactive（綠色）→  Pending（黃色）→  Firing（紅色）
     ↑                                    │
     └──── 條件不再成立 ←──────────────────┘
```

| 狀態 | 意義 | 觸發條件 |
|------|------|---------|
| **Inactive** | 一切正常，PromQL 條件 **不成立** | 預設狀態 |
| **Pending** | PromQL 條件 **成立了**，但還沒超過 `for` 設定的持續時間 | `expr` 為 true，但持續時間 < `for` |
| **Firing** | PromQL 條件 **持續成立超過** `for` 設定的時間，alert 已送到 Alertmanager | `expr` 為 true，持續時間 ≥ `for` |

### 為什麼需要 Pending 狀態？

如果沒有 `for`（或 `for: 0s`），metrics 的瞬間波動就會立刻觸發 alert。例如：

- CPU 使用率瞬間飆到 85%（可能只是一個大的 compilation），1 秒後就降回來了
- 設定 `for: 10m` 可以避免這種誤報，確保是真正的持續性問題才觸發

> 💡 **講師提示：** 可以告訴學生：「Pending 就像是醫院的觀察期。體溫偏高不一定要馬上急救，但如果持續 10 分鐘都很高，那就該處理了。」

---

## 常見問題排解

### 1. Alert Rules 沒有載入

**症狀**：Prometheus Alerts 頁面是空的。

**排查步驟**：

```bash
# 確認 rules 檔案路徑正確
docker compose exec prometheus ls /etc/prometheus/rules/

# 確認 rules 檔案語法正確
docker compose exec prometheus promtool check rules /etc/prometheus/rules/alerts.yml
```

### 2. Alertmanager 設定錯誤

**症狀**：Alertmanager container 啟動後立刻退出。

```bash
# 查看錯誤日誌
docker compose logs alertmanager

# 驗證設定檔
docker compose exec alertmanager amtool check-config /etc/alertmanager/alertmanager.yml
```

### 3. Discord 沒有收到通知

**排查步驟**：

1. 確認 Webhook URL 是正確的
2. 確認 alert 狀態是 **Firing**（不是 Pending）
3. 打開 Alertmanager Web UI 確認 alerts 有被接收
4. 檢查 Alertmanager 的 logs：`docker compose logs alertmanager`

### 4. 想手動測試 alert 但不想等

你可以暫時把 `for` 改成 `0s`，讓 alert 立刻觸發（測試完記得改回來）。

---

## 小結與練習題

### 本章重點回顧

- Alert Rules 定義在 YAML 檔案中，包含 `expr`（條件）、`for`（持續時間）、`labels`（分類）、`annotations`（說明）
- **Severity 等級**：`critical`（立即處理）、`warning`（儘快處理）、`info`（知悉即可）
- **Alertmanager** 負責 routing、grouping、deduplication，最終送出通知
- Alert 有三種狀態：**Inactive** → **Pending** → **Firing**
- `for` 的設計可以有效減少誤報（false positive）

### 練習題

**練習 1：觸發一個 Alert**

手動停掉 Node Exporter 來觸發 `TargetDown` alert：

```bash
docker compose stop node-exporter
```

然後觀察：

1. 在 http://localhost:9090/alerts 上看到 `TargetDown` 從 Inactive → Pending → Firing
2. 在 http://localhost:9093 上看到 Alertmanager 收到了 alert
3. 在 Discord 上看到通知（如果有設定 Webhook）

測試完後恢復 Node Exporter：

```bash
docker compose start node-exporter
```

觀察 alert 自動恢復為 Inactive。

**練習 2：自訂一條 Alert Rule**

在 `config/rules/alerts.yml` 中新增一條 alert rule，監測 Prometheus 自己的 TSDB：

- 名稱：`PrometheusTSDBTooManySeries`
- 條件：`prometheus_tsdb_head_series > 10000`（time series 數量超過 10000 條）
- 持續時間：5 分鐘
- 嚴重等級：warning

> **接下來，我們要加入 Blackbox Exporter 和 Grafana，完成完整的監控系統！**

---

[← 上一章：部署 Prometheus](03-first-prometheus.md) ｜ [下一章：完整監控系統 →](05-full-monitoring-stack.md)
