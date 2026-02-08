package server

import (
	"net/http"

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
	"github.com/gulmix/Social-Network/internal/repository"
	"github.com/gulmix/Social-Network/internal/service"
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

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	resolver := graph.NewResolver(authService, userRepo, cfg)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	router.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))

	router.POST("/query", middleware.Auth(cfg), gin.WrapH(srv))
	router.GET("/query", middleware.Auth(cfg), gin.WrapH(srv))

	return router
}
