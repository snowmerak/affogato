# affogato

Affogato is a library for observing and logging, filtering, and transforming events and logs in Go.

## observer

Observer is an interface for observing events and logs.

```go
type Observer interface {
	WatchLogs(ctx context.Context, identifier string) (line <-chan []byte, err <-chan error)
}
```

Currently, Affogato provides the following observers:
- Docker container logs observer

## parser

Parser is a library set for parsing events and logs.  
Currently, Affogato provides the following parsers:
- Redis/Valkey log parser

## checkpoint

Checkpoint is a library for managing checkpoints.  
Simply, it is a checker for the last processed log line.

```go
type Checkpoint interface {
	Check(ctx context.Context, appName string, lastTimestamp time.Time) (bool, error)
}
```

Currently, Affogato provides the following checkpoints:
- Valkey checkpoint
