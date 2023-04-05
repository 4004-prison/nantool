package snowflake

import (
	"strconv"
	"sync"
	"time"
)

const (
	defaultEpoch     = int64(1577808000000)
	defaultMachineID = int64(0)
)

const (
	timestampBits = int64(41)
	machineBits   = int64(10)
	sequenceBits  = int64(12)

	maxTimestamp = int64(-1 ^ (-1 << timestampBits))
	maxMachine   = int64(-1 ^ (-1 << machineBits))
	maxSequence  = int64(-1 ^ (-1 << sequenceBits))

	timestampBitShift = machineBits + sequenceBits
	machineBitShift   = sequenceBits
)

type SnowFlake struct {
	id int64
}

func (s SnowFlake) Int64() int64 {
	return s.id
}

func (s SnowFlake) String() string {
	return strconv.FormatInt(s.id, 10)
}

type snowFlakeModel struct {
	mu        sync.Mutex
	timestamp int64
	epoch     int64
	machineID int64
	sequence  int64
}

func (s *snowFlakeModel) getSnowFlake() SnowFlake {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()
	t := now - s.epoch
	if t > maxTimestamp {
		return SnowFlake{}
	}
	sequence := int64(0)
	if now == s.timestamp {
		sequence = s.sequence
		s.sequence++
	} else {
		s.sequence = 0
	}
	s.timestamp = now
	return SnowFlake{
		id: t<<timestampBitShift | s.machineID<<machineBitShift | sequence,
	}
}

type Config struct {
	Epoch     int64 `yaml:"epoch"`
	MachineID int64 `yaml:"machineId"`
}

var mu sync.Mutex
var model *snowFlakeModel

func InitSnowFlake(config *Config) {
	if model == nil {
		mu.Lock()
		defer mu.Unlock()
		if model == nil {
			model = &snowFlakeModel{
				timestamp: time.Now().UnixMilli(),
				epoch:     config.Epoch,
				machineID: config.MachineID,
				sequence:  int64(0),
			}
		}
	}
}

func GenID() SnowFlake {
	if model == nil {
		InitSnowFlake(&Config{
			MachineID: defaultMachineID,
			Epoch:     defaultEpoch,
		})
	}
	return model.getSnowFlake()
}
