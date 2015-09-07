// Package main provides ...
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	c "../config"
)

var (
	message    c.Message
	g_upd_port = flag.String("udp_port", "", "port for the udp server.")
)

func main() {
	flag.Usage = show_usage
	flag.Parse()

	setupEnv()
	broadcast()
}

func show_usage() {
	flag.PrintDefaults()
}

func setupEnv() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Panic(fmt.Sprintln("could not get local addrs"))
	}

	message.Cmd = c.CMD_HELLO_CLIENT
	message.Payload = c.G_UDP_PREFIX
	for _, addr := range addrs {
		message.Ips = append(message.Ips, addr.String())
	}
}

/*
Send a udp4 braodcast every n milliseconds.
*/
func broadcast() {
	m, err := json.Marshal(message)
	if err != nil {
		log.Panicf("Message to json failed: %v", message)
	}

	for {
		socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
			IP:   net.IPv4bcast,
			Port: c.G_DEFAULT_UDP_PORT,
		})
		log.Println(c.G_DEFAULT_UDP_PORT)
		if err != nil {
			log.Panicf("DialUDP faild -> %v", err)
		}
		socket.Write(m)
		socket.Close()
		if err != nil {
			log.Panicf("udp broadcast failed cause %v", err)
		}
		time.Sleep(30 * time.Second)
	}
}
