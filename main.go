package main

import (
	"math/rand"
	"time"

	"github.com/darkcat013/pr-food-ordering/domain"
	"github.com/darkcat013/pr-food-ordering/utils"
)

func main() {
	utils.InitializeLogger()
	rand.Seed(time.Now().UnixNano())

	go domain.RestaurantRegistrationsHandler()

	StartServer()
}
