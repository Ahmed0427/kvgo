# kvgo

A really simple in-memory Redis-like server written in Go.  
Implements basic Redis commands like `PING`, `SET`, `GET`, `HSET`, and `HGET` using the RESP (REdis Serialization Protocol) format — because parsing plain text is too mainstream.


### Getting Started
- You'll need Go 1.18 or later. Who knows.
- Then just run:

```bash
git clone <repo-url>
cd kvgo

go build

# My server runs on port 6379 (just like Redis), so make sure Redis isn't running:
sudo systemctl stop redis

# You're good to go
./kvgo
```

----------

### Example Usage

```bash
$ redis-cli
127.0.0.1:6379> PING
"PONG"
127.0.0.1:6379> SET name Ahmed
"OK"
127.0.0.1:6379> GET name
"Ahmed"
127.0.0.1:6379> HSET users u1 Ahmed
"OK"
127.0.0.1:6379> HGET users u1
"Ahmed"
```

### Maybe…

Maybe I’ll add more commands.  
Maybe I’ll turn it into a distributed system.  
Maybe I’ll rewrite it in brainfuck just for fun.  
Maybe I’ll forget this project exists next week.

Who knows?
