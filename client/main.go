package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nhooyr.io/websocket"
)

const addr = "ws://localhost:7000"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, addr, &websocket.DialOptions{
		Subprotocols: []string{"tlsbum-client-to-verifier"},
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	targetUrl := "tcp://127.0.0.1:8080"
	data, err := json.Marshal(targetUrl)
	if err != nil {
		panic(err)
	}

	err = conn.Write(ctx, websocket.MessageText, data)
	if err != nil {
		panic(err)
	}

	const id = 0

	const (
		msgSuccess byte = 0x00
		msgLess    byte = 0x01
		msgMore    byte = 0x02
	)

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

		err := conn.Write(ctx, websocket.MessageBinary, []byte{guess})
		if err != nil {
			panic(err)
		}

		msgType, buf, err := conn.Read(ctx)
		if err != nil {
			panic(err)
		}
		if msgType != websocket.MessageBinary {
			panic(fmt.Errorf("expected websocket.MessageBinary, got: %v", msgType))
		}
		if len(buf) != 1 {
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
