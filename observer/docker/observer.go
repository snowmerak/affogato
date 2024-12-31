package docker

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Observer struct {
	client *client.Client
}

func New(ctx context.Context) (*Observer, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	context.AfterFunc(ctx, func() {
		cli.Close()
	})

	return &Observer{client: cli}, nil
}

func (o *Observer) WatchLogs(ctx context.Context, containerID string) (<-chan []byte, <-chan error) {
	logChan := make(chan []byte, 1024)
	errChan := make(chan error, 1)

	reader, err := o.client.ContainerLogs(context.Background(), containerID, container.LogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "40",
	})
	if err != nil {
		errChan <- fmt.Errorf("failed to get logs: %w", err)
		close(logChan)
		close(errChan)
		return logChan, errChan
	}

	context.AfterFunc(ctx, func() {
		reader.Close()
	})

	go func() {
		defer close(logChan)
		defer close(errChan)

		br := bufio.NewScanner(reader)
		for br.Scan() {
			line := br.Bytes()
			// header := line[:8]
			line = line[8:]
			logChan <- line
		}

		if err := br.Err(); err != nil {
			errChan <- fmt.Errorf("failed to read logs: %e", err)
		}
	}()

	return logChan, errChan
}
