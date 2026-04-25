# 01 — CI/CD 概念介紹

角色設定：

- **Ocean** — SDC SRE Team 新進成員，負責維護社團的基礎設施和部署流程，還在摸索中，常常搞出各種狀況。
- **Andrew** — 整天在喝酒的 SDC 社長，熟悉社團大小事。
- **Snow** — SDC 開發部長，明明很聰明但常常把事情丟給 Ocean 做，口頭禪是「歐不」。

## Table of Contents

- [01 — CI/CD 概念介紹](#01--cicd-概念介紹)
  - [Table of Contents](#table-of-contents)
  - [軟體開發的多人協作](#軟體開發的多人協作)
  - [什麼是 CI/CD？](#什麼是-cicd)
    - [SDC 實際怎麼做？](#sdc-實際怎麼做)
  - [常見的 CI/CD 工具](#常見的-cicd-工具)
  - [小結](#小結)

## 軟體開發的多人協作

Ocean 最近在跟 Andrew 開發一個新專案「Special Ground」，兩個人一起在同一個 GitHub repository 上面做軟體開發。兩個人都直接在 `main` 分支（branch）上推程式碼，有一天 Ocean 推了一版上去，結果把 Andrew 前一天寫好的功能整個蓋掉了，上線中的服務也跟著壞掉。兩個人想把版本 revert 回去，結果發現 git history 已經亂成一鍋粥，用威士忌熬的那種，搞了半天也恢復不了原本的版本。

Snow 看到他們的窘境，決定教他們怎麼好好在 GitHub 上面用版本控制做多人協作，順便介紹一個叫做 CI/CD 的好東西。

Snow 第一件事就是叫他們不要再把程式直接推到 `main` 上面，而是各自開一條分支開發，寫到一個段落之後發 **[Pull Request（PR）](https://shoujhengduan.medium.com/%E4%BB%80%E9%BA%BC%E6%98%AF-pull-request-b476ee3e0217)**。PR 是 GitHub 上的一個功能，用途是告訴其他人「我這邊改好了，請來看一下」。其他人可以在 PR 上看到你改了哪些程式碼、留下討論、提出修改建議。如果需要調整，就繼續在同一條分支上 push 新的 commit，PR 則會自動更新。等大家都確認沒問題了，才由維護者把這條分支合併進 `main`。Snow 自己則負責盯著 `main` 的狀態，確保上面的程式碼是乾淨的。

![Pull Request 頁面](assets/github-pull-request.png)

上圖是一個實際的 PR 頁面（截自 SDC 內部專案）。標題寫了這次改了什麼，下面標示著要把開發分支合併到 `main`。右邊可以指定 **Reviewer**（請誰來看）和 **Assignee**（誰負責這個 PR）。中間的 tab 可以切換看 **Conversation**（討論）、**Commits**（這個 PR 包含哪些 commit）、**Checks**（CI 自動檢查的結果，後面會詳細講）、**Files changed**（改了哪些檔案的哪些地方）。

具體的做法很簡單：`main` 分支永遠保持穩定（隨時都可以拿來部署），每次要寫新功能就從 `main` 開一條新的分支（例如 `feature/add-user-auth`），在上面開發、發 PR、review 通過後合回 `main`。SDC 目前的專案（例如 [Core System backend](https://github.com/NYCU-SDC/core-system-backend)）就是用這種方式協作。

## 什麼是 CI/CD？

**CI/CD** 的全名是 Continuous Integration / Continuous Deployment（持續整合 / 持續部署）。做法是把編譯、測試、部署這些原本靠人手動做的步驟，定義成一條自動化的管線（pipeline）：程式碼一推上去，CI 自動跑檢查確保品質，通過之後 CD 自動部署到對應的環境上。

CI/CD 的核心目標是「維持程式品質」和縮短「寫完程式碼」到「使用者看到新功能」的時間。

品質的部分，每次程式碼變更都會自動跑完整的測試和檢查，每個人的程式碼都經過同一套標準，不會因為本機環境不一樣而出現「在我電腦上可以跑啊」的問題。部署流程也由 pipeline 定義好，每次步驟都一樣，不會因為人為操作漏掉關鍵步驟。

速度的部分，重複性的工作全部交給機器做，工程師把時間和腦力留給產品迭代本身。整個流程自動化之後，團隊也可以更頻繁地發版，每次變更更小、風險更低、出問題也更容易定位。

CI/CD 拆成兩段：CI（Continuous Integration，持續整合）負責「合併前自動跑品質檢查」，CD（Continuous Deployment，持續部署）負責「合併後自動部署到對應環境」。CI 怎麼實作見 [03 章](03-go-ci-pipeline.md)，CD 怎麼實作見 [04 章](04-deployment.md)。

### SDC 實際怎麼做？

SDC 現行的專案（例如 [Core System backend](https://github.com/NYCU-SDC/core-system-backend)、[Clustron backend](https://github.com/NYCU-SDC/clustron-backend)）都具備完整的 CI/CD pipeline，並且他們在不同的環境（正式版 / 測試版）會經過不同的 pipeline 流程。

以 [Core System backend](https://github.com/NYCU-SDC/core-system-backend) 為例，CI 的部分使用 **GitHub Actions**，在每次有程式變更時都會跑一次 lint、test、build 等流程。開發者在處理手上的任務時，當流程全部正常通過，這個 PR 才可以被 merge。CD 的部分則使用由 SRE Team 維護的 n8n deployment workflow 做自動化部署。

回頭對照一下，Ocean 遇到的那些問題，在導入 CI/CD 之後是怎麼被解決的：

| 沒有 CI/CD 的時候 | 有了 CI/CD 之後 |
|------------------|----------------|
| Ocean 忘記跑測試就合併，壞掉的程式碼進到 `main` | 每次合併自動執行測試，沒通過就無法合併 |
| 本機測試環境不一樣，常出現「在我電腦上可以跑啊」的問題 | 每個人的程式碼都在同一套 runner 環境下跑過檢查，品質標準一致 |
| revert 不回去，git history 一團亂 | 每次變更都很小、都經過驗證，出問題也容易定位和復原 |

## 常見的 CI/CD 工具

市面上有許多 CI/CD 工具，以下是最常見的幾個：

- **GitHub Actions** — 本課程的重點！與 GitHub 深度整合，免費且功能強大
- **GitLab CI/CD** — GitLab 內建的 CI/CD 平台，與 GitLab 生態系緊密結合
- **Jenkins** — 老牌的開源 CI/CD 工具，高度可客製化，但設定較複雜
- **CircleCI** — 雲端 CI/CD 服務，以速度與易用性著稱
- **Travis CI** — 早期非常流行的 CI 服務，對開源專案友善

## 小結

- **CI/CD** 是一套自動化實踐，讓程式碼從開發到部署的每個階段都能自動執行，**CI** 著重在頻繁合併與自動測試，**CD** 著重在自動化部署流程。
- CI/CD 能解決手動流程的痛點：忘記測試、環境不一致、部署出錯。
- **GitHub Actions** 是我們今天要學的 CI/CD 工具，與 GitHub 無縫整合。

> 下一章我們就來動手，寫出你的第一個 GitHub Actions workflow。

[下一章：GitHub Actions 基礎 →](02-github-actions-basics.md)
