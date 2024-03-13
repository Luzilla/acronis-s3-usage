# acronis-s3-usage

This is a playground to explore certain APIs provided by ACI (or VHI) to extract usage for the object storage provided.

## pull overall usage (ACI)

Some code to pull S3 storage usage by tenant.

```sh
$ go run ./cmd/tenant-usage/main.go
Got a token: ABC***
Got tenant id: abc-abc-abc-abc-abc

Cyber Infrastructure (Type: hci)
hci_s3_storage -- 11070.81 GB

Cyber Infrastructure (Type: hci)
hci_s3_storage -- 4619.61 GB
```

## extract usage for buckets (ACI &amp; VHI)

```sh
❯ go run ./cmd/ostor/main.go --help
NAME:
   a program to interact with the s3 management APIs in ACI and VHI - A new cli application

USAGE:
   a program to interact with the s3 management APIs in ACI and VHI [global options] command [command options] 

VERSION:
   dev (none, date: unknown)

COMMANDS:
   buckets, b  list buckets
   stats, s    list stats
   users, u    list users
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --s3-endpoint value        [$S3_ENDPOINT]
   --s3-system-key-id value   [$S3_SYSTEM_KEY_ID]
   --s3-system-secret value   [$S3_SYSTEM_SECRET_KEY]
   --help, -h                show help
   --version, -v             print the version
```
