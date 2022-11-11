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
)

//export Callback
func Callback(action C.int, data C.struct_NethogsMonitorRecord) {
	if action == 1 {
		//r.Time = time.Now() // 需要记录时间吗？
		nethogsCollector.Record.Name = C.GoString(data.name)
		if strings.Contains(nethogsCollector.Record.Name, "unknown") {
			return
		}
		nethogsCollector.Record.PID = int(data.pid)
		nethogsCollector.Record.UID = uint32(data.uid)
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
		nethogsCollector.Record.DeviceName = C.GoString(data.device_name)
		nethogsCollector.Record.SentBytes = uint64(data.sent_bytes)
		nethogsCollector.Record.RecvBytes = uint64(data.recv_bytes)
		nethogsCollector.Record.SentKBs = float64(data.sent_kbs)
		nethogsCollector.Record.RecvKBs = float64(data.recv_kbs)
		//atomic.StoreUint64(&g.valBits, math.Float64bits(val))
		//NetCollector.Histogram[r.PID] = r
		fmt.Printf("[tick]\n[Name]%v [PID]%v [UID]%v [DEV]%v [sdbt]%v [rcbt]%v [sdkb]%v [rckb]%v\n%v\n[endtick]\n\n",
			nethogsCollector.Record.Name,
			nethogsCollector.Record.PID,
			nethogsCollector.Record.UID,
			nethogsCollector.Record.DeviceName,
			nethogsCollector.Record.SentBytes,
			nethogsCollector.Record.RecvBytes,
			nethogsCollector.Record.SentKBs,
			nethogsCollector.Record.RecvKBs,
			nethogsCollector.Record)
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
