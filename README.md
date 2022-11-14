# youling-exporter
to take advantage of dynamic library of nethogs to collect network resource usage of specific processes.

因为用到了CGO链接C库，所以目录结构非常重要，不要乱动哦亲

## 编译
```go build```

## 使用
```
Usage of ./youling-exporter:
  -pcap_filter string
        pcap-filter passed to nethogs.
  -web.listen-address string
        Address on which to expose metrics and web interface. (default ":9356")
```

## 示例
```./youling-exporter -web.listen-address :6789 -pcap_filter "udp portrange 4800-5000"```
