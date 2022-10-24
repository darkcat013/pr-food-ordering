package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/darkcat013/pr-food-ordering/config"
	"github.com/darkcat013/pr-food-ordering/domain"
	"github.com/darkcat013/pr-food-ordering/utils"
	"go.uber.org/zap"
)

func StartServer() {
	unhandledRoutes := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		utils.Log.Warn("Path not found", zap.Int("statusCode", http.StatusNotFound))
		http.Error(w, "404 path not found.", http.StatusNotFound)
	}

	register := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		if r.Method != "POST" {
			utils.Log.Warn("Method not allowed", zap.Int("statusCode", http.StatusMethodNotAllowed))
			http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
			return
		}

		var rd domain.RestaurantData
		err := json.NewDecoder(r.Body).Decode(&rd)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			utils.Log.Fatal("Failed to decode restaurand data", zap.String("error", err.Error()))
			return
		}
		utils.Log.Info("Restaurant data decoded", zap.Any("data", rd))

		domain.RestaurantRegisterChan <- rd

		w.WriteHeader(http.StatusOK)
	}

	menu := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		if r.Method != "GET" {
			utils.Log.Warn("Method not allowed", zap.Int("statusCode", http.StatusMethodNotAllowed))
			http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
			return
		}

		utils.Log.Info("Convert menu to JSON ", zap.Any("menu", domain.RestaurantsMenu))

		response, err := json.Marshal(domain.RestaurantsMenu)
		if err != nil {
			utils.Log.Fatal("Failed to convert menu to JSON", zap.String("error", err.Error()), zap.Any("menu", domain.RestaurantsMenu))
		}

		utils.Log.Info("Send back menu")
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}

	order := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		if r.Method != "POST" {
			utils.Log.Warn("Method not allowed", zap.Int("statusCode", http.StatusMethodNotAllowed))
			http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
			return
		}

		var o domain.CliendOrder
		err := json.NewDecoder(r.Body).Decode(&o)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			utils.Log.Fatal("Failed to decode client order", zap.String("error", err.Error()))
			return
		}
		utils.Log.Info("Client order decoded", zap.Any("data", o))

		atomic.AddInt64(&domain.OrderId, 1)

		var orders = make([]domain.ClientOrderResponse, 0)
		for i := 0; i < len(o.Orders); i++ {

			toRestaurant := domain.RestaurantOrderPayload{
				Items:       o.Orders[i].Items,
				Priority:    o.Orders[i].Priority,
				MaxWait:     o.Orders[i].MaxWait,
				CreatedTime: o.Orders[i].CreatedTime,
			}

			body, _ := json.Marshal(toRestaurant)
			resp, _ := http.Post(domain.RestaurantsMenu.RestaurantsData[o.Orders[i].RestaurantId-1].Address+"/v2/order", "application/json", bytes.NewBuffer(body))

			var ord domain.RestaurantOrderResponse

			json.NewDecoder(resp.Body).Decode(&ord)

			order := domain.ClientOrderResponse{
				RestaurantId:         ord.RestaurantId,
				RestaurantAddress:    domain.RestaurantsMenu.RestaurantsData[o.Orders[i].RestaurantId-1].Address,
				OrderId:              ord.OrderId,
				EstimatedWaitingTime: ord.EstimatedWaitingTime,
				CreatedTime:          utils.GetCurrentTimeFloat(),
				RegisteredTime:       ord.RegisteredTime,
			}
			orders = append(orders, order)
		}
		responseObj := domain.ClientResponse{
			OrderId: int(atomic.LoadInt64(&domain.OrderId)),
			Orders:  orders,
		}

		response, err := json.Marshal(responseObj)
		if err != nil {
			utils.Log.Fatal("Failed to convert responseOrder to JSON", zap.String("error", err.Error()), zap.Any("responseOrder", responseObj))
		}

		utils.Log.Info("Send back response order to client", zap.Int("clientId", o.ClientId))

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}

	rating := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		if r.Method != "POST" {
			utils.Log.Warn("Method not allowed", zap.Int("statusCode", http.StatusMethodNotAllowed))
			http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
			return
		}

		utils.Log.Info("food ordering rating POST")
		w.WriteHeader(http.StatusNoContent)
	}

	http.HandleFunc("/", unhandledRoutes)
	http.HandleFunc("/register", register)
	http.HandleFunc("/menu", menu)
	http.HandleFunc("/order", order)
	http.HandleFunc("/rating", rating)

	utils.Log.Info("Started web server at port :" + config.PORT)

	if err := http.ListenAndServe(":"+config.PORT, nil); err != nil {
		utils.Log.Fatal("Could not start web server", zap.String("error", err.Error()))
	}
}
