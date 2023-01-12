#!/bin/sh
cd ../ && go build -o process-net-exporter

mkdir process_net_exporter
mv process-net-exporter process_net_exporter/
tar -zcvf process_net_exporter.tar.gz process_net_exporter/
#rm -rf process_net_exporter/