# router7

> **This is the [0magnet](https://github.com/0magnet) fork of
> [rtr7/router7](https://github.com/rtr7/router7)** (Apache-2.0). It tracks
> upstream and adds one thing: the LAN-serving libraries are *importable*, so
> another Go program can embed router7's DHCPv4 and DNS servers instead of
> shelling out to `dnsmasq`.
>
> Differences from upstream:
>
> - Module path is `github.com/0magnet/router7`.
> - `dhcp4d`, `dhcp6`, `dns`, `oui`, `teelogger` and `multilisten` moved from
>   `internal/` to `pkg/`. The appliance-only packages (`netconfig`, `radvd`,
>   `dhcp4`, `backup`, `diag`, `dyndns`, `notify`) stay in `internal/`.
> - `dhcp4d` no longer imports `internal/netconfig`. It takes an optional
>   `ServerIPFunc` to discover the server's own address; the default
>   (`InterfaceServerIP`) reads it off the live interface, so an embedder needs
>   no `interfaces.json`. `cmd/dhcp4d` passes `netconfig.ServerIP`, keeping the
>   appliance's behavior byte-for-byte. This is what drops the
>   netlink/nftables/wireguard dependency tree from the library's import graph.
> - `dhcp4d`, `dns` and `oui` carry an explicit `//go:build linux` tag so
>   cross-platform consumers can import the tree unconditionally.
> - The daemons under `cmd/` are decoupled from the gokrazy appliance runtime
>   and take a `--perm` state directory, so they also run under a normal host
>   OS (see `contrib/host-os/`).


[![GitHub Actions CI](https://github.com/rtr7/router7/actions/workflows/go.yml/badge.svg)](https://github.com/rtr7/router7/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/rtr7/router7/cmd?status.svg)](https://godoc.org/github.com/rtr7/router7/cmd)
[![Go Report Card](https://goreportcard.com/badge/github.com/rtr7/router7)](https://goreportcard.com/report/github.com/rtr7/router7)

router7 is a pure-Go implementation of a small home internet router. It comes with all the services required to make a [fiber7 internet connection](https://www.init7.net/en/internet/fiber7/) work (DHCPv4, DHCPv6, DNS, etc.).

Note that this project should be considered a (working!) tech demo. Feature requests will likely not be implemented, and see [CONTRIBUTING.md](CONTRIBUTING.md) for details about which contributions are welcome.

**For more details, please see [router7.org](https://router7.org/)**
