package pcap

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestGetClientIP(t *testing.T) {
	if ip, err := GetClientIp(); err != nil {
		panic(err)
	} else {
		log.Println(ip)
	}
}
