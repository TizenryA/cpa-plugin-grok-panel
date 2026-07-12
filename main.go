package main

/*
#include <stdint.h>
#include <stdlib.h>

typedef struct {
	void* ptr;
	size_t len;
} cliproxy_buffer;

typedef int (*cliproxy_host_call_fn)(void*, const char*, const uint8_t*, size_t, cliproxy_buffer*);
typedef void (*cliproxy_host_free_fn)(void*, size_t);

typedef struct {
	uint32_t abi_version;
	void* host_ctx;
	cliproxy_host_call_fn call;
	cliproxy_host_free_fn free_buffer;
} cliproxy_host_api;

typedef int (*cliproxy_plugin_call_fn)(char*, uint8_t*, size_t, cliproxy_buffer*);
typedef void (*cliproxy_plugin_free_fn)(void*, size_t);
typedef void (*cliproxy_plugin_shutdown_fn)(void);

typedef struct {
	uint32_t abi_version;
	cliproxy_plugin_call_fn call;
	cliproxy_plugin_free_fn free_buffer;
	cliproxy_plugin_shutdown_fn shutdown;
} cliproxy_plugin_api;

extern int cliproxyPluginCall(char*, uint8_t*, size_t, cliproxy_buffer*);
extern void cliproxyPluginFree(void*, size_t);
extern void cliproxyPluginShutdown(void);

static const cliproxy_host_api* stored_host;

static void store_host_api(const cliproxy_host_api* host) {
	stored_host = host;
}

static int call_host_api(const char* method, const uint8_t* request, size_t request_len, cliproxy_buffer* response) {
	if (stored_host == NULL || stored_host->call == NULL) {
		return 1;
	}
	return stored_host->call(stored_host->host_ctx, method, request, request_len, response);
}

static void free_host_buffer(void* ptr, size_t len) {
	if (stored_host != NULL && stored_host->free_buffer != NULL && ptr != NULL) {
		stored_host->free_buffer(ptr, len);
	}
}
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unsafe"
)

const abiVersion uint32 = 1

// ---- Envelope ----

type envelope struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *envelopeError  `json:"error,omitempty"`
}

type envelopeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ---- Registration ----

type registration struct {
	SchemaVersion uint32                   `json:"schema_version"`
	Metadata      metadata                 `json:"metadata"`
	Capabilities  registrationCapabilities `json:"capabilities"`
}

type metadata struct {
	Name             string        `json:"Name"`
	Version          string        `json:"Version"`
	Author           string        `json:"Author"`
	GitHubRepository string        `json:"GitHubRepository"`
	Logo             string        `json:"Logo"`
	ConfigFields     []configField `json:"ConfigFields"`
}

type configField struct {
	Key         string `json:"Key"`
	Label       string `json:"Label"`
	Type        string `json:"Type"`
	Required    bool   `json:"Required"`
	Description string `json:"Description"`
}

type registrationCapabilities struct {
	ManagementAPI bool `json:"management_api"`
}

// ---- Management ----

type managementResource struct {
	Path        string `json:"Path"`
	Menu        string `json:"Menu"`
	Description string `json:"Description"`
}

type managementRegistration struct {
	Resources []managementResource `json:"resources,omitempty"`
}

type managementRequest struct {
	Method         string      `json:"Method"`
	Path           string      `json:"Path"`
	Headers        http.Header `json:"Headers"`
	Query          url.Values  `json:"Query"`
	Body           []byte      `json:"Body"`
	HostCallbackID string      `json:"host_callback_id,omitempty"`
}

type managementResponse struct {
	StatusCode int         `json:"StatusCode"`
	Headers    http.Header `json:"Headers"`
	Body       []byte      `json:"Body"`
}

// ---- Host auth list ----

type authListResponse struct {
	Files []authFile `json:"files"`
}

type authFile struct {
	Account        string          `json:"account"`
	AccountType    string          `json:"account_type"`
	AuthIndex      string          `json:"auth_index"`
	CreatedAt      string          `json:"created_at"`
	Disabled       bool            `json:"disabled"`
	Email          string          `json:"email"`
	Failed         int             `json:"failed"`
	ID             string          `json:"id"`
	Label          string          `json:"label"`
	LastRefresh    string          `json:"last_refresh"`
	Name           string          `json:"name"`
	Provider       string          `json:"provider"`
	RecentRequests []recentRequest `json:"recent_requests"`
	Size           int64           `json:"size"`
	Status         string          `json:"status"`
	Success        int             `json:"success"`
	Type           string          `json:"type"`
	Unavailable    bool            `json:"unavailable"`
	UpdatedAt      string          `json:"updated_at"`
}

type recentRequest struct {
	Time    string `json:"time"`
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
}

// ---- Plugin stats (returned to browser) ----

type pluginStats struct {
	TotalFiles    int          `json:"total_files"`
	ActiveFiles   int          `json:"active_files"`
	DisabledNum   int          `json:"disabled_files"`
	TotalSuccess  int          `json:"total_success"`
	TotalFailed   int          `json:"total_failed"`
	Files         []fileStats  `json:"files"`
	RecentBuckets []bucketStat `json:"recent_buckets"`
}

type fileStats struct {
	Email    string `json:"email"`
	Status   string `json:"status"`
	Disabled bool   `json:"disabled"`
	Success  int    `json:"success"`
	Failed   int    `json:"failed"`
}

type bucketStat struct {
	Time    string `json:"time"`
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
}

// ---- Plugin entry points ----

func main() {}

//export cliproxy_plugin_init
func cliproxy_plugin_init(host *C.cliproxy_host_api, plugin *C.cliproxy_plugin_api) C.int {
	if plugin == nil {
		return 1
	}
	C.store_host_api(host)
	plugin.abi_version = C.uint32_t(abiVersion)
	plugin.call = C.cliproxy_plugin_call_fn(C.cliproxyPluginCall)
	plugin.free_buffer = C.cliproxy_plugin_free_fn(C.cliproxyPluginFree)
	plugin.shutdown = C.cliproxy_plugin_shutdown_fn(C.cliproxyPluginShutdown)
	return 0
}

//export cliproxyPluginCall
func cliproxyPluginCall(method *C.char, request *C.uint8_t, requestLen C.size_t, response *C.cliproxy_buffer) C.int {
	if response != nil {
		response.ptr = nil
		response.len = 0
	}
	if method == nil {
		writeResponse(response, errorEnvelope("invalid_method", "method is required"))
		return 1
	}
	var requestBytes []byte
	if request != nil && requestLen > 0 {
		requestBytes = C.GoBytes(unsafe.Pointer(request), C.int(requestLen))
	}
	raw, errHandle := handleMethod(C.GoString(method), requestBytes)
	if errHandle != nil {
		writeResponse(response, errorEnvelope("plugin_error", errHandle.Error()))
		return 1
	}
	writeResponse(response, raw)
	return 0
}

//export cliproxyPluginFree
func cliproxyPluginFree(ptr unsafe.Pointer, len C.size_t) {
	if ptr != nil {
		C.free(ptr)
	}
	_ = len
}

//export cliproxyPluginShutdown
func cliproxyPluginShutdown() {}

// ---- Method dispatch ----

func handleMethod(method string, request []byte) ([]byte, error) {
	switch method {
	case "plugin.register", "plugin.reconfigure":
		return okEnvelope(registration{
			SchemaVersion: 1,
			Metadata: metadata{
				Name:             "grok-panel",
				Version:          "1.0.1",
				Author:           "tizenry",
				GitHubRepository: "https://github.com/TizenryA",
				Logo:             "",
				ConfigFields:     []configField{},
			},
			Capabilities: registrationCapabilities{
				ManagementAPI: true,
			},
		})
	case "management.register":
		return okEnvelope(managementRegistration{
			Resources: []managementResource{
				{
					Path:        "/panel",
					Menu:        "Grok 面板",
					Description: "Grok 账号用量统计面板",
				},
				{
					Path:        "/panel/data",
					Menu:        "",
					Description: "Grok 面板统计数据",
				},
			},
		})

	case "management.handle":
		return handleManagement(request)
	default:
		return errorEnvelope("unknown_method", "unknown method: "+method), nil
	}
}

// ---- Management handler ----

func handleManagement(raw []byte) ([]byte, error) {
	var req managementRequest
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &req); err != nil {
			return nil, fmt.Errorf("decode management request: %w", err)
		}
	}

	path := req.Path
	// Serve data API at /data, HTML at everything else
	if strings.HasSuffix(strings.TrimRight(path, "/"), "/data") {
		return handleData()
	}
	return handleHTML()
}

func handleData() ([]byte, error) {
	// Call host.auth.list to get auth files
	result, err := callHost("host.auth.list", map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("host.auth.list: %w", err)
	}

	var authResp authListResponse
	if err := json.Unmarshal(result, &authResp); err != nil {
		return nil, fmt.Errorf("decode auth list: %w", err)
	}

	// Filter xai files and compute stats
	var xaiFiles []authFile
	for _, f := range authResp.Files {
		if f.Provider == "xai" {
			xaiFiles = append(xaiFiles, f)
		}
	}

	stats := pluginStats{
		TotalFiles:   len(xaiFiles),
		TotalSuccess: 0,
		TotalFailed:  0,
	}

	activeCount := 0
	disabledCount := 0
	for _, f := range xaiFiles {
		if f.Disabled {
			disabledCount++
		} else if f.Status == "active" {
			activeCount++
		}
		stats.TotalSuccess += f.Success
		stats.TotalFailed += f.Failed
		stats.Files = append(stats.Files, fileStats{
			Email:    f.Email,
			Status:   f.Status,
			Disabled: f.Disabled,
			Success:  f.Success,
			Failed:   f.Failed,
		})
	}
	stats.ActiveFiles = activeCount
	stats.DisabledNum = disabledCount

	// Aggregate recent request buckets
	bucketMap := map[string]*bucketStat{}
	for _, f := range xaiFiles {
		for _, r := range f.RecentRequests {
			b, ok := bucketMap[r.Time]
			if !ok {
				b = &bucketStat{Time: r.Time}
				bucketMap[r.Time] = b
			}
			b.Success += r.Success
			b.Failed += r.Failed
		}
	}
	for _, b := range bucketMap {
		stats.RecentBuckets = append(stats.RecentBuckets, *b)
	}
	// Sort buckets by time
	for i := 0; i < len(stats.RecentBuckets); i++ {
		for j := i + 1; j < len(stats.RecentBuckets); j++ {
			if stats.RecentBuckets[i].Time > stats.RecentBuckets[j].Time {
				stats.RecentBuckets[i], stats.RecentBuckets[j] = stats.RecentBuckets[j], stats.RecentBuckets[i]
			}
		}
	}

	jsonBytes, err := json.Marshal(stats)
	if err != nil {
		return nil, fmt.Errorf("marshal stats: %w", err)
	}

	return okEnvelope(managementResponse{
		StatusCode: 200,
		Headers: http.Header{
			"content-type": []string{"application/json; charset=utf-8"},
		},
		Body: jsonBytes,
	})
}

func handleHTML() ([]byte, error) {
	return okEnvelope(managementResponse{
		StatusCode: 200,
		Headers: http.Header{
			"content-type": []string{"text/html; charset=utf-8"},
		},
		Body: []byte(htmlPage),
	})
}

// ---- Host callback ----

func callHost(method string, payload any) (json.RawMessage, error) {
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload %s: %w", method, err)
	}

	cMethod := C.CString(method)
	defer C.free(unsafe.Pointer(cMethod))

	var response C.cliproxy_buffer
	var requestPtr *C.uint8_t
	if len(rawPayload) > 0 {
		cPayload := C.CBytes(rawPayload)
		if cPayload == nil {
			return nil, fmt.Errorf("allocate payload %s", method)
		}
		defer C.free(cPayload)
		requestPtr = (*C.uint8_t)(cPayload)
	}

	callCode := C.call_host_api(cMethod, requestPtr, C.size_t(len(rawPayload)), &response)

	var rawResponse []byte
	if response.ptr != nil && response.len > 0 {
		rawResponse = C.GoBytes(response.ptr, C.int(response.len))
	}
	if response.ptr != nil {
		C.free_host_buffer(response.ptr, response.len)
	}

	if len(rawResponse) == 0 {
		return nil, fmt.Errorf("host callback %s returned no response, code=%d", method, int(callCode))
	}

	var env envelope
	if err := json.Unmarshal(rawResponse, &env); err != nil {
		return nil, fmt.Errorf("decode envelope %s: %w", method, err)
	}
	if !env.OK {
		if env.Error != nil {
			return nil, fmt.Errorf("%s: %s", env.Error.Code, env.Error.Message)
		}
		return nil, fmt.Errorf("host callback %s failed", method)
	}
	return append(json.RawMessage(nil), env.Result...), nil
}

// ---- Helpers ----

func okEnvelope(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(envelope{OK: true, Result: raw})
}

func errorEnvelope(code, message string) []byte {
	raw, _ := json.Marshal(envelope{OK: false, Error: &envelopeError{Code: code, Message: message}})
	return raw
}

func writeResponse(response *C.cliproxy_buffer, raw []byte) {
	if response == nil || len(raw) == 0 {
		return
	}
	ptr := C.CBytes(raw)
	if ptr == nil {
		return
	}
	response.ptr = ptr
	response.len = C.size_t(len(raw))
}
