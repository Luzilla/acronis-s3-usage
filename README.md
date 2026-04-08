# acronis-s3-usage

This project started as a playground to explore certain APIs provided by ACI (or VHI) to extract usage for the object storage provided. It now hosts libraries to interact with ACI (specifically) and Ostor APIs (for Acronis and Virtuozzo Object Storage).

> [!CAUTION]
> The code is still pre `v1.0.0` so some BC breaks are expected. It is used in production already.
> Please **review** the changes particularly from `v0.9.0` to `v0.10.0`.

> [!IMPORTANT]
> Please star the repository if you find it useful.
> For custom feature development or support, please get in touch.

## library

See the [`pkg`](./pkg/) directory for the Golang libraries (or SDKs, if you will) for `acronis` (ACI) and the Acronis/Virtuozzo Ostor API.

## toys

This repository includes ready code that allows you to run code examples against your Acronis or Virtuozzo setups.

Each tool will take arguments for the user/password of required tokens/credentials, but also supports environment variables (see [`.envrc-dist`](./.envrc-dist))

The command line tools (along with the source code), are also available as [release downloads](https://github.com/Luzilla/acronis-s3-usage/releases).

### pull overall usage (ACI)

A command line tool to pull S3 storage usage by tenant from Acronis Cyber Infrastructure (ACI).

```sh
$ go run ./cmd/tenant-usage/main.go
Got a token: ABC***
Got tenant id: abc-abc-abc-abc-abc

Cyber Infrastructure (Type: hci)
hci_s3_storage -- 11070.81 GB

Cyber Infrastructure (Type: hci)
hci_s3_storage -- 4619.61 GB
```

### extract usage for buckets (ACI &amp; VHI)

A command line tool to interact with the Ostor APIs — it allows user management, bucket management and extracting statistics (e.g. number of GET, HEAD, POST and PUT requests and storage used).

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

## adopters

- [Planetary Quantum GmbH](https://www.planetary-quantum.com/)

## contributions

If you end up using this, feel free to let me know by adding yourself to the adopters. All contributions (documentation, bug fixes, feature suggestions) are welcome.