package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"sync"

	c "local/cluster"

	"github.com/gin-gonic/gin"
)

var (
	known_ips map[string]string
	mutex     sync.Mutex
)

func main() {
	messages := make(chan c.Message)
	//go httpServer()
	go udpServer()
	//go tcpServer(messages)

	for {
		cmd := <-messages
		log.Println(cmd)
	}

}

func tcpServer(cmd chan c.Message) {
}

func httpServer() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		if len(known_ips) > 0 {
			var buffer bytes.Buffer
			for k, _ := range known_ips {
				buffer.WriteString(k)
				buffer.WriteString("\n")
			}
			c.String(200, buffer.String())
		} else {
			c.String(200, "no clients found")
		}
	})
	r.Run(":" + strconv.Itoa(c.G_DEFAULT_HTTP_PORT))
}

func udpServer() {
	log.Println(c.G_DEFAULT_UDP_PORT)
	saddr, err := net.ResolveUDPAddr("udp", ":"+string(c.G_DEFAULT_UDP_PORT))
	sconn, err := net.ListenUDP("udp", saddr)
	if err != nil {
		log.Panicf("could not start udp server: %v", err)
	}

	buf := make([]byte, 1024)
	for {
		n, caddr, err := sconn.ReadFromUDP(buf)
		log.Println("asohdoah")
		if err != nil {
			log.Fatalf("error read from client: %d %v", n, caddr)
		} else {
			handleUDPRequest(buf, n)
		}
	}
}

func handleUDPRequest(buf []byte, bufLen int) {
	var m c.Message
	err := json.Unmarshal(buf[:bufLen], m)
	if err != nil {
		log.Fatalf("Fail to Unmarshal Udp Message -> %v", err)
	}

	if m.Cmd == c.CMD_HELLO_CLIENT && m.Payload == c.G_UDP_PREFIX && len(m.Ips) > 0 {
		mutex.Lock()
		for key, _ := range known_ips {
			if key == m.Ips[0] {
				break
			}
			known_ips[m.Ips[0]] = m.Ips[0]
		}
		mutex.Unlock()
	}
}
