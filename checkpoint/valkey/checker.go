package valkey

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
	"unique"

	"github.com/valkey-io/valkey-go"
)

func makeKey(appName string) string {
	return "affogato:checkpoint:" + appName
}

var uniqueKeyMap = sync.Map{}

func getUniqueKey(appName string) string {
	v, ok := uniqueKeyMap.Load(appName)
	if !ok {
		k := makeKey(appName)
		uk := unique.Make[string](k)
		uniqueKeyMap.Store(appName, uk)
	}

	uk, ok := v.(unique.Handle[string])
	if !ok {
		k := makeKey(appName)
		uk = unique.Make[string](k)
		uniqueKeyMap.Store(appName, uk)
	}

	return uk.Value()
}

type Checkpoint struct {
	client valkey.Client
}

func NewCheckpoint(client valkey.Client) *Checkpoint {
	return &Checkpoint{client: client}
}

const checkLocalTTL = 5 * time.Minute

const checkAndSwapScript = `
local key = KEYS[1]
local lastTimestamp = ARGV[1]

local currentTimestamp = redis.call("GET", key)
if currentTimestamp == nil then
	redis.call("SET", key, lastTimestamp)
	return true
end

if currentTimestamp < lastTimestamp then
	redis.call("SET", key, lastTimestamp)
	return true
end

return false
`

func (c *Checkpoint) Check(ctx context.Context, appName string, lastTimestamp time.Time) (bool, error) {
	uk := getUniqueKey(appName)
	lastTimestampEpoch := lastTimestamp.UnixNano()

	latest, err := c.client.DoCache(ctx, c.client.B().Get().Key(uk).Cache(), checkLocalTTL).AsInt64()
	if err != nil {
		return false, fmt.Errorf("failed to get checkpoint: %w", err)
	}

	if lastTimestampEpoch < latest {
		return false, nil
	}

	resp := c.client.Do(ctx, c.client.B().Eval().Script(checkAndSwapScript).Numkeys(1).Key(uk).Arg(strconv.FormatInt(lastTimestampEpoch, 10)).Build())
	if err := resp.Error(); err != nil {
		return false, fmt.Errorf("failed to check checkpoint: %w", err)
	}

	ok, err := resp.AsBool()
	if err != nil {
		return false, fmt.Errorf("failed to check checkpoint: %w", err)
	}

	return ok, nil
}
