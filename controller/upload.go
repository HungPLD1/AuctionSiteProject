package controller

import (
	"hellogorm/model"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

//  @Description Upload one or multiples photos of item , return a JSON message
// @Param Authorization header string true "Session token"
// @Param itemid path string true "Item id to be removed from wishlist"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /upload [POST]
func UploadItemImages(c *gin.Context) {
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

	//Get files data from Multiform
	var uploadform model.UploadItemImageForm
	if err := c.ShouldBind(&uploadform); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Not a valid form!",
		})
		return
	}
	//Copy files to server
	db := GetDBInstance().Db
	for _, file := range uploadform.Images {
		filename := path.Join(".", "view", "images", file.Filename)
		log.Println(filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err,
				"message": "Error while saving photos!",
			})
			return
		}
		//Update filename to database
		images := model.ItemImage{
			ItemID: uploadform.ItemID,
			Images: filename,
		}
		errCreate := db.Table("item_image").Create(&images).Error
		if errCreate != nil {
			log.Println(errCreate)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   errCreate,
				"message": "Error while updating database!",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Success": "Upload Files Successfully!",
	})
	return
}
