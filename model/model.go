package model

import (
	"fmt"
	"log"
	"os"

	"mime/multipart"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	jsoniter "github.com/json-iterator/go"
)

/**************************DATABASE STRUCT FOR GORM*********************************/

//SessionSearch ...Used by gorm and json
type SessionSearch struct {
	SessionID          int             `gorm:"type:bigint(20)" json:"sessionid"`
	ItemID             int             `gorm:"type:bigint(20)" json:"itemid"`
	ItemName           string          `gorm:"type:varchar(255)" json:"itemname" binding:"required"`
	CategoriesID       int             `gorm:"type:int(11)" json:"categoriesid" binding:"required"`
	CategoriesName     string          `gorm:"type:varchar(255)" json:"categoriesname"`
	ItemDescription    string          `gorm:"type:text" json:"itemdescription" binding:"required"`
	ItemCondition      string          `gorm:"type:varchar(30)" json:"itemcondition" binding:"required"`
	SessionStartDate   time.Time       `gorm:"type:datetime" json:"startdate"`
	SessionEndDate     time.Time       `gorm:"type:datetime" json:"enddate"`
	MinimumIncreaseBid int             `gorm:"type:int(11)" json:"minimumincreasebid" binding:"required"`
	UserviewCount      int             `gorm:"type:int(11)" json:"viewcount"`
	Images             []string        `json:"imagelink"`
	SellerID           string          `gorm:"type:varchar(255)" json:"sellerid"`
	SellerName         string          `gorm:"type:varchar(100)" json:"sellername"`
	CurrentBid         int             `gorm:"type:bigint(20)" json:"currentbid" binding:"required"`
	BidLogs            []BidSessionLog `json:"biddingLog"`
	SessionStatus      string          `gorm:"type:varchar30" json:"sessionstatus"`
}

//Items ...Used by gorm and json
type Items struct {
	ItemID          int       `gorm:"type:bigint(20)" json:"itemid"`
	CategoriesID    int       `gorm:"type:int(11)" json:"categoriesid"`
	ItemName        string    `gorm:"type:varchar(255)" json:"itemname"`
	ItemDescription string    `gorm:"type:text" json:"itemdescription"`
	ItemCondition   string    `gorm:"type:varchar(30)" json:"itemcondition"`
	ItemCreateat    time.Time `gorm:"type:datetime" json:"createAt"`
}

//Categories ...Used by gorm and json
type Categories struct {
	CategoriesID   int    `gorm:"type:int(11)" json:"categoriesid"`
	CategoriesName string `gorm:"type:varchar(255)" json:"categoriesname"`
}

//ItemImage ...Used by gorm and json
type ItemImage struct {
	ItemID int    `gorm:"type:bigint(20)" json:"itemid"`
	Images string `gorm:"text" json:"images"`
}

//UserCommon ...Used by gorm and json
type UserCommon struct {
	UserID          string    `gorm:"type:varchar(255)" json:"userid"`
	UserPassword    string    `gorm:"type:varchar(255)" json:"password"`
	UserName        string    `gorm:"type:varchar(100)" json:"name"`
	UserPhone       string    `gorm:"type:varchar(15)" json:"phone"`
	UserEmail       string    `gorm:"type:varchar(255)" json:"email"`
	UserGender      byte      `gorm:"type:char(1)" json:"gender"`
	UserAddress     string    `gorm:"type:varchar(255)" json:"address"`
	UserAvatar      string    `gorm:"type:TEXT" json:"avatarimage"`
	UserAccessLevel int       `gorm:"type:int" json:"accesslevel"`
	UserCreateat    time.Time `gorm:"type:datetime" json:"createAt"`
}

//BidSession ...Used by gorm and json
type BidSession struct {
	SessionID          int       `gorm:"type:bigint(20)" json:"sessionid"`
	ItemID             int       `gorm:"type:bigint(20)" json:"itemid"`
	SellerID           string    `gorm:"type:varchar(255)" json:"sellerid"`
	SessionStartDate   time.Time `gorm:"type:datetime" json:"startdate"`
	SessionEndDate     time.Time `gorm:"type:datetime" json:"enddate"`
	UserviewCount      int       `gorm:"type:int(11)" json:"viewcount"`
	WinnerID           string    `gorm:"type:varchar(255)" json:"winnerid"`
	MinimumIncreaseBid int       `gorm:"type:int(11)" json:"minimumbid"`
	CurrentBid         int       `gorm:"type:bigint(20)" json:"currentbid"`
	SessionStatus      string    `gorm:"type:varchar(30)" json:"status"`
}

//BidSessionLog ...Used by gorm and json
type BidSessionLog struct {
	UserID    string    `gorm:"type:varchar(255)" json:"userid"`
	SessionID int       `gorm:"type:bigint(20)" json:"sessionid"`
	BidAmount int       `gorm:"type:int" json:"amount"`
	BidDate   time.Time `gorm:"type:datetime" json:"createAt"`
}

//UserReview ...Used by gorm and json
type UserReview struct {
	UserWriter    string `gorm:"type:varchar(255)" json:"writerid"`
	UserTarget    string `gorm:"type:varchar(255)" json:"targetid"`
	SessionID     int    `gorm:"type:bigint(20)" json:"sessionid"`
	ReviewContent string `gorm:"type:text" json:"content"`
	ReviewScore   int    `gorm:"type:int(1)" json:"score"`
}

//UserWishlist ...used by gorm and json
type UserWishlist struct {
	UserID  string    `gorm:"type:varchar(255)" json:"userid"`
	ItemID  int       `gorm:"type:bigint(20)" json:"itemid"`
	AddDate time.Time `gorm:"type: datetime" json:"createAt"`
}

/**************************COMMUNICATION STRUCT*********************************/
//SignupLoginResponse ...Respond form
type SignupLoginResponse struct {
	ResponseTime string     `json:"responseTime"`
	Code         int        `json:"code"`
	Message      string     `json:"message"`
	Data         UserCommon `json:"data"`
	SessionToken string     `json:"sessiontoken"`
}

//AuthorizationHeader ...Used to get session token in header
type AuthorizationHeader struct {
	Token string `header:"Authorization"`
}

//UploadItemImageForm ...For uploading item photo
type UploadItemImageForm struct {
	ItemID int                     `form:"itemid" binding:"required"`
	Images []*multipart.FileHeader `form:"imageurl" binding:"required"`
}

//ModifyPassword ...For password update
type ModifyPassword struct {
	OldPassword string `json:"oldpassword" binding:"required"`
	NewPassword string `json:"newpassword" binding:"required"`
}

//NewSession ...Used by controller.CreateBidSession
type NewSession struct {
	ItemName           string   `json:"itemname" binding:"required"`
	ItemDescription    string   `json:"itemdescription" binding:"required"`
	ItemCondition      string   `json:"itemcondition" binding:"required"`
	SessionStartDate   string   `json:"startdate"`
	SessionEndDate     string   `json:"enddate"`
	StartPrice         int      `json:"startprice" binding:"required"`
	MinimumIncreaseBid int      `json:"minimumincreasebid" binding:"required"`
	Images             []string `json:"imagelink"`
	CategoriesID       int      `json:"categoriesid" binding:"required"`
}

//UpdateSession ...Used by controller.UpdateBidSession
type UpdateSession struct {
	SessionID       int    `json:"sessionid" binding:"required"`
	ItemName        string `json:"itemname"`
	ItemDescription string `json:"itemdescription"`
	ItemCondition   string `json:"itemcondition"`
	//Images             []string        `json:"imagelink"`
}

//NewBidLog ...Used when user create a new bid
type NewBidLog struct {
	SessionID int ` json:"sessionid" binding:"required"`
	BidAmount int ` json:"amount" binding:"required"`
}

//BidHistory ...Used by controller.BidSessionHistory
type BidHistory struct {
	SessionID        int       ` json:"sessionid"`
	ItemID           int       `json:"itemid"`
	ItemName         string    ` json:"itemname"`
	Images           []string  `json:"images"`
	BidAmount        int       `json:"bidamount"`
	BidDate          time.Time `json:"biddate"`
	SessionStartDate time.Time `json:"sessionstartdate"`
	SessionEndDate   time.Time `json:"sessionenddate"`
}

//SellHistory ...Used by controller.SellSessionHistory
type SellHistory struct {
	SessionID        int       ` json:"sessionid"`
	ItemID           int       `json:"itemid"`
	ItemName         string    ` json:"itemname"`
	Images           []string  `json:"images"`
	SessionStartDate time.Time `json:"sessionstartdate"`
	SessionEndDate   time.Time `json:"sessionenddate"`
}

//Deletelastlog ...Used by controller.DeleteBidLogs
type Deletelastlog struct {
	SessionID int    ` json:"sessionid" binding:"required"`
	UserID    string `json:"userid" binding:"required"`
}

//NewReview ...Used by controller.CreateReview
type NewReview struct {
	UserTarget	string	`json:"usertargetid" binding:"required"`
	SessionID 	int `json:"sessionid" binding:"required"`
	ReviewContent string `json:"reviewcontent" binding:"required"`
	ReviewScore int `json:"reviewscore" binding:"required"`
}

//Register ...Used by controller.RegisterJSON
type RegisterForm struct {
	UserID       string `json:"userid" binding:"required"`
	UserPassword string `json:"password" binding:"required"`
	UserEmail    string `json:"email" binding:"required"`
}

//LoginForm ...Used by controller.LoginJSON
type LoginForm struct {
	UserID       string `json:"userid" binding:"required"`
	UserPassword string `json:"password" binding:"required"`
}

//Config ...Database login info
type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
		Address  string `json:"address"`
	} `json:"database"`
}

/**************************INTERNAL SECTION*********************************/
var (
	SecretKey = "thonking"
)

func DecodeDataFromJsonFile(f *os.File, data interface{}) error {
	jsonParser := jsoniter.NewDecoder(f)
	err := jsonParser.Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

func SetupConfig() Config {
	var conf Config

	// Đọc file config.dev.json
	configFile, err := os.Open("config.local.json")
	if err != nil {
		// Nếu không có file config.dev.json thì đọc file config.default.json
		configFile, err = os.Open("config.default.json")
		if err != nil {
			panic(err)
		}
		defer configFile.Close()
	}
	defer configFile.Close()

	// Parse dữ liệu JSON và bind vào conf
	err = DecodeDataFromJsonFile(configFile, &conf)
	if err != nil {
		log.Println("Không đọc được file config.")
		panic(err)
	}

	return conf
}

func ConnectDb(user string, password string, database string, address string) *gorm.DB {
	connectionInfo := fmt.Sprintf(`%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local`, user, password, address, database)

	db, err := gorm.Open("mysql", connectionInfo)
	if err != nil {
		panic(err)
	}
	return db
}
