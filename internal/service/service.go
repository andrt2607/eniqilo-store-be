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

	User    *UserService
	Product *ProductService
	// Cat   *CatService
	// Match *MatchService
}

func NewService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *Service {
	service := Service{}
	service.repo = repo
	service.validator = validator
	service.cfg = cfg

	service.User = newUserService(repo, validator, cfg)
	service.Product = newProductService(repo, validator, cfg)
	// service.Cat = newCatService(repo, validator, cfg)
	// service.Match = newMatchService(repo, validator, cfg)

	return &service
}
