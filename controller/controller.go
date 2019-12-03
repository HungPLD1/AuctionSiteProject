package controller

import (
	"hellogorm/model"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

//GetItemByID ...API: Search item by ID, result are JSON form
//  @Description Get the item's informations and pictures by ID
//  @Param id path string the item ID number
//  @Success 200 {object} model.Items
//  @Failure 500 message : Error while fetching item data
//  @Router /item/:id [GET]

func GetItemByID(c *gin.Context) {
	db := GetDBInstance().Db
	itemid := c.Param("id")

	var itemsList []model.Items
	var imagelink []model.ItemImage

	errGetItems := db.Table("item").
		Select("*").
		Where("item_id = ?", itemid).
		Scan(&itemsList).Error
	if errGetItems != nil {
		log.Println(errGetItems)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while fetching item data",
		})
		return
	}
	for i, item := range itemsList {
		db.Table("item_image").
			Where("item_id = ?", item.ItemID).
			Select("*").
			Scan(&imagelink)
		for _, link := range imagelink {
			itemsList[i].ImageLink = append([]string(itemsList[i].ImageLink), link.ImageLink)
		}
	}
	c.JSON(200, itemsList)
}

//GetItemByQuery ...API: Search item by query, result are JSON form
func GetItemByQuery(c *gin.Context) {
	db := GetDBInstance().Db
	itemname := "%" + c.DefaultQuery("name", "all") + "%"
	itemcategories := c.DefaultQuery("categories", "all")

	var itemsList []model.Items
	var imagelink []model.ItemImage

	errGetItems := db.Table("item").
		Select("item.*").
		Where("(item_name LIKE ? OR '%all%' = ?) AND (categories_id = ? OR 'all' = ?)", itemname, itemname, itemcategories, itemcategories).
		Scan(&itemsList).
		Error
	if errGetItems != nil {
		log.Println(errGetItems)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while fetching item data",
		})
		return
	}
	for i, item := range itemsList {
		db.Table("item_image").
			Where("item_id = ?", item.ItemID).
			Select("*").
			Scan(&imagelink)
		for _, link := range imagelink {
			itemsList[i].ImageLink = append([]string(itemsList[i].ImageLink), link.ImageLink)
		}
	}
	c.JSON(200, itemsList)
	return
}

//RegisterJSON ...API: Register new Account by JSON
func RegisterJSON(c *gin.Context) {
	db := GetDBInstance().Db
	var newUser model.UserCommon

	//check if the json form is valid
	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid JSON!",
		})
		return
	}

	//Check for empty field and password length
	if newUser.UserID == "" {
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
	if exist, err := checkUserByID(newUser.UserID); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while fetching user data",
		})
		return
	} else if exist == true {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Tên truy cập đã có người sử dụng!",
		})
		return
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

	//Filling information in struct
	newUser = model.UserCommon{
		UserID:    newUser.UserID,
		UserName:  "",
		UserPhone: "",
		//UserBirth:    	  time.Time{},
		UserGender:       0,
		UserAddress:      "",
		UserPassword:     passwordHash,
		UserAccessLevel:  1,
		UserSessionToken: "",
	}
	newUser.UserSessionToken, _ = tokenGenerate(newUser)

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

//LoginJSON ...API: Login by JSON
func LoginJSON(c *gin.Context) {
	db := GetDBInstance().Db
	var userLogin model.UserCommon

	err := c.BindJSON(&userLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid JSON!",
		})
		return
	}
	//Validing Login Info
	if exist, err := validLoginInfo(userLogin); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while fetching user data",
		})
		return
	} else if exist == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Tên truy cập hoặc mật khẩu không đúng!",
		})
		return
	}
	//Generate token
	if userLogin.UserSessionToken, err = tokenGenerate(userLogin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error,cannot create login session!",
		})
		return
	}
	//db.Save(&userLogin)
	db.Table("user_common").Where("user_id = ?", userLogin.UserID).Update("user_session_token", userLogin.UserSessionToken)
	return
}

//UserProfile ...API: Show user profile stored in jwt session token
func UserProfile(c *gin.Context) {
	//Get the auth key in header
	var headerInfo model.AuthorizationHeader
	if err := c.ShouldBindHeader(&headerInfo); err != nil {
		c.JSON(200, err)
	}
	//check token validation and get userID
	var userID string
	var errtoken error
	if userID, errtoken = checkSessionToken(headerInfo.Token); errtoken != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
		})
		return
	} else if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token không hợp lệ",
		})
		return
	}
	//Get and return user profile
	db := GetDBInstance().Db
	var userprofile model.UserCommon
	errprofile := db.Table("user_common").
		Select("user_name, user_phone,user_birth,user_gender,user_address").
		Where("user_id = ?", userID).
		Scan(&userprofile).Error
	if errprofile != nil {
		log.Println(errprofile)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while fetching user data",
		})
		return
	}
	c.JSON(200, userprofile)
	return
}

//UserProfileUpdate ...Update user info to database
func UserProfileUpdate(c *gin.Context) {
	//Get the auth key in header
	var headerInfo model.AuthorizationHeader
	if err := c.ShouldBindHeader(&headerInfo); err != nil {
		c.JSON(200, err)
	}
	//check token validation and get userID
	var userID string
	var errtoken error
	if userID, errtoken = checkSessionToken(headerInfo.Token); errtoken != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
		})
		return
	} else if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token không hợp lệ",
		})
		return
	}
	//Get user update info
	db := GetDBInstance().Db
	var userUpdate model.UserCommon
	errJSON := c.BindJSON(&userUpdate)
	if errJSON != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid JSON!",
		})
		return
	}
	//Register update
	errUpdate := db.Table("user_common").
		Model(&userUpdate).
		Omit("user_id", "user_password", "user_access_level", "user_session_token").
		Where("user_id = ?", userID).
		Updates(userUpdate).
		Error
	if errUpdate != nil {
		log.Println(errUpdate)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while updating user data",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Update Success!",
	})
	return
}

//ShowWishList ...API: Show user WishList, result are JSON form
func ShowWishList(c *gin.Context) {
	var headerInfo model.AuthorizationHeader
	if err := c.ShouldBindHeader(&headerInfo); err != nil {
		c.JSON(200, err)
	}
	//check token validation and get userID
	var userID string
	var errtoken error
	if userID, errtoken = checkSessionToken(headerInfo.Token); errtoken != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
		})
		return
	} else if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token không hợp lệ",
		})
		return
	}

	db := GetDBInstance().Db
	var wishlist []model.Items
	var imagelink []model.ItemImage
	db.Table("item").
		Joins("JOIN user_wishlist ON item.item_id = user_wishlist.item_id").
		Where("user_wishlist.user_id = ?", userID).
		Select("item.*").
		Scan(&wishlist)

	for i, item := range wishlist {
		db.Table("item_image").
			Where("item_id = ?", item.ItemID).
			Select("*").
			Scan(&imagelink)
		for _, link := range imagelink {
			wishlist[i].ImageLink = append([]string(wishlist[i].ImageLink), link.ImageLink)
		}
	}
	c.JSON(200, wishlist)
	return
}

//BidSession ...API: Search for session id
func BidSession(c *gin.Context) {
	db := GetDBInstance().Db
	sessionid := c.Param("id")
	var session model.BidSession

	errGetSession := db.Table("bid_session").
		Select("*").
		Where("session_id = ?", sessionid).
		Scan(&session).
		Error
	if errGetSession != nil {
		log.Println(errGetSession)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while fetching session data",
		})
		return
	}
	c.JSON(200, session)
	return
}

//BidLogs ...API: Get Bid Session Logs
func BidLogs(c *gin.Context) {
	db := GetDBInstance().Db
	sessionid := c.Param("id")
	var logs []model.BidSessionLog

	errgetLogs := db.Table("bid_session_log").
		Where("session_id = ?", sessionid).
		Select("user_id, bid_amount, bid_date").
		Scan(&logs).
		Error
	if errgetLogs != nil {
		log.Println(errgetLogs)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while fetching bidding data",
		})
		return
	}
	c.JSON(200, logs)
	return
}

//ShowReview ...API: Show review of user
func ShowReview(c *gin.Context) {
	db := GetDBInstance().Db
	userid := c.Param("id")
	var review []model.UserReview

	errgetReview := db.Table("user_review").
		Where("user_target = ?", userid).
		Select("*").
		Scan(&review).
		Error
	if errgetReview != nil {
		log.Println(errgetReview)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while fetching review data",
		})
		return
	}
	c.JSON(200, review)
	return
}

//AddWishList ...API: Search for session id

/**********************************************************************/
/**************************INTERNAL FUNCTIONS**************************/
//check if the username exist in database or not
func checkUserByID(UserID string) (bool, error) {
	db := GetDBInstance().Db
	var usersList model.UserCommon
	err := db.Table("user_common").
		Select("user_id").
		Where("user_id = ?", UserID).
		Scan(&usersList).
		Error
	if err != nil {
		return false, err
	}
	if usersList.UserID != UserID {
		return false, nil
	}
	return true, nil
}

//check if the password is correct or not
func checkUserPassword(userName string, userPassword string) (bool, error) {
	db := GetDBInstance().Db
	var user model.UserCommon
	err := db.Table("user_common").
		Select("user_password").
		Where("user_id = ?", userName).
		Scan(&user).Error
	if err != nil {
		return false, err
	}
	byteHash := []byte(user.UserPassword)
	check := bcrypt.CompareHashAndPassword(byteHash, []byte(userPassword))
	if check != nil {
		return false, nil
	}
	return true, nil
}

/**	Function used by following API: /login
*	Check the user info with those in database.
*	First check if the username exist, then compare the password.
*	Return true if data are matched, false otherwise.
**/
func validLoginInfo(userLogin model.UserCommon) (bool, error) {
	//Check username
	var userCheck, passCheck bool
	var err error
	if userCheck, err = checkUserByID(userLogin.UserID); err != nil {
		return userCheck, err
	}
	//Check password
	if passCheck, err = checkUserPassword(userLogin.UserID, userLogin.UserPassword); err != nil {
		return passCheck, err
	}
	return (userCheck && passCheck), nil
}

/**	Function used by following API: /login , /register
*	Generate a jwt token string to save the login session.
*	Return string: the jwt token, error: Error when generating the token.
**/
func tokenGenerate(user model.UserCommon) (string, error) {
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))

	token.Claims = jwt_lib.MapClaims{
		"userId": user.UserID,
		"Role":   user.UserAccessLevel,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	return token.SignedString([]byte(model.SecretKey))
}

/**	Function used by following API: /profile
*	Check the validation of jwt token session
*	Return userID if token are valid, a empty string otherwise
**/
func checkSessionToken(token string) (string, error) {
	tokenFromHeader := strings.Replace(token, "Bearer ", "", -1)
	claims := jwt_lib.MapClaims{}
	tkn, err := jwt_lib.ParseWithClaims(tokenFromHeader, claims, func(token *jwt_lib.Token) (interface{}, error) {
		return []byte(model.SecretKey), nil
	})
	//In case of error, check for it
	if err != nil {
		if err == jwt_lib.ErrSignatureInvalid {
			log.Println("error 1: Invalid Token")
			return "", nil
		}
		log.Println("error 2: Bad Request", err)
		return "", err
	}
	//Check for token expiration date
	if !tkn.Valid {
		log.Println("error 3: Invalid Token")
		return "", nil
	}
	//Get and userID from the token
	var userID string
	//var roleFromToken int
	for k, v := range claims {
		if k == "userId" {
			userID = v.(string)
		}
		/*if k == "Role" {
			roleFromToken = int(v.(float64))
		}*/
	}
	//Check if the user exist in database
	if chk, _ := checkUserByID(userID); chk == false {
		log.Println("error 1: Invalid Token")
		return "", nil
	}

	log.Println("Success: Token are valid")
	return userID, nil
}

/****************************NOT YET INCLUDED*************************/
