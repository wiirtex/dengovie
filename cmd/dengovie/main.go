package main

import (
	"context"
	_ "dengovie/docs"
	"dengovie/internal/app/dengovie"
	"dengovie/internal/app/middlewares"
	_ "dengovie/internal/config"
	"dengovie/internal/service/debts"
	"dengovie/internal/service/telegram"
	"dengovie/internal/service/users"
	"dengovie/internal/store/postgres"
	"dengovie/internal/utils/env"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host localhost:8080
//	api.dengovie.ingress
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	env.InitEnvs(nil)
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	r.Use(middlewares.PanicCatcher)
	r.Use(middlewares.CORSMiddleware)
	r.Use(middlewares.PrometheusMiddleware())

	storage, err := postgres.New()
	if err != nil {
		log.Fatal(fmt.Errorf("postgres.New: %w", err))
	}

	debtsService := debts.New(storage)
	usersService := users.New(storage)

	bot, err := telegram.NewClient(storage)
	if err != nil {
		log.Fatalf("can not init telegram bot: %v", err)
	}
	ctx := context.Background()
	bot.Start(ctx)

	c := dengovie.NewController(storage, debtsService, usersService, bot)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/request_code", c.RequestCode)
			auth.POST("/login", c.Login)
			auth.POST("/logout", c.Logout)
		}

		groups := v1.Group("/groups")
		{
			groups.Use(middlewares.CheckAuth)
			groups.GET("", c.ListUserGroups)
			groups.GET("/:groupID/users", c.ListUsersInGroup)
		}

		user := v1.Group("/user")
		{
			user.Use(middlewares.CheckAuth)
			user.GET("", c.GetMe)
			user.POST("update_name", c.UpdateName)
			user.DELETE("delete", c.DeleteUser)
		}

		debtsHandler := v1.Group("/debts")
		{
			debtsHandler.Use(middlewares.CheckAuth)
			debtsHandler.GET("", c.ListDebts)
			debtsHandler.POST("share", c.ShareDebt)
			debtsHandler.POST("pay", c.PayDebt)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	err = r.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal("run: %w", err)
	}
}
