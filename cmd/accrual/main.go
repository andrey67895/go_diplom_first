package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.Recoverer, middleware.Logger)
	r.Get("/api/orders/{numbers}", getMockAccrual)
	helpers.TLog.Fatal(http.ListenAndServe(":8080", r))
}

func randBool() bool {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(2) == 1
}

func getMockAccrual(w http.ResponseWriter, req *http.Request) {

	numbers := chi.URLParam(req, "numbers")
	_, err := strconv.Atoi(numbers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var tModel model.OrdersAccrualModel
	if randBool() {
		status := "PROCESSED"
		accrual := rand.Float64() * 100000
		tModel = model.OrdersAccrualModel{
			OrderID: &numbers,
			Status:  &status,
			Accrual: &accrual,
		}
	} else if randBool() {
		status := "INVALID"
		tModel = model.OrdersAccrualModel{
			OrderID: &numbers,
			Status:  &status,
		}
	} else {
		status := "PROCESSING"
		tModel = model.OrdersAccrualModel{
			OrderID: &numbers,
			Status:  &status,
		}
	}
	marshal, err := json.Marshal(tModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
	w.WriteHeader(http.StatusOK)
}
