package utils

import "time"

const (
	layout = "2006-01-02 15:04:05"
)

func ParseLocal(str string) (time.Time, error) {
	return time.ParseInLocation(layout, str, time.Local)
}
