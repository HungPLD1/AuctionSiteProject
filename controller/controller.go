package controller

import (
	"hellogorm/model"
	"log"
	"net/http"
	"sync"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
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

//Showitems ...API: Show item by categories. Show all by default, result are JSON form
func Showitems(c *gin.Context) {
	db := GetDBInstance().Db
	categoriesName := c.Param("categories")
	var itemsList []model.Items

	if categoriesName == "all" {
		errGetItems := db.Table("item").Select("*").Scan(&itemsList).Error
		if errGetItems != nil {
			log.Println(errGetItems)
			return
		}
		c.JSON(200, itemsList)
	} else {
		errGetItems := db.Table("item, categories").
			Select("item.*").
			Where("categories.categories_name = ? AND item.categories_id = categories.categories_id", categoriesName).
			Scan(&itemsList).Error
		if errGetItems != nil {
			log.Println(errGetItems)
			return
		}
		c.JSON(200, itemsList)
	}
}

//SearchItemByName ...API: Search item by name, result are JSON form
func SearchItemByName(c *gin.Context) {
	db := GetDBInstance().Db
	itemname := "%" + c.DefaultQuery("name", "") + "%"
	var itemsList []model.Items

	errGetItems := db.Table("item").
		Select("*").
		Where("item_name LIKE ?", itemname).
		Scan(&itemsList).Error
	if errGetItems != nil {
		log.Println(errGetItems)
		return
	}
	c.JSON(200, itemsList)
}

//SearchItemByID ...API: Search item by ID, result are JSON form
func SearchItemByID(c *gin.Context) {
	db := GetDBInstance().Db
	itemid := c.DefaultQuery("id", "0")
	var itemsList []model.Items

	errGetItems := db.Table("item").
		Select("*").
		Where("item_id = ?", itemid).
		Scan(&itemsList).Error
	if errGetItems != nil {
		log.Println(errGetItems)
		return
	}
	c.JSON(200, itemsList)
}

//RegisterJSON ...API: Register new Account by JSON
func RegisterJSON(c *gin.Context) {
	db := GetDBInstance().Db
	var newUser model.User

	//check if the json form is valid
	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid JSON!",
		})
		return
	}

	//Check for empty field and password length
	if newUser.UserLoginID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Vui lòng nhập tên",
		})
		return
	}
	if len(newUser.UserPassword) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Mật khẩu phải có tối thiểu 4 kí tự",
		})
		return
	}

	//Fetch userdata from database to check for existing username
	var usersList []model.User
	errGetUsers := db.Table("user_common").
		Select("user_phone, user_login_id").
		Scan(&usersList).Error
	if errGetUsers != nil {
		log.Println(errGetUsers)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Cannot connect to user database!",
		})
		return
	}
	for _, user := range usersList {
		if newUser.UserLoginID == user.UserLoginID {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Tên truy cập đã có người sử dụng!",
			})
			return
		}
	}

	//Encrypt the password
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.UserPassword), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Encrypt passsword error",
		})
		return
	}
	passwordHash := string(hash)

	//Generate new userID
	userID := xid.New().String()

	//Filling information in struct
	newUser = model.User{
		UserID:    10,
		UserName:  "",
		UserPhone: "",
		//UserBirth:    	  time.Time{},
		UserGender:       0,
		UserAddress:      "",
		UserLoginID:      newUser.UserLoginID,
		UserPassword:     passwordHash,
		UserAccessLevel:  1,
		UserSessionToken: "",
	}

	// Tạo token với Header lưu thông tin chung:
	// Loại token: JWT
	// Thuật toán mã hoá: HS256
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))

	// Truyền dữ liệu vào phần Claim của token
	// Dữ liệu có kiểu map[string]interface{} mô phỏng một cấu trúc dạng JSON
	token.Claims = jwt_lib.MapClaims{
		"userId": userID,
		"Role":   newUser.UserAccessLevel,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	}

	// Tạo Signature cho token
	// Signature = HS256(Header, Claim, mysupersecretpassword)
	// Sử dụng secretkey như một input đầu vào
	// để thuật toán HS256 tạo ra chuỗi signature
	tokenString, err := token.SignedString([]byte(model.SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while generating token!",
		})
		return
	}
	newUser.UserSessionToken = tokenString

	//Save account info to database
	errInsertDb := db.Table("user_common").Create(newUser).Error
	if errInsertDb != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error: Cannot save new user",
		})
		return
	}

	//Generating success repond
	var rsp = model.SignupLoginResponse{
		ResponseTime: time.Now().String(),
		Code:         0,
		Message:      "Đăng kí thành công",
		Data:         newUser,
	}

	c.JSON(http.StatusOK, rsp)
	return
}

//ShowWishList ...Show user WishList, result are JSON form
func ShowWishList(c *gin.Context) {
	//db := GetDBInstance().Db
	//var itemsList []model.Items
}

/******INTERNAL FUNCTIONS******/
