package redis

import (
	"fmt"
	"strconv"
	"time"
)

type Role byte

const (
	Child      Role = 'C'
	Master     Role = 'M'
	Replica    Role = 'S'
	AppendOnly Role = 'A'
)

type Severity byte

const (
	Info Severity = '*'
	Warn Severity = '#'
)

const TimeFormat = "02 Jan 2006 15:04:05.000"

type Log struct {
	PID      int       `json:"pid"`
	Role     Role      `json:"role"`
	Severity Severity  `json:"severity"`
	Time     time.Time `json:"time"`
	Message  string    `json:"message"`
}

type Step int

const (
	StepStart Step = iota
	StepPID
	StepRole
	StepTime
	StepMessage
)

func ParseLog(line []byte) (log Log, err error) {
	step := StepStart
	prevIdx := 0
	for i := range line {
		switch step {
		case StepStart:
			if line[i] == ':' {
				log.PID, err = strconv.Atoi(string(line[prevIdx:i]))
				if err != nil {
					err = fmt.Errorf("failed to parse PID: %w", err)
					return
				}
				step = StepPID
				prevIdx = i + 1
			}
		case StepPID:
			if line[i] == ' ' {
				switch line[prevIdx] {
				case byte(Child), byte(Master), byte(Replica), byte(AppendOnly):
					log.Role = Role(line[prevIdx])
				default:
					err = fmt.Errorf("unknown role: %c", line[prevIdx])
					return
				}
				step = StepRole
				prevIdx = i + 1
			}
		case StepRole:
			switch line[i] {
			case '#', '*':
				timeValue := line[prevIdx : i-1]
				log.Time, err = time.Parse(TimeFormat, string(timeValue))
				if err != nil {
					err = fmt.Errorf("failed to parse time: %w", err)
					return
				}

				log.Severity = Severity(line[i])
				step = StepTime
				prevIdx = i + 2
			}
		default:
		}
	}

	if prevIdx >= len(line) {
		err = fmt.Errorf("no message found")
		return
	}

	log.Message = string(line[prevIdx:])

	return log, nil
}
