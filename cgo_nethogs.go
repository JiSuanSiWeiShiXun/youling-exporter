package main

/*
#cgo CFLAGS: -I./nethogs/src
#cgo LDFLAGS: -L${SRCDIR}/lib -lnethogs
#include <stdlib.h>
#include "libnethogs.h"
extern void Callback(int action, NethogsMonitorRecord data);

static void WrapperFunc(int action, NethogsMonitorRecord const *data) {
Callback(action, *data);
}

  // normally you will have to define function or variables
  // in another separate C file to avoid the multiple definition
  // errors, however, using "static inline" is a nice workaround
  // for simple functions like this one.
static inline int CallMyFunction(char *filter) {
int rc = nethogsmonitor_loop(WrapperFunc, filter, 0);
	return rc;
}
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	"nethogs-exporter/collector"

	log "github.com/sirupsen/logrus"
)

//export Callback
func Callback(action C.int, data C.struct_NethogsMonitorRecord) {
	if action != 1 {
		log.Warnf("get data from nethogs error: action=%v", action)
	}
	//r.Time = time.Now() // 需要记录时间吗？
	record := new(collector.NethogsMonitorRecord)
	record.Name = C.GoString(data.name)
	if strings.Contains(record.Name, "unknown") {
		return
	}
	record.PID = int(data.pid)
	record.UID = uint32(data.uid)
	record.DeviceName = C.GoString(data.device_name)
	// TODO: 原子操作姿势不对
	//atomic.StoreUint64(&nethogsCollector.RecordMap[pid].SentBytes, uint64(data.sent_bytes))
	//atomic.StoreUint64(&nethogsCollector.RecordMap[pid].RecvBytes, uint64(data.recv_bytes))
	//atomic.StoreUint64(&nethogsCollector.RecordMap[pid].SentKBs, math.Float64bits(float64(data.sent_kbs)))
	//atomic.StoreUint64(&nethogsCollector.RecordMap[pid].RecvKBs, math.Float64bits(float64(data.recv_kbs)))
	record.SentBytes = uint64(data.sent_bytes)
	record.RecvBytes = uint64(data.recv_bytes)
	record.SentKBs = float64(data.sent_kbs)
	record.RecvKBs = float64(data.recv_kbs)
	scrapeChan <- record

	// 别在C里用go的打印
	//log.Debugf("[tick]\n[key]\"%v\" [Name]%v [PID]%v [UID]%v [DEV]%v [sdbt]%v [rcbt]%v [sdkb]%v [rckb]%v\n%v\n[endtick]\n\n",
	//	key,
	//	nethogsCollector.RecordMap[key].Name,
	//	nethogsCollector.RecordMap[key].PID,
	//	nethogsCollector.RecordMap[key].UID,
	//	nethogsCollector.RecordMap[key].DeviceName,
	//	nethogsCollector.RecordMap[key].SentBytes,
	//	nethogsCollector.RecordMap[key].RecvBytes,
	//	nethogsCollector.RecordMap[key].SentKBs,
	//	nethogsCollector.RecordMap[key].RecvKBs,
	//	nethogsCollector.RecordMap[key])
}

func CallNethogs(str string) {
	// 动态链接
	var cstring *C.char = C.CString(str)
	defer C.free(unsafe.Pointer(cstring))
	if record := C.CallMyFunction(cstring); record != 0 {
		fmt.Printf("wrong return\n")
	}
}
