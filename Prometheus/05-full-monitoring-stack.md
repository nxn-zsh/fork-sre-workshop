# 05 — 實戰：完整監控系統

> **時間：30 分鐘**

---

## 目錄

- [學習目標](#學習目標)
- [完整架構回顧](#完整架構回顧)
- [Step 1：加入 Blackbox Exporter](#step-1加入-blackbox-exporter)
  - [什麼是 Blackbox Exporter？](#什麼是-blackbox-exporter)
  - [建立 Blackbox Exporter 設定](#建立-blackbox-exporter-設定)
  - [設定 Prometheus scrape Blackbox Exporter](#設定-prometheus-scrape-blackbox-exporter)
- [Step 2：加入被監控的目標服務](#step-2加入被監控的目標服務)
- [Step 3：新增 Blackbox Alert Rules](#step-3新增-blackbox-alert-rules)
- [Step 4：加入 Grafana](#step-4加入-grafana)
- [Step 5：更新 Docker Compose — 完整版](#step-5更新-docker-compose--完整版)
- [Step 6：啟動完整 Stack](#step-6啟動完整-stack)
- [Step 7：探索 Grafana](#step-7探索-grafana)
  - [建立你的第一個 Dashboard](#建立你的第一個-dashboard)
- [Step 8：端到端測試 — 觸發 Alert](#step-8端到端測試--觸發-alert)
- [五步流程驗證](#五步流程驗證)
- [Next Steps](#next-steps)
- [小結](#小結)

---

## 學習目標

完成本章節後，你將能夠：

- 部署 Blackbox Exporter 來監測 HTTP endpoints
- 部署 Grafana 並連接 Prometheus 作為 data source
- 在 Grafana 中建立基本的 Dashboard
- 從頭到尾走完「蒐集 → 評估 → 告警 → 通知 → 視覺化」的完整流程
- 手動觸發 alert 並觀察整個告警流程

---

## 完整架構回顧

到這一章結束時，你的 monitoring stack 會長這樣：

```
        ┌─────────────────────┐
        │   whoami (範例服務)   │
        │   :80               │
        └──────────┬──────────┘
                   │ HTTP probe
                   ▼
    ┌──────────────────────────┐
    │   Blackbox Exporter      │
    │   :9115                  │
    └─────────────┬────────────┘
                  │ 被 scrape
                  ▼
    ┌──────────────────────────┐     ┌──────────────────┐
    │      Prometheus          │────▶│   Alertmanager   │
    │      :9090               │     │   :9093          │
    └──────┬───────────────────┘     └────────┬─────────┘
           │                                  │
  ┌────────┼────────────┐                     ▼
  │        │            │            ┌──────────────────┐
  ▼        ▼            ▼            │    Discord       │
┌────┐  ┌────────┐  ┌──────────┐    │    (通知)        │
│self│  │ Node   │  │Alertmgr  │    └──────────────────┘
│    │  │Exporter│  │  metrics │
└────┘  │ :9100  │  └──────────┘
        └────────┘

    ┌──────────────────────────┐
    │       Grafana            │ ← 查詢 Prometheus
    │       :3000              │
    └──────────────────────────┘
```

---

## Step 1：加入 Blackbox Exporter

### 什麼是 Blackbox Exporter？

回顧第二章的 Blackbox vs. Whitebox 概念：

- **Node Exporter** 是 whitebox monitoring — 從服務內部回報指標
- **Blackbox Exporter** 是 blackbox monitoring — 從外部探測服務是否正常

Blackbox Exporter 可以做以下類型的探測（probe）：

| Probe 類型 | 檢查什麼 | 範例 |
|-----------|---------|------|
| `http` | HTTP/HTTPS endpoint 是否回應正常 | 網站是否可以連上？回應碼是不是 200？ |
| `tcp` | TCP port 是否開放 | 資料庫的 port 是否可連線？ |
| `icmp` | 主機是否可以 ping 到 | 伺服器是否在線上？ |
| `dns` | DNS 解析是否正常 | 域名能不能正確解析？ |

本課程聚焦在 **HTTP probe**。

### 建立 Blackbox Exporter 設定

建立 `config/blackbox.yml`：

```yaml
modules:
  http_2xx:
    prober: http
    timeout: 5s
    http:
      valid_http_versions: ["HTTP/1.1", "HTTP/2.0"]
      valid_status_codes: [200]
      method: GET
      follow_redirects: true
      preferred_ip_protocol: "ip4"
```

這個設定定義了一個名為 `http_2xx` 的 probe module：
- 使用 HTTP prober
- 超時時間 5 秒
- 期望收到 200 狀態碼
- 使用 GET 方法
- 跟隨 redirect

### 設定 Prometheus scrape Blackbox Exporter

Blackbox Exporter 的 scrape 設定比較特殊——你要告訴 Prometheus：「去 Blackbox Exporter 查詢某個 URL 的探測結果」。

在 `config/prometheus.yml` 的 `scrape_configs` 中加入：

```yaml
  # Blackbox Exporter — HTTP 探測
  - job_name: "blackbox-http"
    metrics_path: /probe
    params:
      module: [http_2xx]
    static_configs:
      - targets:
          - "http://whoami:80"         # 我們要探測的目標 URL
        labels:
          service: "whoami"
          env: "dev"
    relabel_configs:
      # 把 target URL 傳給 Blackbox Exporter 的 `target` 參數
      - source_labels: [__address__]
        target_label: __param_target
      # 保留原始 target URL 作為 `instance` label
      - source_labels: [__param_target]
        target_label: instance
      # 實際要 scrape 的地址改成 Blackbox Exporter
      - target_label: __address__
        replacement: blackbox-exporter:9115
```

> 💡 **講師提示：** `relabel_configs` 這段比較進階。可以簡單解釋為：「我們告訴 Prometheus，去問 Blackbox Exporter 這個 URL 的狀態如何」。不需要讓學生完全理解 relabel 的機制，只要知道這是標準的 Blackbox Exporter 設定模板即可。

---

## Step 2：加入被監控的目標服務

我們需要一個實際的 HTTP 服務來讓 Blackbox Exporter 探測。使用 `traefik/whoami` — 一個極簡的 HTTP 服務，會回傳自己的資訊。

這個服務會在 Docker Compose 中定義（見 Step 5）。

---

## Step 3：新增 Blackbox Alert Rules

在 `config/rules/alerts.yml` 中加入 Blackbox Exporter 相關的 alert rules：

```yaml
  - name: blackbox-alerts
    rules:
      # Alert：HTTP 探測失敗（服務連不上）
      - alert: EndpointDown
        expr: probe_success == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "{{ $labels.service }} ({{ $labels.env }}) 服務無法連線"
          description: "{{ $labels.instance }} 已經超過 1 分鐘無法連線。請立即確認服務狀態。"

      # Alert：HTTP 回應太慢
      - alert: EndpointSlow
        expr: probe_duration_seconds > 3
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "{{ $labels.service }} 回應時間過長"
          description: "{{ $labels.instance }} 的回應時間為 {{ $value | printf \"%.2f\" }} 秒，超過 3 秒門檻。"
```

---

## Step 4：加入 Grafana

Grafana 用來將 Prometheus 的 metrics 視覺化成 Dashboard。

建立 Grafana 的 data source 自動設定，讓 Grafana 啟動時自動連接 Prometheus：

```bash
mkdir -p config/grafana/provisioning/datasources
```

建立 `config/grafana/provisioning/datasources/prometheus.yml`：

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
```

這個設定告訴 Grafana：「自動把 `http://prometheus:9090` 設定成預設的 Prometheus data source」。

---

## Step 5：更新 Docker Compose — 完整版

把所有元件整合到完整的 `docker-compose.yml`：

```yaml
services:
  # === 監控核心 ===
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

  # === Exporters ===
  node-exporter:
    image: prom/node-exporter:v1.9.0
    container_name: node-exporter
    ports:
      - "9100:9100"
    restart: unless-stopped

  blackbox-exporter:
    image: prom/blackbox-exporter:v0.25.0
    container_name: blackbox-exporter
    ports:
      - "9115:9115"
    volumes:
      - ./config/blackbox.yml:/etc/blackbox_exporter/config.yml:ro
    command:
      - "--config.file=/etc/blackbox_exporter/config.yml"
    restart: unless-stopped

  # === 告警 ===
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

  # === 視覺化 ===
  grafana:
    image: grafana/grafana:11.5.2
    container_name: grafana
    ports:
      - "3000:3000"
    volumes:
      - ./config/grafana/provisioning:/etc/grafana/provisioning:ro
      - grafana-data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_AUTH_ANONYMOUS_ENABLED=true
      - GF_AUTH_ANONYMOUS_ORG_ROLE=Admin
    restart: unless-stopped

  # === 被監控的目標服務 ===
  whoami:
    image: traefik/whoami:v1.10
    container_name: whoami
    ports:
      - "8080:80"
    restart: unless-stopped

volumes:
  prometheus-data:
  grafana-data:
```

### 各服務的 Port 一覽

| 服務 | Port | 網址 |
|------|------|------|
| Prometheus | 9090 | http://localhost:9090 |
| Node Exporter | 9100 | http://localhost:9100 |
| Blackbox Exporter | 9115 | http://localhost:9115 |
| Alertmanager | 9093 | http://localhost:9093 |
| Grafana | 3000 | http://localhost:3000 |
| whoami（目標服務） | 8080 | http://localhost:8080 |

---

## Step 6：啟動完整 Stack

### 完整的 prometheus.yml

先確認你的 `config/prometheus.yml` 包含了所有 scrape 設定：

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets: ["alertmanager:9093"]

rule_files:
  - "/etc/prometheus/rules/*.yml"

scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets: ["localhost:9090"]

  - job_name: "node-exporter"
    static_configs:
      - targets: ["node-exporter:9100"]

  - job_name: "alertmanager"
    static_configs:
      - targets: ["alertmanager:9093"]

  - job_name: "blackbox-http"
    metrics_path: /probe
    params:
      module: [http_2xx]
    static_configs:
      - targets:
          - "http://whoami:80"
        labels:
          service: "whoami"
          env: "dev"
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: blackbox-exporter:9115
```

### 啟動所有服務

```bash
cd Prometheus/examples/monitoring-stack
docker compose up -d
```

### 確認所有 container 都正常運行

```bash
docker compose ps
```

你應該看到六個 container 都是 **Up** 狀態。

### 確認 Prometheus Targets

打開 http://localhost:9090/targets，你應該看到四個 scrape jobs：

| Job | Target | 狀態 |
|-----|--------|------|
| `prometheus` | `localhost:9090` | UP |
| `node-exporter` | `node-exporter:9100` | UP |
| `alertmanager` | `alertmanager:9093` | UP |
| `blackbox-http` | `http://whoami:80` | UP |

> 💡 **講師提示：** 這是一個好的停下來確認的時間點。確保所有學生的六個 container 都在跑，四個 targets 都是 UP。如果有 DOWN 的，協助排查。

---

## Step 7：探索 Grafana

打開 http://localhost:3000

因為我們設定了匿名 admin 存取，所以不需要登入。

### 確認 Data Source

1. 點擊左側選單的 ⚙️ → **Data Sources**
2. 你應該看到 **Prometheus** 已經自動設定好了
3. 點擊它，在底部按 **Test** 按鈕——應該顯示 ✅ `Successfully queried the Prometheus API.`

### 用 Explore 試試 PromQL

1. 點擊左側選單的 🔍 **Explore**
2. 確認上方的 Data Source 選的是 **Prometheus**
3. 輸入以下查詢試試看：

```promql
up
```

你應該看到四筆結果，每筆對應一個 scrape target。

### 建立你的第一個 Dashboard

1. 點擊左側選單的 ➕ → **New Dashboard**
2. 點擊 **Add visualization**
3. 選擇 **Prometheus** 作為 data source

**Panel 1：CPU 使用率**

在查詢框輸入：

```promql
(1 - avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m]))) * 100
```

設定：
- Panel title: `CPU Usage (%)`
- 單位：選 **Misc → Percent (0-100)**

點擊右上角的 **Apply** 儲存這個 panel。

**Panel 2：記憶體使用率**

點擊 **Add panel** → **Add a new panel**：

```promql
(1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100
```

設定：
- Panel title: `Memory Usage (%)`
- 單位：**Misc → Percent (0-100)**

**Panel 3：服務狀態**

再新增一個 panel：

```promql
probe_success
```

設定：
- Panel title: `Service Status`
- Visualization type: 改成 **Stat**
- Value mappings: `1` → `UP`（綠色）、`0` → `DOWN`（紅色）

**儲存 Dashboard**

點擊右上角的 💾 儲存按鈕，取名為 `Monitoring Overview`。

> 💡 **講師提示：** 帶著學生一步一步操作。Grafana 的 UI 很直覺，但第一次用可能會有點不知所措。重點是讓學生體驗「把 PromQL 查詢變成視覺化圖表」的過程。

---

## Step 8：端到端測試 — 觸發 Alert

現在來走完整個告警流程！

### 停掉 whoami 服務

```bash
docker compose stop whoami
```

### 觀察 Alert 的狀態變化

1. 打開 http://localhost:9090/alerts

2. 觀察 `EndpointDown` alert 的狀態變化：
   - **0 秒** — 仍然是 Inactive（Blackbox Exporter 還沒偵測到）
   - **~15 秒** — 變成 **Pending**（黃色）— `probe_success == 0` 條件成立了
   - **~1 分 15 秒** — 變成 **Firing**（紅色）— 條件持續超過 `for: 1m`

3. 打開 http://localhost:9093 — 確認 Alertmanager 收到了 alert

4. 如果有設定 Discord Webhook — 你應該在 Discord 頻道看到告警訊息 🚨

### 確認 Grafana Dashboard 的反映

回到 Grafana 的 Dashboard，你應該看到：

- `Service Status` panel 的 whoami 變成 **DOWN**（紅色）
- 其他 metrics 不受影響

### 恢復服務

```bash
docker compose start whoami
```

觀察：
1. Prometheus Alerts 頁面：`EndpointDown` 從 Firing → Inactive
2. Alertmanager：alert 自動 resolved
3. Discord：收到 ✅ resolved 通知
4. Grafana Dashboard：whoami 恢復為 **UP**（綠色）

> 💡 **講師提示：** 這是整個課程最精彩的部分！讓學生親自體驗從「服務掛掉」到「收到通知」到「服務恢復」的完整流程。可以讓學生用自己的手機打開 Discord 看通知，增加臨場感。

---

## 五步流程驗證

回顧第一章提到的五步流程，確認你已經親手建立了每一步：

```
✅ 蒐集 (Collect)
   Prometheus 每 15 秒 scrape：自己、Node Exporter、Alertmanager、Blackbox Exporter

✅ 評估 (Evaluate)
   Prometheus 每 15 秒評估 alert rules：TargetDown、EndpointDown、HighMemoryUsage 等

✅ 告警 (Alert)
   Alert 從 Inactive → Pending → Firing，Prometheus 把 Firing 的 alert 送給 Alertmanager

✅ 通知 (Notify)
   Alertmanager 分組、路由，送出 Discord 通知

✅ 視覺化 (Visualize)
   Grafana 查詢 Prometheus，顯示即時 Dashboard
```

---

## Next Steps

恭喜！你已經從零建立了一套完整的 monitoring 系統。以下是進階的學習方向：

### P0 — 現在就可以做

- [ ] 在 Grafana 中建立更多 panels（硬碟使用率、網路流量等）
- [ ] 修改 alert rules 的 threshold，觀察不同的觸發行為
- [ ] 試著新增另一個 HTTP 服務（例如 `nginx`）作為 Blackbox 探測目標

### P1 — 深入理解

- [ ] 研究 Alertmanager 的 `inhibit_rules`（alert 抑制機制）
- [ ] 學習 Grafana 的 variables 和 templating（動態 Dashboard）
- [ ] 了解 Prometheus 的 service discovery（自動發現 targets）

### P2 — 進階實戰

- [ ] 為你自己的應用程式加上 `/metrics` endpoint（application instrumentation）
- [ ] 設定 Prometheus 的 remote write（長期儲存方案）
- [ ] 建立 Grafana alerting（Grafana 自己的告警機制）
- [ ] 探索 Loki（logs）和 Tempo（traces）完善可觀測性

### 推薦學習資源

- [PromLabs Training](https://training.promlabs.com/) — 免費的 Prometheus 線上教學
- [Prometheus 官方文件](https://prometheus.io/docs/)
- [PromQL Cheat Sheet](https://promlabs.com/promql-cheat-sheet/)
- [Grafana 官方文件](https://grafana.com/docs/)
- [Awesome Prometheus Alerts](https://awesome-prometheus-alerts.grep.to/) — 現成的 alert rule 範本集
- [Alertmanager 官方文件](https://prometheus.io/docs/alerting/latest/alertmanager/)

---

## 小結

- **Blackbox Exporter** 從外部探測 HTTP endpoints 是否正常，是 blackbox monitoring 的核心工具
- **Grafana** 連接 Prometheus data source，把 metrics 視覺化成 Dashboard
- 完整的 monitoring stack = **Prometheus + Node Exporter + Blackbox Exporter + Alertmanager + Grafana**
- 五步流程：**蒐集 → 評估 → 告警 → 通知 → 視覺化**——你已經全部走過一遍了
- 透過手動停掉服務，你親自驗證了 alert 從觸發到通知的完整流程

> **恭喜你完成了 Prometheus Monitoring Workshop！🎉**
>
> 你現在擁有建立和管理一套基本監控系統的能力。繼續探索 Next Steps 中的進階主題，逐步提升你的 SRE 技能！

---

[← 上一章：Alertmanager 告警系統](04-alerting.md)
