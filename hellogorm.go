package main

import (
	model "hellogorm/model"
	"log"
	"sync"

	"github.com/jinzhu/gorm"

	//"time"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/gin-gonic/gin"
)

/******SINGLETON Database Connection******/
var once sync.Once

//DatabaseB ...It hold the pointer to database.
type DatabaseB struct {
	Db *gorm.DB
}

//variance global
var instance *DatabaseB

//GetDBInstance ...Use this function go fetch database instance.
func GetDBInstance() *DatabaseB {
	once.Do(func() { //do not allow repeating
		//thread safe
		instance = &DatabaseB{}
	})

	return instance
}

/**	Items table
*	id	name	bidding_status	item_condition	id_categories	description
**/
type Items struct {
	ItemID            int
	ItemName          string `gorm:"type:varchar(255)"`
	ItemBiddingstatus string `gorm:"type:varchar(20)"`
	ItemCondition     string `gorm:"type:varchar(10)"`
	CategoriesID      int
	ItemDescription   string `gorm:"type:varchar(255)"`
}

func main() {
	//Reference to singleton variance
	databaseB := GetDBInstance()
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

	v1 := router.Group("show")
	{
		v1.GET("/:categories", showitems)
	}

	router.Run(":8080")
}

func showitems(c *gin.Context) {
	db := GetDBInstance().Db
	categoriesName := c.Param("categories")
	var itemsList []Items

	/*Temporaly stuff*/
	//s := new(search)
	//s.Where("name = ?", "jinzhu").Order("name").Attrs("name", "jinzhu").Select("name, age")

	if categoriesName == "all" {
		errGetItems := db.Table("items_auction_info").Select("*").Scan(&itemsList).Error
		if errGetItems != nil {
			log.Println(errGetItems)
			return
		}
		c.JSON(200, itemsList)
	}
}
