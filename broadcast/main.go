package main

import (
	"net"
	"fmt"
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

			_, err = listener.WriteToUDP([]byte("world"), rAddr)
			if err != nil {
				fmt.Printf(err.Error())
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

	n, err := conn.WriteToUDP([]byte("hello"), dstAddr)
	if err != nil {
		fmt.Println(err)
	}

	data := make([]byte, 1024)
	n, _, err = conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Client received:", string(data[:n]))
}
