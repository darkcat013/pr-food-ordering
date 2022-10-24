package domain

type CliendOrder struct {
	ClientId int                  `json:"client_id"`
	Orders   []ClientOrderRequest `json:"orders"`
}

type ClientOrderRequest struct {
	RestaurantId int     `json:"restaurant_id"`
	Items        []int   `json:"items"`
	Priority     int     `json:"priority"`
	MaxWait      float64 `json:"max_wait"`
	CreatedTime  float64 `json:"created_time"`
}

type ClientResponse struct {
	OrderId int                   `json:"order_id"`
	Orders  []ClientOrderResponse `json:"orders"`
}

type ClientOrderResponse struct {
	RestaurantId         int     `json:"restaurant_id"`
	RestaurantAddress    string  `json:"restaurant_address"`
	OrderId              int     `json:"order_id"`
	EstimatedWaitingTime float64 `json:"estimated_waiting_time"`
	CreatedTime          float64 `json:"created_time"`
	RegisteredTime       float64 `json:"registered_time"`
}
