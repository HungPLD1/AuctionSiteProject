package controller

import (
	"hellogorm/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//  @Description Show Session information by ID, return a JSON form
//	@Param id path string true "Session id"
//  @Success 200 {object} model.SessionSearch
//	@Failure 500 {body} string "Error message"
//  @Router /session/:id [GET]
func BidSessionByID(c *gin.Context) {
	sessionid := c.Param("id")
	var session model.SessionSearch

	errGetSession := searchSessionSQL().
		Where("session_id = ?", sessionid).
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

	session.BidLogs = attachSessionLogs(session)
	session.Images = attachSessionImages(session.ItemID, session.Images)
	c.JSON(200, session)
	return
}

//  @Description Search Session by query, return a JSON form
//  @Param name query string true "Name of the item (or part of it)"
// 	@Param categories query string true "Item Categories by number"
//  @Success 200 {object} model.SessionSearch
//	@Failure 500 {body} string "Error message"
//  @Router /session [GET]
func BidSessionByQuery(c *gin.Context) {
	itemname := "%" + c.DefaultQuery("name", "all") + "%"
	itemcategories := c.DefaultQuery("categories", "all")
	var sessionlist []model.SessionSearch
	errGetSession := searchSessionSQL().
		Where("(item.item_name LIKE ? OR '%all%' = ?) AND (item.categories_id = ? OR 'all' = ?)", itemname, itemname, itemcategories, itemcategories).
		Scan(&sessionlist).
		Error
	if errGetSession != nil {
		log.Println(errGetSession)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errGetSession,
			"message": "Error while fetching session data",
		})
		return
	}
	for i, _ := range sessionlist {
		sessionlist[i].BidLogs = attachSessionLogs(sessionlist[i])
		sessionlist[i].Images = attachSessionImages(sessionlist[i].ItemID, sessionlist[i].Images)
	}
	c.JSON(200, sessionlist)
	return
}

//  @Description Create a new bidding session , return a JSON message
//  @Param Authorization header string true "Session token"
//  @Param NewSessionInfo body model.SessionSearch true "Information to be provided"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /session [POST]
func CreateBidSession(c *gin.Context) {
	var headerInfo model.AuthorizationHeader
	if err := c.ShouldBindHeader(&headerInfo); err != nil {
		c.JSON(200, err)
	}
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

	var newsession model.SessionSearch
	errJSON := c.BindJSON(&newsession)
	if errJSON != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errJSON,
			"message": "Not a valid JSON!",
		})
		return
	}

	item := model.Items{
		CategoriesID:    newsession.CategoriesID,
		ItemName:        newsession.ItemName,
		ItemDescription: newsession.ItemDescription,
		ItemCondition:   newsession.ItemCondition,
		ItemCreateat:    time.Now(),
	}
	db := GetDBInstance().Db
	if err := db.Table("item").Create(item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Cannot add item to database",
		})
		return
	}
	var itemID []int
	if err := db.Table("item").
		Where("item_name = ?", newsession.ItemName).
		Pluck("item_id", &itemID).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Cannot get info from database",
		})
		return
	}
	session := model.BidSession{
		ItemID:   itemID[0],
		SellerID: userID,
		//SessionStartDate:   newsession.SessionStartDate,
		//SessionEndDate:     newsession.SessionEndDate,
		SessionStartDate:   time.Now(),
		SessionEndDate:     time.Now().Add(time.Hour * 24),
		UserviewCount:      1,
		MinimumIncreaseBid: newsession.MinimumIncreaseBid,
		CurrentBid:         newsession.CurrentBid,
	}
	if err := db.Table("bid_session").Create(session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Cannot add session to database",
		})
		return
	}
	c.JSON(http.StatusOK, "Successfully create new session!")
	return
}

//  @Description Get bid activity/history of user.
//  @Param Authorization header string true "Session token"
//  @Success 200 {object} model.BidHistory
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /history/bid [GET]
func BidSessionHistory(c *gin.Context) {
	var headerInfo model.AuthorizationHeader
	if err := c.ShouldBindHeader(&headerInfo); err != nil {
		c.JSON(200, err)
	}
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
	var bidhistory []model.BidHistory
	errGetSession := db.Table("bid_session_log a").
		Joins("JOIN bid_session b ON a.session_id = b.session_id").
		Joins("JOIN item c ON b.item_id = c.item_id").
		Where("a.user_id = ?", userID).
		Select("a.session_id, a.bid_amount, a.bid_date, b.session_start_date, b.session_end_date, c.item_id, c.item_name").
		Scan(&bidhistory).
		Error
	if errGetSession != nil {
		log.Println(errGetSession)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errGetSession,
			"message": "Error while fetching session data",
		})
		return
	}
	for i, _ := range bidhistory {
		bidhistory[i].Images = attachSessionImages(bidhistory[i].ItemID, bidhistory[i].Images)
	}
	c.JSON(200, bidhistory)
	return
}

//  @Description Get sell history of user.
//  @Param Authorization header string true "Session token"
//  @Success 200 {object} model.SellHistory
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /history/sell [GET]
func SellSessionHistory(c *gin.Context) {
	var headerInfo model.AuthorizationHeader
	if err := c.ShouldBindHeader(&headerInfo); err != nil {
		c.JSON(200, err)
	}
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
	var sellhistory []model.SellHistory
	errGetSession := db.Table("bid_session a").
		Joins("JOIN item b ON a.item_id = b.item_id").
		Where("a.seller_id = ?", userID).
		Select("a.session_id, a.session_start_date, a.session_end_date, b.item_id, b.item_name").
		Scan(&sellhistory).
		Error
	if errGetSession != nil {
		log.Println(errGetSession)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errGetSession,
			"message": "Error while fetching session data",
		})
		return
	}
	for i, _ := range sellhistory {
		sellhistory[i].Images = attachSessionImages(sellhistory[i].ItemID, sellhistory[i].Images)
	}
	c.JSON(200, sellhistory)
	return
}

/***************************************************************************************/
/***********************************INTERNAL FUNCTION***********************************/
