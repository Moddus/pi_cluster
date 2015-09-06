// Package main provides ...
package cluster

const (
	G_UDP_PREFIX        = "4b8c608f27f76d4bad67aaead8e68bcb710de01f-PI_CLUSTER"
	G_DEFAULT_UDP_PORT  = 43782
	G_DEFAULT_TCP_PORT  = 43783
	G_DEFAULT_HTTP_PORT = 43784

	CMD_HELLO_CLIENT = "Client-Hello"
)

type Message struct {
	Cmd     string
	Payload string
	Ips     []string
}
