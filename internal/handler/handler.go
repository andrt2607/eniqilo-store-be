package handler

import (
	"time"

	"eniqilo-store-be/internal/cfg"
	"eniqilo-store-be/internal/service"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Handler struct {
	router  *chi.Mux
	service *service.Service
	cfg     *cfg.Cfg
}

func NewHandler(router *chi.Mux, service *service.Service, cfg *cfg.Cfg) *Handler {
	handler := &Handler{router, service, cfg}
	handler.registRoute()

	return handler
}

func (h *Handler) registRoute() {

	r := h.router
	var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(h.cfg.JWTSecret), nil, jwt.WithAcceptableSkew(30*time.Second))

	staffHandler := newStaffHandler(h.service.Staff)
	productHandler := newProductHandler(h.service.Product)
	customerHandler := newCustomerHandler(h.service.Customer)
	checkoutHandler := newCheckoutHandler(h.service.Checkout)

	r.Use(middleware.RedirectSlashes)

	r.Post("/v1/staff/register", staffHandler.Register)
	r.Post("/v1/staff/login", staffHandler.Login)

	r.Get("/v1/product/customer", productHandler.GetProductSKU)

	// protected route
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Get("/v1/product", productHandler.GetProduct)
		r.Post("/v1/product", productHandler.Create)
		r.Put("/v1/product/{id}", productHandler.UpdateByID)
		r.Delete("/v1/product/{id}", productHandler.DeleteProduct)

		r.Post("/v1/customer/register", customerHandler.Register)
		r.Get("/v1/customer", customerHandler.GetCustomer)
		r.Post("/v1/product/checkout", checkoutHandler.PostCheckout)
		r.Get("/v1/product/checkout", checkoutHandler.GetCheckout)
	})
}
