
# dnss

dnss is a daemon for using DNS over HTTPS.

It can act as a proxy (the most common use case), and as a server (in case you
want end to end control).


[![Build Status](https://travis-ci.org/albertito/dnss.svg?branch=master)](https://travis-ci.org/albertito/dnss)
[![Go Report Card](https://goreportcard.com/badge/github.com/albertito/dnss)](https://goreportcard.com/report/github.com/albertito/dnss)

[![](https://images.microbadger.com/badges/image/containerfactory/dnss.svg)](https://microbadger.com/images/containerfactory/dnss "Get your own image badge on microbadger.com")


## Features

* Supports the JSON-based protocol as implemented by
  [dns.google.com](https://dns.google.com)
  ([reference](https://developers.google.com/speed/public-dns/docs/dns-over-https)).
* Supports the [DNS Queries over HTTPS
  (DoH)](https://tools.ietf.org/html/draft-ietf-doh-dns-over-https) proposed
  standard (and implemented by [Cloudflare's 1.1.1.1](https://1.1.1.1/)).
* Local cache (optional).
* HTTP(s) proxy support, autodetected from the environment.
* Monitoring HTTP server, with exported variables and tracing to help
  debugging.
* Separate resolution for specific domains, useful for home networks with
  local DNS servers.


## Install

Run as a docker container. It accepts all options like the standalone version on the command line

```shell
docker run -it --rm containerfactory/dnss  --help
```

## Examples

### DNS server (proxy mode)

Listens on port 53 for DNS queries, resolves them using the given HTTPS URL.

```shell
# Use the default HTTPS URL (currently, dns.google.com):
dnss -enable_dns_to_https

# Use Cloudflare's 1.1.1.1:
dnss -enable_dns_to_https -https_upstream="https://1.1.1.1/dns-query"

# Use Google's dns.google.com:
dnss -enable_dns_to_https -https_upstream="https://dns.google.com/resolve"
```

### HTTPS server

Receives DNS over HTTPS requests, resolves them using the machine's configured
DNS servers, and returns the replies.  You will need to have certificates for
the domains you want to serve.

Supports both DoH and JSON modes automatically, and the endpoints are
`/dns-query` and `/resolve`.

```shell
# Serve DNS over HTTPS requests, take certificates from letsencrypt.
DOMAIN=yourdomain.com
dnss -enable_https_to_dns \
  -https_key=/etc/letsencrypt/live/$DOMAIN/privkey.pem \
  -https_cert=/etc/letsencrypt/live/$DOMAIN/fullchain.pem
```

