package main

// 构建的时候要注意 nethogs 库文件目录中有静态链接文件(*.a)时，会报错，因为所依赖的 libpcap 不提供静态链接途径。所以要么提前删除默认生成的静态链接文件，要么在构建的时候指定为动态链接

//#cgo CFLAGS: -I./nethogs/src
//#cgo LDFLAGS: -L${SRCDIR}/nethogs/src -lnethogs
//#include <stdlib.h>
//#include <libnethogs.h>
/*
extern void GoCallback(int action, NethogsPortMonitorRecord data);


static void WrapperFunc(int action, NethogsPortMonitorRecord const *data) {
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
func GoCallback(action C.int, data C.struct_NethogsPortMonitorRecord) {
	if action == 1 {
		r := new(Record)
		r.Time = time.Now()
		r.Name = C.GoString(data.name)
		if strings.Contains(r.Name, "unknown") {
			return
		}
		r.PID = int(data.pid)
		r.UID = int(data.uid)
		r.Port = int(data.port)
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
		fmt.Printf("[tick]%v\n[endtick]\n\n", r)
	}
}

func main() {
	// 动态链接
	var cstring *C.char = C.CString("")
	defer C.free(unsafe.Pointer(cstring))
	if record := C.CallMyFunction(cstring); record != 0 {
		fmt.Printf("wrong return\n")
	}
}
