package main

import (
	"os"
	"io"
	"bufio"
	"sync"
	"time"
)

type AOF struct {
	file *os.File
	reader *bufio.Reader
	lock sync.Mutex
}

func newAOF(path string) (*AOF, error) {
	file, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR, 0644)	
	if err != nil {
		return nil, err
	}

	aof := &AOF{
		file: file,
		reader: bufio.NewReader(file),
		lock: sync.Mutex{},
	}

	go func() {
		for {
			aof.lock.Lock()
			aof.file.Sync()
			aof.lock.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *AOF) Close() error {
	aof.lock.Lock()
	defer aof.lock.Unlock()
	return aof.file.Close()
}

func (aof *AOF) Write(value Value) error {
	aof.lock.Lock()
	defer aof.lock.Unlock()
	_, err := aof.file.Write(value.Encode())
	if err != nil {
		return err
	}
	return nil
}

func (aof *AOF) Read(callback func(value Value)) error {
	aof.lock.Lock()
	defer aof.lock.Unlock()
	dec := NewDecoder(aof.file)

	for {
		val, err := dec.Decode()	
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		callback(val)
	}
	return nil
}
