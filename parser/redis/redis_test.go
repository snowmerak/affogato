package redis_test

import (
	"testing"
	"time"

	"github.com/snowmerak/affogato/parser/redis"
)

func TestParseLog(t *testing.T) {
	testInput := []string{
		"19538:C 31 Dec 2024 19:04:28.664 # WARNING Memory overcommit must be enabled! Without it, a background save or replication may fail under low memory condition. Being disabled, it can also cause failures without low memory condition, see https://github.com/jemalloc/jemalloc/issues/1328. To fix this issue add 'vm.overcommit_memory = 1' to /etc/sysctl.conf and then reboot or run the command 'sysctl vm.overcommit_memory=1' for this to take effect.",
		"19538:C 31 Dec 2024 19:04:28.664 * oO0OoO0OoO0Oo Valkey is starting oO0OoO0OoO0Oo",
		"19538:C 31 Dec 2024 19:04:28.664 * Valkey version=8.0.1, bits=64, commit=00000000, modified=0, pid=19538, just started",
		"19538:C 31 Dec 2024 19:04:28.664 # Warning: no config file specified, using the default config. In order to specify a config file use valkey-server /path/to/valkey.conf",
		"19538:M 31 Dec 2024 19:04:28.664 * Increased maximum number of open files to 10032 (it was originally set to 1024).",
		"19538:M 31 Dec 2024 19:04:28.664 * monotonic clock: POSIX clock_gettime",
		"19538:M 31 Dec 2024 19:04:28.665 * Server initialized",
		"19538:M 31 Dec 2024 19:04:28.665 * Ready to accept connections tcp",
	}

	expectedOutput := []redis.Log{
		{
			PID:      19538,
			Role:     redis.Child,
			Time:     time.Date(2024, time.December, 31, 19, 4, 28, 664000000, time.UTC),
			Severity: redis.Warn,
			Message:  "WARNING Memory overcommit must be enabled! Without it, a background save or replication may fail under low memory condition. Being disabled, it can also cause failures without low memory condition, see https://github.com/jemalloc/jemalloc/issues/1328. To fix this issue add 'vm.overcommit_memory = 1' to /etc/sysctl.conf and then reboot or run the command 'sysctl vm.overcommit_memory=1' for this to take effect.",
		},
		{
			PID:      19538,
			Role:     redis.Child,
			Time:     time.Date(2024, time.December, 31, 19, 4, 28, 664000000, time.UTC),
			Severity: redis.Info,
			Message:  "oO0OoO0OoO0Oo Valkey is starting oO0OoO0OoO0Oo",
		},
		{
			PID:      19538,
			Role:     redis.Child,
			Time:     time.Date(2024, time.December, 31, 19, 4, 28, 664000000, time.UTC),
			Severity: redis.Info,
			Message:  "Valkey version=8.0.1, bits=64, commit=00000000, modified=0, pid=19538, just started",
		},
		{
			PID:      19538,
			Role:     redis.Child,
			Time:     time.Date(2024, time.December, 31, 19, 4, 28, 664000000, time.UTC),
			Severity: redis.Warn,
			Message:  "Warning: no config file specified, using the default config. In order to specify a config file use valkey-server /path/to/valkey.conf",
		},
		{
			PID:      19538,
			Role:     redis.Master,
			Time:     time.Date(2024, time.December, 31, 19, 4, 28, 664000000, time.UTC),
			Severity: redis.Info,
			Message:  "Increased maximum number of open files to 10032 (it was originally set to 1024).",
		},
		{
			PID:      19538,
			Role:     redis.Master,
			Time:     time.Date(2024, time.December, 31, 19, 4, 28, 664000000, time.UTC),
			Severity: redis.Info,
			Message:  "monotonic clock: POSIX clock_gettime",
		},
		{
			PID:      19538,
			Role:     redis.Master,
			Time:     time.Date(2024, time.December, 31, 19, 4, 28, 665000000, time.UTC),
			Severity: redis.Info,
			Message:  "Server initialized",
		},
		{
			PID:      19538,
			Role:     redis.Master,
			Time:     time.Date(2024, time.December, 31, 19, 4, 28, 665000000, time.UTC),
			Severity: redis.Info,
			Message:  "Ready to accept connections tcp",
		},
	}

	for i, input := range testInput {
		output, err := redis.ParseLog(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if output.PID != expectedOutput[i].PID {
			t.Errorf("PID mismatch: expected %d, got %d", expectedOutput[i].PID, output.PID)
		}

		if output.Role != expectedOutput[i].Role {
			t.Errorf("Role mismatch: expected %c, got %c", expectedOutput[i].Role, output.Role)
		}

		if output.Severity != expectedOutput[i].Severity {
			t.Errorf("Severity mismatch: expected %c, got %c", expectedOutput[i].Severity, output.Severity)
		}

		if output.Time != expectedOutput[i].Time {
			t.Errorf("Time mismatch: expected %v, got %v", expectedOutput[i].Time, output.Time)
		}

		if output.Message != expectedOutput[i].Message {
			t.Errorf("Message mismatch: expected %s, got %s", expectedOutput[i].Message, output.Message)
		}
	}
}
