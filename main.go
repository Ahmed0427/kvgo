package main

import (
	"fmt"
	"net"
	"os"
	"io"
	"bytes"
	"strconv"
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
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}
		codec := NewCodec(bytes.NewReader(buf[:n]))
		value, err := codec.Decode()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(strconv.Quote(string(buf[:n]))) 
		fmt.Println(strconv.Quote(string(value.Encode()))) 

		conn.Write([]byte("+OK\r\n"))
	}
}
