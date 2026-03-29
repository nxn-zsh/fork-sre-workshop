# Docker 入門指南

> **預計閱讀時間**：2 小時
> **適用對象**：具備 Linux 基礎，初次接觸 Docker 的開發者
> **範例語言**：Go
> **前提條件**：已安裝 Docker Engine 及 Docker Compose

---

## 目錄

- [Docker 入門指南](#docker-入門指南)
  - [目錄](#目錄)
  - [Part 1：Docker 核心概念與基本操作](#part-1docker-核心概念與基本操作)
    - [1.1 什麼是容器化？](#11-什麼是容器化)
      - [問題情境：「在我電腦上可以跑」](#問題情境在我電腦上可以跑)
      - [所以 Docker 如何解決這個問題？？？](#所以-docker-如何解決這個問題)
      - [容器化的適用場景](#容器化的適用場景)
    - [1.2 容器 vs 虛擬機器](#12-容器-vs-虛擬機器)
      - [虛擬機器的運作方式](#虛擬機器的運作方式)
      - [容器的運作方式](#容器的運作方式)
      - [架構比較](#架構比較)
      - [詳細比較](#詳細比較)
      - [選擇指引](#選擇指引)
      - [容器的底層技術](#容器的底層技術)
    - [1.3 Docker 架構](#13-docker-架構)
      - [架構圖](#架構圖)
      - [各元件說明](#各元件說明)
      - [一個指令的完整旅程](#一個指令的完整旅程)
    - [1.4 映像檔與容器的關係](#14-映像檔與容器的關係)
      - [概念對比](#概念對比)
      - [映像檔的分層架構（Layer）](#映像檔的分層架構layer)
      - [容器的可寫層（Container Layer）](#容器的可寫層container-layer)
      - [多容器共用映像檔](#多容器共用映像檔)
      - [映像檔命名與 Tag](#映像檔命名與-tag)
      - [概念總結](#概念總結)
    - [1.5 安裝與驗證](#15-安裝與驗證)
    - [1.6 映像檔基本操作](#16-映像檔基本操作)
      - [搜尋映像檔](#搜尋映像檔)
      - [下載映像檔](#下載映像檔)
      - [列出本機映像檔](#列出本機映像檔)
      - [刪除映像檔](#刪除映像檔)
      - [查看映像檔資訊](#查看映像檔資訊)
    - [1.7 容器基本操作](#17-容器基本操作)
      - [建立並啟動容器](#建立並啟動容器)
      - [常用 `docker run` 參數](#常用-docker-run-參數)
      - [列出容器](#列出容器)
      - [容器生命週期管理](#容器生命週期管理)
    - [1.8 Port Mapping](#18-port-mapping)
    - [1.9 Volume](#19-volume)
    - [1.10 容器除錯](#110-容器除錯)
      - [查看日誌](#查看日誌)
      - [進入容器內部](#進入容器內部)
      - [查看資源使用狀況](#查看資源使用狀況)
      - [查看容器程序](#查看容器程序)
      - [複製檔案](#複製檔案)
    - [1.11 練習 1：執行你的第一個容器](#111-練習-1執行你的第一個容器)
  - [Part 2：Dockerfile 深入實作](#part-2dockerfile-深入實作)
    - [2.1 什麼是 Dockerfile？](#21-什麼是-dockerfile)
    - [2.2 Dockerfile 基礎指令](#22-dockerfile-基礎指令)
      - [FROM — 基礎映像檔](#from--基礎映像檔)
      - [COPY vs ADD](#copy-vs-add)
      - [RUN — 建置階段執行命令](#run--建置階段執行命令)
      - [CMD vs ENTRYPOINT](#cmd-vs-entrypoint)
      - [EXPOSE — 宣告連接埠](#expose--宣告連接埠)
      - [HEALTHCHECK — 健康檢查](#healthcheck--健康檢查)
    - [2.3 .dockerignore](#23-dockerignore)
    - [2.4 為 Go 應用撰寫 Dockerfile](#24-為-go-應用撰寫-dockerfile)
      - [基本版 Dockerfile（單階段建置）](#基本版-dockerfile單階段建置)
    - [2.5 Multi-stage Build](#25-multi-stage-build)
      - [使用 scratch 的極致精簡版](#使用-scratch-的極致精簡版)
      - [建置與執行](#建置與執行)
    - [2.6 映像檔最佳實踐](#26-映像檔最佳實踐)
      - [1. 善用 Layer Cache](#1-善用-layer-cache)
      - [2. 減少 Layer 數量](#2-減少-layer-數量)
      - [3. 使用非 root 使用者](#3-使用非-root-使用者)
      - [4. 安裝套件時不保留快取](#4-安裝套件時不保留快取)
      - [5. 建置參數化](#5-建置參數化)
    - [2.7 練習 2：建置你的 Go 應用映像檔](#27-練習-2建置你的-go-應用映像檔)
  - [Part 3：Docker Compose 多容器編排](#part-3docker-compose-多容器編排)
    - [3.1 為什麼需要 Docker Compose？](#31-為什麼需要-docker-compose)
    - [3.2 基礎語法](#32-基礎語法)
    - [3.3 服務定義詳解](#33-服務定義詳解)
      - [使用映像檔 vs 建置](#使用映像檔-vs-建置)
      - [depends\_on — 服務相依](#depends_on--服務相依)
      - [restart — 重啟策略](#restart--重啟策略)
      - [資源限制](#資源限制)
    - [3.4 Network](#34-network)
      - [自訂網路（進階）](#自訂網路進階)
    - [3.5 Volume](#35-volume)
    - [3.6 環境變數管理](#36-環境變數管理)
      - [方式 1：直接在 YAML 中定義](#方式-1直接在-yaml-中定義)
      - [方式 2：使用 .env 檔案](#方式-2使用-env-檔案)
      - [方式 3：使用 env\_file](#方式-3使用-env_file)
    - [3.7 常用 Compose 指令](#37-常用-compose-指令)
    - [3.8 實戰：Go API + PostgreSQL + Redis](#38-實戰go-api--postgresql--redis)
      - [目錄結構](#目錄結構)
      - [.env](#env)
      - [docker-compose.yml](#docker-composeyml)
      - [db/init.sql](#dbinitsql)
      - [api/Dockerfile](#apidockerfile)
      - [api/main.go（完整版）](#apimaingo完整版)
      - [啟動與測試](#啟動與測試)
    - [3.9 練習 3：用 Compose 部署完整服務](#39-練習-3用-compose-部署完整服務)
  - [Part 4：綜合演練與延伸學習](#part-4綜合演練與延伸學習)
    - [4.1 綜合練習](#41-綜合練習)
    - [4.2 常見問題排查](#42-常見問題排查)
      - [容器啟動失敗](#容器啟動失敗)
      - [容器間無法連線](#容器間無法連線)
      - [磁碟空間不足](#磁碟空間不足)
      - [映像檔建置快取問題](#映像檔建置快取問題)
      - [常見錯誤訊息](#常見錯誤訊息)
    - [4.3 延伸學習資源](#43-延伸學習資源)
      - [官方文件](#官方文件)
      - [進階主題](#進階主題)
  - [附錄 A：Docker 指令速查表](#附錄-adocker-指令速查表)
    - [映像檔 (Image)](#映像檔-image)
    - [容器 (Container)](#容器-container)
    - [資料卷 (Volume)](#資料卷-volume)
    - [網路 (Network)](#網路-network)
    - [系統](#系統)
  - [附錄 B：Dockerfile 指令速查表](#附錄-bdockerfile-指令速查表)
  - [附錄 C：Docker Compose 速查表](#附錄-cdocker-compose-速查表)

---

## Part 1：Docker 核心概念與基本操作

### 1.1 什麼是容器化？

容器化（Containerization）是一種將「應用程式」及其「完整執行環境」包括程式語言、相依套件、設定檔與環境變數，封裝為方便轉移的元件的技術，而這個元件稱為**容器（Container）**。無論把容器部署到開發者的筆電、測試伺服器，還是雲端正式環境，執行結果都會保持一致。

這帶來幾個明顯的好處：團隊成員不再需要各自手動建置環境、不同專案的相依套件可以完全隔離、開發與正式環境的差異大幅縮小。容器化特別適合用在多人協作開發、CI/CD 流水線，以及需要快速複製環境的場景。

#### 問題情境：「在我電腦上可以跑」

想像一個情境，SDC 有很多個新鮮的肝正在協作開發一個由 Go 撰寫的網頁應用，專案依賴 PostgreSQL 資料庫與特定版本的函式庫。今天，身為 Project Team 的 Andrew 想要請 SRE Team 的 Ocean 把這個服務上線，兩個人的電腦環境差異很大：

| 開發者 | 作業系統 | Go 版本 | PostgreSQL 版本 | 函式庫版本 |
|--------|---------|---------|----------------|-----------|
| Andrew | macOS   | 1.21    | 15             | v2.3.1    |
| Ocean    | Ubuntu  | 1.19    | 14             | v2.1.0    |

Andrew 的程式原本在 macOS 上跑得順順的，但當程式轉移到 Ocean 的電腦之後就看到終端機噴錯了。Ocean 看著空空如也的 README 不知道要從哪裡改設定，在跟 Claude Code 奮戰兩個小時之後大喝了兩口 Whiskey，一不小心睡著了。

Andrew 因為答應了計中明天要上線，只好自己登入到系辦的機器上部署服務。因為正式環境又是另一套設定，程式剛一執行就跑出了 2000 行的 errors。可憐的 Andrew 只好一邊咒罵 SRE Team，一邊喝著自己的 Triple espresso with Vodka，跟這些 errors 奮鬥到天亮。

還好這些不是真實故事，~~因為 Ocean 會乖乖地部署完服務再睡覺~~，聰明的 Docker 也會避免這件事情發生。

#### 所以 Docker 如何解決這個問題？？？

開發者將執行環境的完整規格撰寫為 **Dockerfile**，定義以下資訊：

- 使用哪個版本的 Go
- 需要安裝哪些套件
- 需要哪些環境變數
- 程式如何啟動

任何人在任何機器上，都能透過 Docker 依照 Dockerfile 自動建置出**完全相同的容器**。新成員加入專案不再需要逐步照著手冊設定環境，執行一個指令即可開始工作。

#### 容器化的適用場景

- 多人協作開發：不同人使用不同作業系統與不同版本的工具，容器確保每個人面對完全相同的執行環境。

- 解決相依性衝突：不同專案可能依賴同一套件的不同版本。容器讓每個專案在隔離的空間中獨立運作，彼此不受干擾。

- 標準化的專案啟動程序：Dockerfile 是可執行的環境文件，明確記錄專案需求與啟動方式，取代容易過時的手動安裝說明。

- 一致的部署流程：開發、測試、正式環境使用同一個設定流程，大幅降低環境差異導致的問題。

- 快速實驗：想試 PostgreSQL 17 或 Redis 新版本？一行指令就可以快速啟動，並且用完直接回收，不需要安裝也不會汙染系統環境。

> **核心概念**：容器化將「環境」從各自的機器中抽離，轉變為可版本控制、可共享、可重現的成品。

### 1.2 容器 vs 虛擬機器

容器與虛擬機器（VM）都旨在解決「隔離」與「資源管理」的問題，但實作方式截然不同。

#### 虛擬機器的運作方式

虛擬機器透過 **Hypervisor**（虛擬機管理程式）在硬體上模擬出多台獨立電腦。每台虛擬機器都擁有完整的作業系統——從核心到使用者空間，一應俱全。

> **類比說明**：虛擬機器如同大樓裡的獨立公寓——各自擁有完整的水電管路與門牌。隔離性強，但每一間都需要重複建設相同的基礎設施，資源利用率較低。

常見的 Hypervisor 包括 VMware ESXi、KVM、Microsoft Hyper-V。

#### 容器的運作方式

容器不模擬硬體，而是利用 Linux 核心本身的功能（Namespace + Cgroup）隔離程序。所有容器共用同一個宿主機核心，每個容器只包含應用程式及其所需的函式庫。

> **類比說明**：容器如同大樓裡的共用辦公空間（co-working space）——共享水電、網路、電梯等基礎設施，各個團隊擁有獨立的工作區域，互不干擾。省空間、啟動快，但隔離強度不如獨立公寓。

#### 架構比較

```
┌──────────────────────────────┬──────────────────────────────┐
│        虛擬機器 (VM)          │        容器 (Container)       │
├──────────────────────────────┼──────────────────────────────┤
│                              │                              │
│  ┌────────┐  ┌────────┐     │  ┌────────┐  ┌────────┐     │
│  │ App A  │  │ App B  │     │  │ App A  │  │ App B  │     │
│  ├────────┤  ├────────┤     │  ├────────┤  ├────────┤     │
│  │ Bins/  │  │ Bins/  │     │  │ Bins/  │  │ Bins/  │     │
│  │ Libs   │  │ Libs   │     │  │ Libs   │  │ Libs   │     │
│  ├────────┤  ├────────┤     │  └────┬───┘  └───┬────┘     │
│  │Guest OS│  │Guest OS│     │       │          │           │
│  └────┬───┘  └───┬────┘     │  ┌────┴──────────┴────┐     │
│  ┌────┴──────────┴────┐     │  │   Docker Engine     │     │
│  │    Hypervisor      │     │  ├─────────────────────┤     │
│  ├────────────────────┤     │  │     Host OS         │     │
│  │     Host OS        │     │  ├─────────────────────┤     │
│  ├────────────────────┤     │  │    Infrastructure   │     │
│  │   Infrastructure   │     │  └─────────────────────┘     │
│  └────────────────────┘     │                              │
│                              │                              │
└──────────────────────────────┴──────────────────────────────┘
```

容器架構中沒有 Guest OS 這一層，這是容器輕量化的根本原因。

#### 詳細比較

| 比較項目 | 虛擬機器 (VM) | 容器 (Container) |
|---------|--------------|-----------------|
| **隔離方式** | 硬體層虛擬化 (Hypervisor) | 作業系統層虛擬化 (Namespace + Cgroup) |
| **Guest OS** | 每台 VM 都有完整的 OS | 共用宿主機核心，無 Guest OS |
| **啟動時間** | 分鐘級（需要開機） | 秒級（啟動一個程序） |
| **映像檔大小** | GB 級（含完整 OS） | MB 級（App + Libs） |
| **資源佔用** | 高（每台都需分配 CPU、記憶體給 OS） | 低（共用核心，按需分配） |
| **隔離強度** | 強（完整硬體隔離） | 較弱（共用核心，靠 Namespace 隔離） |
| **部署密度** | 一台伺服器通常跑數台至數十台 VM | 一台伺服器可跑數十至數百個容器 |
| **適用場景** | 強隔離需求、不同 OS、遺留系統 | 微服務、CI/CD、快速部署、開發環境 |

#### 選擇指引

在實務中，兩者經常搭配使用，例如：在雲端平台上啟動一台 VM（EC2 / GCE），在 VM 中執行 Docker 容器。

**適合使用 VM 的場景：**
- 需要執行不同的作業系統（例如在 Linux 上執行 Windows）
- 需要硬體等級的隔離（多租戶環境、金融法規要求）
- 遺留系統無法容器化

**適合使用容器的場景：**
- 微服務架構
- CI/CD 流水線
- 開發與測試環境
- 需要快速橫向擴展（scale out）的應用

#### 容器的底層技術

容器依賴 Linux 核心的三大機制實現隔離：

**Namespace — 隔離可見範圍**

Namespace 讓每個容器以為自己是系統上唯一的程序，每個容器只能看到屬於自己的資源。

| Namespace 類型 | 隔離內容 | 效果 |
|---------------|---------|------|
| **PID** | 程序 ID | 容器內的程序以為自己是 PID 1 |
| **Network** | 網路介面、IP、路由 | 每個容器有獨立的網路堆疊與 IP 位址 |
| **Mount** | 檔案系統掛載點 | 容器只能看到自己的檔案系統 |
| **UTS** | 主機名稱 | 每個容器可以有自己的 hostname |
| **IPC** | 跨程序通訊 | 容器間的 IPC 相互隔離 |
| **User** | 使用者/群組 ID | 容器內的 root 不等同於宿主機的 root |
| **Cgroup** | Cgroup 層級 | 容器只能看到自己的資源限制 |

在宿主機上執行 `ps aux` 可以看到所有程序，但在容器內執行同樣的命令只能看到容器自己的程序，這就是 PID Namespace 的效果。

**Cgroup（Control Group）— 限制資源用量**

Cgroup 限制每個容器可以使用的系統資源：

| 資源 | 說明 | 範例 |
|------|------|------|
| **CPU** | 限制 CPU 使用量 | 最多使用 1.5 個 CPU 核心 |
| **Memory** | 限制記憶體使用量 | 最多使用 512MB，超過則 OOM Kill |
| **I/O** | 限制磁碟讀寫速度 | 讀寫速度上限 100MB/s |
| **Network** | 限制網路頻寬 | 上行/下行各 100Mbps |

**Union Filesystem — 管理檔案系統**

Union Filesystem（如 OverlayFS）讓多個目錄「疊加」呈現為統一的檔案系統，是 Docker 映像檔分層架構的基礎（詳見 [1.4 節](#14-映像檔與容器的關係)）。

```
┌───────────────────────────────────────────────────────┐
│              容器的三大底層技術                          │
├───────────────────────────────────────────────────────┤
│                                                       │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────┐  │
│  │  Namespace  │  │   Cgroup    │  │ Union FS     │  │
│  │             │  │             │  │              │  │
│  │ 隔離可見性   │  │ 限制資源量   │  │ 管理檔案系統  │  │
│  │ 「看到什麼」 │  │ 「用多少」   │  │ 「存什麼」   │  │
│  └─────────────┘  └─────────────┘  └──────────────┘  │
│         │                │                │           │
│         └────────────────┼────────────────┘           │
│                          ▼                            │
│                 ┌─────────────────┐                   │
│                 │   Linux Kernel  │                   │
│                 └─────────────────┘                   │
│                                                       │
│  這些都是 Linux 核心的原生功能，Docker 的貢獻是將它們   │
│  封裝成簡單易用的工具。                                  │
│                                                       │
└───────────────────────────────────────────────────────┘
```

> **注意**：容器並非完全獨立的虛擬機器，而是透過核心功能隔離出來的程序。在宿主機上執行 `ps aux` 可以看到容器內的程序——它們本質上是普通的 Linux 程序，只是被 Namespace 與 Cgroup 圍起來。

---

### 1.3 Docker 架構

Docker 採用 **Client-Server 架構**，由多個元件協作完成容器管理。理解這些元件有助於排查問題並預測 Docker 的行為。

#### 架構圖

```
┌──────────────────────────────────────────────────────────────────┐
│                         Docker 架構                                │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐                                                │
│  │ Docker CLI   │  終端機輸入的指令                                │
│  │ (Client)     │                                                │
│  │              │                                                │
│  │ docker run   │                                                │
│  │ docker build │                                                │
│  │ docker pull  │                                                │
│  └──────┬───────┘                                                │
│         │  REST API（通常透過 Unix Socket）                        │
│         ▼                                                        │
│  ┌──────────────────────────────────────────┐                    │
│  │          Docker Daemon (dockerd)         │                    │
│  │                                          │                    │
│  │  接收 Client 指令，協調所有 Docker 操作     │                    │
│  │                                          │                    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ │                    │
│  │  │ Image    │ │ Network  │ │ Volume   │ │                    │
│  │  │ 管理     │ │ 管理     │ │ 管理     │ │                    │
│  │  └──────────┘ └──────────┘ └──────────┘ │                    │
│  └──────────────────┬───────────────────────┘                    │
│                     │                                            │
│                     ▼                                            │
│  ┌──────────────────────────────────────────┐                    │
│  │            containerd                    │                    │
│  │                                          │                    │
│  │  高階容器執行時期：管理容器生命週期          │                    │
│  │  (拉取映像檔、建立容器、管理儲存和網路)      │                    │
│  └──────────────────┬───────────────────────┘                    │
│                     │                                            │
│                     ▼                                            │
│  ┌──────────────────────────────────────────┐                    │
│  │              runc                        │                    │
│  │                                          │                    │
│  │  低階容器執行時期（OCI 標準實作）           │                    │
│  │  實際呼叫 Linux Namespace/Cgroup 建立容器  │                    │
│  └──────────────────────────────────────────┘                    │
│                                                                  │
│                     ▲                                            │
│                     │ 拉取映像檔                                  │
│                     ▼                                            │
│  ┌──────────────────────────────────────────┐                    │
│  │          Docker Registry                 │                    │
│  │          (Docker Hub 等)                  │                    │
│  │                                          │                    │
│  │  存放和分發映像檔的倉庫                     │                    │
│  │  nginx, postgres, redis, golang, alpine  │                    │
│  └──────────────────────────────────────────┘                    │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

#### 各元件說明

**Docker Client（用戶端）**

`docker` 命令列工具。本身不執行任何容器管理，負責將指令轉換為 REST API 請求並傳送給 Docker Daemon。

Client 與 Daemon 之間預設透過 Unix Socket（`/var/run/docker.sock`）通訊，通常在同一台機器上，但也可透過 TCP 連線至遠端 Daemon。

```bash
# 查看 Docker Client 與 Server（Daemon）的版本資訊
docker version
```

**Docker Daemon（dockerd）**

Docker 的核心服務，常駐於背景。接收 Client 的 API 請求，負責協調映像檔管理、容器建立、網路設定與儲存管理等所有操作。

執行 `docker run nginx` 時，Daemon 依序執行以下步驟：
1. 檢查本機是否有 `nginx` 映像檔
2. 若無，從 Registry 下載
3. 請 containerd 建立並啟動容器
4. 設定網路（分配 IP、建立 bridge）
5. 掛載 Volume（若有指定）
6. 回報結果給 Client

**containerd（容器執行時期）**

Docker Daemon 底下的容器生命週期管理者，負責拉取與推送映像檔、管理容器的建立與刪除，以及管理映像檔儲存。

containerd 是 CNCF（Cloud Native Computing Foundation）的畢業專案，Kubernetes 也可直接使用 containerd，無需 Docker。

**runc（低階容器執行時期）**

實際建立容器的元件。根據 OCI（Open Container Initiative）規格，呼叫 Linux 核心的 Namespace 與 Cgroup 建立隔離的程序。runc 完成工作後即退出，容器直接由作業系統管理。

**Docker Registry（映像檔倉庫）**

存放與分發映像檔的服務，類似程式碼的 GitHub。

| Registry | 說明 |
|----------|------|
| **Docker Hub** | 官方公開倉庫，有大量社群映像檔（預設） |
| **GitHub Container Registry (ghcr.io)** | GitHub 提供的映像檔倉庫 |
| **Amazon ECR** | AWS 的私有映像檔倉庫 |
| **Google Artifact Registry** | GCP 的映像檔倉庫 |
| **Harbor** | 開源的自建映像檔倉庫，常用於企業內部 |

#### 一個指令的完整旅程

以下追蹤 `docker run nginx` 從輸入到容器啟動的完整流程：

```
┌─────────────────────────────────────────────────────────────┐
│  docker run nginx 的完整流程                                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. 使用者輸入 docker run nginx                              │
│     │                                                       │
│     ▼                                                       │
│  2. Docker CLI 將指令轉換為 REST API 請求                    │
│     POST /containers/create + POST /containers/{id}/start   │
│     │                                                       │
│     ▼                                                       │
│  3. Docker Daemon 收到請求                                   │
│     ├─→ 本機有 nginx 映像檔嗎？                               │
│     │   ├─ 有 → 直接使用                                     │
│     │   └─ 沒有 → 從 Docker Hub 下載                         │
│     │                                                       │
│     ▼                                                       │
│  4. Daemon 請 containerd 建立容器                             │
│     │                                                       │
│     ▼                                                       │
│  5. containerd 請 runc 建立容器程序                           │
│     │                                                       │
│     ▼                                                       │
│  6. runc 呼叫核心 API：                                      │
│     ├─ 建立 Namespace（PID、Network、Mount...）              │
│     ├─ 設定 Cgroup（資源限制）                                │
│     ├─ 掛載 Union Filesystem                                 │
│     └─ 啟動容器內的程序（nginx）                              │
│     │                                                       │
│     ▼                                                       │
│  7. 容器啟動完成，nginx 開始運作                               │
│     runc 退出，容器由核心直接管理                              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

> **macOS 與 Windows 的特殊情況**：Docker 容器需要 Linux 核心。在 macOS 與 Windows 上，Docker Desktop 會在背景執行一台輕量 Linux VM（使用 Apple Hypervisor Framework 或 Hyper-V），Docker Daemon 執行於此 VM 中。因此在 macOS 上，容器實際運行於一台隱藏的 Linux VM 內。

---

### 1.4 映像檔與容器的關係

映像檔（Image）與容器（Container）是 Docker 最核心的兩個概念。

#### 概念對比

| 比喻 | 映像檔（Image） | 容器（Container） |
|------|----------------|-------------------|
| **程式設計** | Class（類別定義） | Instance（實體物件） |
| **烘焙** | 食譜 + 模具 | 依食譜做出來的蛋糕 |
| **建築** | 建築藍圖 | 蓋好的房子 |
| **作業系統** | ISO 安裝映像檔 | 安裝並執行中的系統 |

**關鍵區別**：映像檔是靜態的、唯讀的「定義」；容器是動態的、可讀寫的「執行實體」。一個映像檔可以產生無限多個容器，如同一個 Class 可以 `new()` 出無限多個 Instance。

#### 映像檔的分層架構（Layer）

Docker 映像檔不是單一巨大的檔案，而是由多個**唯讀的 Layer** 堆疊而成。每一個 Dockerfile 指令（`FROM`、`RUN`、`COPY` 等）都會建立一個新的 Layer。

```
┌────────────────────────────────────────────────────────────┐
│                映像檔的分層結構                               │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Dockerfile 指令                    對應的 Layer             │
│                                                            │
│  FROM golang:1.24-alpine    →    ┌──────────────────────┐  │
│                                  │ Layer 1: 基礎映像檔    │  │
│                                  │ (Alpine + Go 工具鏈)   │  │
│                                  │ 大小: ~270MB           │  │
│  COPY go.mod go.sum ./      →    ├──────────────────────┤  │
│                                  │ Layer 2: 相依定義檔    │  │
│                                  │ 大小: ~5KB             │  │
│  RUN go mod download        →    ├──────────────────────┤  │
│                                  │ Layer 3: 下載的套件    │  │
│                                  │ 大小: ~50MB            │  │
│  COPY . .                   →    ├──────────────────────┤  │
│                                  │ Layer 4: 原始碼        │  │
│                                  │ 大小: ~1MB             │  │
│  RUN go build -o server .   →    ├──────────────────────┤  │
│                                  │ Layer 5: 編譯產物      │  │
│                                  │ 大小: ~15MB            │  │
│                                  └──────────────────────┘  │
│                                                            │
│  所有 Layer 都是唯讀的（Read-Only）                          │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

分層架構帶來三個重要優勢：

**建置快取**：若某一層的指令與輸入沒有改變，Docker 直接使用上次的快取，無需重新執行。僅修改一行 Go 原始碼時，Layer 1–3 可直接使用快取。

**磁碟共用**：若 10 個映像檔都基於 `alpine:3.21`，該層在磁碟上只需存一份。`docker images` 顯示的總大小可能大於實際佔用空間，原因正是共用的 Layer 只計算一次。

**傳輸效率**：推送或拉取映像檔時，只需傳輸對方尚未擁有的 Layer。

#### 容器的可寫層（Container Layer）

從映像檔建立容器時，Docker 在所有唯讀 Layer 之上加入一個**可讀寫的容器層**：

```
┌────────────────────────────────────────────────────────────┐
│                Container = Image + 可寫層                    │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  ┌──────────────────────────────────────────────┐          │
│  │  Container Layer（可讀寫）                     │          │
│  │                                              │          │
│  │  容器執行期間所有的修改都寫在這一層：            │          │
│  │  - 新建立的檔案                               │          │
│  │  - 修改的設定檔                               │          │
│  │  - 應用程式寫入的日誌                          │          │
│  │  - 暫存資料                                   │          │
│  │                                              │          │
│  │  ⚠ 容器刪除時，這一層隨之消失！               │          │
│  ├──────────────────────────────────────────────┤          │
│  │  Layer 5: go build（唯讀）                    │          │
│  ├──────────────────────────────────────────────┤          │
│  │  Layer 4: COPY . .（唯讀）                    │   Image  │
│  ├──────────────────────────────────────────────┤  Layers  │
│  │  Layer 3: go mod download（唯讀）             │ (共用的)  │
│  ├──────────────────────────────────────────────┤          │
│  │  Layer 2: COPY go.mod（唯讀）                 │          │
│  ├──────────────────────────────────────────────┤          │
│  │  Layer 1: golang:1.24-alpine（唯讀）          │          │
│  └──────────────────────────────────────────────┘          │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

**Copy-on-Write（CoW）機制**：當容器需要修改映像檔 Layer 中的檔案時，Docker 不直接修改唯讀 Layer，而是先將該檔案複製至可寫層，再於可寫層修改。此機制確保映像檔的不可變性——無論多少容器正在執行，底層的映像檔 Layer 永遠不會被修改。

#### 多容器共用映像檔

```
┌────────────────────────────────────────────────────────────┐
│                 多容器共用同一個映像檔                        │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Container A          Container B          Container C     │
│  ┌──────────┐        ┌──────────┐        ┌──────────┐     │
│  │ 可寫層 A  │        │ 可寫層 B  │        │ 可寫層 C  │     │
│  │ (日誌、   │        │ (不同的   │        │ (各自獨   │     │
│  │  暫存等)  │        │  資料)    │        │  立的)    │     │
│  └─────┬────┘        └─────┬────┘        └─────┬────┘     │
│        │                   │                   │           │
│        └───────────────────┼───────────────────┘           │
│                            ▼                               │
│              ┌──────────────────────┐                      │
│              │  共用的映像檔 Layer   │                      │
│              │  (唯讀，只存一份)      │                      │
│              │                      │                      │
│              │  Layer 5: app        │                      │
│              │  Layer 4: deps       │                      │
│              │  Layer 3: runtime    │                      │
│              │  Layer 2: libs       │                      │
│              │  Layer 1: base OS    │                      │
│              └──────────────────────┘                      │
│                                                            │
│  三個容器共用同一份映像檔（~300MB），各自的可寫層僅數 MB      │
│  總磁碟佔用 ≈ 300MB + 數 MB × 3，而非 300MB × 3            │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

#### 映像檔命名與 Tag

每個映像檔可以有多個 Tag，用來區分不同版本。命名格式如下：

```
┌────────────────────────────────────────────────────────────┐
│                    映像檔命名格式                             │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  完整格式：                                                  │
│  [registry/][username/]repository[:tag]                     │
│                                                            │
│  範例：                                                     │
│  docker.io/library/nginx:1.27-alpine                       │
│  ├────────┘ ├─────┘ ├───┘ ├─────────┘                      │
│  │          │       │     └─ Tag（版本標記）                  │
│  │          │       └─ Repository（映像檔名稱）              │
│  │          └─ Username（官方映像檔為 library，通常省略）      │
│  └─ Registry（預設 docker.io，通常省略）                     │
│                                                            │
│  常見 Tag 慣例：                                             │
│  nginx:latest        → 最新版（不建議正式環境使用）           │
│  nginx:1.27          → 主要版本                              │
│  nginx:1.27.3        → 精確版本（正式環境建議使用）           │
│  nginx:1.27-alpine   → 基於 Alpine 的精簡版                 │
│  nginx:1.27-bookworm → 基於 Debian Bookworm 的版本          │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

> **重要**：`latest` 只是一個普通的 Tag 名稱，不保證是最新版，且會隨映像檔更新而指向不同版本。正式環境中**務必使用明確的版本號**（如 `nginx:1.27.3`）。

#### 概念總結

| 概念 | 映像檔（Image） | 容器（Container） |
|------|----------------|-------------------|
| **性質** | 唯讀模板，由多個 Layer 組成 | 可讀寫的執行實體，映像檔 + 可寫層 |
| **生命週期** | 建置後不可變（Immutable） | 可建立、啟動、暫停、停止、刪除 |
| **儲存** | 分層結構，可跨映像檔共用 | 可寫層存放執行期間的變更 |
| **可重複性** | 同一 Dockerfile 永遠建出相同的映像檔 | 每個容器的可寫層內容可能不同 |
| **數量關係** | 一個映像檔可產生多個容器 | 一個容器只對應一個映像檔 |
| **資料持久性** | 除非主動刪除，否則永久存在 | 可寫層隨容器刪除而消失（需 Volume 持久化） |

---

### 1.5 安裝與驗證

**macOS（使用 Docker Desktop 或 OrbStack）**

```bash
# 驗證安裝
docker version

# 查看 Docker 系統資訊
docker info

# 執行測試容器
docker run hello-world
```

**預期輸出：**

```
Hello from Docker!
This message shows that your installation appears to be working correctly.
...
```

出現此訊息表示 Docker 已正確安裝。

`docker run hello-world` 背後的完整流程：

```
┌──────────────────────────────────────────────────────┐
│  docker run hello-world 的完整流程                     │
├──────────────────────────────────────────────────────┤
│                                                      │
│  1. Docker Client 傳送指令給 Docker Daemon            │
│                    │                                 │
│                    ▼                                 │
│  2. Daemon 在本機尋找 hello-world 映像檔               │
│                    │                                 │
│                    ▼ （找不到）                        │
│  3. 從 Docker Hub 下載 hello-world 映像檔              │
│                    │                                 │
│                    ▼                                 │
│  4. 用映像檔建立一個新的容器                           │
│                    │                                 │
│                    ▼                                 │
│  5. 容器執行程式，輸出訊息                             │
│                    │                                 │
│                    ▼                                 │
│  6. 程式結束，容器停止（不會自動刪除）                  │
│                                                      │
└──────────────────────────────────────────────────────┘
```

---

### 1.6 映像檔基本操作

#### 搜尋映像檔

```bash
# 從 Docker Hub 搜尋映像檔
docker search nginx

# 建議直接前往 Docker Hub 網站搜尋，資訊更完整
# https://hub.docker.com
```

#### 下載映像檔

```bash
# 下載最新版（:latest 標籤）
docker pull nginx

# 下載特定版本
docker pull nginx:1.27

# 下載特定平台版本
docker pull --platform linux/amd64 nginx:1.27

# 下載 Alpine 精簡版（體積更小，建議優先使用）
docker pull nginx:1.27-alpine
```

> **注意**：Tag 是映像檔的版本標記。`latest` 不保證是最新版，僅是預設標籤。正式環境務必指定明確版本號。

#### 列出本機映像檔

```bash
docker images

# 輸出範例：
# REPOSITORY   TAG           IMAGE ID       CREATED       SIZE
# nginx        1.27-alpine   a2bd6dc6e5e6   2 weeks ago   43.3MB
# nginx        1.27          39286ab8a5e1   2 weeks ago   192MB
# hello-world  latest        d2c94e258dcb   9 months ago  13.3kB
```

`nginx:1.27-alpine`（43MB）與 `nginx:1.27`（192MB）的大小差異，體現了 Alpine 基礎映像檔的輕量優勢。

#### 刪除映像檔

```bash
# 刪除指定映像檔
docker rmi nginx:1.27

# 強制刪除（即使有已停止的容器正在使用）
docker rmi -f nginx:1.27

# 刪除所有未使用的映像檔（dangling images）
docker image prune

# 刪除所有未被容器使用的映像檔
docker image prune -a
```

#### 查看映像檔資訊

```bash
# 查看映像檔詳細資訊（層數、環境變數等）
docker inspect nginx:1.27-alpine

# 查看映像檔的建置歷史（每一層的指令）
docker history nginx:1.27-alpine
```

---

### 1.7 容器基本操作

#### 建立並啟動容器

```bash
# 基本語法
docker run [OPTIONS] IMAGE [COMMAND] [ARG...]

# 前景模式執行（Ctrl+C 可停止）
docker run nginx:1.27-alpine

# 背景模式執行（-d = detach）
docker run -d nginx:1.27-alpine

# 指定容器名稱（--name）
docker run -d --name my-nginx nginx:1.27-alpine

# 執行後自動刪除容器（--rm，適合一次性任務）
docker run --rm nginx:1.27-alpine nginx -v
```

#### 常用 `docker run` 參數

| 參數 | 說明 | 範例 |
|------|------|------|
| `-d` | 背景執行（detach） | `docker run -d nginx` |
| `--name` | 命名容器 | `docker run --name web nginx` |
| `-p` | Port mapping（host:container） | `docker run -p 8080:80 nginx` |
| `-v` | 掛載 Volume | `docker run -v /data:/app/data nginx` |
| `-e` | 設定環境變數 | `docker run -e DB_HOST=db nginx` |
| `--rm` | 停止後自動刪除 | `docker run --rm nginx` |
| `-it` | 互動式終端 | `docker run -it ubuntu bash` |
| `--restart` | 重啟策略 | `docker run --restart=unless-stopped nginx` |
| `--network` | 指定網路 | `docker run --network=mynet nginx` |
| `--memory` | 記憶體限制 | `docker run --memory=512m nginx` |
| `--cpus` | CPU 限制 | `docker run --cpus=1.5 nginx` |

#### 列出容器

```bash
# 列出執行中的容器
docker ps

# 列出所有容器（包含已停止的）
docker ps -a

# 只顯示容器 ID
docker ps -q

# 輸出範例：
# CONTAINER ID   IMAGE               STATUS          PORTS                  NAMES
# a1b2c3d4e5f6   nginx:1.27-alpine   Up 10 minutes   0.0.0.0:8080->80/tcp   my-nginx
```

#### 容器生命週期管理

```bash
# 停止容器（發送 SIGTERM，10 秒後 SIGKILL）
docker stop my-nginx

# 指定等待時間（秒）
docker stop -t 30 my-nginx

# 啟動已停止的容器
docker start my-nginx

# 重啟容器
docker restart my-nginx

# 強制停止容器（發送 SIGKILL，不建議使用）
docker kill my-nginx

# 刪除已停止的容器
docker rm my-nginx

# 強制刪除執行中的容器
docker rm -f my-nginx

# 刪除所有已停止的容器
docker container prune
```

**容器生命週期圖：**

```
┌─────────────────────────────────────────────────────────────┐
│                    容器生命週期                               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  docker create         docker start         docker stop     │
│  ┌──────────┐  ───→   ┌──────────┐  ───→   ┌──────────┐   │
│  │ Created  │         │ Running  │         │ Stopped  │    │
│  └──────────┘  ←───   └──────────┘  ←───   └──────────┘   │
│                        docker stop          docker start    │
│                                                             │
│  docker run = docker create + docker start                  │
│                                                             │
│              docker rm                                      │
│  任何狀態 ──────────→ 刪除                                   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

### 1.8 Port Mapping

容器預設處於隔離狀態，外部無法直接存取容器內的服務。**Port Mapping** 將宿主機的連接埠對映至容器的連接埠，使外部得以存取。

```
┌─────────────────────────────────────────────────┐
│                    Port Mapping                  │
├─────────────────────────────────────────────────┤
│                                                 │
│  宿主機 (Host)                                   │
│  ┌─────────────────────────────────────┐        │
│  │                                     │        │
│  │   瀏覽器 → http://localhost:8080    │        │
│  │                    │                │        │
│  │                    ▼                │        │
│  │          Host Port 8080             │        │
│  │                    │                │        │
│  │         ┌──────────┼──────────┐     │        │
│  │         │ Container│          │     │        │
│  │         │          ▼          │     │        │
│  │         │  Container Port 80  │     │        │
│  │         │          │          │     │        │
│  │         │     ┌────▼────┐     │     │        │
│  │         │     │  Nginx  │     │     │        │
│  │         │     └─────────┘     │     │        │
│  │         └─────────────────────┘     │        │
│  │                                     │        │
│  └─────────────────────────────────────┘        │
│                                                 │
└─────────────────────────────────────────────────┘
```

```bash
# 語法：-p <宿主機連接埠>:<容器連接埠>
docker run -d --name web -p 8080:80 nginx:1.27-alpine

# 對映多個連接埠
docker run -d -p 8080:80 -p 8443:443 nginx:1.27-alpine

# 指定綁定的 IP（預設 0.0.0.0，即所有介面）
docker run -d -p 127.0.0.1:8080:80 nginx:1.27-alpine

# 隨機分配宿主機連接埠（用 docker ps 或 docker port 查看）
docker run -d -P nginx:1.27-alpine

# 查看連接埠對映
docker port web
```

```bash
# 測試
curl http://localhost:8080
# 應看到 Nginx 歡迎頁面的 HTML
```

---

### 1.9 Volume

容器內的資料在容器刪除後即消失。**Volume** 將宿主機目錄或 Docker 管理的儲存掛載至容器內，實現資料持久化。

**三種掛載方式：**

```
┌───────────────────────────────────────────────────────────────┐
│                     Volume 掛載方式                             │
├───────────────────────────────────────────────────────────────┤
│                                                               │
│  1. Bind Mount（綁定掛載）                                     │
│     宿主機目錄 → 容器目錄                                       │
│     -v /host/path:/container/path                             │
│     用途：開發時同步原始碼                                      │
│                                                               │
│  2. Named Volume（具名資料卷）                                  │
│     Docker 管理的儲存 → 容器目錄                                │
│     -v mydata:/container/path                                 │
│     用途：資料庫等需要持久化的資料                               │
│                                                               │
│  3. tmpfs Mount（記憶體掛載）                                   │
│     記憶體 → 容器目錄                                           │
│     --tmpfs /container/path                                   │
│     用途：暫存敏感資料，容器停止即消失                           │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

```bash
# Bind Mount：將宿主機目錄掛載至容器
docker run -d --name web \
  -p 8080:80 \
  -v $(pwd)/html:/usr/share/nginx/html:ro \
  nginx:1.27-alpine
# :ro 表示容器內唯讀（read-only）

# Named Volume：建立 Docker 管理的資料卷
docker volume create pgdata
docker run -d --name db \
  -v pgdata:/var/lib/postgresql/data \
  postgres:17-alpine

# 列出所有 Volume
docker volume ls

# 查看 Volume 詳細資訊
docker volume inspect pgdata

# 刪除未使用的 Volume
docker volume prune
```

> **注意**：Bind Mount 使用絕對路徑（或 `$(pwd)`），Named Volume 只使用名稱。若路徑以 `/` 或 `./` 開頭，Docker 視為 Bind Mount。

---

### 1.10 容器除錯

#### 查看日誌

```bash
# 查看容器日誌
docker logs my-nginx

# 持續追蹤日誌（類似 tail -f）
docker logs -f my-nginx

# 查看最近 100 行日誌
docker logs --tail 100 my-nginx

# 附帶時間戳記
docker logs -t my-nginx

# 查看特定時間之後的日誌
docker logs --since 2024-01-01T00:00:00 my-nginx

# 組合使用
docker logs -f --tail 50 -t my-nginx
```

#### 進入容器內部

```bash
# 在容器內開啟 shell
docker exec -it my-nginx sh
# -i: 保持標準輸入開啟（interactive）
# -t: 分配偽終端（tty）

# 若容器有 bash（Alpine 通常只有 sh）
docker exec -it my-nginx bash

# 執行單一命令（不進入互動模式）
docker exec my-nginx cat /etc/nginx/nginx.conf

# 以 root 身份進入
docker exec -it -u root my-nginx sh
```

#### 查看資源使用狀況

```bash
# 即時顯示所有容器的資源使用
docker stats

# 只看特定容器
docker stats my-nginx

# 輸出範例：
# CONTAINER ID   NAME       CPU %   MEM USAGE / LIMIT     MEM %   NET I/O
# a1b2c3d4e5f6   my-nginx   0.00%   3.441MiB / 7.667GiB   0.04%   1.45kB / 0B
```

#### 查看容器程序

```bash
docker top my-nginx
```

#### 複製檔案

```bash
# 從宿主機複製至容器
docker cp ./index.html my-nginx:/usr/share/nginx/html/

# 從容器複製至宿主機
docker cp my-nginx:/etc/nginx/nginx.conf ./nginx.conf
```

---

### 1.11 練習 1：執行你的第一個容器

**目標**：啟動一個 Nginx 容器，掛載自訂頁面，並透過瀏覽器存取。

```bash
# 1. 建立工作目錄
mkdir -p ~/docker-lab/lab1
cd ~/docker-lab/lab1

# 2. 建立自訂 HTML 頁面
cat << 'EOF' > index.html
<!DOCTYPE html>
<html>
<head>
    <title>Docker Lab 1</title>
</head>
<body>
    <h1>Hello from Docker!</h1>
    <p>This page is served from a Docker container.</p>
    <p>Nginx is running inside the container, serving files from the host machine.</p>
</body>
</html>
EOF

# 3. 啟動 Nginx 容器
docker run -d \
  --name lab1-nginx \
  -p 8080:80 \
  -v $(pwd):/usr/share/nginx/html:ro \
  nginx:1.27-alpine

# 4. 確認容器正在執行
docker ps

# 5. 測試
curl http://localhost:8080

# 6. 查看容器日誌
docker logs lab1-nginx

# 7. 進入容器內部檢視
docker exec -it lab1-nginx sh
# 在容器內執行：
#   ls /usr/share/nginx/html/
#   cat /etc/nginx/nginx.conf
#   exit

# 8. 清理
docker stop lab1-nginx
docker rm lab1-nginx
```

**延伸練習：**

1. 修改 `index.html` 的內容，不重啟容器，直接重新整理瀏覽器查看變化（思考：為什麼修改會立即生效？）
2. 移除 `-v` 中的 `:ro` 限制，在容器內建立新檔案，觀察宿主機是否能看到該檔案。

---

## Part 2：Dockerfile 深入實作

### 2.1 什麼是 Dockerfile？

Dockerfile 是一個**純文字檔**，包含一系列指令，定義 Docker 如何建置（build）一個映像檔。

```
┌────────────────────────────────────────────────┐
│                                                │
│  Dockerfile ─── docker build ──→ Docker Image  │
│  （建置腳本）      （建置過程）      （成品）      │
│                                                │
└────────────────────────────────────────────────┘
```

---

### 2.2 Dockerfile 基礎指令

| 指令 | 說明 | 範例 |
|------|------|------|
| `FROM` | 指定基礎映像檔（必須是第一行） | `FROM golang:1.24-alpine` |
| `WORKDIR` | 設定工作目錄（不存在會自動建立） | `WORKDIR /app` |
| `COPY` | 複製檔案/目錄至映像檔 | `COPY . .` |
| `ADD` | 類似 COPY，但可解壓縮 tar 與下載 URL | `ADD app.tar.gz /app` |
| `RUN` | 建置階段執行命令 | `RUN go build -o server .` |
| `ENV` | 設定環境變數 | `ENV GIN_MODE=release` |
| `ARG` | 定義建置時的參數 | `ARG GO_VERSION=1.24` |
| `EXPOSE` | 宣告容器監聽的連接埠（僅文件用途） | `EXPOSE 8080` |
| `CMD` | 容器啟動時的預設命令 | `CMD ["./server"]` |
| `ENTRYPOINT` | 容器啟動時的入口命令 | `ENTRYPOINT ["./server"]` |
| `USER` | 指定執行命令的使用者 | `USER nonroot` |
| `HEALTHCHECK` | 定義健康檢查 | `HEALTHCHECK CMD curl -f http://localhost/` |
| `LABEL` | 新增映像檔中繼資料 | `LABEL version="1.0"` |

#### FROM — 基礎映像檔

```dockerfile
# 使用特定版本（建議）
FROM golang:1.24-alpine

# 使用最小映像檔（僅有基本工具）
FROM alpine:3.21

# 從空白開始（適合靜態編譯的 Go 程式）
FROM scratch
```

**常見基礎映像檔比較：**

| 基礎映像檔 | 大小 | 特點 |
|-----------|------|------|
| `ubuntu:24.04` | ~78MB | 完整套件管理，偵錯方便 |
| `debian:bookworm-slim` | ~74MB | 精簡版 Debian |
| `alpine:3.21` | ~7MB | 極小，使用 musl libc + apk |
| `scratch` | 0MB | 空映像檔，適合靜態編譯的二進位 |
| `gcr.io/distroless/static` | ~2MB | Google Distroless，無 shell |

#### COPY vs ADD

```dockerfile
# COPY — 單純複製（建議使用，行為可預期）
COPY go.mod go.sum ./
COPY . .

# ADD — 具備額外功能，但行為較不直覺
ADD app.tar.gz /app/          # 自動解壓縮
ADD https://example.com/f /f  # 可下載 URL（不建議，應改用 RUN curl）
```

> **最佳實踐**：除非需要自動解壓縮 tar 檔案，否則一律使用 `COPY`。

#### RUN — 建置階段執行命令

```dockerfile
# 每個 RUN 會建立一個新的 Layer
RUN apk add --no-cache git

# 合併多個命令以減少 Layer 數量
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*
```

> **最佳實踐**：以 `&&` 串聯相關命令，減少不必要的 Layer。每個 Layer 都會增加映像檔體積。

#### CMD vs ENTRYPOINT

```dockerfile
# CMD — 定義預設命令，可被 docker run 的參數覆蓋
CMD ["./server"]
# docker run myapp              → 執行 ./server
# docker run myapp /bin/sh      → 執行 /bin/sh（CMD 被覆蓋）

# ENTRYPOINT — 定義固定入口，不會被覆蓋
ENTRYPOINT ["./server"]
# docker run myapp              → 執行 ./server
# docker run myapp --port 9090  → 執行 ./server --port 9090（參數被附加）

# 最佳搭配：ENTRYPOINT 定義程式，CMD 定義預設參數
ENTRYPOINT ["./server"]
CMD ["--port", "8080"]
# docker run myapp                     → ./server --port 8080
# docker run myapp --port 9090         → ./server --port 9090
```

**Shell Form vs Exec Form：**

```dockerfile
# Exec Form（建議）— 直接執行，程式成為 PID 1
CMD ["./server", "--port", "8080"]

# Shell Form — 透過 /bin/sh -c 執行，server 不是 PID 1
CMD ./server --port 8080
# 等同於：/bin/sh -c "./server --port 8080"
```

> **重要**：使用 Exec Form，讓程式成為 PID 1，才能正確接收 SIGTERM 訊號（`docker stop` 時）。Shell Form 會讓 shell 成為 PID 1，應用程式可能無法優雅關閉（Graceful Shutdown）。

#### EXPOSE — 宣告連接埠

```dockerfile
# EXPOSE 只是文件宣告，不會實際開放連接埠
# 實際的連接埠對映需在 docker run 時以 -p 指定
EXPOSE 8080
```

#### HEALTHCHECK — 健康檢查

```dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

| 參數 | 說明 | 預設值 |
|------|------|--------|
| `--interval` | 檢查間隔 | 30s |
| `--timeout` | 逾時時間 | 30s |
| `--start-period` | 啟動寬限期 | 0s |
| `--retries` | 失敗重試次數 | 3 |

---

### 2.3 .dockerignore

`.dockerignore` 類似 `.gitignore`，用來排除不需要傳送至 Docker Daemon 的檔案，可**加速建置**並**避免洩漏敏感資訊**。

```gitignore
# .dockerignore

# 版本控制
.git
.gitignore

# IDE 設定
.idea/
.vscode/
*.swp
*.swo

# Go 建置產物
/bin/
/dist/
vendor/

# Docker 相關
Dockerfile
docker-compose.yml
.dockerignore

# 文件
*.md
docs/
LICENSE

# 測試
*_test.go
testdata/

# 敏感檔案
.env
.env.*
*.pem
*.key
```

> **為何重要？** `docker build` 會將整個 build context 傳送給 Docker Daemon。若目錄中有 1GB 的 `vendor/` 或 `.git/`，建置將會非常緩慢。

---

### 2.4 為 Go 應用撰寫 Dockerfile

以下是一個簡單的 Go HTTP 伺服器範例：

**main.go**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		fmt.Fprintf(w, "Hello from Go! Hostname: %s\n", hostname)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

#### 基本版 Dockerfile（單階段建置）

```dockerfile
FROM golang:1.24-alpine

WORKDIR /app

# 先複製 go.mod, go.sum 以利用 Layer Cache
COPY go.mod go.sum ./
RUN go mod download

# 複製原始碼並編譯
COPY . .
RUN go build -o server .

EXPOSE 8080
CMD ["./server"]
```

**問題**：最終映像檔包含整個 Go 編譯工具鏈（~300MB+），但執行時只需要編譯好的二進位檔。

---

### 2.5 Multi-stage Build

Multi-stage Build（多階段建置）是 Docker 的重要功能，讓你在**建置階段**使用完整工具鏈，但**最終映像檔**只包含執行所需的最小內容。

```dockerfile
# ==========================================
# Stage 1: Build（建置階段）
# ==========================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 先複製相依套件定義，利用 Layer Cache
COPY go.mod go.sum ./
RUN go mod download

# 複製原始碼
COPY . .

# 編譯：靜態連結，停用 CGO
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o server .

# ==========================================
# Stage 2: Runtime（執行階段）
# ==========================================
FROM alpine:3.21

# 安裝 CA 憑證（若需要對外 HTTPS 請求）與時區資料
RUN apk add --no-cache ca-certificates tzdata

# 建立非 root 使用者
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# 從 builder 階段複製編譯好的二進位檔
COPY --from=builder /app/server .

# 使用非 root 使用者執行
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./server"]
```

**映像檔大小比較：**

| 方式 | 映像檔大小 |
|------|-----------|
| 單階段 `golang:1.24-alpine` | ~300MB |
| 多階段 `alpine:3.21` | ~15MB |
| 多階段 `scratch`（最精簡） | ~8MB |

#### 使用 scratch 的極致精簡版

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o server .

# scratch = 完全空白的映像檔
FROM scratch

# 從 builder 複製 CA 憑證
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 複製時區資料
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /app/server /server

EXPOSE 8080
ENTRYPOINT ["/server"]
```

> **注意**：`scratch` 沒有 shell，無法用 `docker exec` 進入容器偵錯，也無法使用 `wget`/`curl` 進行 HEALTHCHECK。適合已穩定上線的正式環境。

#### 建置與執行

```bash
# 建置映像檔
docker build -t my-go-app:v1 .

# 查看映像檔大小
docker images my-go-app

# 執行容器
docker run -d --name my-app -p 8080:8080 my-go-app:v1

# 測試
curl http://localhost:8080
curl http://localhost:8080/health

# 查看容器健康狀態
docker inspect --format='{{.State.Health.Status}}' my-app
```

---

### 2.6 映像檔最佳實踐

#### 1. 善用 Layer Cache

Docker 會快取每一個 Layer。Layer 快取失效後，其後所有 Layer 都必須重新建置。因此，**將不常變動的指令放在 Dockerfile 前段**，可最大化快取命中率。

```dockerfile
# 不建議：任何原始碼變更都會導致 go mod download 重新執行
COPY . .
RUN go mod download
RUN go build -o server .

# 建議：只有 go.mod/go.sum 改變時才重新下載相依套件
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .
```

```
┌───────────────────────────────────────────────────────┐
│              Layer Cache 最佳順序                       │
├───────────────────────────────────────────────────────┤
│                                                       │
│  FROM golang:1.24-alpine     ← 幾乎不變               │
│  RUN apk add ...             ← 很少變                 │
│  COPY go.mod go.sum ./       ← 偶爾變（加新套件）      │
│  RUN go mod download         ← 跟上一層連動            │
│  COPY . .                    ← 常常變（修改原始碼）     │
│  RUN go build ...            ← 跟上一層連動            │
│                                                       │
│  不常變 ──────────────────────────────── 常常變         │
│  （放上面）                              （放下面）      │
│                                                       │
└───────────────────────────────────────────────────────┘
```

#### 2. 減少 Layer 數量

```dockerfile
# 不建議：3 個 Layer
RUN apk add --no-cache git
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata

# 建議：1 個 Layer
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata
```

#### 3. 使用非 root 使用者

以 root 執行容器會帶來安全風險。應建立專用的非特權使用者：

```dockerfile
# 建立專用使用者
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 切換使用者
USER appuser
```

#### 4. 安裝套件時不保留快取

```dockerfile
# 建議：不留下套件快取
RUN apk add --no-cache ca-certificates

# 不建議：套件快取會增加映像檔體積
RUN apk add ca-certificates
```

#### 5. 建置參數化

```dockerfile
ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-alpine AS builder

ARG APP_VERSION=dev
RUN go build -ldflags="-X main.version=${APP_VERSION}" -o server .
```

```bash
# 建置時傳入參數
docker build --build-arg GO_VERSION=1.24 --build-arg APP_VERSION=v1.2.3 -t my-app:v1.2.3 .
```

---

### 2.7 練習 2：建置你的 Go 應用映像檔

**目標**：為一個簡單的 Go HTTP 伺服器建置多階段映像檔。

```bash
# 1. 建立專案目錄
mkdir -p ~/docker-lab/lab2
cd ~/docker-lab/lab2

# 2. 初始化 Go 模組
go mod init lab2

# 3. 建立 main.go
cat << 'GOEOF' > main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var startTime = time.Now()

type InfoResponse struct {
	Message    string `json:"message"`
	Hostname   string `json:"hostname"`
	GoVersion  string `json:"go_version"`
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	Uptime     string `json:"uptime"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/info", handleInfo)
	http.HandleFunc("/health", handleHealth)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Hello from Go in Docker! Hostname: %s\n", hostname)
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	info := InfoResponse{
		Message:   "Hello from Docker Workshop!",
		Hostname:  hostname,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Uptime:    time.Since(startTime).Round(time.Second).String(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
GOEOF

# 4. 建立 .dockerignore
cat << 'EOF' > .dockerignore
.git
*.md
Dockerfile
.dockerignore
EOF

# 5. 建立多階段 Dockerfile
cat << 'DEOF' > Dockerfile
# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# Stage 2: Runtime
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/server .
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
ENTRYPOINT ["./server"]
DEOF

# 6. 建置映像檔
docker build -t lab2-app:v1 .

# 7. 查看映像檔大小
docker images lab2-app

# 8. 執行容器
docker run -d --name lab2 -p 8080:8080 lab2-app:v1

# 9. 測試各端點
curl http://localhost:8080/
curl http://localhost:8080/info | jq .
curl http://localhost:8080/health

# 10. 清理
docker stop lab2
docker rm lab2
```

**延伸練習：**

1. 將 `FROM alpine:3.21` 改為 `FROM scratch`，重新建置並比較映像檔大小差異。
2. 將 `COPY go.mod` 與 `COPY . .` 合併為單一 `COPY`，修改一行程式碼後重新建置，觀察建置時間差異。

---

## Part 3：Docker Compose 多容器編排

### 3.1 為什麼需要 Docker Compose？

在實際專案中，一個服務通常需要多個容器協同運作：

```
┌──────────────────────────────────────────────────────┐
│                    典型的應用架構                       │
├──────────────────────────────────────────────────────┤
│                                                      │
│                   ┌──────────┐                       │
│                   │  Client  │                       │
│                   └────┬─────┘                       │
│                        │                             │
│                        ▼                             │
│                 ┌──────────────┐                     │
│                 │   Go API     │  ← Container 1     │
│                 │   (Port 8080)│                     │
│                 └──┬───────┬──┘                      │
│                    │       │                         │
│              ┌─────▼──┐  ┌─▼────────┐               │
│              │PostgreSQL│  │  Redis   │               │
│              │(Port 5432)│  │(Port 6379)│              │
│              │Container 2│  │Container 3│              │
│              └──────────┘  └──────────┘              │
│                                                      │
└──────────────────────────────────────────────────────┘
```

若不使用 Compose，需要手動執行一系列指令：

```bash
# 建立網路
docker network create myapp

# 啟動 PostgreSQL
docker run -d --name db --network myapp \
  -e POSTGRES_PASSWORD=secret \
  -v pgdata:/var/lib/postgresql/data \
  postgres:17-alpine

# 啟動 Redis
docker run -d --name cache --network myapp \
  redis:7-alpine

# 啟動 Go API
docker run -d --name api --network myapp \
  -p 8080:8080 \
  -e DATABASE_URL=postgres://postgres:secret@db:5432/mydb \
  -e REDIS_URL=redis://cache:6379 \
  my-go-app:v1
```

每次啟動都需要執行這些指令，容易出錯且難以維護。**Docker Compose** 讓你用一個 YAML 檔案定義所有服務，以單一指令啟動整個應用。

---

### 3.2 基礎語法

```yaml
# docker-compose.yml（或 compose.yml，兩者皆可）

services:
  # 服務名稱（同時也是容器在網路中的 hostname）
  service-name:
    image: image:tag           # 使用現有映像檔
    # 或
    build: ./path              # 從 Dockerfile 建置
    ports:
      - "host:container"       # 連接埠對映
    volumes:
      - name:/path             # 資料掛載
    environment:
      KEY: value               # 環境變數
    depends_on:
      - other-service          # 相依關係
    networks:
      - network-name           # 所屬網路

# 資料卷定義
volumes:
  name:

# 網路定義（通常不需要，Compose 會自動建立）
networks:
  network-name:
```

---

### 3.3 服務定義詳解

#### 使用映像檔 vs 建置

```yaml
services:
  # 方式 1：直接使用映像檔
  redis:
    image: redis:7-alpine

  # 方式 2：從 Dockerfile 建置
  api:
    build: .
    # 等同於 docker build .

  # 方式 3：指定 Dockerfile 路徑與 build context
  api:
    build:
      context: .
      dockerfile: Dockerfile.prod
      args:
        GO_VERSION: "1.24"
        APP_VERSION: "v1.0.0"
    image: my-go-app:v1   # 建置後標記映像檔名稱
```

#### depends_on — 服務相依

```yaml
services:
  api:
    build: .
    depends_on:
      db:
        condition: service_healthy  # 等待 db 健康檢查通過
      redis:
        condition: service_started  # 只等待 redis 啟動

  db:
    image: postgres:17-alpine
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
```

> **注意**：`depends_on` 預設只等待容器啟動，不等待服務就緒。搭配 `condition: service_healthy` 與 `healthcheck` 才能確保相依服務真正可用。

#### restart — 重啟策略

```yaml
services:
  api:
    restart: unless-stopped
    # no            — 不自動重啟（預設）
    # always        — 總是重啟
    # on-failure    — 非正常退出時重啟
    # unless-stopped — 除非手動停止，否則重啟
```

#### 資源限制

```yaml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: "1.0"
          memory: 512M
        reservations:
          cpus: "0.25"
          memory: 128M
```

---

### 3.4 Network

Docker Compose 自動建立一個預設網路（`<專案名>_default`），所有服務加入此網路後，**可以用服務名稱作為 hostname 互相存取**。

```yaml
services:
  api:
    build: .
    environment:
      # 直接以服務名稱 "db" 作為 hostname
      DATABASE_URL: postgres://postgres:secret@db:5432/mydb
      # 直接以服務名稱 "redis" 作為 hostname
      REDIS_URL: redis://redis:6379

  db:
    image: postgres:17-alpine

  redis:
    image: redis:7-alpine
```

```
┌──────────────────────────────────────────────────┐
│         預設網路 (myapp_default)                    │
│                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐       │
│  │   api    │  │    db    │  │  redis   │       │
│  │          │──│          │  │          │       │
│  │          │  │ :5432    │  │ :6379    │       │
│  │          │──│──────────│──│          │       │
│  │ :8080    │  │          │  │          │       │
│  └────┬─────┘  └──────────┘  └──────────┘       │
│       │                                          │
└───────┼──────────────────────────────────────────┘
        │
   Port Mapping
   8080:8080
        │
   ┌────▼──────────┐
   │  Host Machine  │
   │  localhost:8080 │
   └────────────────┘
```

#### 自訂網路（進階）

```yaml
services:
  api:
    networks:
      - frontend
      - backend

  db:
    networks:
      - backend     # db 只在 backend 網路，frontend 無法直接存取

  nginx:
    networks:
      - frontend

networks:
  frontend:
  backend:
```

---

### 3.5 Volume

```yaml
services:
  db:
    image: postgres:17-alpine
    volumes:
      # Named Volume — Docker 管理的持久化儲存
      - pgdata:/var/lib/postgresql/data

      # Bind Mount — 掛載宿主機目錄（開發常用）
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro

  api:
    build: .
    volumes:
      # 開發模式：掛載原始碼，搭配 hot reload
      - .:/app

# 宣告 Named Volume
volumes:
  pgdata:
    # driver: local  # 預設值，可省略
```

> Named Volume 的資料由 Docker 管理，即使容器刪除也不會消失。需使用 `docker volume rm` 明確刪除。

---

### 3.6 環境變數管理

#### 方式 1：直接在 YAML 中定義

```yaml
services:
  db:
    environment:
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: mydb
```

#### 方式 2：使用 .env 檔案

```bash
# .env（與 docker-compose.yml 同目錄）
POSTGRES_USER=myuser
POSTGRES_PASSWORD=secret
POSTGRES_DB=mydb
DB_PORT=5432
```

```yaml
services:
  db:
    image: postgres:17-alpine
    ports:
      - "${DB_PORT:-5432}:5432"    # 使用 .env 中的變數，預設值 5432
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
```

#### 方式 3：使用 env_file

```yaml
services:
  api:
    env_file:
      - .env          # 載入 .env 中所有變數
      - .env.local    # 可載入多個檔案，後者覆蓋前者
```

> **安全提醒**：務必將 `.env` 加入 `.gitignore`，避免將密碼提交至版本控制。

---

### 3.7 常用 Compose 指令

```bash
# 啟動所有服務（背景模式）
docker compose up -d

# 啟動並強制重新建置映像檔
docker compose up -d --build

# 查看服務狀態
docker compose ps

# 查看服務日誌
docker compose logs

# 追蹤特定服務的日誌
docker compose logs -f api

# 停止所有服務
docker compose down

# 停止並刪除 Volume（注意：資料將消失）
docker compose down -v

# 重啟特定服務
docker compose restart api

# 進入特定服務的容器
docker compose exec api sh

# 在服務中執行一次性命令
docker compose run --rm api go test ./...

# 查看服務設定（展開所有變數後的結果）
docker compose config

# 拉取所有映像檔
docker compose pull

# 只建置映像檔（不啟動）
docker compose build
```

| 指令 | 說明 |
|------|------|
| `up -d` | 啟動所有服務（背景） |
| `down` | 停止並移除容器、網路 |
| `down -v` | 同上，並刪除 Volume |
| `ps` | 列出服務狀態 |
| `logs -f [service]` | 追蹤日誌 |
| `exec [service] [cmd]` | 在執行中的容器內執行命令 |
| `run --rm [service] [cmd]` | 建立新容器執行一次性命令 |
| `build` | 建置映像檔 |
| `restart [service]` | 重啟服務 |
| `config` | 驗證並顯示設定 |

---

### 3.8 實戰：Go API + PostgreSQL + Redis

以下是一個完整的多容器應用範例。

#### 目錄結構

```
myapp/
├── docker-compose.yml
├── .env
├── api/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── .dockerignore
└── db/
    └── init.sql
```

#### .env

```bash
# .env
POSTGRES_USER=appuser
POSTGRES_PASSWORD=secretpassword
POSTGRES_DB=myapp
REDIS_PASSWORD=redispassword
API_PORT=8080
```

#### docker-compose.yml

```yaml
services:
  # ─── Go API 服務 ───
  api:
    build:
      context: ./api
      dockerfile: Dockerfile
    ports:
      - "${API_PORT:-8080}:8080"
    environment:
      DATABASE_URL: postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable
      REDIS_URL: redis://:${REDIS_PASSWORD}@redis:6379/0
      PORT: "8080"
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      start_period: 10s
      retries: 3

  # ─── PostgreSQL 資料庫 ───
  db:
    image: postgres:17-alpine
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./db/init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    ports:
      - "5432:5432"    # 開發時可從宿主機直接連線；正式環境可移除
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # ─── Redis 快取 ───
  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redisdata:/data
    ports:
      - "6379:6379"    # 開發時可從宿主機直接連線；正式環境可移除
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${REDIS_PASSWORD}", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  pgdata:
  redisdata:
```

#### db/init.sql

```sql
-- Initialize database schema
CREATE TABLE IF NOT EXISTS visits (
    id SERIAL PRIMARY KEY,
    path TEXT NOT NULL,
    visited_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_visits_path ON visits(path);
```

#### api/Dockerfile

```dockerfile
# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# Stage 2: Runtime
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/server .

USER appuser
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./server"]
```

#### api/main.go（完整版）

```go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	db          *sql.DB
	redisClient *redis.Client
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to PostgreSQL
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Connect to Redis
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	redisClient = redis.NewClient(opt)
	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/stats", handleStats)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
	log.Println("Server stopped")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Record visit to PostgreSQL
	_, err := db.ExecContext(ctx,
		"INSERT INTO visits (path) VALUES ($1)", r.URL.Path)
	if err != nil {
		log.Printf("Failed to record visit: %v", err)
	}

	// Increment counter in Redis
	count, err := redisClient.Incr(ctx, "visit_count").Result()
	if err != nil {
		log.Printf("Failed to increment Redis counter: %v", err)
		count = -1
	}

	hostname, _ := os.Hostname()
	response := map[string]any{
		"message":     "Hello from Go + Docker Compose!",
		"hostname":    hostname,
		"visit_count": count,
		"path":        r.URL.Path,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check PostgreSQL
	if err := db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "db: unhealthy: %v", err)
		return
	}

	// Check Redis
	if err := redisClient.Ping(ctx).Err(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "redis: unhealthy: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Query recent visits from PostgreSQL
	rows, err := db.QueryContext(ctx,
		"SELECT path, COUNT(*) as count FROM visits GROUP BY path ORDER BY count DESC LIMIT 10")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type PathStat struct {
		Path  string `json:"path"`
		Count int    `json:"count"`
	}

	var stats []PathStat
	for rows.Next() {
		var s PathStat
		if err := rows.Scan(&s.Path, &s.Count); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stats = append(stats, s)
	}

	// Get total count from Redis
	totalCount, _ := redisClient.Get(ctx, "visit_count").Int64()

	response := map[string]any{
		"total_visits": totalCount,
		"top_paths":   stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

#### 啟動與測試

```bash
# 啟動所有服務
docker compose up -d --build

# 查看服務狀態
docker compose ps

# 查看日誌
docker compose logs -f api

# 測試 API
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/stats | jq .

# 停止服務
docker compose down

# 停止服務並刪除資料
docker compose down -v
```

---

### 3.9 練習 3：用 Compose 部署完整服務

**目標**：在上方範例的基礎上，新增 Adminer（資料庫管理介面）服務。

在 `docker-compose.yml` 的 `services` 區塊新增：

```yaml
  adminer:
    image: adminer:4
    ports:
      - "8081:8080"
    depends_on:
      - db
    restart: unless-stopped
```

**步驟：**

1. 在 `docker-compose.yml` 新增 `adminer` 服務
2. 執行 `docker compose up -d`
3. 開啟瀏覽器前往 `http://localhost:8081`
4. 登入資訊：
   - System: PostgreSQL
   - Server: `db`
   - Username、Password、Database：使用 `.env` 中設定的值
5. 查看 `visits` 資料表的內容

**延伸練習：**

1. 新增一個 `nginx` 服務作為反向代理，將 `/` 導向 `api:8080`，將 `/adminer` 導向 `adminer:8080`。
2. 嘗試執行 `docker compose scale api=3` 啟動 3 個 API 實體，觀察不同請求的 hostname 變化。

---

## Part 4：綜合演練與延伸學習

### 4.1 綜合練習

**情境**：部署一個包含以下元件的微服務系統：

1. **前端**：Nginx 靜態網站
2. **後端 API**：Go HTTP 伺服器
3. **資料庫**：PostgreSQL
4. **快取**：Redis

**要求：**

- [ ] 所有服務以 Docker Compose 管理
- [ ] Go API 使用多階段建置
- [ ] PostgreSQL 資料使用 Named Volume 持久化
- [ ] 服務間設定正確的 `depends_on` + `healthcheck`
- [ ] 以 `.env` 檔案管理敏感資訊
- [ ] Nginx 反向代理 Go API

**提示 — Nginx 反向代理設定（nginx.conf）：**

```nginx
upstream api {
    server api:8080;
}

server {
    listen 80;

    # Static files
    location / {
        root /usr/share/nginx/html;
        index index.html;
        try_files $uri $uri/ @api;
    }

    # Reverse proxy to Go API
    location @api {
        proxy_pass http://api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Direct API routing
    location /api/ {
        proxy_pass http://api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

### 4.2 常見問題排查

#### 容器啟動失敗

```bash
# 查看容器日誌
docker logs <container-name>

# 查看容器詳細資訊
docker inspect <container-name>

# 檢查容器退出碼
docker inspect --format='{{.State.ExitCode}}' <container-name>
# 0   = 正常退出
# 1   = 應用程式錯誤
# 137 = OOM Killed 或 SIGKILL
# 143 = SIGTERM（正常停止）
```

#### 容器間無法連線

```bash
# 確認容器在同一網路
docker network inspect <network-name>

# 進入容器測試連線
docker exec -it api sh
# 在容器內執行：
#   ping db
#   wget -qO- http://api:8080/health
```

#### 磁碟空間不足

```bash
# 查看 Docker 佔用空間
docker system df

# 清理所有未使用的資源（映像檔、容器、網路、Volume）
docker system prune -a --volumes
# ⚠ 注意：此操作將刪除所有未使用的資料，請謹慎執行
```

#### 映像檔建置快取問題

```bash
# 強制不使用快取重新建置
docker build --no-cache -t my-app .

# Compose 中強制重新建置
docker compose build --no-cache
docker compose up -d --build --force-recreate
```

#### 常見錯誤訊息

| 錯誤訊息 | 原因 | 解決方式 |
|---------|------|---------|
| `port is already allocated` | 宿主機連接埠被佔用 | 更換連接埠或停止佔用的程序 |
| `no space left on device` | 磁碟空間不足 | 執行 `docker system prune` |
| `network not found` | 網路不存在 | 檢查 `docker-compose.yml` 的 `networks` 設定 |
| `exec format error` | 映像檔平台不符 | 加上 `--platform linux/amd64` 重新建置 |
| `permission denied` | 權限不足 | 檢查 `USER` 指令或目錄權限 |
| `connection refused` | 服務尚未就緒 | 新增 `depends_on` + `healthcheck` |

---

### 4.3 延伸學習資源

#### 官方文件

- [Docker 官方文件](https://docs.docker.com/)
- [Dockerfile 最佳實踐](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Docker Compose 規格](https://docs.docker.com/compose/compose-file/)

#### 進階主題

| 主題 | 說明 |
|------|------|
| **Docker Network 進階** | Bridge、Host、Overlay 等網路模式 |
| **Docker Security** | 映像檔掃描、Rootless 模式、Secret 管理 |
| **Container Registry** | Harbor、ECR、GCR 等私有映像檔倉庫 |
| **Container Orchestration** | Kubernetes、Docker Swarm |
| **CI/CD 整合** | GitHub Actions、GitLab CI 中使用 Docker |
| **多平台建置** | `docker buildx` 建置 ARM/AMD64 映像檔 |
| **日誌管理** | 日誌驅動程式、集中式日誌收集 |
| **監控** | Prometheus + Grafana 監控容器 |

---

## 附錄 A：Docker 指令速查表

### 映像檔 (Image)

```bash
docker pull <image>:<tag>           # 下載映像檔
docker images                        # 列出映像檔
docker rmi <image>                   # 刪除映像檔
docker build -t <name>:<tag> .       # 建置映像檔
docker tag <src> <dest>              # 標記映像檔
docker push <image>:<tag>            # 推送映像檔
docker save -o file.tar <image>      # 匯出映像檔
docker load -i file.tar              # 匯入映像檔
docker history <image>               # 查看建置歷史
docker inspect <image>               # 查看詳細資訊
docker image prune                   # 清理 dangling 映像檔
docker image prune -a                # 清理所有未使用映像檔
```

### 容器 (Container)

```bash
docker run [opts] <image> [cmd]      # 建立並啟動容器
docker create [opts] <image> [cmd]   # 建立容器（不啟動）
docker start <container>             # 啟動容器
docker stop <container>              # 停止容器
docker restart <container>           # 重啟容器
docker kill <container>              # 強制停止
docker rm <container>                # 刪除容器
docker ps                            # 列出執行中容器
docker ps -a                         # 列出所有容器
docker logs [-f] <container>         # 查看日誌
docker exec -it <container> <cmd>    # 在容器內執行命令
docker cp <src> <container>:<dest>   # 複製檔案
docker inspect <container>           # 查看詳細資訊
docker stats                         # 即時資源使用
docker top <container>               # 查看程序
docker port <container>              # 查看連接埠對映
docker container prune               # 清理已停止容器
```

### 資料卷 (Volume)

```bash
docker volume create <name>          # 建立資料卷
docker volume ls                     # 列出資料卷
docker volume inspect <name>         # 查看詳細資訊
docker volume rm <name>              # 刪除資料卷
docker volume prune                  # 清理未使用資料卷
```

### 網路 (Network)

```bash
docker network create <name>         # 建立網路
docker network ls                    # 列出網路
docker network inspect <name>        # 查看詳細資訊
docker network rm <name>             # 刪除網路
docker network connect <net> <ctn>   # 將容器加入網路
docker network disconnect <net> <c>  # 將容器移出網路
docker network prune                 # 清理未使用網路
```

### 系統

```bash
docker system df                     # 磁碟使用量
docker system prune                  # 清理未使用資源
docker system prune -a --volumes     # 深度清理（含映像檔和資料卷）
docker info                          # 系統資訊
docker version                       # 版本資訊
```

---

## 附錄 B：Dockerfile 指令速查表

```dockerfile
# Base image
FROM <image>:<tag> [AS <name>]

# Working directory
WORKDIR /path

# Copy files
COPY [--from=<stage>] <src> <dest>
ADD <src> <dest>

# Execute commands (build stage)
RUN <command>

# Environment variables
ENV KEY=value
ARG KEY=default

# Declare port
EXPOSE <port>

# Container startup command
CMD ["executable", "param1"]
ENTRYPOINT ["executable", "param1"]

# User
USER <user>[:<group>]

# Health check
HEALTHCHECK [OPTIONS] CMD <command>

# Metadata
LABEL key="value"

# Stop signal
STOPSIGNAL signal

# Shell
SHELL ["executable", "parameters"]
```

---

## 附錄 C：Docker Compose 速查表

```bash
docker compose up -d                 # 啟動服務（背景）
docker compose up -d --build         # 啟動並重新建置
docker compose down                  # 停止並移除
docker compose down -v               # 停止、移除、刪除 Volume
docker compose ps                    # 服務狀態
docker compose logs [-f] [service]   # 服務日誌
docker compose exec <svc> <cmd>      # 在服務內執行命令
docker compose run --rm <svc> <cmd>  # 一次性命令
docker compose build [--no-cache]    # 建置映像檔
docker compose pull                  # 拉取映像檔
docker compose restart [service]     # 重啟服務
docker compose config                # 驗證設定
docker compose top                   # 查看服務程序
```

---

> 完成本指南後，你已掌握以下能力：
>
> - 理解 Docker 核心概念（容器化、映像檔、容器、Layer 架構）
> - 使用 Dockerfile 建置 Go 應用映像檔
> - 以 Multi-stage Build 大幅縮減映像檔體積
> - 以 Docker Compose 管理多容器應用
>
> **建議下一步**：將現有專案容器化，並以 Docker Compose 管理開發環境，親身體驗容器化帶來的一致性。
