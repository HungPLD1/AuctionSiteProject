package main

import (
	"hellogorm/controller"
	model "hellogorm/model"

	//"log"

	//"time"
	//jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
)

func main() {
	//Reference to singleton variance
	databaseB := controller.GetDBInstance()
	//Open Database from JSON config
	config := model.SetupConfig()
	databaseB.Db = model.ConnectDb(config.Database.User,
		config.Database.Password,
		config.Database.Database,
		config.Database.Address)
	defer databaseB.Db.Close()
	databaseB.Db.LogMode(true)

	//Create Router
	/*the 3 following line of code are equivalent to router := gin.Default() */
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	//router.Use(gin.LoggerWithFormatter(loggerformat))

	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to hellogorm")
	})

	//API: Show items by categories and. Show all by default
	router.GET("/categories/:categories", controller.Showitems)

	v2 := router.Group("item")
	//API: Search item by name
	v2.GET("name/", controller.SearchItemByName)
	//API: Search item by id
	v2.GET("id/", controller.SearchItemByID)

	//API: Register new Account by JSON
	router.POST("/signup", controller.RegisterJSON)

	//API: Login by JSON
	router.POST("/login", controller.LoginJSON)
	//API: Show user profile
	router.GET("/profile", jwt.Auth(model.SecretKey), controller.UserProfile)
	//API: Modify user profile
	router.PUT("/profile", jwt.Auth(model.SecretKey), controller.UserProfileUpdate)
	//API: Show user wishlish
	router.GET("/wishlist", controller.ShowWishList)
	//API: Bid Session
	//router.GET("/session/:id", controller.BidSession)

	router.Run(":8080")
}
