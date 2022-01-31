package cpustat

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

type Stat struct {
	path    string
	bufSize int
	snap    []Snapshot
	pos     int

	ModeStatAverage Snapshot
	LoadAverage     float32
	AverageDone     chan struct{}
}

func NewCPUStat(path string, bufSize int) *Stat {
	return &Stat{
		path:        path,
		bufSize:     bufSize,
		snap:        make([]Snapshot, bufSize),
		AverageDone: make(chan struct{}, 1),
	}
}

func (s *Stat) MakeSnapshot() (*Snapshot, error) {
	f, err := os.OpenFile(s.path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	b, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}

	columns := strings.Fields(string(b))
	if len(columns) != 11 {
		return nil, errors.New("/proc/stat is broken")
	}
	columns = columns[1:]

	for i := range columns {
		val, err := strconv.ParseUint(columns[i], 10, 64)
		if err != nil {
			return nil, err
		}
		switch i {
		case 0:
			s.snap[s.pos].User = val
		case 1:
			s.snap[s.pos].Nice = val
		case 2:
			s.snap[s.pos].System = val
		case 3:
			s.snap[s.pos].Idle = val
		}
	}
	s.snap[s.pos].TimeStamp = time.Now()
	snap := s.snap[s.pos]
	s.pos++
	if s.pos == s.bufSize {
		s.ModeStatAverage = Snapshot{
			User:      s.snap[s.bufSize-1].User - s.snap[0].User,
			Nice:      s.snap[s.bufSize-1].Nice - s.snap[0].Nice,
			System:    s.snap[s.bufSize-1].System - s.snap[0].System,
			Idle:      s.snap[s.bufSize-1].Idle - s.snap[0].Idle,
			TimeStamp: s.snap[s.bufSize-1].TimeStamp,
		}

		s.LoadAverage = 100.0 * float32((s.ModeStatAverage.User+s.ModeStatAverage.Nice+s.ModeStatAverage.System)/(s.ModeStatAverage.User+s.ModeStatAverage.Nice+s.ModeStatAverage.System+s.ModeStatAverage.Idle))
		s.AverageDone <- struct{}{}
		s.pos = 0
	}

	return &snap, nil
}
