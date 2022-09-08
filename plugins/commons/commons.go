package commons

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/stangirard/yatas/internal/yatas"
)

// Yatas is the interface that we're exposing as a plugin.
type Yatas interface {
	Run(c *yatas.Config) []yatas.Tests
}

// Here is an implementation that talks over RPC
type YatasRPC struct{ client *rpc.Client }

func (g *YatasRPC) Run(c *yatas.Config) []yatas.Tests {
	var resp []yatas.Tests
	err := g.client.Call("Plugin.Run", c, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that YatasRPC talks to, conforming to
// the requirements of net/rpc
type YatasRPCServer struct {
	// This is the real implementation
	Impl Yatas
}

func (s *YatasRPCServer) Run(c *yatas.Config, resp *[]yatas.Tests) error {
	*resp = s.Impl.Run(c)
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a YatasRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return YatasRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type YatasPlugin struct {
	// Impl Injection
	Impl Yatas
}

func (p *YatasPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &YatasRPCServer{Impl: p.Impl}, nil
}

func (YatasPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &YatasRPC{client: c}, nil
}
