package main

import (
	"hellogorm/controller"
	model "hellogorm/model"

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
	url := ginSwagger.URL("http://site.ap.loclx.io/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to hellogorm")
	})

	v2 := router.Group("item")
	v2.GET("/:id", controller.GetItemByID)                 //API: Search item by id
	v2.GET("/", controller.GetItemByQuery)                 //API: Search item by query (name, categories)
	router.GET("/categories", controller.SearchCategories) //API: Search categories by id, return all by default

	router.GET("/session/:id", controller.BidSession) //API: Get Bid Session
	router.GET("/logs/:id", controller.BidLogs)       //APT: Get Bid Session Logs

	router.POST("/signup", controller.RegisterJSON)                                 //API: Register new Account by JSON
	router.POST("/login", controller.LoginJSON)                                     //API: Login by JSON
	router.GET("/profile", jwt.Auth(model.SecretKey), controller.UserProfile)       //API: Show user profile
	router.PUT("/profile", jwt.Auth(model.SecretKey), controller.UserProfileUpdate) //API: Modify user profile
	router.PUT("/password", jwt.Auth(model.SecretKey), controller.UpdatePassword)   //API: Change user password

	router.GET("/wishlist", jwt.Auth(model.SecretKey), controller.ShowWishList)                  //API: Show user wishlist
	router.POST("/wishlist/:id", jwt.Auth(model.SecretKey), controller.AddItemToWishList)        //API: Add new item to wishlist
	router.DELETE("/wishlist/:id", jwt.Auth(model.SecretKey), controller.RemoveItemFromWishList) //API: Delete item from wishlist

	router.GET("/review/:id", controller.ShowReview) //API: Show user review

	router.POST("/upload", jwt.Auth(model.SecretKey), controller.UploadItemImages) //API: Upload image of item to database
	//router.DELETE("/upload", jwt.Auth(model.SecretKey), controller.UploadItemImages) //API: Delete image of item from database

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
