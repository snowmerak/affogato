package observer

import "context"

type Observer interface {
	WatchLogs(ctx context.Context, identifier string) (line <-chan []byte, err <-chan error)
}
