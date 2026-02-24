package graph

import (
	"github.com/gulmix/Social-Network/internal/config"
	"github.com/gulmix/Social-Network/internal/pubsub"
	"github.com/gulmix/Social-Network/internal/repository"
	"github.com/gulmix/Social-Network/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	AuthService    *service.AuthService
	UserService    *service.UserService
	PostService    *service.PostService
	CommentService *service.CommentService
	LikeService    *service.LikeService
	FollowService  *service.FollowService
	UserRepo       *repository.UserRepository
	PostRepo       *repository.PostRepository
	CommentRepo    *repository.CommentRepository
	LikeRepo       *repository.LikeRepository
	FollowRepo     *repository.FollowRepository
	Config         *config.Config
	PubSub         *pubsub.PubSub
}

func NewResolver(
	authService *service.AuthService,
	userService *service.UserService,
	postService *service.PostService,
	commentService *service.CommentService,
	likeService *service.LikeService,
	followService *service.FollowService,
	userRepo *repository.UserRepository,
	postRepo *repository.PostRepository,
	commentRepo *repository.CommentRepository,
	likeRepo *repository.LikeRepository,
	followRepo *repository.FollowRepository,
	cfg *config.Config,
	ps *pubsub.PubSub,
) *Resolver {
	return &Resolver{
		AuthService:    authService,
		UserService:    userService,
		PostService:    postService,
		CommentService: commentService,
		LikeService:    likeService,
		FollowService:  followService,
		UserRepo:       userRepo,
		PostRepo:       postRepo,
		CommentRepo:    commentRepo,
		LikeRepo:       likeRepo,
		FollowRepo:     followRepo,
		Config:         cfg,
		PubSub:         ps,
	}
}
