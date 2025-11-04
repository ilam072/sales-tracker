package main

import (
	"context"
	analyticsrepo "github.com/ilam072/sales-tracker/internal/analytics/repo/postgres"
	analyticsrest "github.com/ilam072/sales-tracker/internal/analytics/rest"
	analyticsservice "github.com/ilam072/sales-tracker/internal/analytics/service"
	categoryrepo "github.com/ilam072/sales-tracker/internal/category/repo/postgres"
	categoryrest "github.com/ilam072/sales-tracker/internal/category/rest"
	categoryservice "github.com/ilam072/sales-tracker/internal/category/service"
	"github.com/ilam072/sales-tracker/internal/config"
	itemrepo "github.com/ilam072/sales-tracker/internal/item/repo/postgres"
	itemrest "github.com/ilam072/sales-tracker/internal/item/rest"
	itemservice "github.com/ilam072/sales-tracker/internal/item/service"
	"github.com/ilam072/sales-tracker/internal/middlewares"
	"github.com/ilam072/sales-tracker/internal/validator"
	"github.com/ilam072/sales-tracker/pkg/db"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize logger
	zlog.Init()

	// Context
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Initialize config
	cfg := config.MustLoad()

	// Connect to DB
	DB, err := db.OpenDB(cfg.DB)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to connect to DB")
	}

	// Initialize validator
	v := validator.New()

	// Initialize category, item and analytics repositories
	categoryRepo := categoryrepo.New(DB)
	itemRepo := itemrepo.New(DB)
	analyticsRepo := analyticsrepo.New(DB)

	// Initialize category, item and analytics services
	category := categoryservice.New(categoryRepo)
	item := itemservice.New(itemRepo)
	analytics := analyticsservice.New(analyticsRepo)

	// Initialize category, item and analytics handlers
	categoryHandler := categoryrest.NewCategoryHandler(category, v)
	itemHandler := itemrest.NewItemHandler(item, v)
	analyticsHandler := analyticsrest.NewAnalyticsHandler(analytics, v)

	// Initialize Gin engine and set routes
	engine := ginext.New("")
	engine.Use(ginext.Logger())
	engine.Use(ginext.Recovery())
	engine.Use(middlewares.CORS())

	api := engine.Group("/api")

	// categories
	api.POST("/categories", categoryHandler.CreateCategory)
	api.GET("/categories/:id", categoryHandler.GetCategoryByID)
	api.GET("/categories", categoryHandler.GetAllCategories)
	api.PUT("/categories/:id", categoryHandler.UpdateCategory)
	api.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	// items
	api.POST("/items", itemHandler.CreateItem)
	api.GET("/items/:id", itemHandler.GetItemByID)
	api.GET("/items", itemHandler.GetAllItems) // query параметры ?from=...&to=...&category_id=...&type=...
	api.PUT("/items/:id", itemHandler.UpdateItem)
	api.DELETE("/items/:id", itemHandler.DeleteItem)

	// analytics
	api.GET("/analytics/sum", analyticsHandler.Sum)                        // query параметры ?from=...&to=...&category_id=...&type=...
	api.GET("/analytics/avg", analyticsHandler.Avg)                        // query параметры ?from=...&to=...&category_id=...&type=...
	api.GET("/analytics/count", analyticsHandler.Count)                    // query параметры ?from=...&to=...&category_id=...&type=...
	api.GET("/analytics/median", analyticsHandler.Median)                  // query параметры ?from=...&to=...&category_id=...&type=...
	api.GET("/analytics/percentile", analyticsHandler.PercentileNinetieth) // query параметры ?from=...&to=...&category_id=...&type=...

	// Initialize and start http server
	server := &http.Server{
		Addr:    cfg.Server.HTTPPort,
		Handler: engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to listen start http server")
		}
	}()

	<-ctx.Done()

	// Graceful shutdown
	withTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(withTimeout); err != nil {
		zlog.Logger.Error().Err(err).Msg("server shutdown failed")
	}

	if err := DB.Master.Close(); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to close master database")
	}
}
