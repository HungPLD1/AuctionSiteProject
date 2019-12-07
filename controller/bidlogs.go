package controller

import (
	"hellogorm/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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

//  @Description Create a new bidding session , return a JSON message
//  @Param Authorization header string true "Session token"
//  @Param NewBidInfo body model.NewBidLog true "Information to be provided"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /logs [POST]
func NewBid(c *gin.Context) {
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

	var newbid model.NewBidLog
	if err := c.BindJSON(&newbid); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Not a valid JSON!",
		})
		return
	}
	newbidlog := model.BidSessionLog{
		UserID:    userID,
		SessionID: newbid.SessionID,
		BidAmount: newbid.BidAmount,
		BidDate:   time.Now(),
	}

	db := GetDBInstance().Db
	//Check if the bidder userid are different than seller userid
	//Check if the new bid amount are highest and respect the minimum bid increase
	var currentBid []int
	var minimumIncreaseBid []int
	var sellerID []string
	if err := db.Table("bid_session").
		Where("session_id = ?", newbid.SessionID).
		Pluck("seller_id", &sellerID).
		Pluck("current_bid", &currentBid).
		Pluck("minimum_increase_bid", &minimumIncreaseBid).
		Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err,
			"message": "Error while getting session bid data",
		})
		return
	}
	if newbid.BidAmount < (currentBid[0] + minimumIncreaseBid[0]) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bidding amount too low!",
		})
		return
	}
	if userID == sellerID[0] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Seller cannot bid",
		})
		return
	}

	//Create new Bid Logs
	if err := db.Table("bid_session_log").Create(newbidlog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Cannot create new Bid to database",
		})
		return
	}

	//update session current price
	session := model.BidSession{
		CurrentBid: newbid.BidAmount,
	}
	errUpdate := db.Table("bid_session").
		Model(&session).
		Where("session_id = ?", newbid.SessionID).
		Updates(session).
		Error
	if errUpdate != nil {
		log.Println(errUpdate)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   errUpdate,
			"message": "Error while updating session data",
		})
		return
	}

	c.JSON(http.StatusOK, "Successfully create new Bid to database!")
	return
}
