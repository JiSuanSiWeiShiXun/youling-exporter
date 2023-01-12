# pcap-exporter

## 编译
编译环境需要安装libpcap-dev
因为我的目标环境是ubuntu 所以执行 ```sudo apt install -y libpcap-dev```

## 部署
deploy/ 目录下的docker-compose.yml定义了测试环境，其中包含一个udp-client容器 和 一个包含udp-server以及pcap-exporter的容器

## 运行
./process-net-exporter -web.listen-address :6789 -pcap_filter "udp portrange 4800-5000"