# rscat

rscat is a [kafkacat](https://github.com/edenhill/kafkacat) like command-line
utility written in Go, that can produce and consume items from
[Redis Streams](https://redis.io/topics/streams-intro). Think of it like netcat
for Redis Streams. This is incredibly useful for quick debugging, building POCs
,demos and aims to embrace Unix pipeline architecture.

In **produce** mode rscat reads messages from stdin, delimited by newline and
produces them to provided Redis endpoint (--endpoint) and stream (--stream).

In **consume** mode rscat reads messages from a redis stream, and prints them
to stdout. This gives a `tail -f` like functionality.

**Installation:**

```sh
$ go install github.com/ppai-plivo/rscat
```

**Usage:**

```sh
$ rscat -h
Usage: rscat [--endpoint ENDPOINT] --stream STREAM --mode MODE [--fmt FMT] [--id ID] [--silent] [--block BLOCK]

Options:
  --endpoint ENDPOINT    redis endpoint [default: 127.0.0.1:6379]
  --stream STREAM        key name of redis stream
  --mode MODE            mode: produce or consume
  --fmt FMT              produce mode: source format (only csv supported) [default: csv]
  --id ID                produce mode: possible values - auto or linenum [default: auto]
  --silent               produce mode: suppress printing id
  --block BLOCK          consume mode: time in seconds to block for [default: 1]
  --help, -h             display this help and exit
```
