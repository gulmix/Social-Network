package graph

import (
	"github.com/gulmix/Social-Network/internal/config"
	"github.com/gulmix/Social-Network/internal/repository"
	"github.com/gulmix/Social-Network/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	AuthService *service.AuthService
	UserRepo    *repository.UserRepository
	Config      *config.Config
}

func NewResolver(authService *service.AuthService, userRepo *repository.UserRepository, cfg *config.Config) *Resolver {
	return &Resolver{
		AuthService: authService,
		UserRepo:    userRepo,
		Config:      cfg,
	}
}
