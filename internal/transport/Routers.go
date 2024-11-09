package transport

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/openapi-ui/go-openapi-ui/pkg/doc"

	"github.com/andrey67895/go_diplom_first/internal/middlewareproject"
)

func SwaggerHandler(filePath string) http.HandlerFunc {

	var spec []byte
	spec, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(spec)
	}
}

func GetRoutersGophermart() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP, middleware.Recoverer, middleware.Logger, middlewareproject.AuthValidate)
	docSpec := doc.Doc{
		Title:    "Asset manager API",
		SpecFile: "/app/app.yaml",
		SpecPath: "/openapi.yaml",
		DocsPath: "/apidocs",
		Theme:    "light",
	}
	r.Handle(docSpec.DocsPath, docSpec.Handler())
	r.Handle(docSpec.SpecPath, SwaggerHandler(docSpec.SpecFile))
	r.Post("/api/user/register", UserRegister)
	r.Post("/api/user/login", AuthUser)
	r.Post("/api/user/orders", SaveOrders)
	r.Get("/api/user/orders", GetOrders)
	r.Get("/api/user/balance", GetBalance)
	r.Post("/api/user/balance/withdraw", WithdrawBalance)
	r.Get("/api/user/withdrawals", GetWithdrawalsHistory)
	r.Get("/api/ping", Ping)
	return r
}
