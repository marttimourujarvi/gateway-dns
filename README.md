```
                o                                            8
                8                                            8
.oPYo. .oPYo.  o8P .oPYo. o   o   o .oPYo. o    o       .oPYo8 odYo. .oPYo.
8    8 .oooo8   8  8oooo8 Y. .P. .P .oooo8 8    8 ooooo 8    8 8' ´8 Yb..
8    8 8    8   8  8.     ´b.d'b.d' 8    8 8    8       8    8 8   8   'Yb.
ˋYooP8 ˋYooP8   8  ˋYooo'  ˋY' ˋY'  ˋYooP8 ˋYooP8       ˋYooP' 8   8 ˋYooP'
:....8 :.....:::..::.....:::..::..:::.....::....8 :::::::.....:..::..:.....:
::ooP'.::::::::::::::::::::::::::::::::::::::ooP'.::::::::::::::::::::::::::
::...::::::::::::::::::::::::::::::::::::::::...::::::::::::::::::::::::::::
```

# gateway-dns

A lightweight in-cluster DNS server that automatically resolves hostnames from Kubernetes Gateway API resources (`HTTPRoute`, `GRPCRoute`) to the IPv4 addresses of their parent `Gateway`.

## How it works

1. Watches `HTTPRoute` and `GRPCRoute` resources via the Kubernetes controller-runtime.
2. Extracts hostnames from route specs and looks up the parent Gateway's `.status.addresses`.
3. Stores the resulting A records in an in-memory `DNSStore`.
4. Serves DNS A queries over UDP on port `5353` (configurable).

## Build

```bash
make build   # docker build . -t gateway-dns
```

## Run locally (requires kubeconfig)

```bash
go run ./cmd/gateway-dns/main.go
```

The DNS server listens on `:5353` by default. Use `dig` to test:

```bash
dig @127.0.0.1 -p 5353 <hostname-from-httproute>
```

## Deploy

A Dockerfile and deploy manifests are included. See `deploy/` and `Makefile` for kind-based workflow.

## Project layout

| Path                          | Purpose                                   |
|-------------------------------|-------------------------------------------|
| `cmd/gateway-dns/main.go`     | Entrypoint (DNS server + controller mgr)  |
| `internal/server/server.go`   | UDP DNS server using `miekg/dns`          |
| `internal/dnsstore/store.go`  | Thread-safe in-memory A-record store      |
| `internal/controller/`        | Gateway API reconciler                    |
| `deploy/`                     | Kubernetes manifests                      |

## Release process

This repo follows **trunk-based development**:

1. Merge feature / fix PRs directly to `main` as often as needed.
2. When you are ready to cut a release, go to **Actions → Release Please** and click **Run workflow** (or wait for the weekly Monday cron).
3. Review the generated Release PR, then merge it.
4. Merging the Release PR creates a GitHub Release + tag, which triggers the **Release** workflow to build and push the multi-arch Docker image.
5. If you ever need to rebuild a Docker image for an existing tag, go to **Actions → Release** and run it manually with the desired tag.
