package pcap

/*
在linux环境使用libpcap库采集网络流量数据
*/
import (
	"fmt"
	"net"
	"strings"

	"github.com/JiSuanSiWeiShiXun/pcap_exporter/collector"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	log "github.com/sirupsen/logrus"
)

var (
	iface  = "eth0"
	buffer = int32(1600)
	ipAddr string
)

func init() {
	var err error
	ipAddr, err = GetClientIp()
	if err != nil {
		log.Panicf("get ip addr failed: %v", err)
	}
	log.Infof("host ip addr is %v", ipAddr)
}

// GetClientIp 获取运行环境的ip地址
func GetClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("cannot find the client ip address")
}

func isWriteOrRead(srcIP, dstIP string) (isWrite, isRead bool) {
	if strings.Contains(srcIP, ipAddr) {
		isWrite = true
	}
	if strings.Contains(dstIP, ipAddr) {
		isRead = true
	}
	return
}

// deviceExists 查找设备网卡中是否包含目标网卡
func deviceExists(name string) bool {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panic(err)
	}
	for _, device := range devices {
		if device.Name == name {
			return true
		}
	}
	return false
}

// harvestData 获取目标数据
func harvestData(pc *collector.PcapCollector, packet gopacket.Packet) {
	app := packet.ApplicationLayer()
	if app != nil {
		// App is present
		srcIP := packet.NetworkLayer().NetworkFlow().Src()
		dstIP := packet.NetworkLayer().NetworkFlow().Dst()
		srcPort := packet.TransportLayer().TransportFlow().Src()
		dstPort := packet.TransportLayer().TransportFlow().Dst()

		srcAddr := fmt.Sprintf("%v:%v", srcIP, srcPort)
		dstAddr := fmt.Sprintf("%v:%v", dstIP, dstPort)
		isWrite, isRead := isWriteOrRead(srcAddr, dstAddr)
		networkProtocol := packet.NetworkLayer().LayerType().String()
		transportProtocol := packet.TransportLayer().LayerType().String()
		byteLen := uint64(len(app.LayerContents()))

		// 直接写入PcapCollector.NetRecordCache[name]
		var key string
		// 发送
		if isWrite {
			key = fmt.Sprintf("%v|%v|%v", srcAddr, networkProtocol, transportProtocol)
			r, ok := pc.NetRecordCache[key]
			if !ok {
				// 第一次记录
				pc.NetRecordCache[key] = collector.NetRecord{
					Name:     key,
					Protocol: transportProtocol,
				}
				r = pc.NetRecordCache[key]
			}
			r.SendBytesTotal += byteLen

			log.Tracef("[pcap send]\n [name]%v [protocol]%v [len]%v\n[pcap send]\n\n",
				key, transportProtocol, byteLen,
			)
		}

		// 接收
		if isRead {
			key = fmt.Sprintf("%v|%v|%v", dstAddr, networkProtocol, transportProtocol)
			r, ok := pc.NetRecordCache[key]
			if !ok {
				// 第一次记录
				pc.NetRecordCache[key] = collector.NetRecord{
					Name: key,
				}
				r = pc.NetRecordCache[key]
			}
			r.RecvBytesTotal += byteLen

			log.Tracef("[pcap recv]\n [name]%v [protocol]%v [len]%v\n[pcap recv]\n\n",
				key, transportProtocol, byteLen,
			)
		}
	}
}

func PacketCapture(pc *collector.PcapCollector, filter string) {
	if !deviceExists(iface) {
		log.Fatal("Unable to open device ", iface)
	}

	handler, err := pcap.OpenLive(iface, buffer, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handler.Close()

	if err := handler.SetBPFFilter(filter); err != nil {
		log.Fatal(err)
	}

	source := gopacket.NewPacketSource(handler, handler.LinkType())
	for packet := range source.Packets() {
		// fmt.Println(packet)
		harvestData(pc, packet)
	}
}
