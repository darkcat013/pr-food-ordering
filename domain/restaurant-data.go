package domain

type RestaurantMenuData struct {
	RestaurantId int     `json:"restaurant_id"`
	Name         string  `json:"name"`
	MenuItems    int     `json:"menu_items"`
	Menu         []Food  `json:"menu"`
	Rating       float64 `json:"rating"`
}

type Food struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	PreparationTime  float64 `json:"preparation-time"`
	Complexity       int     `json:"complexity"`
	CookingApparatus string  `json:"cooking-apparatus"`
}

type RestaurantData struct {
	RestaurantId   int     `json:"restaurant_id"`
	Name           string  `json:"name"`
	Address        string  `json:"address"`
	MenuItems      int     `json:"menu_items"`
	Menu           []Food  `json:"menu"`
	Rating         float64 `json:"rating"`
	RegisteredTime float64
}

type RestaurantPayloadRating struct {
	OrderId              int     `json:"order_id"`
	Rating               int     `json:"rating"`
	EstimatedWaitingTime float64 `json:"estimated_waiting_time"`
	WaitingTime          float64 `json:"waiting_time"`
}

type RestaurantResponseRating struct {
	RestaurantId        int     `json:"restaurant_id"`
	RestaurantAvgRating float64 `json:"restaurant_avg_rating"`
	PreparedOrders      int     `json:"prepared_orders"`
}
