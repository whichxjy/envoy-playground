package main

import (
	"net/http"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
}

func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &authContext{}
}

type authContext struct {
	types.DefaultHttpContext
}

func (ctx *authContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	reqHeaders, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogErrorf("Fail to get request headers: %v", err)
		return types.ActionPause
	}

	if _, err := proxywasm.DispatchHttpCall("auth_cluster", reqHeaders, nil, nil, 50000, authCallback); err != nil {
		proxywasm.LogErrorf("Fail to call auth: %v", err)
		return types.ActionPause
	}

	return types.ActionPause
}

func authCallback(numHeaders, bodySize, numTrailers int) {
	respBody, err := proxywasm.GetHttpCallResponseBody(0, bodySize)
	if err != nil {
		proxywasm.LogErrorf("Fail to get response body: %v", err)
		return
	}

	proxywasm.LogInfof("Auth response body: %v", string(respBody))

	if string(respBody) != "OK" {
		if err := proxywasm.SendHttpResponse(http.StatusForbidden, nil, nil, -1); err != nil {
			proxywasm.LogErrorf("Fail to send local response: %v", err)
		}
		return
	}

	proxywasm.ResumeHttpRequest()
}

func main() {
	proxywasm.SetVMContext(&vmContext{})
}
