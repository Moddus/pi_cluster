package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"

	c "local/cluster"

	"github.com/gin-gonic/gin"
)

var (
	REGEX_IPV4 = regexp.MustCompile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")
	known_ips  = make(map[string]string)
	mutex      sync.Mutex
)

func main() {
	messages := make(chan c.Message)
	go httpServer()
	go udpServer()
	go tcpServer(messages)

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
	saddr, err := net.ResolveUDPAddr("udp", ":43782") //+string(c.G_DEFAULT_UDP_PORT))
	sconn, err := net.ListenUDP("udp", saddr)
	if err != nil {
		log.Panicf("could not start udp server: %v", err)
	}

	buf := make([]byte, 1024)
	for {
		n, caddr, err := sconn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalf("error read from client: %d %v", n, caddr)
		} else {
			handleUDPRequest(buf, n)
		}
	}
}

func handleUDPRequest(buf []byte, bufLen int) {
	var m c.Message
	err := json.Unmarshal(buf[:bufLen], &m)
	if err != nil {
		log.Fatalf("Fail to Unmarshal Udp Message -> %v", err)
	}

	if m.Cmd == c.CMD_HELLO_CLIENT && strings.HasPrefix(m.Payload, c.G_UDP_PREFIX) && len(m.Ips) > 0 {
		var ip string
		for _, n := range m.Ips {
			if REGEX_IPV4.MatchString(n) {
				ip = n
			}
		}
		mutex.Lock()
		if len(known_ips) > 0 {
			for key, _ := range known_ips {
				if key == ip {
					break
				}
				known_ips[ip] = ip
			}
		} else {
			known_ips[ip] = ip
		}
		mutex.Unlock()
	}
}
