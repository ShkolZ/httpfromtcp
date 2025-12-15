package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	raddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalln(err)
		log.Fatalln(err)

	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	sin := os.Stdin
	rd := bufio.NewReader(sin)

	for {
		fmt.Print(">")

		str, err := rd.ReadString(10)
		if err != nil {
			log.Fatal(err)
		}

		n, err := conn.Write([]byte(str))
		if err != nil {
			log.Printf("Couldn't write to connection: %v", err)
		}
		if n > 0 {
			fmt.Printf("%v bytes was written to connection\n", n)
		}
	}

}
