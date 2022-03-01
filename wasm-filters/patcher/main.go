package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/sjson"
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
	return &patcherContext{}
}

type patcherContext struct {
	types.DefaultHttpContext
	totalResponseBodySize int
}

func (ctx *patcherContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	if err := proxywasm.RemoveHttpResponseHeader("Content-Length"); err != nil {
		proxywasm.LogErrorf("Fail to remove content-length response header: %v", err)
		return types.ActionPause
	}

	return types.ActionContinue
}

func (ctx *patcherContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	ctx.totalResponseBodySize += bodySize
	if !endOfStream {
		return types.ActionPause
	}

	originalRespBody, err := proxywasm.GetHttpResponseBody(0, ctx.totalResponseBodySize)
	if err != nil {
		proxywasm.LogErrorf("Fail to get response body: %v", err)
		return types.ActionPause
	}
	proxywasm.LogInfof("Original response body: %v", string(originalRespBody))

	newRespBody, err := sjson.SetRawBytes(originalRespBody, "hi", []byte("\"hi\""))
	if err != nil {
		proxywasm.LogErrorf("Fail to set response body: %v", err)
		return types.ActionPause
	}
	proxywasm.LogInfof("New response body: %v", string(newRespBody))

	if err := proxywasm.ReplaceHttpResponseBody(newRespBody); err != nil {
		proxywasm.LogErrorf("Fail to replace response: %v", err)
		return types.ActionPause
	}

	return types.ActionContinue
}

func main() {
	proxywasm.SetVMContext(&vmContext{})
}
