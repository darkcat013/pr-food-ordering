package domain

import "github.com/darkcat013/pr-food-ordering/utils"

type Menu struct {
	Restaurants     int                    `json:"restaurants"`
	RestaurantsData map[int]RestaurantData `json:"restaurants_data"`
}

func RestaurantRegistrationsHandler() {
	for {
		rd := <-RestaurantRegisterChan
		RestaurantsMenu.Restaurants++

		restaurantMenu := RestaurantData{
			Address:        rd.Address,
			RestaurantId:   rd.RestaurantId,
			Name:           rd.Name,
			MenuItems:      rd.MenuItems,
			Menu:           rd.Menu,
			Rating:         rd.Rating,
			RegisteredTime: utils.GetCurrentTimeFloat(),
		}

		RestaurantsMenu.RestaurantsData[rd.RestaurantId] = restaurantMenu
	}
}
