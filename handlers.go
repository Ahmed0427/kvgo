package main

import (
	"sync"
)

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{kind:"string", str:"PONG"}
	}
	return Value{kind:"string", str:args[0].bulk}
}

var setMap = map[string]string{}
var setMapLock = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{
			str: "ERR wrong number of arguments for 'set' command",
			kind: "error",
		}
	}

	setMapLock.Lock()
	setMap[args[0].bulk] = args[1].bulk
	setMapLock.Unlock()

	return Value{kind: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{
			str: "ERR wrong number of arguments for 'get' command",
			kind: "error",
		}
	}

	setMapLock.RLock()
	val, ok := setMap[args[0].bulk]
	setMapLock.RUnlock()

	if ok {
		return Value{kind: "bulk", bulk: val}
	}
	return Value{kind: "null"}
}

var hsetMap = map[string]map[string]string{}
var hsetMapLock = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{
			str: "ERR wrong number of arguments for 'hset' command",
			kind: "error",
		}
	}

	hkey := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	hsetMapLock.Lock()
	if _, ok := hsetMap[hkey]; !ok {
		hsetMap[hkey] = map[string]string{}
	}
	hsetMap[hkey][key] = value
	hsetMapLock.Unlock()

	return Value{kind: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{
			str: "ERR wrong number of arguments for 'hget' command",
			kind: "error",
		}
	}

	hkey := args[0].bulk
	key := args[1].bulk

	var ok bool
	var val string

	hsetMapLock.RLock()
	if _, ok = hsetMap[hkey]; ok {
		val, ok = hsetMap[hkey][key]
	}
	hsetMapLock.RUnlock()

	if ok {
		return Value{kind: "bulk", bulk: val}
	}
	return Value{kind: "null"}
}

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET": set,
	"GET": get,
	"HSET": hset,
	"HGET": hget,
}
