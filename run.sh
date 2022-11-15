#!/bin/sh
# 脚本调用方式：run.sh /opt/process_net_exporter/process-net-exporter :6789 "udp port 4869"
# 前提是可执行文件{process_net_exporter}，已经在目标地址${1}了
# 目录结构e.g ${1} = /opt/process_net_exporter/
# /opt/process_net_exporter
# ---lib/
#    ---libnethogs.so.
#    ---libnethogs.so
# ---process_net_exporter
# ---run.sh
# ---...

# debian/ubuntu环境 预装nethogs外部依赖
# sudo apt-get install build-essential libncurses5-dev libpcap-dev

# 切换为sudo权限，启动exporter
sudo su -
echo "可执行文件路径 ${1}"
echo "exporter监听端口 ${2}"
echo "pcap_filter ${3}"

# nohup /opt/process_net_exporter/process-net-exporter -web.listen-address :6789 -pcap_filter "udp port 4869" > output.log 2>&1 &
echo "sudo nohup ${1} -web.listen-address ${2} -pcap_filter ${3} > output.log 2>&1 &"
sudo nohup ${1} -web.listen-address ${2} -pcap_filter ${3} > output.log 2>&1 &