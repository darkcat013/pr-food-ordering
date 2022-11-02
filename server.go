package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
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

		oId := atomic.AddInt64(&domain.OrderId, 1)

		var orders = make([]domain.ClientOrderResponse, 0)
		var mu = &sync.Mutex{}
		var wg = &sync.WaitGroup{}

		for i := 0; i < len(o.Orders); i++ {

			rAddr := domain.RestaurantsMenu.RestaurantsData[o.Orders[i].RestaurantId].Address

			toRestaurant := domain.RestaurantOrderPayload{
				Items:       o.Orders[i].Items,
				Priority:    o.Orders[i].Priority,
				MaxWait:     o.Orders[i].MaxWait,
				CreatedTime: o.Orders[i].CreatedTime,
			}
			wg.Add(1)

			go func(restaurantAddress string, payload domain.RestaurantOrderPayload) {

				utils.Log.Info("Send order to restaurant", zap.Any("data", payload))

				body, err := json.Marshal(payload)

				if err != nil {
					utils.Log.Fatal("Failed to convert restaurant order to json", zap.String("error", err.Error()), zap.Any("order", payload))
				}

				resp, err := http.Post(restaurantAddress+"/v2/order", "application/json", bytes.NewBuffer(body))
				if err != nil {
					utils.Log.Error("Failed to send order to restaurant"+restaurantAddress, zap.String("error", err.Error()), zap.Any("order", payload))
					http.Error(w, "Failed to send order to restaurant "+restaurantAddress, http.StatusBadRequest)
				}

				var ord domain.RestaurantOrderResponse

				json.NewDecoder(resp.Body).Decode(&ord)

				utils.Log.Info("Decoded Restautant response", zap.Any("data", ord))

				order := domain.ClientOrderResponse{
					RestaurantId:         ord.RestaurantId,
					RestaurantAddress:    restaurantAddress,
					OrderId:              ord.OrderId,
					EstimatedWaitingTime: ord.EstimatedWaitingTime,
					CreatedTime:          utils.GetCurrentTimeFloat(),
					RegisteredTime:       ord.RegisteredTime,
				}

				utils.Log.Info("Created order for client", zap.Any("order", order))

				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
				wg.Done()
			}(rAddr, toRestaurant)
		}
		wg.Wait()

		responseObj := domain.ClientResponse{
			OrderId: int(oId),
			Orders:  orders,
		}
		utils.Log.Info("Created response object for client", zap.Any("data", responseObj))

		response, err := json.Marshal(responseObj)
		if err != nil {
			utils.Log.Fatal("Failed to convert responseOrder to JSON", zap.String("error", err.Error()), zap.Any("responseOrder", responseObj))
		}

		utils.Log.Info("Send back response order to client", zap.Int("clientId", o.ClientId), zap.Int("orderId", responseObj.OrderId))

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}

	ratingMutex := &sync.Mutex{}
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

		var ratings domain.ClientRating
		err := json.NewDecoder(r.Body).Decode(&ratings)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			utils.Log.Fatal("Failed to decode client ratings", zap.String("error", err.Error()))
			return
		}
		utils.Log.Info("Client ratings decoded", zap.Any("data", ratings))

		newRatings := make(map[int]float64)
		var mu = &sync.Mutex{}
		var wg = &sync.WaitGroup{}

		for i := 0; i < len(ratings.Orders); i++ {

			r := domain.RestaurantsMenu.RestaurantsData[ratings.Orders[i].RestaurantId]

			toRestaurant := domain.RestaurantPayloadRating{
				OrderId:              ratings.Orders[i].OrderId,
				Rating:               ratings.Orders[i].Rating,
				EstimatedWaitingTime: ratings.Orders[i].EstimatedWaitingTime,
				WaitingTime:          ratings.Orders[i].WaitingTime,
			}
			wg.Add(1)

			go func(restaurant domain.RestaurantData, payload domain.RestaurantPayloadRating) {

				utils.Log.Info("Send ratings to restaurant", zap.Any("data", payload))

				body, err := json.Marshal(payload)

				if err != nil {
					utils.Log.Fatal("Failed to convert restaurant rating to json", zap.String("error", err.Error()), zap.Any("rating", payload))
				}

				resp, err := http.Post(restaurant.Address+"/v2/rating", "application/json", bytes.NewBuffer(body))
				if err != nil {
					utils.Log.Error("Failed to send rating to restaurant"+restaurant.Address, zap.String("error", err.Error()), zap.Any("rating", payload))
					http.Error(w, "Failed to send rating to restaurant "+restaurant.Address, http.StatusBadRequest)
				}

				var rating domain.RestaurantResponseRating

				json.NewDecoder(resp.Body).Decode(&rating)

				utils.Log.Info("Decoded Restautant rating response", zap.Any("data", rating))

				mu.Lock()

				newRatings[restaurant.RestaurantId] = rating.RestaurantAvgRating

				mu.Unlock()
				wg.Done()
			}(r, toRestaurant)
		}
		wg.Wait()

		ratingMutex.Lock()
		sumRatings := 0.0
		for i := 1; i <= len(domain.RestaurantsMenu.RestaurantsData); i++ {
			if entry, ok := newRatings[i]; ok {
				field := domain.RestaurantsMenu.RestaurantsData[i]
				field.Rating = entry
				domain.RestaurantsMenu.RestaurantsData[i] = field
			}
			sumRatings += domain.RestaurantsMenu.RestaurantsData[i].Rating
		}

		avg := sumRatings / float64(domain.RestaurantsMenu.Restaurants)

		utils.LogRep.Info("AVG RATING", zap.Float64("rating", avg))
		ratingMutex.Unlock()

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
