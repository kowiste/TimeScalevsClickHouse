package model

import "time"

type Measure struct {
	ID    string
	Asset string
	Value int
	Time  time.Time
}
type DataEntry struct {
	ID    string
	Asset string
}

type IDatabase interface {
	Connect() error
	Name() string
	IsPopulate() (int64, error)
	Save([]Measure) error
	Delete() error
	GetByAssets(assetID []string) ([]Measure, error)
	GetAssets() (data []string, err error)
	GetOldNewTime() (old, new time.Time, err error)
	GetByIntervalAndAssets(start, end time.Time, assets []string) (data []Measure, err error)
}
