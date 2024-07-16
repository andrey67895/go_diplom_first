package transport

import (
	"github.com/andrey67895/go_diplom_first/internal/middlewareProject"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func GetRoutersGophermart() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP, middleware.Recoverer, middleware.Logger, middlewareProject.AuthValidate)
	r.Post("/api/user/register", UserRegister)
	r.Post("/api/user/login", AuthUser)
	r.Post("/api/user/orders", SaveOrders)
	r.Get("/api/user/orders", GetOrders)
	r.Get("/api/user/balance", GetBalance)
	r.Post("/api/user/balance/withdraw", WithdrawBalance)
	r.Get("/api/user/withdrawals", GetWithdrawalsHistory)
	return r
}
