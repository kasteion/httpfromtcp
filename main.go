package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main()  {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not listen on %s - %s\n", port, err)
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", port)
	fmt.Println("===================")

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not accept connection: %s", err)
		}

		fmt.Println("Connection accepted")
		linesChan := getLinesChannel(con)
		for line := range linesChan {
			fmt.Println(line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8)
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
				return
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	
	return lines
}
