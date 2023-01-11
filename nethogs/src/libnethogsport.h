//include <stdlib.h>
//include "libnethogs.h"

#ifndef LIBNETHOGSPORT_H_
#define LIBNETHOGSPORT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>

#define NETHOGS_DSO_VISIBLE __attribute__((visibility("default")))
#define NETHOGS_DSO_HIDDEN __attribute__((visibility("hidden")))

#define NETHOGS_APP_ACTION_SET 1
#define NETHOGS_APP_ACTION_REMOVE 2

#define NETHOGS_STATUS_OK 0
#define NETHOGS_STATUS_FAILURE 1
#define NETHOGS_STATUS_NO_DEVICE 2

typedef struct NethogsPortMonitorRecord {
	int record_id;
	const char* name;
	int pid;
	uint32_t uid;
	int port;
	const char* device_name;
	uint64_t sent_bytes;
	uint64_t recv_bytes;
	float sent_kbs;
	float recv_kbs;
} NethogsPortMonitorRecord;


typedef void (*NethogsPortMonitorCallback)(int action,
	NethogsPortMonitorRecord const* data);


NETHOGS_DSO_VISIBLE int nethogsportmonitor_loop(NethogsPortMonitorCallback cb,
	char* filter,
	int to_ms);

NETHOGS_DSO_VISIBLE int nethogsportmonitor_loop_devices(NethogsPortMonitorCallback cb,
	char* filter, int devc,
	char** devicenames,
	bool all,
	int to_ms);


NETHOGS_DSO_VISIBLE void nethogsportmonitor_breakloop();

#undef NETHOGS_DSO_VISIBLE
#undef NETHOGS_DSO_HIDDEN

#ifdef __cplusplus
}
#endif

#endif // LIBNETHOGSPORT_H_