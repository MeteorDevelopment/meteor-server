package core

import (
	"fmt"
	"time"
)

func GetDate() string {
	dt := time.Now()
	return fmt.Sprintf("%02d-%02d-%d", dt.Day(), dt.Month(), dt.Year())
}
