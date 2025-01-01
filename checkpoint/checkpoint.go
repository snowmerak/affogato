package checkpoint

import (
	"context"
	"time"
)

type Checkpoint interface {
	Check(ctx context.Context, appName string, lastTimestamp time.Time) (bool, error)
}
