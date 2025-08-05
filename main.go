package main

import (
	"fmt"
	"net"
	"strings"
	"os"
	"io"
)

func main() {
	lst, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer lst.Close()

	conn, err := lst.Accept()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	aof, err := newAOF("kvgodb.aof")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer aof.Close()

	aof.Read(func(value Value) {
		cmd := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[cmd]
		if !ok {
			fmt.Println("Invalid command:", cmd)
			return
		}
		handler(args)	
	})

	for {
		dec := NewDecoder(conn)
		value, err := dec.Decode()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		if value.kind != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		cmd := strings.ToUpper(value.array[0].bulk)
		if cmd == "COMMAND" {
			conn.Write(Value{kind: "string", str: ""}.Encode())
			continue
		}
		args := value.array[1:]

		handler, ok := Handlers[cmd]
		if !ok {
			fmt.Println("Invalid command:", cmd)
			conn.Write(Value{kind: "string", str: ""}.Encode())
			continue
		}

		if cmd == "SET" || cmd == "HSET" {
			aof.Write(value)
		}

		var res Value = handler(args)	
		conn.Write(res.Encode())
	}
}
