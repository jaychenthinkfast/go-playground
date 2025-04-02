package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/playground/share-service/pkg/api"
	"github.com/playground/share-service/pkg/storage/mongo"
)

func main() {
	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 获取环境变量
	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "playground"
	}

	// 初始化存储
	storage, err := mongo.NewMongoStorage(ctx, mongoURI, mongoDB, "shares")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer storage.Close(ctx)

	// 创建 Gin 路由
	router := gin.Default()

	// 添加中间件
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// 创建 API 处理器
	handler := api.NewHandler(storage)

	// 注册路由
	router.GET("/health", handler.HealthCheck)
	router.POST("/api/share", handler.CreateShare)
	router.GET("/api/share/:id", handler.GetShare)
	router.POST("/api/share/:id/view", handler.IncrementViews)
	router.POST("/api/execute", handler.ExecuteCode)

	// 启动服务器
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置关闭超时
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
