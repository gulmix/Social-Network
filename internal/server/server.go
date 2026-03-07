package server

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gulmix/Social-Network/internal/config"
	"github.com/gulmix/Social-Network/internal/database"
	"github.com/gulmix/Social-Network/internal/graph"
	"github.com/gulmix/Social-Network/internal/middleware"
	"github.com/gulmix/Social-Network/internal/pubsub"
	"github.com/gulmix/Social-Network/internal/repository"
	"github.com/gulmix/Social-Network/internal/service"
	"github.com/gulmix/Social-Network/internal/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	router.Use(middleware.CORS())

	db, err := database.InitPostgres(cfg)
	if err != nil {
		panic(err)
	}

	if err := database.InitRedis(cfg); err != nil {
		panic(err)
	}
	ps := pubsub.New(database.RedisClient)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	followRepo := repository.NewFollowRepository(db)
	convRepo := repository.NewConversationRepository(db)
	msgRepo := repository.NewMessageRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	eventRepo := repository.NewEventRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	storyRepo := repository.NewStoryRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo, followRepo)
	postService := service.NewPostService(postRepo, userRepo, likeRepo, commentRepo, followRepo)
	commentService := service.NewCommentService(commentRepo, postRepo, userRepo)
	likeService := service.NewLikeService(likeRepo, postRepo, userRepo)
	followService := service.NewFollowService(followRepo, userRepo)
	convService := service.NewConversationService(convRepo, userRepo)
	msgService := service.NewMessageService(msgRepo, convRepo)
	groupService := service.NewGroupService(groupRepo)
	eventService := service.NewEventService(eventRepo)
	notifService := service.NewNotificationService(notifRepo)
	storyService := service.NewStoryService(storyRepo, followRepo)

	resolver := graph.NewResolver(
		authService, userService, postService,
		commentService, likeService, followService,
		convService, msgService, groupService,
		eventService, notifService, storyService,
		userRepo, postRepo, commentRepo, likeRepo, followRepo,
		cfg, ps,
	)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{
		MaxMemory:     32 << 20, // 32MB
		MaxUploadSize: 10 << 20, // 10MB
	})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, *transport.InitPayload, error) {
			if _, ok := middleware.GetUserIDFromContext(ctx); ok {
				return ctx, nil, nil
			}

			authVal := initPayload.Authorization()
			if authVal != "" {
				tokenString := utils.ExtractTokenFromHeader(authVal)
				if tokenString == "" && !strings.HasPrefix(authVal, "Bearer ") {
					tokenString = authVal
				}
				if tokenString != "" {
					claims, err := utils.ValidateToken(tokenString, cfg)
					if err != nil {
						log.Printf("ws: invalid token in connection_init: %v", err)
					} else {
						ctx = context.WithValue(ctx, middleware.UserIDKey, claims.UserID)
					}
				}
			}

			return ctx, nil, nil
		},
	})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// GraphQL endpoints
	router.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))
	router.POST("/query", middleware.Auth(cfg), gin.WrapH(srv))
	router.GET("/query", middleware.Auth(cfg), gin.WrapH(srv))

	// Media upload endpoint
	uploadDir := cfg.Upload.Dir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	router.Static("/uploads", uploadDir)
	router.POST("/upload", middleware.Auth(cfg), func(c *gin.Context) {
		_, exists := middleware.GetUserIDFromContext(c.Request.Context())
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
			return
		}

		filename, _, err := utils.SaveUploadedFile(file, uploadDir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		url := "/uploads/" + filename
		c.JSON(http.StatusOK, gin.H{
			"url":      url,
			"filename": filename,
		})
	})

	return router
}
