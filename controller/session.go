package controller

import (
	"hellogorm/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//  @Description Get the item's informations and pictures by ID, return a JSON form
//  @Param id path string true "the item ID number"
//  @Success 200 {object} model.Items
//	@Failure 500 {body} string "Error message"
//  @Router /item/:id [GET]
func GetItemByID(c *gin.Context) {
	db := GetDBInstance().Db
	itemid := c.Param("id")

	var itemsList []model.Items
	var images []model.ItemImage

	errGetItems := db.Table("item").
		Where("item_id = ?", itemid).
		Select("*").
		Scan(&itemsList).Error
	if errGetItems != nil {
		log.Println(errGetItems)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errGetItems,
			"message": "Error while fetching item data",
		})
		return
	}
	for i, item := range itemsList {
		db.Table("item_image").
			Where("item_id = ?", item.ItemID).
			Select("*").
			Scan(&images)
		for _, link := range images {
			itemsList[i].Images = append([]string(itemsList[i].Images), link.Images)
		}
	}
	c.JSON(200, itemsList)
}

//  @Description Search item by query, return a JSON form
//  @Param name query string true "Name of the item (or part of it)"
// 	@Param categories query string true "Item Categories by number"
//  @Success 200 {object} model.Items
//	@Failure 500 {body} string "Error message"
//  @Router /item [GET]
func GetItemByQuery(c *gin.Context) {
	db := GetDBInstance().Db
	itemname := "%" + c.DefaultQuery("name", "all") + "%"
	itemcategories := c.DefaultQuery("categories", "all")

	var itemsList []model.Items
	var images []model.ItemImage

	errGetItems := db.Table("item").
		Where("(item_name LIKE ? OR '%all%' = ?) AND (categories_id = ? OR 'all' = ?)", itemname, itemname, itemcategories, itemcategories).
		Select("item.*").
		Scan(&itemsList).
		Error
	if errGetItems != nil {
		log.Println(errGetItems)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errGetItems,
			"message": "Error while fetching item data",
		})
		return
	}
	for i, item := range itemsList {
		db.Table("item_image").
			Where("item_id = ?", item.ItemID).
			Select("*").
			Scan(&images)
		for _, link := range images {
			itemsList[i].Images = append([]string(itemsList[i].Images), link.Images)
		}
	}
	c.JSON(200, itemsList)
	return
}

//  @Description Show Session information by ID, return a JSON form
//	@Param id path string true "Session id"
//  @Success 200 {object} model.BidSession
//	@Failure 500 {body} string "Error message"
//  @Router /session/:id [GET]
func BidSession(c *gin.Context) {
	db := GetDBInstance().Db
	sessionid := c.Param("id")
	var session model.BidSession

	errGetSession := db.Table("bid_session").
		Where("session_id = ?", sessionid).
		Select("*").
		Scan(&session).
		Error
	if errGetSession != nil {
		log.Println(errGetSession)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errGetSession,
			"message": "Error while fetching session data",
		})
		return
	}
	c.JSON(200, session)
	return
}

//  @Description Get Bid Session Logs by session ID, return a JSON form
//	@Param id path string true "Session id"
//  @Success 200 {object} model.BidSessionLog
//	@Failure 500 {body} string "Error message"
//  @Router /logs/:id [GET]
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
			"error":   errgetLogs,
			"message": "Error while fetching bidding data",
		})
		return
	}
	c.JSON(200, logs)
	return
}
