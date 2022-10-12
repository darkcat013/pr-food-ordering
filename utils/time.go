package utils

import (
	"time"

	"github.com/darkcat013/pr-food-ordering/config"
)

func GetCurrentTimeFloat() float64 {
	if config.TIME_UNIT >= time.Millisecond && config.TIME_UNIT < time.Second {
		return float64(time.Now().UnixMilli())
	} else {
		return float64(time.Now().Unix())
	}
}
