package cpu

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Snapshot struct {
	User   uint64
	Nice   uint64
	System uint64
	Idle   uint64

	TimeStamp time.Time
}

func (c *CPUStat) Collect() (*Snapshot, error) {
	file, err := os.OpenFile(c.path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	columns := strings.Fields(string(line))
	if len(columns) != 11 {
		return nil, errors.New(ErrSourceFileIsInvalid)
	}

	snapshot := &Snapshot{
		TimeStamp: time.Now(),
	}

	for i := 1; i < 5; i++ {
		val, err := strconv.ParseUint(columns[i], 10, 64)
		if err != nil {
			return nil, err
		}
		switch i {
		case 1:
			snapshot.User = val
		case 2:
			snapshot.Nice = val
		case 3:
			snapshot.System = val
		case 4:
			snapshot.Idle = val
		}
	}

	return snapshot, nil
}

const (
	ErrSourceFileIsInvalid = "invalid file"

	DefaultSourceFilePath = "/proc/stat"
)
