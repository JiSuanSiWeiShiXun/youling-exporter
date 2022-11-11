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
	"math"
	"strings"
	"sync/atomic"
	"unsafe"
	"youling-exporter/collector"
)

//export Callback
func Callback(action C.int, data C.struct_NethogsMonitorRecord) {
	if action == 1 {
		pid := int(data.pid)
		if _, ok := nethogsCollector.RecordMap[pid]; !ok {
			nethogsCollector.RecordMap[pid] = new(collector.NethogsMonitorRecord)
		}
		//r.Time = time.Now() // 需要记录时间吗？
		nethogsCollector.RecordMap[pid].Name = C.GoString(data.name)
		if strings.Contains(nethogsCollector.RecordMap[pid].Name, "unknown") {
			return
		}
		nethogsCollector.RecordMap[pid].PID = int(data.pid)
		nethogsCollector.RecordMap[pid].UID = uint32(data.uid)
		//r.Port = int(data.port)
		// 筛选数据
		//flag := false
		//for _, port := range NetCollector.Ports {
		//	if port == r.Port {
		//		flag = true
		//	}
		//}
		//if !flag {
		//	return
		//}
		// TODO: 处理竞争
		nethogsCollector.RecordMap[pid].DeviceName = C.GoString(data.device_name)
		atomic.StoreUint64(&nethogsCollector.RecordMap[pid].SentBytes, uint64(data.sent_bytes))
		atomic.StoreUint64(&nethogsCollector.RecordMap[pid].RecvBytes, uint64(data.recv_bytes))
		atomic.StoreUint64(&nethogsCollector.RecordMap[pid].SentKBs, math.Float64bits(float64(data.sent_kbs)))
		atomic.StoreUint64(&nethogsCollector.RecordMap[pid].RecvKBs, math.Float64bits(float64(data.recv_kbs)))
		//atomic.StoreUint64(&g.valBits, math.Float64bits(val))
		//NetCollector.Histogram[r.PID] = r
		fmt.Printf("[tick]\n[Name]%v [PID]%v [UID]%v [DEV]%v [sdbt]%v [rcbt]%v [sdkb]%v [rckb]%v\n%v\n[endtick]\n\n",
			nethogsCollector.RecordMap[pid].Name,
			nethogsCollector.RecordMap[pid].PID,
			nethogsCollector.RecordMap[pid].UID,
			nethogsCollector.RecordMap[pid].DeviceName,
			nethogsCollector.RecordMap[pid].SentBytes,
			nethogsCollector.RecordMap[pid].RecvBytes,
			nethogsCollector.RecordMap[pid].SentKBs,
			nethogsCollector.RecordMap[pid].RecvKBs,
			nethogsCollector.RecordMap[pid])
	}
}

func CallNethogs() {
	// 动态链接
	var cstring *C.char = C.CString("")
	defer C.free(unsafe.Pointer(cstring))
	if record := C.CallMyFunction(cstring); record != 0 {
		fmt.Printf("wrong return\n")
	}
}
