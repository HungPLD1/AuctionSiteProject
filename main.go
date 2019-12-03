package main

import (
	"hellogorm/controller"
	model "hellogorm/model"

	//"log"

	//"time"
	//jwt_lib "github.com/dgrijalva/jwt-go"
	_ "hellogorm/docs"

	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.Use(CORSMiddleware()) //CORS config

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	//swagger init
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to hellogorm")
	})

	v2 := router.Group("item")
	//API: Search item by id
	v2.GET("/:id", controller.GetItemByID)
	//API: Search item by query (name, categories)
	v2.GET("/", controller.GetItemByQuery)
	//API: Get Bid Session
	router.GET("/session/:id", controller.BidSession)
	//APT: Get Bid Session Logs
	router.GET("/logs/:id", controller.BidLogs)
	//API: Show user profile
	router.GET("/profile", jwt.Auth(model.SecretKey), controller.UserProfile)
	//API: Show user wishlist
	router.GET("/wishlist", jwt.Auth(model.SecretKey), controller.ShowWishList)
	//API: Show user review
	router.GET("/review/:id", controller.ShowReview)
	//API: Search categories by id, return all by default
	router.GET("/categories", controller.SearchCategories)

	//API: Register new Account by JSON
	router.POST("/signup", controller.RegisterJSON)

	//API: Login by JSON
	router.POST("/login", controller.LoginJSON)
	//API: Modify user profile
	router.PUT("/profile", jwt.Auth(model.SecretKey), controller.UserProfileUpdate)
	//API: Add item to wishlist
	//router.POST("/wishlist", jwt.Auth(model.SecretKey), controller.AddWishlist)

	router.Run(":8080")
}

//CORSMiddleware ...Allow ACAO for all request and for all methods
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
