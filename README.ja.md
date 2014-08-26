# RDS TestRunner

RDS TestRunnerは、RDSの最新SnapshotからDBを複製し、それに対して各種クエリを実行し、完了までの時間を計測します。

もう少し細かく言うと、以下の作業をします。

* RDSの最新SnapshotからDBを複製する
    * インスタンス名は実行ユーザー名を元に自動で付与される
* DB Warm up用のSQLを投入する
* 列挙したSQLクエリを投入、各クエリの実行時間を計測する
* チャットに実行完了通知を送る
    * 現在はHipChatのみ

DBの複製機能以外はOptionalです。また、複製したDBの一括破棄もできます。

## 使い方

RDS TestRunnerの使い方を紹介します。

### DBの複製とクエリの計測

クエリの計測には、`open`コマンドを使い、`-q`でクエリを列挙したファイルを渡して実行します。
(クエリファイルは、現状では1行毎に評価・計測するようになっているので注意してください。)

```
rds-testrunner open -q test-query.sql
```

`default`以外に設定した`resource`があれば、下記のようにパラメータで切替できます。

```
rds-testrunner open -r "fooresource" -q foo-test-query.sql
```

また、DBのInstance Typeを変更したい場合にも`-t`で指定できます(デフォルトでは`db.m3.medium`)

```
rds-testrunner open -t db.m3.large -q foo-test-query.sql
```

`rds-testrunnner --help` および `rds-testrunner open --help`を参照してください。

### DBの破棄

利用し終わったDB群をまとめて破棄します。

```
rds-testrunner close
```

### 利用中のDB一覧

現在、自分が起動しているミラーDBの一覧を表示します。

```
rds-testrunner list
```

## Install

バイナリを[Download](https://github.com/crocos/rds-testrunner/releases)してください。
適当にpathの通ったところに`rds-testrunner`のバイナリを置けば利用できます。


## 設定ファイル

`rds-testrunner`は設定ファイルとして`/etc/rds-testrunner.conf`もしくは`$HOME/.rtrrc`を読み込みます。両方ある場合は後者が優先されます。

設定ファイルの書式にはHashiCorpが先日リリースした[HCL](https://github.com/hashicorp/hcl)を利用しています。

設定例は、`rds-testrunner.conf.example`にあります。

```
resource "default" {
  instance = "crocos-test-db"
  region = "ap-northeast-1"
  uesr = "dbadmin"
  password = "dbpassword"
  warmup = "/tmp/warmup-query.sql"

}

resource "other..." {
  ..
}

aws {
  key = "AWS_ACCESS_KEY_ID"
  secret = "AWS_SECRET_ACCESS_KEY"
}

notify {
  type = "hipchat"
  token = "hipchat token"
  room = "notify room id or name"
}
```

### resource

`resource`ディレクティブにはDBインスタンスの情報を記載します。
複数定義することができ、標準では`default`が利用されます。

コマンド実行時に`-r`で対象とする`resource`の種類を指定することができます。

### aws

`aws`ディレクティブにはCredential情報を記載します。

この項目を設定する代わりに、ENVに`AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY`をexportしておくか、IAM credentialを使える環境でも実行できます。

### notify (optional)

チャットへの通知に関するディレクティブです。
現在は、`type`として`hipchat`のみが指定可能です。

`token`にアクセストークン、`room`に通知先のroom_id or room名を指定してください。


## Development

ライブラリの依存関係の管理には`godep`を利用しています。

```
godep get
```

にて必要なライブラリを取得できます。

また、buildにはパッケージのディレクトリに移動して

```
go build
```

もしくは

```
gox -output="build/{{.OS}}/{{.Arch}}/{{.Dir}}"
```

でbuildできます。



