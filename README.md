# youling-exporter
to take advantage of dynamic library of nethogs to collect network resource usage of specific processes.

因为用到了CGO链接C库，所以目录结构非常重要，不要乱动哦亲

## 编译
```
./build.sh
```

## 使用
```
# 首先需要将.so所在文件夹加入LD_LIBRARY_PATH，以root权限执行
export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:<your_path>/lib/"
```
```
Usage of ./youling-exporter:
  -pcap_filter string
        pcap-filter passed to nethogs.
  -web.listen-address string
        Address on which to expose metrics and web interface. (default ":9356")
```

## 示例
- **请务必以root权限执行**
- ```./process-net-exporter -web.listen-address :6789 -pcap_filter "udp portrange 4800-5000"```
- ```nohup ./process-net-exporter -pcap_filter "udp port 4869" > output.log 2>&1 &```
