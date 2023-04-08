package commons

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

func (g *YatasRPC) Run(c *Config) []Tests {
	var resp []Tests
	err := g.client.Call("Plugin.Run", c, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

func (s *YatasRPCServer) Run(c *Config, resp *[]Tests) error {
	*resp = s.Impl.Run(c)
	return nil
}

func (p *YatasPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &YatasRPCServer{Impl: p.Impl}, nil
}

func (YatasPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &YatasRPC{client: c}, nil
}
