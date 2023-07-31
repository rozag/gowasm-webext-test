package main

import (
	"fmt"
	"net"
)

func main() {
	const id = 0

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	const (
		msgSuccess byte = 0x00
		msgLess    byte = 0x01
		msgMore    byte = 0x02
	)

	buf := make([]byte, 1)

	min := byte(0x00)
	max := byte(0xFF)

	hasKey := false

	for !hasKey && min <= max {
		sum := uint16(min) + uint16(max)

		var guess byte
		if sum%2 == 0 {
			guess = byte(sum / 2)
		} else {
			guess = byte((sum + 1) / 2)
		}

		n, err := conn.Write([]byte{guess})
		if err != nil {
			panic(err)
		}
		if n != 1 {
			panic("Written bytes count != 1")
		}

		n, err = conn.Read(buf)
		if err != nil {
			panic(err)
		}
		if n != 1 {
			panic("Read bytes count != 1")
		}

		answer := buf[0]
		switch answer {
		case msgSuccess:
			fmt.Printf("ID: %d, SUCCESS, answer is: %d\n", id, guess)
			hasKey = true
		case msgLess:
			fmt.Printf("ID: %d, Guess is LESS than the answer: %d\n", id, guess)
			min = guess + 1
		case msgMore:
			fmt.Printf("ID: %d, Guess is MORE than the answer: %d\n", id, guess)
			max = guess - 1
		default:
			panic(fmt.Errorf("Unknown answer: %v", answer))
		}
	}
}
