# go-admin

  <img align="right" width="320" src="https://doc-image.zhangwj.com/img/go-admin.svg">


[![Build Status](https://github.com/go-admin-team/go-admin/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/go-admin-team/go-admin)
[![Release](https://img.shields.io/github/release/go-admin-team/go-admin.svg?style=flat-square)](https://github.com/go-admin-team/go-admin/releases)
[![License](https://img.shields.io/github/license/go-admin-team/go-admin.svg)](https://github.com/go-admin-team/go-admin)

[English](https://github.com/go-admin-team/go-admin/blob/master/README.md) | [简体中文](https://github.com/go-admin-team/go-admin/blob/master/README.Zh-cn.md) | 繁體中文 | [日本語](https://github.com/go-admin-team/go-admin/blob/master/README.ja-JP.md)

基於 Gin + Vue + Element UI OR Arco Design OR Ant Design 的前後端分離權限管理系統。系統初始化極為簡單，只需在設定檔中修改資料庫連線資訊即可。系統支援多指令操作：遷移指令讓資料庫初始化變得更簡單，服務指令則能輕鬆啟動 API 服務。

[線上文件](https://www.go-admin.pro)

[前端專案](https://github.com/go-admin-team/go-admin-ui)

[影片教學](https://space.bilibili.com/565616721/channel/detail?cid=125737)

## 🎬 線上體驗

Element Plus vue3 體驗：[https://vue.go-admin.pro](https://vue.go-admin.pro/#/login)
> ⚠️⚠️⚠️ 帳號 / 密碼： admin / 123456

antd 體驗（go-admin-pro）：[https://antd.go-admin.pro](https://antd.go-admin.pro/)
> ⚠️⚠️⚠️ 帳號 / 密碼： admin / 123456

## ✨ 特性

- 遵循 RESTful API 設計規範

- 基於 GIN WEB API 框架，提供豐富的中介軟體支援（使用者認證、跨域、存取日誌、追蹤 ID 等）

- 基於 Casbin 的 RBAC 存取控制模型

- JWT 認證

- 支援 Swagger 文件（基於 swaggo）

- 基於 GORM 的資料庫儲存，可擴充多種類型資料庫

- 設定檔簡單的模型映射，快速取得所需設定

- 程式碼產生工具

- 表單建構工具

- 多指令模式

- 多租戶的支援

- TODO: 單元測試

## 🎁 內建

1. 多租戶：系統預設支援多租戶，按資料庫分離，一個資料庫一個租戶。
1. 使用者管理：使用者是系統操作者，該功能主要完成系統使用者設定。
2. 部門管理：設定系統組織架構（公司、部門、小組），以樹狀結構呈現並支援資料權限。
3. 職位管理：設定系統使用者所擔任的職務。
4. 選單管理：設定系統選單、操作權限、按鈕權限標識、介面權限等。
5. 角色管理：角色選單權限分配、設定角色按機構進行資料範圍權限劃分。
6. 字典管理：對系統中經常使用且較為固定的資料進行維護。
7. 參數管理：對系統動態設定常用參數。
8. 操作日誌：系統正常操作的日誌記錄與查詢；系統異常資訊的日誌記錄與查詢。
9. 登入日誌：系統登入日誌記錄查詢，包含登入異常。
1. 介面文件：根據業務程式碼自動產生相關的 API 介面文件。
1. 程式碼產生：根據資料表結構產生對應的增刪改查業務，全程視覺化操作，讓基本業務可以零程式碼實現。
1. 表單建構：自訂頁面樣式，拖拉放實現頁面佈局。
1. 服務監控：檢視伺服器的基本資訊。
1. 內容管理：demo 功能，下設分類管理、內容管理，可參考使用以快速入門。
1. 排程任務：自動化任務，目前支援介面呼叫與函式呼叫。

## 準備工作

你需要在本機安裝 [go] [gin] [node](http://nodejs.org/) 和 [git](https://git-scm.com/)

同時配套了系列教學（含影片與文件），說明如何從下載到熟練使用。強烈建議先看完這些教學再來實作本專案！！！

### 輕鬆用 go-admin 寫出第一個應用 - 文件教學

[步驟一 - 基礎內容介紹](https://www.go-admin.pro/guide/intro/tutorial01.html)

[步驟二 - 實際應用 - 撰寫增刪改查](https://www.go-admin.pro/guide/intro/tutorial02.html)

### 手把手教你從入門到放棄 - 影片教學

[如何啟動 go-admin](https://www.bilibili.com/video/BV1z5411x7JG)

[使用產生工具輕鬆實現業務](https://www.bilibili.com/video/BV1Dg4y1i79D)

[v1.1.0 版本程式碼產生工具 - 釋放雙手](https://www.bilibili.com/video/BV1N54y1i71P) [進階]

[多指令啟動方式講解以及 IDE 設定](https://www.bilibili.com/video/BV1Fg4y1q7ph)

[go-admin 選單的設定說明](https://www.bilibili.com/video/BV1Wp4y1D715) [必看]

[如何設定選單資訊以及介面資訊](https://www.bilibili.com/video/BV1zv411B7nG) [必看]

[go-admin 權限設定使用說明](https://www.bilibili.com/video/BV1rt4y197d3) [必看]

[go-admin 資料權限使用說明](https://www.bilibili.com/video/BV1LK4y1s71e) [必看]

**如有問題請先參閱上述文件與文章，若仍無法解決，歡迎提出 issue 與 pr。影片教學與文件持續更新中**

## 📦 本機開發

### 環境需求

go 1.26.5

node 版本: v22+（建議 v24 LTS）

套件管理器: pnpm v9+（UI 專案使用 pnpm）

### 建立開發目錄

```bash

# 建立開發目錄
mkdir goadmin
cd goadmin
```

### 取得程式碼

> 重點注意：兩個專案必須放在同一資料夾下；

```bash
# 取得後端程式碼
git clone https://github.com/go-admin-team/go-admin.git

# 取得前端程式碼
git clone https://github.com/go-admin-team/go-admin-ui.git

```

### 啟動說明

#### 伺服器端啟動說明

```bash
# 進入 go-admin 後端專案
cd ./go-admin

# 更新整理相依套件
go mod tidy

# 編譯專案
go build

# 修改設定
# 檔案路徑  go-admin/config/settings.yml
vi ./config/settings.yml

# 1. 在設定檔中修改資料庫資訊
# 注意: settings.database 下對應的設定資料
# 2. 確認 log 路徑
```

⚠️注意 在 Windows 環境若未安裝 CGO，會出現這個問題；

```bash
E:\go-admin>go build
# github.com/mattn/go-sqlite3
cgo: exec /missing-cc: exec: "/missing-cc": file does not exist
```

or

```bash
D:\Code\go-admin>go build
# github.com/mattn/go-sqlite3
cgo: exec gcc: exec: "gcc": executable file not found in %PATH%
```

[解決 cgo 問題請進入](https://www.go-admin.pro/zh-CN/guide/faq#cgo-%E7%9A%84%E9%97%AE%E9%A2%98)


#### 初始化資料庫，以及服務啟動

``` bash
# 首次設定需要初始化資料庫資源資訊
# macOS or linux 下使用
$ ./go-admin migrate -c config/settings.dev.yml

# ⚠️注意:windows 下使用
$ go-admin.exe migrate -c config/settings.dev.yml


# 啟動專案，也可以用 IDE 進行除錯
# macOS or linux 下使用
$ ./go-admin server -c config/settings.yml


# ⚠️注意:windows 下使用
$ go-admin.exe server -c config/settings.yml
```

#### sys_api 表的資料如何新增

在專案啟動時，使用 `-a true` 系統會自動新增缺少的介面資料
```bash
./go-admin server -c config/settings.yml -a true
```

#### 使用 docker 編譯啟動

```shell
# 編譯映像檔
docker build -t go-admin .

# 啟動容器，第一個 go-admin 是容器名稱，第二個 go-admin 是映像檔名稱
# -v 映射設定檔 本機路徑：容器路徑
docker run --name go-admin -p 8000:8000 -v /config/settings.yml:/config/settings.yml -d go-admin-server
```

#### 文件產生

```bash
go generate
```

#### 交叉編譯

```bash
# windows
env GOOS=windows GOARCH=amd64 go build main.go

# or
# linux
env GOOS=linux GOARCH=amd64 go build main.go
```

### UI 互動端啟動說明

```bash
# 安裝 pnpm（若未安裝）
npm install -g pnpm

# 安裝相依套件
pnpm install

# 中國大陸網路可指定鏡像來源加速
pnpm install --registry=https://registry.npmmirror.com

# 啟動服務
pnpm dev
```

## 📨 互動

<table>
   <tr>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/wx.png" width="180px"></td>
    <td><img src="https://doc-image.zhangwj.com/img/qrcode_for_gh_b798dc7db30c_258.jpg" width="180px"></td>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/qq2.png" width="200px"></td>
    <td><a href="https://space.bilibili.com/565616721">wenjianzhang</a></td>
  </tr>
  <tr>
    <td>微信</td>
    <td>公众号🔥🔥🔥</td>
    <td><a target="_blank" href="https://shang.qq.com/wpa/qunwpa?idkey=0f2bf59f5f2edec6a4550c364242c0641f870aa328e468c4ee4b7dbfb392627b"><img border="0" src="https://pub.idqqimg.com/wpa/images/group.png" alt="go-admin技术交流乙号" title="go-admin技术交流乙号"></a></td>
    <td>哔哩哔哩🔥🔥🔥</td>
  </tr>
</table>

## 💎 貢獻者


<span style="margin: 0 5px;" ><a href="https://github.com/wenjianzhang" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/3890175?v=4&h=60&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/G-Akiraka" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/45746659?s=64&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/lwnmengjing" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/12806223?s=64&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/bing127" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/31166183?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/chengxiao" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/1379545?s=64&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/NightFire0307" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/19854086?v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/appleboy" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/21979?s=64&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/ninstein" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/580303?v=4&h=60&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/kikiyou" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/17959053?s=64&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/horizonzy" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/22524871?s=64&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/Cassuis" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/48005724?s=64&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/hqcchina" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/5179057?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/nodece" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/16235121?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/stephenzhang0713" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/18169290?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/zhouxixi-dev" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/100399679?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/Jalins" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/31172582?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/wkf928592" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/6063351?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/wxxiong6" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/6983441?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/Silicon-He" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/52478309?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/GizmoOAO" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/20385106?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/bestgopher" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/36840497?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/wxb1207" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/20775558?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/misakichan" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/16569274?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/zhuxuyang" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/19301024?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/mss-boot" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/109259065?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/AuroraV" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/37330199?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/Vingurzhou" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/57127283?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/haimait" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/40926384?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/zyd" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/3446278?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/infnan" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/38274826?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/d1y" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/45585937?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/qlijin" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/515900?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/logtous
" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/88697234?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/stepway
" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/9927079?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/NaturalGao
" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/43291304?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/DemoLiang
" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/23476007?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/jfcg
" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/1410597?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/Nicole0724
" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/10487328?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
## JetBrains 開源證書支援

`go-admin` 專案一直以來都是在 JetBrains 公司旗下的 GoLand 整合開發環境中進行開發，基於 **free JetBrains Open Source license(s)** 正版免費授權，在此表達我的謝意。

<a href="https://www.jetbrains.com/?from=kubeadm-ha" target="_blank"><img src="https://raw.githubusercontent.com/panjf2000/illustrations/master/jetbrains/jetbrains-variant-4.png" width="250" align="middle"/></a>

## 🤝 特別感謝

1. [ant-design](https://github.com/ant-design/ant-design)
2. [ant-design-pro](https://github.com/ant-design/ant-design-pro)
2. [arco-design](https://github.com/arco-design/arco-design)
2. [arco-design-pro](https://github.com/arco-design/arco-design-pro)
4. [gin](https://github.com/gin-gonic/gin)
5. [casbin](https://github.com/casbin/casbin)
6. [spf13/viper](https://github.com/spf13/viper)
7. [gorm](https://github.com/go-gorm/gorm)
8. [gin-swagger](https://github.com/swaggo/gin-swagger)
9. [golang-jwt](https://github.com/golang-jwt/jwt)
10. [vue-element-admin](https://github.com/PanJiaChen/vue-element-admin)
11. [ruoyi-vue](https://gitee.com/y_project/RuoYi-Vue)
12. [form-generator](https://github.com/JakHuang/form-generator)

## 🤟 贊助

> 如果你覺得這個專案幫助到了你，可以幫作者買一杯果汁表示鼓勵 :tropical_drink:

<img class="no-margin" src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/pay.png"  height="200px" >

## 🤝 連結

- [mss-boot-io](https://docs.mss-boot-io.top/)

## 🔑 License

[MIT](https://github.com/go-admin-team/go-admin/blob/master/LICENSE.md)

Copyright (c) 2026 wenjianzhang
