# 03 — 動手做：Metrics 與 PromQL

---

## 目錄

- [03 — 動手做：Metrics 與 PromQL](#03--動手做metrics-與-promql)
  - [目錄](#目錄)
  - [Metric 的四種類型](#metric-的四種類型)
    - [1. Counter：只增不減的計數器](#1-counter只增不減的計數器)
    - [2. Gauge：可升可降的儀表](#2-gauge可升可降的儀表)
    - [3. Histogram：分布直方圖](#3-histogram分布直方圖)
    - [4. Summary：摘要](#4-summary摘要)
    - [該用哪一種？](#該用哪一種)
    - [在 `/metrics` 看到的四種類型](#在-metrics-看到的四種類型)
  - [Labels — 資料的維度標籤](#labels--資料的維度標籤)
    - [為什麼需要 Labels？](#為什麼需要-labels)
    - [常見的 Label 設計](#常見的-label-設計)
    - [Label 的注意事項](#label-的注意事項)
  - [Metrics 的傳輸格式](#metrics-的傳輸格式)
    - [Exposition format 長什麼樣子](#exposition-format-長什麼樣子)
    - [對照原始碼：metric 怎麼被吐出來](#對照原始碼metric-怎麼被吐出來)
  - [Prometheus 怎麼抓這些 metrics？](#prometheus-怎麼抓這些-metrics)
    - [Targets 頁面](#targets-頁面)
    - [prometheus.yml](#prometheusyml)
  - [PromQL 入門](#promql-入門)
    - [基本查詢](#基本查詢)
    - [Label 標籤](#label-標籤)
    - [Range Vector — 時間範圍查詢](#range-vector--時間範圍查詢)
    - [常用函數](#常用函數)
    - [試試看：三個查詢](#試試看三個查詢)
  - [到這邊我們有](#到這邊我們有)
  - [常見問題排解](#常見問題排解)
  - [小結與練習題](#小結與練習題)

---

## Metric 的四種類型

Prometheus 定義了四種 metric 類型，每種適用於不同的場景。

### 1. Counter：只增不減的計數器

**特性**：只會往上增加（或在重啟時歸零），永遠不會減少。

**適用場景**：累計事件的總次數。

```
# 範例：總共處理了多少個 HTTP 請求？
http_requests_total{method="GET", path="/api/users"} = 15234

# 範例：總共發生了多少個錯誤？
http_errors_total{service="api", code="500"} = 42
```

可以想成汽車的里程表。只會往上跳，不會倒退。

### 2. Gauge：可升可降的儀表

**特性**：數值可以增加也可以減少，反映某個當下的狀態。

**適用場景**：記錄目前的數值

```
# 範例：目前的記憶體使用量
memory_usage_bytes{instance="server-01"} = 4294967296

# 範例：目前的 CPU 溫度
node_hwmon_temp_celsius{chip="coretemp"} = 65.0

# 範例：目前有多少個活躍的連線？
active_connections{service="api"} = 127
```

**比喻**：溫度計——溫度會上升也會下降

### 3. Histogram：分布直方圖

**特性**：把觀測值分配到預先定義的 buckets（區間）裡，同時記錄總數和總和。

**適用場景**：分析數值的分布情況（例如：延遲分布）。

```
# 範例：HTTP 請求回應時間的分布
# 有多少請求在 100ms 以內完成？ 500ms 以內？ 1s 以內？
http_request_duration_seconds_bucket{le="0.1"}  = 8000   # ≤ 100ms
http_request_duration_seconds_bucket{le="0.5"}  = 9500   # ≤ 500ms
http_request_duration_seconds_bucket{le="1.0"}  = 9900   # ≤ 1s
http_request_duration_seconds_bucket{le="+Inf"} = 10000  # 全部
http_request_duration_seconds_count             = 10000  # 總請求數
http_request_duration_seconds_sum               = 3250.5 # 總時間（秒）
```

**比喻**：考試成績分布圖——多少人 60 分以下、60-70、70-80、80-90、90 以上。

### 4. Summary：摘要

**特性**：在客戶端直接計算百分位數（quantile），結果不能跨實例聚合

**適用場景**：需要精確百分位數時使用

```
# 範例：HTTP 請求回應時間的百分位數
http_request_duration_seconds{quantile="0.5"}  = 0.15   # 中位數是 150ms
http_request_duration_seconds{quantile="0.9"}  = 0.45   # P90 是 450ms
http_request_duration_seconds{quantile="0.99"} = 1.2    # P99 是 1.2s
```

### 該用哪一種？

| 情境 | 建議使用 | 原因 |
|------|---------|------|
| 計算請求總數 | **Counter** | 累計值，只增不減 |
| 記錄記憶體使用量 | **Gauge** | 數值會上下變動 |
| 分析回應時間的分布 | **Histogram** | 可以在 server-side 計算百分位，可跨實例聚合 |
| 需要精確百分位數 | **Summary** | 在 client-side 直接計算 |

### 在 `/metrics` 看到的四種類型

先對 app 製造一點流量，`/metrics` 才有東西可看。這個服務有三個業務 endpoint：`/`、`/health`、`/version`：

```bash
for i in {1..10}; do curl -s http://localhost:8000/; done
curl -s http://localhost:8000/health
curl -s http://localhost:8000/version
```

打開 http://localhost:8000/metrics，對照上面四種類型：

| 在 `/metrics` 看到的 | 對應到 |
|---------------------|--------|
| `http_requests_total{path,method,status}` | **Counter**（只會增加） |
| `http_active_requests` | **Gauge**（上下都會動） |
| `http_request_duration_seconds_bucket` / `_sum` / `_count` | **Histogram**（分桶的分佈） |
| `http_request_summary_seconds{quantile}` | **Summary**（client 端算好 quantile） |

---

## Labels — 資料的維度標籤

### 為什麼需要 Labels？

如果沒有 Labels，你只能知道「總共有 15234 個 HTTP 請求」。但有了 Labels，你可以回答更細緻的問題：

```
# 沒有 labels — 只有一個數字
http_requests_total = 15234

# 有 labels — 可以從不同維度分析
http_requests_total{service="api", method="GET", env="prod"}    = 10000
http_requests_total{service="api", method="POST", env="prod"}   = 3000
http_requests_total{service="api", method="GET", env="staging"} = 2234
```

Labels 讓你可以：

- **Filtering**：只看 production 環境的資料（`env="prod"`）
- **Grouping**：按 service 分組統計（`by (service)`）
- **Routing**：把 critical alerts 送到特定的 Discord 頻道（`severity="critical"`）

### 常見的 Label 設計

| Label | 用途 | 範例值 |
|-------|------|--------|
| `service` | 哪個服務 | `grafana`, `api`, `prometheus` |
| `env` | 哪個環境 | `prod`, `staging`, `dev` |
| `instance` | 哪個instance | `10.0.1.5:8080` |
| `method` | HTTP method | `GET`, `POST`, `PUT` |
| `status` | HTTP stautus | `200`, `404`, `500` |
| `severity` | Alert 嚴重度 | `critical`, `warning`, `info` |

### Label 的注意事項

- 每一組唯一的 label 組合 = 一條獨立的 time series
- Label 的值不要用「高基數」的東西（如 user ID、request ID），否則會產生爆量的 time series
- Label 名稱用小寫加底線（例：`response_code`），值用小寫（例：`success`）

---

## Metrics 的傳輸格式

因為 Prometheus 是 **pull-based** 的，每個被監控的服務都要 expose 一個 HTTP endpoint（慣例是 `/metrics`），用純文字回傳 metrics。這個文字格式叫 **Prometheus exposition format**。

### Exposition format 長什麼樣子

```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/api/users"} 15234
http_requests_total{method="POST",path="/api/users"} 3021

# HELP memory_usage_bytes Current memory usage in bytes
# TYPE memory_usage_bytes gauge
memory_usage_bytes 4294967296

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.1"} 8000
http_request_duration_seconds_bucket{le="0.5"} 9500
http_request_duration_seconds_bucket{le="1.0"} 9900
http_request_duration_seconds_bucket{le="+Inf"} 10000
http_request_duration_seconds_count 10000
http_request_duration_seconds_sum 3250.5
```

每一行包含：

- `# HELP` — metric 的說明文字
- `# TYPE` — metric 的類型（counter、gauge、histogram、summary）
- metric 名稱 + labels + 數值

剛剛 http://localhost:8000/metrics 你看到的，就是這個格式。

### 對照原始碼：metric 怎麼被吐出來

服務要做三件事，就能讓 Prometheus 抓得到：

1. **宣告**要追蹤哪些 metric（名字、類型、有哪些 label）
2. 每次有事情發生時，**更新**對應的數字
3. 開一個 `/metrics` endpoint，讓 Prometheus 來讀

在 `cicd/examples/sample-app/` 裡，這三件事各自對應到：

| 做什麼 | 在哪個檔案 | 長什麼樣子 |
|--------|-----------|-----------|
| 宣告 metric | `metrics.go` | `prometheus.NewCounterVec(...)`、`NewGauge(...)` 這些函式 |
| 更新 metric | `main.go` 的 `record()` | 每次 request 進來就 `+1`、量時間 |
| 開 `/metrics` | `main.go` 的 `main()` | `mux.Handle("GET /metrics", ...)` 一行 |

有興趣可以打開檔案翻翻看，但這章不用讀懂 Go 程式碼
重點是「服務主動把自己的狀態寫下來，放在一個 HTTP endpoint 讓 Prometheus 來拿」。不管用什麼程式語言，流程都一樣。

---

## Prometheus 怎麼抓這些 metrics？

打開 http://localhost:9090。這是 Prometheus 自己的 Web UI。

### Targets 頁面

點上方導覽列的 **Status → Target health**。你會看到一張表，每一列是一個 Prometheus 正在 scrape 的 target：

- `prometheus` — Prometheus 自己（self-monitoring）
- `go_practice` — 剛剛那個 Go app
- `alertmanager` — Alertmanager 自己的 metrics

每個 target 應該都是 **UP**（綠色）。如果 `go_practice` 是 DOWN，表示 Prometheus 連不到 app container——通常是網路設定或 service name 寫錯。

### prometheus.yml

打開 `examples/prometheus/prometheus.yml`：

```yaml
global:
  scrape_interval: 30s       # 每 30 秒去 scrape 一次
  evaluation_interval: 15s   # 每 15 秒評估一次 alert rules

rule_files:
  - "/alert_rules.yml"       # 告警規則檔的路徑

alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - "alertmanager:9093"   # 告警要送去哪

scrape_configs:
  - job_name: "go_practice"
    static_configs:
      - targets:
          - "app:8000"
```

幾個細節：

- **`targets: "app:8000"`** 不是 `localhost:8000`——因為 Prometheus 和 app 都是 container，彼此用 Docker network 的 **service name** 通訊。
- **`scrape_interval: 30s`** 是 Global 預設值，個別 job 可以覆蓋。
- **`alerting:` 和 `rule_files:`** 下一章 Alertmanager 用的

---

## PromQL 入門

**PromQL（Prometheus Query Language）** 是 Prometheus 的查詢語言。你可以用它在 Prometheus Web UI 或 Grafana 中查詢和分析 metrics。

### 基本查詢

```promql
# 查詢一個 metric 的所有 time series
up

# 查詢特定 metric
prometheus_tsdb_head_samples_appended_total
```

`up` 是一個特殊的 metric，Prometheus 自動為每個 scrape target 生成：
- `up == 1` 表示 target 正常
- `up == 0` 表示 target 掛了

### Label 標籤

```promql
# 精確匹配
http_requests_total{method="GET"}

# 不等於
http_requests_total{method!="DELETE"}
```

### Range Vector — 時間範圍查詢

在 metric 後面加上 `[duration]`，可以查詢一段時間範圍內的資料：

```promql
# 過去 5 分鐘的資料
http_requests_total[5m]

# 過去 1 小時的資料
http_requests_total[1h]
```

常用的時間單位：

| 單位 | 意思 |
|------|------|
| `s` | 秒 |
| `m` | 分鐘 |
| `h` | 小時 |
| `d` | 天 |
| `w` | 週 |

### 常用函數

| 函數 | 用途 | 範例 |
|------|------|------|
| `rate()` | 計算 Counter 在一段時間內的每秒增長率 | `rate(http_requests_total[5m])` |
| `increase()` | 計算 Counter 在一段時間內的總增長量 | `increase(http_requests_total[1h])` |
| `avg()` | 計算平均值 | `avg(node_cpu_seconds_total)` |
| `sum()` | 計算總和 | `sum(http_requests_total)` |
| `max()` / `min()` | 最大值 / 最小值 | `max(memory_usage_bytes)` |
| `histogram_quantile()` | 從 Histogram 計算百分位數 | `histogram_quantile(0.95, rate(http_duration_seconds_bucket[5m]))` |

>`rate()` 只能用在 Counter 類型的 metric 上。對 Gauge 用 `rate()` 是沒有意義的。

### 試試看：三個查詢

回到 Prometheus UI 首頁（Graph 頁面），輸入框就是 PromQL 查詢的地方。

**看誰還活著**

```promql
up
```

按 **Execute**，你會看到每個 target 一筆結果，`1` 代表 scrape 成功、`0` 代表失敗。切到 **Graph** 分頁可以看到時間軸上的走勢——理想狀況下應該是一條直直的 `1`。

再產生一波流量：

```bash
for i in {1..15}; do curl -s http://localhost:8000/; done
```

在 Prometheus 查：

```promql
rate(http_requests_total[1m])
```

`http_requests_total` 是 counter（只會漲），直接看它沒意義——你會只看到一條持續往上爬的線。`rate(...[1m])` 算的是 **過去 1 分鐘內每秒的平均增長率**，這才是真正的「每秒請求數」。

因為 metric 帶了 `path` label，可以直接按 path 分組：

```promql
sum by (path) (rate(http_requests_total[1m]))
```

> Counter 一定要搭配 `rate()` / `irate()` / `increase()` 才有意義。

**現在有幾個 request？**

```promql
http_active_requests
```

這是 gauge，可以直接看。它大部分時候會是 0

```bash
# 10 個並行連線 30 秒
ab -c 10 -t 30 http://localhost:8000/
# 沒有 ab 的話用這個也行：
# for i in {1..200}; do curl -s http://localhost:8000/ & done; wait
```

---

## 到這邊我們有

- 一個會 expose `/metrics` 的服務
- 一個把它抓下來、存進 TSDB 的 Prometheus
- 一個可以用 PromQL 查詢畫圖的 Web UI

但 Prometheus Web UI 的 Graph 很陽春，看完 query 就沒了，沒辦法長期觀察、沒辦法拼 dashboard；而且它只會「顯示」資料，不會「通知」你出事。


---

## 常見問題排解

**1. Port 被占用**

```
Error starting userland proxy: listen tcp4 0.0.0.0:9090: bind: address already in use
```

找出誰在用：

```bash
lsof -i :9090
# 或者：改 docker-compose.yml 的 ports 對應，例如 "9091:9090"
```

**2. `app` target 是 DOWN**

```bash
# 從 prometheus container 裡面 ping 一下 app
docker compose exec prometheus wget -qO- http://app:8000/metrics | head
```

如果上面這行也連不到，通常是 `app` container 沒起來，或 service name 寫錯。

**3. `/metrics` 看不到 `http_requests_total`**

Counter/Histogram/Summary 都是帶 label 的，要至少有一次 observation 才會出現對應的 time series。先 `curl http://localhost:8000/` 幾次再回去刷 `/metrics`。

**4. 設定檔改了 Prometheus 沒反應**

Prometheus 不會自動 reload。兩種方式：

```bash
docker compose restart prometheus           # 重啟 container
docker compose kill -s SIGHUP prometheus    # 送 SIGHUP（不中斷服務）
```

改完後先用 `promtool` 驗語法再套用：

```bash
docker run --rm \
  -v "$PWD/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro" \
  -v "$PWD/prometheus/alert_rules.yml:/alert_rules.yml:ro" \
  --entrypoint promtool prom/prometheus:v3.11.2 \
  check config /etc/prometheus/prometheus.yml
```

---

## 小結與練習題

**本章重點回顧**

- 四種 metric 類型：**Counter**（累計）、**Gauge**（當下值）、**Histogram**（分布）、**Summary**（百分位）
- **Labels** 是用來區分標記 metrics 的 key-value pairs
- 服務 expose metrics 的方式：用 client library 宣告 metric → 在程式邏輯裡更新 → 開 `/metrics` HTTP endpoint
- `prometheus.yml` 的重點是 `scrape_configs`，用 `job_name` 分組，target 在同一個 Docker 網路裡要用 **service name** 而不是 `localhost`
- **PromQL** 的核心：`up`、`rate()` / `increase()`（Counter 必套）、label 篩選、range vector



---

[← 上一章：Prometheus 核心概念](02-prometheus-fundamentals.md) ｜ [下一章：Alertmanager 告警系統 →](04-alerting.md)
