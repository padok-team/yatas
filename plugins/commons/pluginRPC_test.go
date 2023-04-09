package commons

import (
	"bytes"
	"net/rpc"
	"testing"
)

type mockYatasImpl struct{}

func (m *mockYatasImpl) Run(c *Config) []Tests {
	return []Tests{
		{
			Account: "testAccount",
			Checks:  []Check{},
		},
	}
}

type MockRPCConn struct {
	server *YatasRPCServer
	buffer bytes.Buffer
}

func (m *MockRPCConn) Read(p []byte) (n int, err error) {
	return m.buffer.Read(p)
}

func (m *MockRPCConn) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func (m *MockRPCConn) Close() error {
	return nil
}

func TestYatasRPCServer_Run(t *testing.T) {
	impl := &mockYatasImpl{}
	server := &YatasRPCServer{Impl: impl}
	cfg := &Config{}
	var resp []Tests

	err := server.Run(cfg, &resp)
	if err != nil {
		t.Errorf("Error in YatasRPCServer Run: %v", err)
	}

	if len(resp) != 1 {
		t.Errorf("Expected 1 test result, got %d", len(resp))
	}

	if resp[0].Account != "testAccount" {
		t.Errorf("Expected account name 'testAccount', got '%s'", resp[0].Account)
	}
}

func TestYatasPlugin_Server(t *testing.T) {
	impl := &mockYatasImpl{}
	plugin := &YatasPlugin{Impl: impl}
	server, err := plugin.Server(nil)

	if err != nil {
		t.Errorf("Error in YatasPlugin Server: %v", err)
	}

	_, ok := server.(*YatasRPCServer)
	if !ok {
		t.Error("YatasPlugin Server did not return a YatasRPCServer")
	}
}

func TestYatasPlugin_Client(t *testing.T) {
	plugin := &YatasPlugin{}
	client, err := plugin.Client(nil, rpc.NewClient(&MockRPCConn{}))

	if err != nil {
		t.Errorf("Error in YatasPlugin Client: %v", err)
	}

	_, ok := client.(*YatasRPC)
	if !ok {
		t.Error("YatasPlugin Client did not return a YatasRPC")
	}
}
