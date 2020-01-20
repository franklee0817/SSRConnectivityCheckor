# SSRConnectivityCheckor

This A SSR node connectivity checkor application.[中文版说明](https://github.com/franklee0817/SSRConnectivityCheckor/blob/master/README.zh.md)

# Use Guide

Compile go executable binary. Move the binary to `/usr/local/ssr_connectivity_scanner` and move `check-ssr` to `/usr/local/bin/` after you modified the content of `check-ssr.sh`. Default config file is in `~/.ssr_scanner_conf` like:

```json
{
  "config_file": "",
  "subscribe_url": ""
}
```

config_file: the exported json config file path of ShadowSocksR
subscribe_url: the subscribe url of ssr server.

You can use `check-ssr -f ssrJsonConfigFileName` and `check-ssr -s subscribeURL` as well.

Enjoy your SSR connectivity checkor
