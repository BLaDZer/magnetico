# magnetico
*Autonomous (self-hosted) BitTorrent DHT search engine suite.*

[![Go](https://github.com/tgragnato/magnetico/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/tgragnato/magnetico/actions/workflows/go.yml)
[![CodeQL](https://github.com/tgragnato/magnetico/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/tgragnato/magnetico/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tgragnato/magnetico)](https://goreportcard.com/report/github.com/tgragnato/magnetico)

![Flow of Operations](/doc/operations.svg)

magnetico is the first autonomous (self-hosted) BitTorrent DHT search engine suite that is *designed for end-users*. The suite consists of a single binary with two components:

- a crawler for the BitTorrent DHT network, which discovers info hashes and fetches metadata from the peers.
- a lightweight web interface for searching and browsing the torrents discovered by its counterpart.

This allows anyone with a decent Internet connection to access the vast amount of torrents waiting to be discovered within the BitTorrent DHT space, *without relying on any central entity*.

**magnetico** liberates BitTorrent from the yoke of centralised trackers & web-sites and makes it
*truly decentralised*. Finally!

## Easy Run and Compilation

The easiest way to run magnetico on amd64 platforms is to use the OCI image built within the CI pipeline:
- `docker pull ghcr.io/tgragnato/magnetico:latest`
- `docker run --rm -it ghcr.io/tgragnato/magnetico:latest --help`
- `docker run --rm -it -v <your_data_dir>:/data -p 8080:8080/tcp ghcr.io/tgragnato/magnetico:latest --database=sqlite3:///data/magnetico.sqlite3 --max-rps=1000 --addr=0.0.0.0:8080`

To compile using the standard Golang toolchain:
- Download the latest golang release from [the official website](https://go.dev/dl/)
- Follow the [installation instructions for your platform](https://go.dev/doc/install)
- Run `go install --tags fts5 .`
- The `magnetico` binary is now available in your `$GOBIN` directory

## Features

Easy installation & minimal requirements:
  - Easy to build golang static binaries.
  - Root access is *not* required to install or to use.

**magnetico** trawls the BitTorrent DHT by "going" from one node to another, and fetches the metadata using the nodes without using trackers. No reliance on any centralised entity!

Unlike client-server model that web applications use, P2P networks are *chaotic* and **magnetico** is designed to handle all the operational errors accordingly.

High performance implementation in Go: **magnetico** utilizes every bit of your resources to discover as many infohashes & metadata as possible.

**magnetico** features a lightweight web interface to help you access the database without getting on your way.

If you'd like to password-protect the access to **magnetico**, you need to store the credentials
in file. The `credentials` file must consist of lines of the following format: `<USERNAME>:<BCRYPT HASH>`.

- `<USERNAME>` must start with a small-case (`[a-z]`) ASCII character, might contain non-consecutive underscores except at the end, and consists of small-case a-z characters and digits 0-9.
- `<BCRYPT HASH>` is the output of the well-known bcrypt function.

You can use `htpasswd` (part of `apache2-utils` on Ubuntu) to create lines:

```
$  htpasswd -bnBC 12 "USERNAME" "PASSWORD"
USERNAME:$2y$12$YE01LZ8jrbQbx6c0s2hdZO71dSjn2p/O9XsYJpz.5968yCysUgiaG
```

### Screenshots

| ![The Homepage](/doc/homepage.png) | ![Searching for torrents](/doc/search.png) | ![Search result](/doc/result.png) |
|:-------------------------------------------------------------------------------------------------------------------------------------------------------:|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------:|:---------------------------------------------------------------------------------------------------------------------------------------------:|
|                                                                     __The Homepage__                                                                    |                                                                     __Searching for torrents__                                                                    |                                                     __Viewing the metadata of a torrent__                                                     |

## Why?
BitTorrent, being a distributed P2P file sharing protocol, has long suffered because of the
centralised entities that people depended on for searching torrents (websites) and for discovering
other peers (trackers). Introduction of DHT (distributed hash table) eliminated the need for
trackers, allowing peers to discover each other through other peers and to fetch metadata from the
leechers & seeders in the network. **magnetico** is the finishing move that allows users to search
for torrents in the network, hence removing the need for centralised torrent websites.
