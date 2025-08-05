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

	for {
		dec := NewDecoder(conn)
		value, err := dec.Decode()
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println(err)
			return
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
		var res Value = handler(args)	
		conn.Write(res.Encode())
	}
}
