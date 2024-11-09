package main

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.Recoverer, middleware.Logger)
	r.Get("/api/orders/{numbers}", getMockAccrual)
	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           r,
	}

	helpers.TLog.Fatal(server.ListenAndServe())
}

func randBool() bool {
	b, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		return false
	}
	return b.Int64() == 1
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
		accrual := 500.0
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
	_, err = w.Write(marshal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
