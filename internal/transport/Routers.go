package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func GetRouters() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP, middleware.Recoverer)
	r.Post("/api/user/register", UserRegister)
	r.Post("/api/user/login", AuthUser)
	//r.Post("/api/user/orders", handlers.SaveArraysMetDataForJSON(iStorage))
	//r.Get("/api/user/orders", handlers.GetDataForJSON(iStorage))
	//
	//r.Get("/api/user/balance", handlers.GetDataByPathParams(iStorage))
	//r.Post("/api/user/balance/withdraw", handlers.GetPing(iStorage))
	//r.Get("/api/user/withdrawals", handlers.GetAllData(iStorage))
	return r
}
