package main

import (
	"github.com/darkcat013/pr-food-ordering/domain"
	"github.com/darkcat013/pr-food-ordering/utils"
)

func main() {
	utils.InitializeLogger()
	go domain.RestaurantRegistrationsHandler()

	StartServer()
}
