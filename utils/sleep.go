package utils

import (
	"time"

	"github.com/darkcat013/pr-food-ordering/config"
)

func SleepFor(t float64) {
	time.Sleep(time.Duration(t) * config.TIME_UNIT)
}
