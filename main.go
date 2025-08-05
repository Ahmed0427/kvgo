package main

import (
	"fmt"
	"net"
	"strings"
	"os"
	"io"
)

var aof *AOF

func handle(conn net.Conn) {
	defer conn.Close()
	for {
		dec := NewDecoder(conn)
		value, err := dec.Decode()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %s disconnected\n", conn.RemoteAddr())
				return
			}
			fmt.Println("Decode error:", err)
			return
		}

		if value.kind != "array" {
			errMsg := "-ERR expected array"
			conn.Write(Value{str: errMsg, kind: "error"}.Encode())
			continue
		}

		if len(value.array) == 0 {
			errMsg := "-ERR empty array"
			conn.Write(Value{str: errMsg, kind: "error"}.Encode())
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
			errMsg := "-ERR unknown command"
			conn.Write(Value{str: errMsg, kind: "error"}.Encode())
			continue
		}

		if cmd == "SET" || cmd == "HSET" {
			aof.Write(value)
		}

		var res Value = handler(args)	
		conn.Write(res.Encode())
	}
}

func main() {
	var err error
	aof, err = newAOF("kvgodb.aof")
	if err != nil {
		fmt.Println("AOF init error:", err)
		os.Exit(1)
	}
	defer aof.Close()
	
	// Replay AOF
	aof.Read(func(value Value) {
		cmd := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[cmd]
		if !ok {
			fmt.Println("Invalid command in AOF:", cmd)
			return
		}
		handler(args)
	})

	lst, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Listen error:", err)
		os.Exit(1)
	}
	defer lst.Close()

	fmt.Println("Server is listening on port 6379")

	for {
		conn, err := lst.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handle(conn)
	}
}
