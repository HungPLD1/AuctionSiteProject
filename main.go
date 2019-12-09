package main

import (
	"AuctionSiteProject/controller"
	model "AuctionSiteProject/model"

	_ "AuctionSiteProject/docs"

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
	url := ginSwagger.URL("http://siteb.ap.loclx.io/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to hellogorm")
	})

	router.GET("/categories", controller.SearchCategories)
	router.POST("/categories", jwt.Auth(model.SecretKey), controller.NewCategories)

	router.GET("/session/:sessionid", controller.BidSessionByID)
	router.GET("/session", controller.BidSessionByQuery)
	router.POST("/session", jwt.Auth(model.SecretKey), controller.CreateBidSession)
	router.PUT("/session", jwt.Auth(model.SecretKey), controller.UpdateBidSession)
	router.GET("/awaitpayment", jwt.Auth(model.SecretKey), controller.UnpaidSession)
	router.PUT("lock/:sessionid", controller.LockSession)

	router.GET("/logs/:sessionid", controller.BidLogs)
	router.POST("/logs", jwt.Auth(model.SecretKey), controller.NewBid)

	router.GET("/history/bid", jwt.Auth(model.SecretKey), controller.BidSessionHistory)
	router.GET("/history/sell", jwt.Auth(model.SecretKey), controller.SellSessionHistory)

	router.POST("/signup", controller.RegisterJSON)
	router.POST("/login", controller.LoginJSON)
	router.GET("/profile", jwt.Auth(model.SecretKey), controller.UserProfile)
	router.PUT("/profile", jwt.Auth(model.SecretKey), controller.UserProfileUpdate)
	router.PUT("/password", jwt.Auth(model.SecretKey), controller.UpdatePassword)

	router.GET("/wishlist", jwt.Auth(model.SecretKey), controller.ShowWishList)
	router.POST("/wishlist/:itemid", jwt.Auth(model.SecretKey), controller.AddItemToWishList)
	router.DELETE("/wishlist/:itemid", jwt.Auth(model.SecretKey), controller.RemoveItemFromWishList)

	router.GET("/review/:userid", controller.ShowReview)
	router.POST("/review", jwt.Auth(model.SecretKey), controller.CreateReview)

	//router.POST("/upload", jwt.Auth(model.SecretKey), controller.UploadItemImages)
	//router.GET("/upload", controller.SendImages)

	//ADMIN ONLY
	router.DELETE("/session/:sessionid", jwt.Auth(model.SecretKey), controller.DeleteBidSession)
	router.DELETE("/logs/last", jwt.Auth(model.SecretKey), controller.DeleteLastBidLogs)
	router.DELETE("/review/:sessionid", jwt.Auth(model.SecretKey), controller.DeleteReview)

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
