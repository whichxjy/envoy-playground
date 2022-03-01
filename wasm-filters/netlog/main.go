package main

import (
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

func (*pluginContext) NewTcpContext(contextID uint32) types.TcpContext {
	return &netlogContext{}
}

type netlogContext struct {
	types.DefaultTcpContext
}

func (ctx *netlogContext) OnDownstreamData(dataSize int, endOfStream bool) types.Action {
	data, err := proxywasm.GetDownstreamData(0, dataSize)
	if err != nil {
		proxywasm.LogErrorf("failed to get downstream data: %v", err)
		return types.ActionContinue
	}

	proxywasm.LogInfof("Get downstream data (size: %v): %v", dataSize, string(data))
	return types.ActionContinue
}

func main() {
	proxywasm.SetVMContext(&vmContext{})
}
