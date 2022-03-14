package main

import (
	"fmt"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &metricPluginContext{}
}

type metricPluginContext struct {
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (ctx *metricPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &metricHttpContext{}
}

type metricHttpContext struct {
	types.DefaultHttpContext
}

var counters = map[string]proxywasm.MetricCounter{}
var requiredHeaders []string

// Override types.DefaultPluginContext.
func (ctx metricPluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
	}

	configMap := parseConfig(data)
	requiredHeaders = strings.Split(configMap["required_headers"], ",")
	prepareCounters(requiredHeaders)

	return types.OnPluginStartStatusOK
}

func parseConfig(data []byte) map[string]string {
	configMap := make(map[string]string)
	configString := string(data)
	lines := strings.Split(configString, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			keyValuePair := strings.Split(line, ":")
			configMap[keyValuePair[0]] = keyValuePair[1]
		}
	}
	return configMap
}

func prepareCounters(requiredHeaders []string) {
	for _, headerName := range requiredHeaders {
		fqn := fmt.Sprintf("envoy_missing_required_headers_header=%s", headerName)
		counter := proxywasm.DefineCounterMetric(fqn)
		counters[headerName] = counter
	}
}

// Override types.DefaultHttpContext.
func (ctx *metricHttpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	for _, headerName := range requiredHeaders {
		_, err := proxywasm.GetHttpRequestHeader(headerName)
		if err != nil {
			counters[headerName].Increment(1)
		}
	}

	return types.ActionContinue
}
