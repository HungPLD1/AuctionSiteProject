package controller

import (
	"AuctionSiteProject/model"
	"log"
	"strings"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"

	"golang.org/x/crypto/bcrypt"
)

const administrator int = 9

/**********************************************************************/
/**************************INTERNAL FUNCTIONS**************************/
//check if the username exist in database or not
func checkUserByID(userID string) (bool, error) {
	db := GetDBInstance().Db
	//var usersList model.UserCommon
	if NotExist := db.Table("user_common").
		Where("user_id = ?", userID).
		First(&model.UserCommon{}).
		RecordNotFound(); NotExist == true {
		return false, nil
	}
	return true, nil
}

//check if the password is correct or not
func checkUserPassword(userName string, userPassword string) (bool, error) {
	db := GetDBInstance().Db
	var user model.UserCommon
	err := db.Table("user_common").
		Where("user_id = ?", userName).
		Select("user_password").
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

/** Check if user are Administrator*/
func checkAdministrator(userID string) (bool, error) {
	db := GetDBInstance().Db
	var accesslevel []int
	if err := db.Table("user_common").
		Where("user_id = ?", userID).
		Pluck("user_access_level", &accesslevel).
		Error; err != nil {
		return false, err
	}
	if accesslevel[0] != administrator {
		return false, nil
	}
	return true, nil
}

/** SQL query for session searching*/
func searchSessionSQL() *gorm.DB {
	db := GetDBInstance().Db
	return db.Table("bid_session").
		Joins("JOIN item ON item.item_id = bid_session.item_id").
		Joins("JOIN user_common ON user_common.user_id = bid_session.seller_id").
		Joins("JOIN categories ON categories.categories_id = item.categories_id").
		Select("bid_session.session_id, bid_session.item_id, bid_session.session_start_date, bid_session.session_end_date, bid_session.minimum_increase_bid, bid_session.seller_id, bid_session.current_bid, item.item_name, item.item_description, user_common.user_name AS seller_name, categories.categories_name")
}

/** Attach the bidding log data to session model */
func attachSessionLogs(session model.SessionSearch) []model.BidSessionLog {
	db := GetDBInstance().Db
	var result []model.BidSessionLog
	db.Table("bid_session_log").
		Where("session_id = ?", session.SessionID).
		Select("user_id, bid_amount, bid_date").
		Scan(&result)
	return result
}

/** Attach the images data to session model */
func attachSessionImages(itemid int, imageslink []string) []string {
	db := GetDBInstance().Db
	var images []string
	db.Table("item_image").
		Where("item_id = ?", itemid).
		Select("*").
		Pluck("images", &images)
	imageslink = append([]string(imageslink), images...)
	return imageslink
}

/****************************NOT YET INCLUDED*************************/

/****************************TO DO LIST*************************/
/*
GENERAL:
Tách controller.go ra làm nhiều files nhỏ để dễ tìm kiếm Chẳng hản như get.go, post.go, put.go, delete.go
Phải tạo model error message thay vì sử dụng gin.H{} . Về sau có thể tích hợp vào swagger
nghiên cứu cách update session tự động
Sua lai Status errors

API:
API delete wishlist kiểm tra item id có tồn tại trong database hay không?
Làm upload API

*/
