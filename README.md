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

