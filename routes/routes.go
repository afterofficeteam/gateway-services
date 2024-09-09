package routes

import (
	"gateway-service/config"
	"gateway-service/handlers/cart"
	"gateway-service/handlers/order"
	"gateway-service/handlers/payment"
	"gateway-service/handlers/users"
	"gateway-service/util/middleware"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Routes struct {
	Router *http.ServeMux
}

func URLRewriter(baseURLPath string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, baseURLPath)

		next.ServeHTTP(w, r)
	}
}

func (r *Routes) SetupBaseURL() {
	baseURL := viper.GetString("BASE_URL_PATH")
	if baseURL != "" && baseURL != "/" {
		r.Router.HandleFunc(baseURL+"/", URLRewriter(baseURL, r.Router))
	}
}

func (r *Routes) cartRoutes() {
	r.Router.HandleFunc("GET /cart/{user_id}", middleware.ApplyMiddleware(cart.GetByUserID, middleware.EnabledCors, middleware.LoggerMiddleware()))
}

func (r *Routes) orderRoutes() {
	r.Router.HandleFunc("POST /order/{user_id}", middleware.ApplyMiddleware(order.CreateOrder, middleware.EnabledCors, middleware.LoggerMiddleware()))
	r.Router.HandleFunc("POST /order/callback", middleware.ApplyMiddleware(order.UpdateOrder, middleware.EnabledCors, middleware.LoggerMiddleware()))
}

func (r *Routes) userRoutes() {
	r.Router.HandleFunc("POST /login", middleware.ApplyMiddleware(users.Login, middleware.EnabledCors, middleware.LoggerMiddleware()))
}

func (r *Routes) paymentRoutes() {
	r.Router.HandleFunc("POST /payment", middleware.ApplyMiddleware(payment.CreatePayment, middleware.EnabledCors, middleware.LoggerMiddleware()))
}

func (r *Routes) SetupRouter() {
	r.Router = http.NewServeMux()
	r.SetupBaseURL()
	r.cartRoutes()
	r.orderRoutes()
	r.userRoutes()
	r.paymentRoutes()
}

func (r *Routes) Run(port string) {
	r.SetupRouter()

	log.Printf("[Running-Success] clients on localhost on port :%s", port)
	srv := &http.Server{
		Handler:      r.Router,
		Addr:         "localhost:" + port,
		WriteTimeout: config.WriteTimeout() * time.Second,
		ReadTimeout:  config.ReadTimeout() * time.Second,
	}

	log.Panic(srv.ListenAndServe())
}
