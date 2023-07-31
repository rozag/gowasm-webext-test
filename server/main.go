package main

import (
	"fmt"
	"math/rand"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	fmt.Println("Listening on", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("Connection accepted")

		go func(c net.Conn) {
			defer c.Close()

			key := byte(rand.Uint32())
			fmt.Println("Key generated:", key)

			const (
				msgSuccess byte = 0x00
				msgLess    byte = 0x01
				msgMore    byte = 0x02
			)

			buf := make([]byte, 1)
			for {
				n, err := c.Read(buf)
				if err != nil {
					panic(err)
				}
				if n != 1 {
					panic("Read bytes count != 1")
				}

				guess := buf[0]
				if guess == key {
					fmt.Printf("Key:\t%d,\tGuess:\t%d,\tSUCCESS\n", key, guess)
					n, err := c.Write([]byte{msgSuccess})
					if err != nil {
						panic(err)
					}
					if n != 1 {
						panic("Success, written bytes count != 1")
					}
					break
				} else if guess < key {
					fmt.Printf("Key:\t%d,\tGuess:\t%d,\tLESS\n", key, guess)
					n, err := c.Write([]byte{msgLess})
					if err != nil {
						panic(err)
					}
					if n != 1 {
						panic("Less, written bytes count != 1")
					}
				} else {
					fmt.Printf("Key:\t%d,\tGuess:\t%d,\tMORE\n", key, guess)
					n, err := c.Write([]byte{msgMore})
					if err != nil {
						panic(err)
					}
					if n != 1 {
						panic("More, written bytes count != 1")
					}
				}
			}
		}(conn)
	}
}
