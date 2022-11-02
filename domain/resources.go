package domain

var RestaurantsMenu Menu = Menu{
	Restaurants:     0,
	RestaurantsData: make(map[int]RestaurantData),
}

var OrderId int64
