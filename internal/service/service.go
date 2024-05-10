package service

import (
	"github.com/go-playground/validator/v10"

	"eniqilo-store-be/internal/cfg"
	"eniqilo-store-be/internal/repo"
)

type Service struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg

	Staff   *StaffService
	Product *ProductService
	// Cat   *CatService
	// Match *MatchService
	Customer *CustomerService
	Checkout *CheckoutService
}

func NewService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *Service {
	service := Service{}
	service.repo = repo
	service.validator = validator
	service.cfg = cfg

	service.Staff = newStaffService(repo, validator, cfg)
	service.Product = newProductService(repo, validator, cfg)
	service.Customer = newCustomerService(repo, validator, cfg)
	service.Checkout = newCheckoutService(repo, validator, cfg)
	// service.Match = newMatchService(repo, validator, cfg)

	return &service
}
