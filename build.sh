#!/bin/bash

# 编译nethogs项目，移动动态链接库到lib文件夹
cp lib/libnethogs.so lib/libnethogs.so.

# 编译golang项目（需要预装go 1.19环境）
go build -o process-net-exporter

# 将lib依赖库文件 和 可执行文件打包到一起，所以最终可执行文件在process_net_exporter.tar.gz里
mkdir process_net_exporter
mv process-net-exporter process_net_exporter
cp -r lib/ process_net_exporter/
tar -zcvf process_net_exporter.tar.gz process_net_exporter/
# rm -rf process_net_exporter/
