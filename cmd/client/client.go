package main

import (
	"fmt"
	"github.com/bluntenpassant/word-of-wisdom-tcp/transport"
	"net"
	"time"
)

const maxQuoteSize = 1024

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:8081")
	if err != nil {
		panic(err)
	}

	err = conn.SetDeadline(time.Now().Add(15 * time.Minute))
	if err != nil {
		panic(err)
	}

	powConn := transport.NewPowClient(conn)

	err = powConn.EstablishSecureConnection()
	if err != nil {
		panic(err)
	}

	quoteBuf := make([]byte, maxQuoteSize)

	n, err := powConn.Read(quoteBuf)
	if err != nil {
		panic(err)
	}

	if n == 0 {
		panic("zero bytes received from server during reading quote")
	}

	fmt.Println(string(quoteBuf))
}
