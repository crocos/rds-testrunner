# RDS TestRunner

This is Query execution time measurement tool.
It is running query and auto-warmup for created mirror DB from RDS(MySQL) Snapshot.

## How to use

* Create DB and execute query.

```
rds-testrunner open -q test-query.sql
  # test-query.sql is SQL Query list file.
```

It is used default resource. If you want useing other resource:

```
rds-testrunner open -r "fooresource" -q foo-test-query.sql
```

And If you want to change instance type(default: `db.m3.medium`):

```
rds-testrunner open -t db.m3.large -q foo-test-query.sql
```

see more: `rds-testrunner --help` and `rds-testrunner open --help`


* Delete end of use DB's

```
rds-testrunner close
```

* check opened DB list

```
rds-testrunner list
```

## Installation

Download [Binary](https://github.com/crocos/rds-testrunner/releases).


## Configuration

Create `/etc/rds-testrunner.conf` or `$HOME/.rtrrc` config file.

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


(required) `resource "default"` is used default resource at `-r` parameter not specified.

(optional) `aws` config is AWS API Credentials.
You can use ENV `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY instead of this config.

(optional) `notify` config is notify finished execution. (Now, notify `type` can set only `hipchat`.)

see other configuration: `rds-testrunner.conf.example` and `warmul.sql.example`


Config style adopted [HCL](https://github.com/hashicorp/hcl).
You can choice hcl or json style as you like.


## Development

Dependency packages get by `godep`

```
godep get
```

and build follow:

```
go build
```

or

```
gox -output="build/{{.OS}}/{{.Arch}}/{{.Dir}}"
```


