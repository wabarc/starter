# Starter

Starter is a tool for packaging and storing Chrome extensions are used with [Xvfb](https://en.wikipedia.org/wiki/Xvfb).

## Background

Websites may contain challenges such as paywalls and various CAPTCHAs, and quite well solutions including 
[`puppeteer-extra`](https://github.com/berstend/puppeteer-extra) do not work in all cases. It is hoped that this project will 
provide a solution that, to some extent, reproduces the real browser workplace.

## How it works

1. Packaging Chrome extensions into the binary.
2. Launching Chrome to load extensions with a specified directory to stores user data.
3. Launch a remote debugging browser with the flag `--user-data-dir={workspace}/UserDataDir`.

Note: remote debugging address is default to listen on `0.0.0.0:9222`.

## Prerequisite

- [Go](https://golang.org/) (requires Go 1.16 or later)
- [GNU Make](https://www.gnu.org/software/make/manual/make.html)
- Docker/Podman

## Building

File the secrets of buster to `misc/secrets.json`

```shell
make build
```

## Installation

```sh
sh <(wget https://github.com/wabarc/starter/raw/main/install.sh -O-)
```

## Extensions

- [x] [bypass-paywalls](https://github.com/iamadamdev/bypass-paywalls-chrome)
- [ ] [dessant/buster](https://github.com/dessant/buster)
- [ ] [gorhill/uBlock](https://github.com/gorhill/uBlock)
- [ ] [privacypass/challenge-bypass-extension](https://github.com/privacypass/challenge-bypass-extension)
- [ ] [n4cr/cookiepopupblocker](https://github.com/n4cr/cookiepopupblocker)
- [ ] [Sainan/Universal-Bypass](https://github.com/Sainan/Universal-Bypass)

## License

This software is released under the terms of the GNU General Public License v3.0. See the [LICENSE](https://github.com/wabarc/starter/blob/main/LICENSE) file for details.
