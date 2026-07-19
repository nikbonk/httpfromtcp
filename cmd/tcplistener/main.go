package main

import (
	"fmt"
	"log"
	"net"

	"github.com/nikbonk/httpfromtcp/internal/request"
)

const inputFilePath = "messages.txt"
const listenerPort = ":42069"

func main() {
	listener, err := net.Listen("tcp", listenerPort)
	if err != nil {
		log.Fatalf("could not listen: %s\n", err)
	}
	defer listener.Close()

	fmt.Printf("Accepting connections on %v\n", listenerPort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("accepted connection from %v\n", conn.RemoteAddr())
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Print(err)
			return
		}
		fmt.Print("Request line:")
		fmt.Printf(
			"\n- Method: %v\n- Target: %v\n- Version: %v\n",
			req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion,
		)

		fmt.Println("Headers:")
		for key, values := range req.Headers {
			fmt.Printf("- %s: %s\n", key, values)
		}

		fmt.Println("Body:")
		fmt.Println(string(req.Body))

		fmt.Printf("closed connection from %v\n", conn.RemoteAddr())
	}
}

// Request line:
// - Method: METHOD
// - Target: TARGET
// - Version: VERSION
// Headers:
// - KEY: VALUE
// - KEY: VALUE
// Body:
// BODY_STRING
