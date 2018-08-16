package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	go func() {
		listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 9981})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Server listen at:", listener.LocalAddr().String())

		data := make([]byte, 1024)
		for {
			n, rAddr, err := listener.ReadFromUDP(data)
			if err != nil {
				fmt.Println("Error during read:", err)
			}
			fmt.Println("Server received:", string(data[:n]))

			tcpAddr := strings.Split(rAddr.String(), ":")[0] + ":9982"
			if c, e := net.Dial("tcp", tcpAddr); e == nil {
				defer c.Close()
				if _, e := c.Write([]byte("world")); e != nil {
					fmt.Println("tcp write error", e)
				}
			}
		}
	}()

	time.Sleep(time.Second * 2)

	ip := net.ParseIP("255.255.255.255")
	//srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 9982}
	dstAddr := &net.UDPAddr{IP: ip, Port: 9981}
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Client listen at:", conn.LocalAddr().String())

	_, err = conn.WriteToUDP([]byte("hello"), dstAddr)
	if err != nil {
		fmt.Println(err)
	}

	tcpListen, tcpE := net.Listen("tcp", "0.0.0.0:9982")
	if tcpE != nil {
		fmt.Println("tcp listen error", tcpE)
	}
	defer tcpListen.Close()
	tcpCon, tcpCE := tcpListen.Accept()
	if tcpCE != nil {
		fmt.Println("tcp accept error", tcpCE)
	}
	buf := make([]byte, 1024)
	if tcpn, rE := tcpCon.Read(buf); rE == nil {
		fmt.Println("Client received:", string(buf[:tcpn]))
	}
}
