package controller

import (
	"hellogorm/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//  @Description Show user profile, return user general profile in JSON form
//  @Param Authorization header string true "Session token"
//  @Success 200 {object} model.UserCommon
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /profile [GET]
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
			"error":   errtoken,
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
		Where("user_id = ?", userID).
		Select("user_id, user_name, user_email, user_phone,user_gender,user_address").
		Scan(&userprofile).Error
	if errprofile != nil {
		log.Println(errprofile)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errprofile,
			"message": "Error while fetching user data",
		})
		return
	}
	c.JSON(200, userprofile)
	return
}

//  @Description Modify/Update user profile, return message in JSON form
//  @Param Authorization header string true "Session token"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /profile [PUT]
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
			"error":   errtoken,
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
			"error":   errJSON,
			"message": "Not a valid JSON!",
		})
		return
	}
	//Register update
	errUpdate := db.Table("user_common").
		Model(&userUpdate).
		Omit("user_id", "user_password", "user_access_level", "user_createAt", "user_avatar").
		Where("user_id = ?", userID).
		Updates(userUpdate).
		Error
	if errUpdate != nil {
		log.Println(errUpdate)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errUpdate,
			"message": "Error while updating user data",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Update Success!",
	})
	return
}

// @Description Add new item to wishlist, return a JSON message
// @Param Authorization header string true "Session token"
// @Param oldpassword body string true "Old Password"
// @Param newpassword body string true "New Password"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /password [PUT]
func UpdatePassword(c *gin.Context) {
	var headerInfo model.AuthorizationHeader
	if err := c.ShouldBindHeader(&headerInfo); err != nil {
		c.JSON(200, err)
	}
	//check token validation and get userID
	var userid string
	var errtoken error
	if userid, errtoken = checkSessionToken(headerInfo.Token); errtoken != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errtoken,
			"message": "Bad request",
		})
		return
	} else if userid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token không hợp lệ",
		})
		return
	}
	//Get user update info
	var passwordinfo model.ModifyPassword
	errJSON := c.BindJSON(&passwordinfo)
	if errJSON != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errJSON,
			"message": "Not a valid JSON!",
		})
		return
	}
	//check the old password
	if passCheck, err := checkUserPassword(userid, passwordinfo.OldPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Error while checking password!",
		})
		return
	} else if passCheck == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid JSON!",
		})
		return
	}

	//Encrypt the new password
	hash, errcrypt := bcrypt.GenerateFromPassword([]byte(passwordinfo.NewPassword), bcrypt.MinCost)
	if errcrypt != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errcrypt,
			"message": "Encrypt passsword error",
		})
		return
	}
	passwordHash := string(hash)

	//Update the new password
	db := GetDBInstance().Db
	passwordUpdate := model.UserCommon{
		UserPassword: passwordHash,
	}
	errUpdate := db.Table("user_common").
		Model(&passwordUpdate).
		Select("user_password").
		Where("user_id = ?", userid).
		Updates(passwordUpdate).
		Error
	if errUpdate != nil {
		log.Println(errUpdate)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   errUpdate,
			"message": "Error while updating user data",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Update Success!",
	})
	return
}
