package gorpc

import (
	"imooc.com/resk/infra"
	"imooc.com/resk/infra/base"
)

type GoRpcApiStarter struct {
	infra.BaseStarter
}

func (g *GoRpcApiStarter) Init(ctx infra.StarterContext) {
	base.RpcRegister(new(EnvelopeRpc))
}
