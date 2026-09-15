# go-admin

  <img align="right" width="320" src="https://doc-image.zhangwj.com/img/go-admin.svg">


[![Build Status](https://github.com/go-admin-team/go-admin/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/go-admin-team/go-admin)
[![Release](https://img.shields.io/github/release/go-admin-team/go-admin.svg?style=flat-square)](https://github.com/go-admin-team/go-admin/releases)
[![License](https://img.shields.io/github/license/go-admin-team/go-admin.svg)](https://github.com/go-admin-team/go-admin)

[English](https://github.com/go-admin-team/go-admin/blob/master/README.md) | [简体中文](https://github.com/go-admin-team/go-admin/blob/master/README.Zh-cn.md) | [繁體中文](https://github.com/go-admin-team/go-admin/blob/master/README.zh-TW.md) | 日本語

Gin + Vue + Element UI / Arco Design / Ant Design による、フロントエンドとバックエンドを分離した権限管理システムです。初期化は非常に簡単で、設定ファイルのデータベース接続情報を変更するだけで動作します。複数のコマンドに対応しており、マイグレーションコマンドでデータベースの初期化が容易になり、サーバーコマンドで API を手軽に起動できます。

[オンラインドキュメント](https://www.go-admin.pro)

[フロントエンドプロジェクト](https://github.com/go-admin-team/go-admin-ui)

[動画チュートリアル](https://space.bilibili.com/565616721/channel/detail?cid=125737)

## 🎬 オンラインデモ

Element Plus vue3 デモ：[https://vue.go-admin.pro](https://vue.go-admin.pro/#/login)
> ⚠️⚠️⚠️ アカウント / パスワード： admin / 123456

antd デモ（go-admin-pro）：[https://antd.go-admin.pro](https://antd.go-admin.pro/)
> ⚠️⚠️⚠️ アカウント / パスワード： admin / 123456

## ✨ 特徴

- RESTful API の設計規約に準拠

- GIN WEB API フレームワークをベースに、豊富なミドルウェアを提供（ユーザー認証、CORS、アクセスログ、トレース ID など）

- Casbin による RBAC アクセス制御モデル

- JWT 認証

- Swagger ドキュメントに対応（swaggo ベース）

- GORM によるデータベース永続化、複数種類のデータベースに拡張可能

- 設定ファイルからモデルへの単純なマッピングで、必要な設定をすぐに取得

- コード生成ツール

- フォームビルダー

- マルチコマンド方式

- マルチテナント対応

- TODO: ユニットテスト

## 🎁 標準機能

1. マルチテナント：デフォルトで対応。データベース単位で分離し、1 データベースにつき 1 テナント。
1. ユーザー管理：システムの操作者であるユーザーの設定を行います。
2. 部門管理：組織構造（会社・部門・グループ）を設定します。ツリー構造で表示し、データ権限に対応します。
3. 役職管理：ユーザーが担当する職務を設定します。
4. メニュー管理：メニュー、操作権限、ボタン権限識別子、API 権限などを設定します。
5. ロール管理：ロールへのメニュー権限の割り当て、および組織単位でのデータ範囲権限の設定を行います。
6. 辞書管理：システム内で頻繁に使う固定的なデータを管理します。
7. パラメータ管理：よく使うパラメータを動的に設定します。
8. 操作ログ：正常系の操作ログと異常情報のログを記録・検索します。
9. ログインログ：ログイン履歴を記録・検索します。ログイン異常も含みます。
1. API ドキュメント：業務コードから API ドキュメントを自動生成します。
1. コード生成：テーブル定義から CRUD 業務を生成します。すべて画面上で操作でき、基本的な業務をコードなしで実現できます。
1. フォームビルダー：ページのスタイルをカスタマイズし、ドラッグ＆ドロップでレイアウトを作成します。
1. サービス監視：サーバーの基本情報を確認します。
1. コンテンツ管理：デモ機能。カテゴリ管理とコンテンツ管理を含み、入門用の参考実装として利用できます。
1. スケジュールタスク：自動実行タスク。現在は API 呼び出しと関数呼び出しに対応しています。

## 事前準備

ローカルに [go] [gin] [node](http://nodejs.org/) と [git](https://git-scm.com/) をインストールしてください。

ダウンロードから使いこなすまでを解説した動画とドキュメントのチュートリアルを用意しています。本プロジェクトを試す前に、まずこれらに目を通すことを強くおすすめします。

### go-admin で最初のアプリケーションを作る - ドキュメント

[ステップ 1 - 基礎の紹介](https://www.go-admin.pro/guide/intro/tutorial01.html)

[ステップ 2 - 実践 - CRUD を書く](https://www.go-admin.pro/guide/intro/tutorial02.html)

### 動画チュートリアル

[go-admin の起動方法](https://www.bilibili.com/video/BV1z5411x7JG)

[生成ツールで業務を手軽に実装する](https://www.bilibili.com/video/BV1Dg4y1i79D)

[v1.1.0 のコード生成ツール](https://www.bilibili.com/video/BV1N54y1i71P) [応用]

[マルチコマンドでの起動方法と IDE 設定](https://www.bilibili.com/video/BV1Fg4y1q7ph)

[go-admin のメニュー設定](https://www.bilibili.com/video/BV1Wp4y1D715) [必見]

[メニュー情報と API 情報の設定方法](https://www.bilibili.com/video/BV1zv411B7nG) [必見]

[go-admin の権限設定](https://www.bilibili.com/video/BV1rt4y197d3) [必見]

[go-admin のデータ権限](https://www.bilibili.com/video/BV1LK4y1s71e) [必見]

**不明点はまず上記のドキュメントと記事をご確認ください。解決しない場合は issue や pr をお寄せください。動画とドキュメントは継続的に更新しています**

## 📦 ローカル開発

### 動作要件

go 1.26.5

node バージョン: v22 以上（v24 LTS 推奨）

パッケージマネージャー: pnpm v9 以上（UI プロジェクトは pnpm を使用）

### 開発ディレクトリの作成

```bash

# 開発ディレクトリを作成
mkdir goadmin
cd goadmin
```

### コードの取得

> 重要：2 つのプロジェクトは同じディレクトリに配置してください。

```bash
# バックエンドのコードを取得
git clone https://github.com/go-admin-team/go-admin.git

# フロントエンドのコードを取得
git clone https://github.com/go-admin-team/go-admin-ui.git

```

### 起動方法

#### サーバーの起動

```bash
# go-admin バックエンドプロジェクトへ移動
cd ./go-admin

# 依存関係を整理
go mod tidy

# ビルド
go build

# 設定を変更
# ファイルパス  go-admin/config/settings.yml
vi ./config/settings.yml

# 1. 設定ファイル内のデータベース情報を変更
# 注意: settings.database 配下の設定項目
# 2. log のパスを確認
```

⚠️注意 Windows 環境で CGO が未導入の場合、次のエラーが発生します。

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

[cgo の問題の解決方法はこちら](https://www.go-admin.pro/zh-CN/guide/faq#cgo-%E7%9A%84%E9%97%AE%E9%A2%98)


#### データベースの初期化とサービス起動

``` bash
# 初回はデータベースのリソース情報を初期化する必要があります
# macOS または linux の場合
$ ./go-admin migrate -c config/settings.dev.yml

# ⚠️注意: windows の場合
$ go-admin.exe migrate -c config/settings.dev.yml


# プロジェクトを起動します。IDE からデバッグ実行することもできます
# macOS または linux の場合
$ ./go-admin server -c config/settings.yml


# ⚠️注意: windows の場合
$ go-admin.exe server -c config/settings.yml
```

#### sys_api テーブルへのデータ追加方法

起動時に `-a true` を付けると、不足している API データが自動的に追加されます。
```bash
./go-admin server -c config/settings.yml -a true
```

#### docker でのビルドと起動

```shell
# イメージをビルド
docker build -t go-admin .

# コンテナを起動します。1 つ目の go-admin はコンテナ名、2 つ目はイメージ名です
# -v は設定ファイルのマウント ローカルパス：コンテナ内パス
docker run --name go-admin -p 8000:8000 -v /config/settings.yml:/config/settings.yml -d go-admin-server
```

#### ドキュメント生成

```bash
go generate
```

#### クロスコンパイル

```bash
# windows
env GOOS=windows GOARCH=amd64 go build main.go

# or
# linux
env GOOS=linux GOARCH=amd64 go build main.go
```

### UI 側の起動方法

```bash
# pnpm をインストール（未導入の場合）
npm install -g pnpm

# 依存関係をインストール
pnpm install

# 中国本土のネットワークではミラーを指定すると高速化できます
pnpm install --registry=https://registry.npmmirror.com

# 開発サーバーを起動
pnpm dev
```

## 📨 コミュニティ

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

## 💎 コントリビューター


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
## JetBrains のオープンソースライセンス支援

`go-admin` は一貫して JetBrains 社の GoLand 統合開発環境で開発されています。**free JetBrains Open Source license(s)** による正規の無償ライセンス提供に、この場を借りて感謝を申し上げます。

<a href="https://www.jetbrains.com/?from=kubeadm-ha" target="_blank"><img src="https://raw.githubusercontent.com/panjf2000/illustrations/master/jetbrains/jetbrains-variant-4.png" width="250" align="middle"/></a>

## 🤝 謝辞

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

## 🤟 支援

> このプロジェクトがお役に立ちましたら、作者にジュースを一杯おごる形で応援いただけます :tropical_drink:

<img class="no-margin" src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/pay.png"  height="200px" >

## 🤝 関連リンク

- [mss-boot-io](https://docs.mss-boot-io.top/)

## 🔑 License

[MIT](https://github.com/go-admin-team/go-admin/blob/master/LICENSE.md)

Copyright (c) 2026 wenjianzhang
