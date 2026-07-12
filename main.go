package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const inputFilePath = "messages.txt"
const listenerPort = ":42069"

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)

		currentLineContents := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				lines <- currentLineContents + parts[i]
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}

	}()
	return lines
}

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
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Printf("closed connection from %v\n", conn.RemoteAddr())
	}
}
