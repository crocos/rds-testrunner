# RDS TestRunner

This is Query execution time measurement tool.
It is running query and auto-warmup for created mirror DB from RDS Snapshot.

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

see more: `rds-testrunner open --help`


* Delete DB

```
rds-testrunner close
```

* check opened DB list

```
rds-testrunner list
```

## Installation

Download Binary.


## Configuration

Create `/etc/rds-testrunner.conf` or `$HOME/.rtrrc` config file.

```
aws {
  key = "AWS_ACCESS_KEY_ID"
  secret = "AWS_SECRET_ACCESS_KEY"
}

resource "default" {
  instance = "crocos-test-db"
  region = "ap-northeast-1"
  uesr = "dbadmin"
  password = "dbpassword"
  warmup = "/tmp/warmup-query.sql"

}

resource "other..." {
  ..
```

`aws` config is AWS API Credentials. You can use ENV `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY instead of this config.

`resource "default"` is used default resource at `-r` parameter not specified.

see other configuration: `rds-testrunner.conf.example` and `warmul.sql.example`


Config style adopted [HCL](https://github.com/hashicorp/hcl).
You can choice hcl or json style as you like.


