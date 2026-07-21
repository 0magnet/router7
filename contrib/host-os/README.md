# router7 under a host OS (standalone port)

Upstream router7 is a [Gokrazy](https://gokrazy.org/) appliance — it *replaces*
the operating system with just the Go router programs. This `feat/host-os-standalone`
port removes that coupling so the daemons run as **ordinary Go programs under an
existing Linux distro** (Debian, Armbian, etc.) alongside whatever else the box
does — no reflash, no appliance image.

## What changed vs. upstream

1. **No gokrazy runtime dependency.** The daemons imported
   `github.com/gokrazy/gokrazy` only for `PrivateInterfaceAddrs()`,
   `IsInPrivateNet()` (both live in the pure-Go `gokrazy/ifaddr` subpackage —
   swapped to those) and `DontStartOnBoot()` (appliance-only — dropped; the
   service manager's restart policy governs a missing-config exit).
2. **Configurable state dir.** Every daemon takes `--perm` (default `/perm`),
   so state (DHCP leases, DNS aliases, DUID, prefixes) lives at a normal path
   like `/var/lib/router7` instead of gokrazy's persistent partition.
3. **Configurable interfaces.** The hardcoded `lan0`/`uplink0` names are now
   flags (`--interface`, `--lan`). Interface *roles* remain file-driven via
   `interfaces.json` under `--perm`, which is exactly what lets the router
   serve a *chosen* interface and coexist with the host's own networking.

## The two daemons that matter for a Skywire VPN router

For a Skywire VPN-router the router only needs to hand out addresses and resolve
names on the downstream LAN — Skywire itself owns the uplink (the vpn-client
`tun`) and the NAT. That's just:

- **`dhcp4d`** — the DHCPv4 server (leases + a status page + Prometheus metrics)
- **`dnsd`**  — the DNS server (resolves LAN hostnames from the leases)

These replace the `dnsmasq` the Skywire `vpn-router` app currently shells out
to — pure Go, embeddable, one fewer external dependency. `radvd`/`dhcp6`/
`netconfigd` are there for a full dual-stack standalone router but aren't needed
for the Skywire slice.

## Run it standalone

```sh
go build -o /usr/local/bin/router7-dhcp4d ./cmd/dhcp4d
go build -o /usr/local/bin/router7-dnsd   ./cmd/dnsd
sudo mkdir -p /var/lib/router7

# serve DHCP+DNS on eth1 (the LAN side); state under /var/lib/router7
sudo /usr/local/bin/router7-dhcp4d --interface eth1 --perm /var/lib/router7 &
sudo /usr/local/bin/router7-dnsd   --lan eth1       --perm /var/lib/router7 &
```

Or via the systemd templates in this directory (`%i` = the LAN interface):

```sh
sudo cp router7-dhcp4d@.service router7-dnsd@.service /etc/systemd/system/
sudo systemctl enable --now router7-dhcp4d@eth1 router7-dnsd@eth1
```

## Config vs. GUI

router7 has **no configuration GUI** by design — it's configured by writing
JSON files under `--perm` (`interfaces.json`, `dnsd/aliases.json`,
`radvd/prefixes.json`, …) and the daemons reload. The HTTP endpoints they serve
are **read-only status** (leases table, DNS stats, Prometheus `/metrics`) plus a
couple of narrow operational POSTs (`/sethostname`, `/dyndns`). That's ideal for
embedding: Skywire configures the router by writing those files, with zero GUI
coupling, and can surface the status pages in its own hypervisor UI if useful.
