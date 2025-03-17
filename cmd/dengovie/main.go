package main

import (
	"dengovie/internal/app/dengovie"
	"dengovie/internal/app/middlewares"
	"dengovie/internal/store/postgres"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"

	_ "dengovie/docs"
	_ "dengovie/internal/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {

	r := gin.Default()

	connString, found := os.LookupEnv("POSTGRES_CONN_STRING")
	if !found {
		log.Fatal("POSTGRES_CONN_STRING environment variable not found")
	}
	fmt.Println(connString)

	storage, err := postgres.New(connString)
	if err != nil {
		log.Fatal(fmt.Errorf("postgres.New: %w", err))
	}

	c := dengovie.NewController(storage)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/request_code", c.RequestCode)
			auth.POST("/login", c.Login)
		}

		groups := v1.Group("/groups")
		{
			// TODO: check authorization
			groups.Use(middlewares.CheckAuth)
			groups.GET("", c.ListUserGroups)
			groups.GET("/:groupID/users", c.ListUsersInGroup)
		}

		user := v1.Group("/user")
		{
			user.GET("", func(context *gin.Context) { /**/ })
		}

		debts := v1.Group("/debts")
		{
			debts.GET("", func(context *gin.Context) {})
			debts.POST("share", func(context *gin.Context) {})
			debts.POST("pay", func(context *gin.Context) {})
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err = r.Run("127.0.0.1:8080")
	if err != nil {
		log.Fatal("run: %w", err)
	}
}
