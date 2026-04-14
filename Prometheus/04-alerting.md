# 04 — Alertmanager 告警系統

---

## 目錄

- [04 — Alertmanager 告警系統](#04--alertmanager-告警系統)
  - [目錄](#目錄)
  - [Alert 的完整流程](#alert-的完整流程)
  - [Step 1： Alert Rules](#step-1-alert-rules)
    - [Alert Rule 結構解說](#alert-rule-結構解說)
    - [Severity 等級設計](#severity-等級設計)
  - [Step 2：設定 Alertmanager](#step-2設定-alertmanager)
    - [Alertmanager 設定檔解說](#alertmanager-設定檔解說)
      - [全域設定](#全域設定)
      - [Route 路由設計](#route-路由設計)
      - [`webhook_url_file` 而非 `webhook_url`](#webhook_url_file-而非-webhook_url)
  - [Step 3：建立 Discord Webhook](#step-3建立-discord-webhook)
    - [在 Discord 建立 Webhook](#在-discord-建立-webhook)
    - [把 Webhook 存成檔案](#把-webhook-存成檔案)
  - [Step 4：啟動並驗證](#step-4啟動並驗證)
    - [確認 Alert Rules 已載入](#確認-alert-rules-已載入)
    - [確認 Alertmanager 正常運作](#確認-alertmanager-正常運作)
    - [確認 Targets](#確認-targets)
  - [Alert 的三種狀態](#alert-的三種狀態)
    - [為什麼需要 Pending 狀態？](#為什麼需要-pending-狀態)
  - [常見問題排解](#常見問題排解)
    - [1. Alert Rules 沒有載入](#1-alert-rules-沒有載入)
    - [2. Alertmanager 設定錯誤](#2-alertmanager-設定錯誤)
  - [小結與練習題](#小結與練習題)
    - [本章重點回顧](#本章重點回顧)
    - [練習題](#練習題)


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

---

## Step 1： Alert Rules

Alert rules 告訴 Prometheus「什麼情況下要發出警報」。

 `Prometheus/examples/prometheus/alert_rules.yml`

```yaml
groups:
  - name: alert_manager_rules
    rules:
      # Alert 1：Target 掛掉
      - alert: InstanceDown
        expr: up == 0
        for: 5m
        labels:
          severity: critical
          component: alertmanager
        annotations:
          summary: "Instance {{ $labels.instance }} down"
          description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes."

  - name: node_exporter_alerts
    interval: 30s
    rules:
      # Alert 2：記憶體快用完
      - alert: HostOutOfMemory
        expr: (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100) < 10
        for: 2m
        labels:
          severity: critical
          component: node_exporter
        annotations:
          summary: "Host out of memory (instance {{ $labels.instance }})"
          description: "Node memory is filling up (< 10% available)"

      # Alert 3：硬碟空間不足
      - alert: HostOutOfDiskSpace
        expr: (node_filesystem_avail_bytes{mountpoint="/",fstype!="rootfs"} / node_filesystem_size_bytes{mountpoint="/",fstype!="rootfs"} * 100) < 10
        for: 2m
        labels:
          severity: critical
          component: node_exporter
        annotations:
          summary: "Host out of disk space (instance {{ $labels.instance }})"
          description: "Disk is almost full (< 10% left)"

      # Alert 4：CPU 使用率過高
      - alert: HostHighCpuLoad
        expr: (100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[2m])) * 100)) > 90
        for: 5m
        labels:
          severity: critical
          component: node_exporter
        annotations:
          summary: "Host high CPU load (instance {{ $labels.instance }})"
          description: "CPU load is > 90%"
```

### Alert Rule 結構解說

以 `InstanceDown` 為例，逐欄位解說：

```yaml
- alert: InstanceDown                  # Alert 名稱（PascalCase）
  expr: up == 0                        # PromQL 條件式
  for: 5m                              # 條件必須持續成立多久才真的觸發
  labels:                              # 附加到 alert 上的 labels
    severity: critical                 # 決定告警的嚴重等級
    component: alertmanager            # 告警來源元件，方便之後做 routing / filtering
  annotations:                         # 給人看的說明文字
    summary: "Instance {{ $labels.instance }} down"
    description: "{{ $labels.instance }} of job {{ $labels.job }} ..."
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


---

## Step 2：設定 Alertmanager

Alertmanager 負責接收 Prometheus 發來的 alerts，進行 **分組（grouping）、去重複（deduplication）、路由（routing）**，最後送到指定的通知管道。

建立 `Prometheus/examples/prometheus/alertmanager.yml`：

```yaml
global:
  resolve_timeout: 5m    # 如果 alert 停止觸發超過 5 分鐘，自動標記為 resolved

# 路由設定：所有 alert 都走到同一個 default receiver
route:
  receiver: 'default'
  group_by: ['alertname', 'severity', 'env', 'service']
  group_wait: 1m                     # 收到新 alert 後等待多久才發送（等待更多同組的 alerts 一起發）
  group_interval: 10m                # 同一組 alert 重複通知的最短間隔
  repeat_interval: 4h                # 已發送的 alert 重複提醒的間隔

# 通知接收者定義
receivers:
  - name: 'default'
    discord_configs:
      - webhook_url_file: '/etc/alertmanager/secrets/discord/webhook-default'
        send_resolved: true
        title: '{{ if eq .Status "firing" }}⚠️{{ else }}✅{{ end }} [{{ .Status | toUpper }}] {{ (index .Alerts 0).Labels.alertname }}'
        message: |
          {{- range .Alerts }}
          **{{ .Labels.alertname }}** ({{ .Labels.severity }})
          {{ .Annotations.summary }}
          {{ .Annotations.description }}
          {{- end }}
        username: 'AlertBot'
```

### Alertmanager 設定檔解說

#### 全域設定

| 設定 | 說明 |
|------|------|
| `resolve_timeout: 5m` | 如果 Prometheus 在 5 分鐘內沒有再發送同一個 alert，Alertmanager 會自動認為這個 alert 已經 resolved |

#### Route 路由設計

這裡所有 alert 都送到同一個 `default` receiver，沒有按 severity 做分流。實際上線時常見的做法是加上子路由把 `critical` 送到警急頻道、`warning` 送到一般頻道，你可以在 Alertmanager 文件裡看到 `routes:` 和 `match:` 的用法。

| 設定 | 說明 |
|------|------|
| `group_by` | 用哪些 labels 來分組。相同組的 alerts 會被合併成一封通知 |
| `group_wait` | 收到第一個 alert 後等多久才發送，目的是等待同組的其他 alerts 一起送出 |
| `group_interval` | 同一組如果有新的 alert 加入，至少間隔多久才發送更新通知 |
| `repeat_interval` | 已經發送過的 alert 還在 firing，多久再提醒一次 |

#### `webhook_url_file` 而非 `webhook_url`

Webhook URL 是機密資訊——一旦 commit 進 git 或推到公開 repo，任何人都可以拿去送假訊息。這裡用 `webhook_url_file` 指向一個檔案路徑，我們把實際的 URL 放在 `prometheus/secrets/discord/webhook-default`，透過 docker volume 掛進 container，並把整個 `prometheus/secrets/` 加到 `.gitignore`。這樣設定檔可以放心 commit，機密資料留在本機。


---

## Step 3：建立 Discord Webhook


### 在 Discord 建立 Webhook

1. 打開 Discord，進入你的伺服器
2. 選擇一個頻道（或建立一個新的 `#monitoring-alerts` 頻道）
3. 點擊頻道名稱旁的 ⚙️（頻道設定）
4. 選擇 **Integrations** → **Webhooks** → **New Webhook**
5. 取個名稱（例如 `Prometheus Alerts`）
6. 點擊 **Copy Webhook URL**

### 把 Webhook 存成檔案

不要把 URL 直接貼進 `alertmanager.yml`。改成建立一個獨立的檔案：

```bash
mkdir -p Prometheus/examples/prometheus/secrets/discord
echo -n 'https://discord.com/api/webhooks/xxxxx/yyyyy' \
  > Prometheus/examples/prometheus/secrets/discord/webhook-default
```
---


## Step 4：啟動並驗證
### 確認 Alert Rules 已載入

1. 打開 http://localhost:9090/alerts
2. 你應該看到 Alert rules 都已載入，狀態大多是 **Inactive**（綠色）——代表目前沒有問題
3. 每條 rule 旁邊會顯示 `expr`、`for`、`labels`、`annotations`

### 確認 Alertmanager 正常運作

1. 打開 http://localhost:9093
2. 你會看到 Alertmanager 的 Web UI
3. 目前應該沒有任何 alerts（因為一切正常）

### 確認 Targets

1. 打開 http://localhost:9090/targets
2. 你應該看到三個 targets 都是 **UP**：`prometheus`、`node-exporter`、`alertmanager`


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



---

## 常見問題排解

### 1. Alert Rules 沒有載入

**症狀**：Prometheus Alerts 頁面是空的。

**排查步驟**：

```bash
# 確認 rules 檔案有成功掛進 container
docker compose exec prometheus ls /alert_rules.yml

# 確認 rules 檔案語法正確
docker compose exec prometheus promtool check rules /alert_rules.yml
```

### 2. Alertmanager 設定錯誤

**症狀**：Alertmanager container 啟動後立刻退出。

```bash
# 查看錯誤日誌
docker compose logs alertmanager

# 驗證設定檔
docker compose exec alertmanager amtool check-config /etc/alertmanager/alertmanager.yml
```


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

手動停掉 Node Exporter 來觸發 `InstanceDown` alert：

```bash
docker compose stop node-exporter
```

然後觀察：

1. 在 http://localhost:9090/alerts 上看到 `InstanceDown` 從 Inactive → Pending → Firing
2. 在 http://localhost:9093 上看到 Alertmanager 收到了 alert
3. 把 `for` 改成 `0s`，讓 alert 立刻觸發（測試完記得改回來）
4. 用 `docker log alertmanager`
```time=2026-04-14T17:27:12.936s+08:00 level=DEBUG source=notify.go:975 msg="Notify success" component=dispatcher receiver=default integration=discord[0] aggrGroup="{}:{alertname=\"InstanceDown\"}" attempts=1 duration=401.501667ms numAlerts=1 alerts="InstanceDown: 1"```
1. 在 Discord 上看到通知

測試完後恢復 Node Exporter：

```bash
docker compose start node-exporter
```

觀察 alert 自動恢復為 Inactive。

**練習 2：自訂一條 Alert Rule**

在 `prometheus/alert_rules.yml` 中新增一條 alert rule，監測 Prometheus 自己的 TSDB：

- 名稱：`PrometheusTSDBTooManySeries`
- 條件：`prometheus_tsdb_head_series > 10000`（time series 數量超過 10000 條）
- 持續時間：5 分鐘
- 嚴重等級：warning

新增完 reload Prometheus：

```bash
docker compose kill -s SIGHUP prometheus
```


---

[← 上一章：部署 Prometheus](03-first-prometheus.md) ｜ [下一章：完整監控系統 →](05-full-monitoring-stack.md)
