package main

/*
#cgo CFLAGS: -I./nethogs/src
#cgo LDFLAGS: -L${SRCDIR}/lib -lnethogs
#include <stdlib.h>
#include "libnethogs.h"
extern void GoCallback(int action, NethogsMonitorRecord data);

static void WrapperFunc(int action, NethogsMonitorRecord const *data) {
	GoCallback(action, *data);
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
	"time"
	"unsafe"
)

type Record struct {
	Time       time.Time
	Name       string
	PID        int
	UID        int
	Port       int
	DeviceName string
	SentBytes  uint64
	RecvBytes  uint64
	SentKbs    float32
	RecvKbs    float32
}

//export GoCallback
func GoCallback(action C.int, data C.struct_NethogsMonitorRecord) {
	if action == 1 {
		r := new(Record)
		r.Time = time.Now()
		r.Name = C.GoString(data.name)
		if strings.Contains(r.Name, "unknown") {
			return
		}
		r.PID = int(data.pid)
		r.UID = int(data.uid)
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
		r.DeviceName = C.GoString(data.device_name)
		r.SentBytes = uint64(data.sent_bytes)
		r.RecvBytes = uint64(data.recv_bytes)
		r.SentKbs = float32(data.sent_kbs)
		r.RecvKbs = float32(data.recv_kbs)
		//NetCollector.Histogram[r.PID] = r
		fmt.Printf("[tick]\n[Name]%v [PID]%v [UID]%v [DEV]%v [sdbt]%v [rcbt]%v [sdkb]%v [rckb]%v\n%v\n[endtick]\n\n", r.Name, r.PID, r.UID, r.DeviceName, r.SentBytes, r.RecvBytes, r.SentKbs, r.RecvKbs, r)
	}
}

func callNethogs() {
	// 动态链接
	var cstring *C.char = C.CString("")
	defer C.free(unsafe.Pointer(cstring))
	if record := C.CallMyFunction(cstring); record != 0 {
		fmt.Printf("wrong return\n")
	}
}
