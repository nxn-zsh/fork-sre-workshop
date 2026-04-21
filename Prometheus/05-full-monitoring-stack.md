# 05 — Grafana + Node Exporter + 真實告警


前一章我們讓 Alert rule 觸發，Alertmanager 會把訊息送到 Discord。但 Prometheus Web UI 的 Graph 很陽春，沒辦法長時間盯著看。這一章我們做三件事：

1. 打開 **Grafana**，直接看 **Node Exporter Full** 這張社群經典 dashboard
2. 花幾分鐘把幾個關鍵 panel 看懂——CPU、記憶體、硬碟、Load、網路
3. 用一行 `docker run` 把 VM 的 CPU 榨到 90% 以上，**實際觸發** `HostHighCpuLoad` alert，看著它從 Pending → Firing → 出現在 Discord

---

## 目錄

- [05 — Grafana + Node Exporter + 真實告警](#05--grafana--node-exporter--真實告警)
  - [目錄](#目錄)
  - [Step 1：Grafana 裡有什麼？](#step-1grafana-裡有什麼)
    - [登入](#登入)
    - [看一下 Data source 和 Dashboards 是怎麼來的](#看一下-data-source-和-dashboards-是怎麼來的)
  - [Step 2：Node Exporter Full](#step-2node-exporter-full)
    - [CPU Basic](#cpu-basic)
    - [Memory Basic](#memory-basic)
    - [Disk Space Used / Disk IOps](#disk-space-used--disk-iops)
    - [Network Traffic / Load Average](#network-traffic--load-average)
    - [所有 panel 的底層都是 PromQL](#所有-panel-的底層都是-promql)
  - [Step 3：Prometheus 2.0 Stats dashboard](#step-3prometheus-20-stats-dashboard)
  - [Step 4：用 Docker 把 VM CPU 榨乾](#step-4用-docker-把-vm-cpu-榨乾)
    - [Firing](#firing)
  - [Step 5：跟著 alert 走一遍完整流程](#step-5跟著-alert-走一遍完整流程)
    - [0:00 — Grafana：CPU panel 衝上去](#000--grafanacpu-panel-衝上去)
    - [~0:30 — Prometheus /alerts：Pending](#030--prometheus-alertspending)
    - [~5:00 — Firing + Discord 通知](#500--firing--discord-通知)
    - [停止測試](#停止測試)
  - [五步流程驗證](#五步流程驗證)
  - [常見問題排解](#常見問題排解)
  - [小結與 Next Steps](#小結與-next-steps)

---

## Step 1：Grafana 裡有什麼？

### 登入

打開 http://localhost:3000。預設帳密都是 `admin` / `admin`（由 `docker-compose.yml` 的 `ADMIN_USER` / `ADMIN_PASSWORD` 環境變數控制）。登入後 Grafana 會問要不要改密碼，直接 Skip 就好。

### 看一下 Data source 和 Dashboards 是怎麼來的

打開左側選單 → **Connections → Data sources**：

- 你會看到一個叫 **Prometheus** 的 data source 已經存在了
- URL 是 `http://prometheus:9090`（service name，不是 localhost）
- UID 是 `prometheus-main`

它是被 **provisioning** 機制自動帶進來的——`examples/grafana/provisioning/datasources/datasources.yml` 這個檔案在 Grafana 啟動時會被讀進去。

打開 **Dashboards**：

- **Node Exporter Full**
- **Prometheus 2.0 Stats**

這兩張也是 provisioning 帶進來的——JSON 檔擺在 `examples/grafana/provisioning/dashboards/` 下，Grafana 掃描那個資料夾自動載入。你不需要手動匯入。

---

## Step 2：Node Exporter Full

點進 **Node Exporter Full** dashboard。

**第一件事**：上方有 `job` 和 `instance` 兩個下拉——把 `job` 選成 `node`，`instance` 選成 `node-exporter:9100`（通常是唯一的選項）。沒選對的話 panel 會空空的。

這張 dashboard 有幾十個 panel，我們只看最上面的 **Quick CPU / Mem / Disk** 以及幾個常用區塊。

### CPU Basic

上方綠色那一堆大數字就是 **Quick CPU / Mem / Disk**——一眼看狀態。往下有 **CPU Basic** 這一區：

- **CPU Busy** — 一個大數字，CPU 使用率的百分比
- **Sys Load (5m avg)** — 5 分鐘平均 load，大於 CPU 核心數就代表排隊了
- **Usage by mode** — 把 CPU 時間拆成 user / system / iowait / idle / steal …

| 你會在圖上看到 | 代表什麼 |
|--------------|---------|
| `user` 高 | 應用程式本身在吃 CPU |
| `system` 高 | kernel 在忙（大量 syscall、context switch） |
| `iowait` 高 | CPU 在等 I/O（通常是硬碟或網路），CPU 其實沒事做 |
| `steal` 高 | 你跑在 VM 上、hypervisor 被其他 VM 搶走了時間片 |

等一下我們要把 `user` / `system` 那一段衝上去。

### Memory Basic

往下是 **Memory Basic**：

- **RAM Used** — 實際被用掉的 memory（扣掉 cache / buffer 後）
- **RAM Cache + Buffer** — Linux 會把閒置 memory 拿去做磁碟 cache。**這是好事**，不是 memory leak，空出來時隨時會還。這也是為什麼 Prometheus 的 `MemAvailable_bytes` 比 `MemFree_bytes` 更值得看。
- **SWAP Used** — 開始用 swap 通常代表 memory 壓力很大了

對應的 PromQL：`node_memory_MemAvailable_bytes` / `node_memory_MemTotal_bytes` — 這就是 Ch03 那行「記憶體使用率」query 的資料來源。

### Disk Space Used / Disk IOps

**Disk Space Used** — 每個 mount point 的使用率，紅色就該清了。注意 dashboard 會幫你過濾掉 tmpfs / overlay 這些假 filesystem。

**Disk IOps / Disk R/W Data** — 每秒幾次 I/O、每秒讀寫多少 bytes。搭配 **Time Spent Doing I/Os**（aka iowait）一起看，可以判斷是不是 I/O 瓶頸。

### Network Traffic / Load Average

**Network Traffic by Interface** — 每張網卡每秒收/送多少 bytes。`eth0` 通常是對外的主介面；`lo` 是 loopback、`docker0` 之類的可以忽略。

**Sys Load** — `load1` / `load5` / `load15`。一個直覺判讀：`load5 > CPU 核心數` 就代表「可以跑的 process 比 CPU 多」，系統在排隊。

### 所有 panel 的底層都是 PromQL

任何 panel 右上角點 **⋮ → Inspect → Panel JSON** 或 **Edit**，你就能看到它在查什麼。例如 CPU Busy 其實就是：

```promql
(((count(count(node_cpu_seconds_total{instance="$instance",job="$job"}) by (cpu))) -
  avg(sum by (mode)(rate(node_cpu_seconds_total{mode="idle",instance="$instance",job="$job"}[5m])))) * 100) /
  count(count(node_cpu_seconds_total{instance="$instance",job="$job"}) by (cpu))
```

看起來嚇人，但拆開來就是 Ch03 那個 `(1 - avg(rate(...idle[5m]))) * 100`，只是 Grafana 幫你用 `$instance` / `$job` 變數把「給人選的下拉」接到 query 裡。

---

## Step 3：Prometheus 2.0 Stats dashboard

回到 Dashboards 清單，打開 **Prometheus 2.0 Stats**。這張看的是 **Prometheus 自己**——不是它監測的對象，而是它自己這個 service 健不健康。

幾個值得看的 panel：

| Panel | 回答的問題 |
|-------|-----------|
| **Scrape duration** | 每次 scrape 花多久？太慢表示被 scrape 的 target 回應慢 |
| **Samples appended per second** | 每秒 TSDB 進來多少筆資料 |
| **Head series** | 目前活躍的 time series 數量——cardinality 爆炸的第一個線索 |
| **Memory usage** | Prometheus 自己吃了多少 RAM |

這張 dashboard 會讓你未來在「Prometheus 自己變慢/掛掉」時有線索可以看。

---

## Step 4：用 Docker 把 VM CPU 榨乾

現在我們要 **真的觸發** 第 4 章設的那個 `HostHighCpuLoad` alert：

```
expr: (100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[2m])) * 100)) > 90
for: 5m
```

也就是 **CPU 使用率大於 90%、撐 5 分鐘以上**。

Docker Desktop 在 macOS / Windows 底下其實是一個 LinuxKit VM，Node Exporter 讀的是 **VM 的** `/proc`，不是你 Mac/Windows 主機的。

- LinuxKit VM 通常只開 2~4 vCPU，想讓 **VM** 的 CPU 炸掉很簡單——在任何 container 裡跑壓測就行
- 因為 Node Exporter 看的是 VM，VM CPU 炸了 = Node Exporter 回報 CPU 炸了 = `HostHighCpuLoad` alert fires

（Linux 主機上實測需要更多 worker；下面命令裡的 `--cpu 4` 通常夠了。如果 5 分鐘後沒 fire，把數字調高。）

### Firing

打開一個新的 terminal，跑：

```bash
docker run --rm -it --cpus 2 alpine sh -c "apk add --no-cache stress-ng && stress-ng --cpu 4 --timeout 600s"
```

這一行做的事：

- 起一個 `alpine` container，限制它最多用 2 顆 vCPU
- 在 container 裡裝 `stress-ng`（Linux 經典壓測工具）
- 叫 `stress-ng` 開 4 個 CPU worker，跑 600 秒（10 分鐘，足夠撐過 `for: 5m`）

跑起來後 **不要關這個 terminal**。

---

## Step 5：跟著 alert 走一遍完整流程

### 0:00 — Grafana：CPU panel 衝上去

回到 Grafana 的 **Node Exporter Full** dashboard。十幾秒內你就會看到：

- **CPU Busy** 從個位數衝到 90%+
- **CPU Basic** 的 usage by mode 那條線，`user` / `system` 佔滿整片
- **Sys Load** 的 load average 也開始爬升

這就是 monitoring 的第一個價值——你**看得到**系統在發生什麼事，不用 `ssh` 進去跑 `top`。

### ~0:30 — Prometheus /alerts：Pending

打開 http://localhost:9090/alerts，找到 `HostHighCpuLoad`：

- 前 2 分鐘左右它還是 **Inactive**（rate 需要 2 分鐘視窗才能算穩）
- 接著變 **Pending**（黃色）—— `expr` 成立了，但 `for: 5m` 還沒到
- Pending 期間 Prometheus 會顯示「已經等了多久、還剩多少秒」

這段 Pending 時間是 alert 設計的精華：防止 CPU 短暫飆一下就驚動值班工程師。

### ~5:00 — Firing + Discord 通知

`for` 倒數結束，狀態變 **Firing**（紅色）。這時候：

1. Prometheus 把 firing alert **push** 給 Alertmanager（AM `9093`）
2. Alertmanager 按 `group_wait: 1m` 等 1 分鐘（方便把相關 alert 合併）
3. 接著發 Discord——你應該在頻道裡看到 ⚠️ 開頭的訊息，內容包含 alert name、instance、severity、觸發時間

全程大概 6~7 分鐘從「開始壓測」到「Discord 響」——中間每一段都有設計過的延遲。

### 停止測試

回到跑 `stress-ng` 的 terminal，`Ctrl+C` 停掉。

看著它走另一條路徑：

1. CPU 掉下來，Grafana 的 CPU 圖馬上回落
2. `rate(...)` 窗口 2 分鐘後，`expr` 不再成立，Prometheus 把 alert 從 Firing → Inactive
3. Alertmanager 依 `send_resolved: true` 再發一則訊息到 Discord

---

## 五步流程驗證

第 1 章提到的五步流程——到這裡你 **全部** 親手走過了：

```
蒐集 (Collect)
   Prometheus 每 30 秒 scrape：自己、Go app、Node Exporter、Alertmanager

評估 (Evaluate)
   Prometheus 每 15 秒評估 rules（HostHighCpuLoad、HostOutOfMemory、InstanceDown…）

告警 (Alert)
   stress-ng 壓測 → 條件成立 → Pending → Firing

通知 (Notify)
   Alertmanager 分組、等 group_wait、送 Discord webhook

視覺化 (Visualize)
   Grafana 拿 Prometheus 做 data source，Node Exporter Full 即時畫圖
```

---

## 常見問題排解

**1. Grafana panel 全部空白 / "No data"**

- 上方 `job` / `instance` 下拉沒選對。選 `job=node` + `instance=node-exporter:9100`
- Prometheus Targets 頁面（http://localhost:9090/targets）確認 `node` target 是 UP
- Node Exporter 掛了：`docker compose logs node-exporter`

**2. 5 分鐘過了 `HostHighCpuLoad` 還是 Pending**

- 看 Grafana 上的 CPU Busy 確實有 >90%，還是只有 70~80%？後者的話把 stress-ng 的 `--cpu` 數字加大
- 或在跑壓測的同時開第二個 terminal 再加一輪 `stress-ng`
- 記住 `rate([2m])` 需要 2 分鐘的資料才會穩，「5 分鐘 for」實際上是從 rate 穩定後才開始算

**3. Discord 沒收到通知**

- Ch04 的 webhook 設對了嗎？`cat prometheus/secrets/discord/webhook-default` 看看是不是還是 dummy URL
- `docker compose logs alertmanager` 看有沒有 send 失敗的錯誤
- Alertmanager UI（http://localhost:9093）能看到 firing alert 嗎？有 = Prometheus→AM 沒問題，問題在 AM→Discord

**4. Dashboard 面板顯示 "Datasource not found"**

通常是社群 JSON 裡 `${DS_PROMETHEUS}` placeholder 沒替換。我們的 JSON 已經預先替換成 `prometheus-main` 了；如果你自己從 grafana.com 匯入新 dashboard，記得做同樣替換，或在匯入時手動選 data source。

---

## 小結與 Next Steps

**本章重點回顧**

- Grafana 用 **provisioning** 自動帶 data source 和 dashboard JSON——不用手動點設定
- **Node Exporter Full**是看主機健康；**Prometheus 2.0 Stats**用來監測 Prometheus 自己
- 所有 Grafana panel 本質上都是 PromQL query，右鍵 Edit 就能看到底層 query
- Docker Desktop 的 LinuxKit VM 架構讓「用 container 觸發 host-level alert」變得很容易
- 一個 alert 的真實生命週期：`expr` 成立 → Pending（等 `for`）→ Firing → Alertmanager group_wait → Discord → （問題解決）→ Resolved 通知

**推薦學習資源**

- [PromLabs Training](https://training.promlabs.com/)
- [Prometheus 官方文件](https://prometheus.io/docs/)
- [PromQL Cheat Sheet](https://promlabs.com/promql-cheat-sheet/)
- [Grafana 官方文件](https://grafana.com/docs/)
- [Awesome Prometheus Alerts](https://awesome-prometheus-alerts.grep.to/) — 現成 alert rule 範本集
---

[← 上一章：Alertmanager 告警系統](04-alerting.md)
