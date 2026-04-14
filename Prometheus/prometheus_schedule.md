# **Prometheus & Alertmanager 入門教學**

## **課程概要**

目標對象：從未接觸過 Prometheus 和 Alertmanager 的初學者

課程目標：理解 Monitoring 基本概念，並在本機成功部署一套能自動告警的基礎監控系統

課程時長：60 分鐘（講解 25 分鐘 \+ 休息 5 分鐘 \+ 實作 27 分鐘 \+ Q\&A 3 分鐘）

前置需求：Docker、make、Git 已安裝，已 clone sre-workshop repo

## 

## 課程設計

| 教學內容 | 時間 | 備註 |
| ----- | ----- | ----- |
| **什麼是 Monitoring？**   \- 醫生類比：監測生命徵象   \- Metrics / Logs / Traces 三大概念   \- 簡介課堂目標與最終成果 | 8 分鐘 | 用醫生監測病人的比喻切入，讓學生理解 Monitoring 的價值 展示最終成果：一個能自動發 Alert 的監控系統 |
| **Prometheus 核心概念**   \- Pull-based vs Push-based   \- Time Series 儲存方式   \- 架構總覽（四大元件） | 10 分鐘 | 搭配架構圖講解 重點：Prometheus、Blackbox Exporter、Alertmanager、Grafana 各自的角色。 用「蒐集 → 評估 → 告警 → 通知 → 視覺化」五步流程串起來 |
| **Alert Rules 簡介**   \- YAML 結構解說   \- Severity 分級   \- Alert → Discord 流程 | 7 分鐘 | 用 EndpointDown 範例講解 alert rule 結構。 只需理解：expr（條件）、for（持續時間）、severity（等級）。 展示 Discord alert 通知截圖 |
| **休息** | 5 分鐘 | 確認學生環境已安裝 |
| **（實作）啟動 Monitoring Stack**   \- 啟動 Traefik   \- 啟動 Observe Stack   \- 確認所有服務正常運作 | 10 分鐘 | 跟著 PDF Step 1\~3 操作： make core-up-dev → make run SERVICE="observe whoami" 打開 Prometheus、Grafana、Alertmanager Web UI 確認 |
| **（實作）探索 Web UI & PromQL**   \- Prometheus Targets 頁面   \- 嘗試基本 PromQL query   \- Grafana Explore 介面 | 10 分鐘 | 引導學生操作： 1\. 查看 Targets（綠色=正常、紅色=掛掉） 2\. 輸入 up、process\_resident\_memory\_bytes 3\. 在 Grafana Explore 查詢同樣的 metric |
| **（實作）觸發 Alert & 觀察流程**   \- 手動關掉 whoami 服務   \- 觀察 Alert 從 Pending → Firing   \- 確認 Alertmanager 收到通知 | 7 分鐘 | docker stop whoami → 等待 1\~2 分鐘 在 Prometheus Alerts 頁面觀察狀態變化 在 Alertmanager 頁面確認 alert 被接收 |
| **結尾與 Q\&A**   \- 回顧五步流程   \- Next Steps（P0/P1/P2）   \- Q\&A | 3 分鐘 | 快速回顧：蒐集 → 評估 → 告警 → 通知 → 視覺化 |

## **時間分配**

講解（概念介紹）：25 分鐘

休息（環境確認）：5 分鐘

實作（動手操作）：27 分鐘

Q\&A：3 分鐘

**合計：60 分鐘**