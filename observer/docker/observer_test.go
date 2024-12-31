package docker_test

import (
	"context"
	"os"
	"os/signal"
	"testing"

	"github.com/snowmerak/affogato/observer/docker"
)

func TestObserver_WatchLogs(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	dc, err := docker.New(ctx)
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}

	containerID := "some-redis"

	logChan, errChan := dc.WatchLogs(ctx, containerID)
	for {
		select {
		case log, ok := <-logChan:
			if !ok {
				if len(errChan) > 0 {
					err := <-errChan
					t.Fatalf("failed to read logs: %v", err)
				}
				return
			}
			t.Logf("%s", log)
		case err := <-errChan:
			t.Fatalf("failed to read logs: %v", err)
		}
	}
}
