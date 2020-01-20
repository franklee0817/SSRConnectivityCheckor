# SSRConnectivityCheckor

这是一个检查 SSR 节点和端口联通性的程序，帮你找到连通性最好的前 20 个节点

# 使用说明

编译代码，然后将可执行文件重命名后放到`/usr/local/ssr_connectivity_scanner/check-ssr`，将`check-ssr` 放到 `/usr/local/bin/`。程序有个自定义的配置文件 `~/.ssr_scanner_conf`， 文件格式如下：

```json
{
  "config_file": "",
  "subscribe_url": ""
}
```

config_file: ssr json 配置文件地址，格式为 ShardowSocksR 导出的配置文件格式
subscribe_url: ssr 订阅地址

使用配置文件则可以通过 check-ssr 命令直接执行扫描，否则请执行 check-ssr -h 查看帮助
