//go:build proxytest

package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func TestUpdateMetricWhenHeaderIsMissing(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithVMContext(&vmContext{}).
		WithPluginConfiguration([]byte("required_headers:x-site"))
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())
	require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

	// Initialize http context.
	headers := [][2]string{{"my-custom-header", "foo"}}
	contextID := host.InitializeHttpContext()

	action := host.CallOnRequestHeaders(contextID, headers, false)
	require.Equal(t, types.ActionContinue, action)

	// Check metrics.
	value, err := host.GetCounterMetric("envoy_missing_required_headers_header=x-site")
	require.NoError(t, err)
	require.Equal(t, uint64(1), value)
}

func TestMissingMetricWhenAllRequiredHeadersArePresent(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithVMContext(&vmContext{}).
		WithPluginConfiguration([]byte("required_headers:x-site"))
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())
	require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

	// Initialize http context.
	headers := [][2]string{{"x-site", "foo"}}
	contextID := host.InitializeHttpContext()

	action := host.CallOnRequestHeaders(contextID, headers, false)
	require.Equal(t, types.ActionContinue, action)

	// Check metrics.
	value, err := host.GetCounterMetric("envoy_missing_required_headers_header=x-site")
	require.NoError(t, err)
	require.Equal(t, uint64(0), value)
}

func TestShouldNotPublishMetricsWhenPluginConfigIsNotProvided(t *testing.T) {
	opt := proxytest.NewEmulatorOption().
		WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnVMStart.
	require.Equal(t, types.OnVMStartStatusOK, host.StartVM())
	require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())

	// Initialize http context.
	headers := [][2]string{{"x-site", "foo"}}
	contextID := host.InitializeHttpContext()

	action := host.CallOnRequestHeaders(contextID, headers, false)
	require.Equal(t, types.ActionContinue, action)

	// Check metrics.
	_, err := host.GetCounterMetric("envoy_missing_required_headers_header=x-site")
	require.Error(t, err)
}
