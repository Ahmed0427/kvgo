package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	kind  string
	str   string
	bulk  string
	array []Value
	num   int
}

func printValue(val Value, pad int) {
	for i := 0; i < pad; i++ {
		fmt.Printf("  ")
	}
	switch val.kind {
	case "bulk":
		fmt.Println("BULK:", val.bulk)
	case "array":
		fmt.Printf("ARRAY: [\n")
		for _, v := range val.array {
			printValue(v, pad+1)
		}
		for i := 0; i < pad; i++ {
			fmt.Printf("  ")
		}
		fmt.Printf("]\n")
	default:
		return
	}
}

type Codec struct {
	reader *bufio.Reader
}

func NewCodec(ioReader io.Reader) *Codec {
	return &Codec{reader: bufio.NewReader(ioReader)}
}

func (codec *Codec) decodeArray() (Value, error) {
	line, err := codec.reader.ReadSlice('\r')
	if err != nil {
		return Value{}, err
	}
	codec.reader.ReadByte()

	line = line[:len(line)-1]
	length, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, err
	}

	value := Value{kind: "array", array: make([]Value, length)}
	for i := 0; i < length; i++ {
		val, err := codec.Decode()
		if err != nil {
			return Value{}, err
		}
		value.array[i] = val
	}

	return value, nil
}

func (codec *Codec) decodeBulk() (Value, error) {
	line, err := codec.reader.ReadSlice('\r')
	if err != nil {
		return Value{}, err
	}
	codec.reader.ReadByte()

	line = line[:len(line)-1]
	length, err := strconv.Atoi(string(line))

	bulk := make([]byte, length)
	codec.reader.Read(bulk)

	codec.reader.ReadByte()
	codec.reader.ReadByte()

	return Value{kind: "bulk", bulk: string(bulk)}, nil
}

func (codec *Codec) Decode() (Value, error) {
	kind, err := codec.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch kind {
	case BULK:
		return codec.decodeBulk()
	case ARRAY:
		return codec.decodeArray()
	default:
		return Value{}, fmt.Errorf("Unkonwn Type %s\n", string(kind))
	}
}

func (v Value) encodeArray() []byte {
	var ret []byte
	ret = append(ret, ARRAY)
	arrayLen := strconv.Itoa(len(v.array))
	ret = append(ret, []byte(arrayLen)...)
	ret = append(ret, '\r', '\n')

	for i := range v.array {
		ret = append(ret, v.array[i].Encode()...)
	}
	return ret
}

func (v Value) encodeBulk() []byte {
	var ret []byte
	ret = append(ret, BULK)
	bulkLen := strconv.Itoa(len(v.bulk))
	ret = append(ret, []byte(bulkLen)...)
	ret = append(ret, '\r', '\n')
	ret = append(ret, []byte(v.bulk)...)
	ret = append(ret, '\r', '\n')
	return ret
}

func (v Value) encodeString() []byte {
	var ret []byte
	ret = append(ret, STRING)
	ret = append(ret, []byte(v.str)...)
	ret = append(ret, '\r', '\n')
	return ret
}

func (v Value) encodeError() []byte {
	var ret []byte
	ret = append(ret, ERROR)
	ret = append(ret, []byte(v.str)...)
	ret = append(ret, '\r', '\n')
	return ret
}

func (v Value) encodeNull() []byte {
	return []byte("$-1\r\n")
}

func (v Value) Encode() []byte {
	switch v.kind {
	case "bulk":
		return v.encodeBulk()
	case "array":
		return v.encodeArray()
	case "string":
		return v.encodeString()
	case "null":
		return v.encodeNull()
	case "error":
		return v.encodeError()
	default:
		return []byte{}
	}
}
