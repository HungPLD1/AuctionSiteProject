package controller

import (
	"hellogorm/model"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//  @Description Show user WishList, return a JSON form
//  @Param Authorization header string true "Session token"
//  @Success 200 {object} model.Items
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /wishlist [GET]
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

	db := GetDBInstance().Db
	var wishlist []model.Items
	var images []model.ItemImage

	db.Table("item").
		Joins("JOIN user_wishlist ON item.item_id = user_wishlist.item_id").
		Where("user_wishlist.user_id = ?", userID).
		Select("item.*").
		Scan(&wishlist)

	for i, item := range wishlist {
		db.Table("item_image").
			Where("item_id = ?", item.ItemID).
			Select("*").
			Scan(&images)
		for _, link := range images {
			wishlist[i].Images = append([]string(wishlist[i].Images), link.Images)
		}
	}
	c.JSON(200, wishlist)
	return
}

// @Description Add new item to wishlist, return a JSON message
// @Param Authorization header string true "Session token"
// @Param itemid path string true "Item id to be added to wishlist"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /wishlist/:id [POST]
func AddItemToWishList(c *gin.Context) {
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

	db := GetDBInstance().Db
	id := c.Param("id")
	//check if the item exist in database
	if NotExist := db.Table("item").
		Where("item_id = ?", id).
		First(&model.Items{}).
		RecordNotFound(); NotExist == true {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Item ID!",
		})
		return
	}
	//Attach item ID to wishlist
	itemid, _ := strconv.Atoi(id)
	wishlist := model.UserWishlist{
		UserID:  userid,
		ItemID:  itemid,
		AddDate: time.Now(),
	}
	errCreate := db.Table("user_wishlist").Create(&wishlist).Error
	if errCreate != nil {
		log.Println(errCreate)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errCreate,
			"message": "Server Error: Cannot add item to wishlist!",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Item successfully added to wishlist!",
	})
	return
}

// @Description Remove item from wishlist, return a JSON message
// @Param Authorization header string true "Session token"
// @Param itemid path string true "Item id to be removed from wishlist"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /wishlist/:id [DELETE]
func RemoveItemFromWishList(c *gin.Context) {
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

	db := GetDBInstance().Db
	id := c.Param("id")
	//Delete item from wishlist
	errDelete := db.Table("user_wishlist").
		Where("item_id = ?", id).
		Delete(&model.UserWishlist{}).
		Error
	if errDelete != nil {
		log.Println(errDelete)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errDelete,
			"message": "Server Error: Cannot remove item from wishlist!",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Item successfully removed from wishlist!",
	})
	return
}
