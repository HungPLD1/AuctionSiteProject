package controller

import (
	"AuctionSiteProject/model"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//	@Tags BidSession Controller
//	@Summary Tìm phiên đấu giá bằng id
//  @Description Tìm phiên đáu giá bằng id dưới dạng path, trả về JSON form
//	@Param id path string true "Session id"
//  @Success 200 {object} model.SessionSearch
//	@Failure 500 {body} string "Error message"
//  @Router /session/:sessionid [GET]
func BidSessionByID(c *gin.Context) {
	sessionid := c.Param("sessionid")
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

//	@Tags BidSession Controller
//	@Summary Tìm phiên đấu giá bằng tên mặt hàng và/hoặc categories
//  @Description Tìm phiên đấu giá bằng tên mặt hàng và/hoặc categories. Tên mặt hàng không cần viết đầy đủ. Trả về tất cả nếu để trống. Trả về JSON form
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

//	@Tags BidSession Controller
//	@Summary Hiển thị phiên đấu giá chưa thanh toán
//  @Description Hiển thị phiên đấu giá chưa thanh toán của người dùng. Cần phải đăng nhập để sử dụng. Trả về JSON form
//  @Param Authorization header string true "Session token"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /awaitpayment [GET]
func UnpaidSession(c *gin.Context) {
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
	var sessionlist []model.SessionSearch
	errGetSession := searchSessionSQL().
		Where("bid_session.winner_id = ? AND bid_session.session_status = ?", userID, AWAITINGPAYMENT).
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

//	@Tags BidSession Controller
//	@Summary Tạo phiên đấu giá mới
//  @Param Authorization header string true "Session token"
//  @Param NewSessionInfo body model.NewSession true "Information to be provided"
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

	var newsession model.NewSession
	errJSON := c.BindJSON(&newsession)
	if errJSON != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errJSON,
			"message": "Not a valid JSON!",
		})
		return
	}
	//create item
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

	//add image Url
	for _, image := range newsession.Images {
		newimage := model.ItemImage{
			ItemID: itemID[0],
			Images: image,
		}
		if err := db.Table("item_image").Create(newimage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err,
				"message": "Cannot add image link to database",
			})
			return
		}
	}

	//create session

	session := model.BidSession{
		ItemID:             itemID[0],
		SellerID:           userID,
		UserviewCount:      1,
		WinnerID:           userID,
		MinimumIncreaseBid: newsession.MinimumIncreaseBid,
		CurrentBid:         newsession.StartPrice,
		SessionStatus:      RUNNING,
	}
	var errtime error
	session.SessionStartDate, errtime = time.Parse(time.RFC3339, newsession.SessionStartDate)
	if errtime != nil {
		log.Println(errtime)
	}
	session.SessionEndDate, errtime = time.Parse(time.RFC3339, newsession.SessionEndDate)
	if errtime != nil {
		log.Println(errtime)
	}
	log.Println(newsession.SessionStartDate, newsession.SessionEndDate)
	log.Println(session.SessionStartDate, session.SessionEndDate)

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

//	@Tags BidSession Controller
//	@Summary Thay đổi thông tin phiên đấu giá.
//  @Param Authorization header string true "Session token"
//  @Param NewSessionInfo body model.UpdateSession true "Information to be provided"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /session [PUT]
func UpdateBidSession(c *gin.Context) {
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
	var updatedata model.UpdateSession
	errJSON := c.BindJSON(&updatedata)
	if errJSON != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errJSON,
			"message": "Not a valid JSON!",
		})
		return
	}

	db := GetDBInstance().Db
	var itemID []int
	var sellerid []string
	db.Table("bid_session").
		Where("session_id = ?", updatedata.SessionID).
		Pluck("item_id", &itemID).
		Pluck("seller_id", &sellerid)

	//check if user are administrator, if not then check if user are seller
	if check, err := checkAdministrator(userID); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Error while fetching user data",
		})
		return
	} else if check == false {
		if userID != sellerid[0] {
			c.JSON(http.StatusUnauthorized, "Unauthorized : Not the session owner")
			return
		}
	}

	//update the session/item info
	newitem := model.Items{
		ItemName:        updatedata.ItemName,
		ItemDescription: updatedata.ItemDescription,
		ItemCondition:   updatedata.ItemCondition,
	}
	errUpdate := db.Table("item").
		Model(&model.Items{}).
		Where("item_id = ?", itemID[0]).
		Updates(newitem).
		Error
	if errUpdate != nil {
		log.Println(errUpdate)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errUpdate,
			"message": "Error while updating session",
		})
		return
	}

	c.JSON(http.StatusOK, "Successfully update bid session!")
	return
}

//	@Tags BidSession Controller
//	@Summary Khóa phiên đấu giá (API nội bộ, không public cho user)
//  @Param SessionID path string true "session ID"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /lock/:sessionid [PUT]
func LockSession(c *gin.Context) {
	sessionid := c.Param("sessionid")
	var session model.BidSession
	var winnerid []string
	var sellerid []string
	db := GetDBInstance().Db
	db.Table("bid_session").
		Where("session_id = ?", sessionid).
		Pluck("seller_id", &sellerid).
		Pluck("winner_id", &winnerid)
	if sellerid[0] == winnerid[0] {
		session.SessionStatus = FINISHED
	} else {
		session.SessionStatus = AWAITINGPAYMENT
	}
	if errLock := db.Table("bid_session").
		Model(&model.BidSession{}).
		Where("session_id = ?", sessionid).
		Updates(session).
		Error; errLock != nil {
		log.Println(errLock)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   errLock,
			"message": "Error while locking session",
		})
	}

	c.JSON(http.StatusOK, "Successfully lock session!")
	return
}

//	@Tags BidSession Controller
//	@Summary Xóa phiên đấu giá (Administrator only)
//  @Description Xóa toàn bộ phiên đấu giá và các dữ liệu liên quan như lịch sử đấu giá, thông tin mặt hàng của phiên đấu giá.
//  @Param Authorization header string true "Session token"
//  @Param sessionid path string true "session id to be deleted"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /session/:sessionid [DELETE]
func DeleteBidSession(c *gin.Context) {
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
	//check if user are adminitrator
	if check, err := checkAdministrator(userID); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Error while fetching user data",
		})
		return
	} else if check == false {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Only Administrator can use this API!",
		})
		return
	}

	sessionid, _ := strconv.Atoi(c.Param("sessionid"))
	db := GetDBInstance().Db
	var itemid []int
	errGetItemID := db.Table("bid_session").Where("session_id = ?", sessionid).Pluck("item_id", &itemid).Error
	errDeleteSession := db.Table("bid_session").Where("session_id = ?", sessionid).Delete(&model.BidSession{}).Error
	errDeleteItem := db.Table("item").Where("item_id = ?", itemid[0]).Delete(&model.Items{}).Error
	if errGetItemID != nil {
		log.Println(errGetItemID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errGetItemID,
			"message": "Error while fetching session data",
		})
		return
	}
	if errDeleteSession != nil {
		log.Println(errDeleteSession)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errDeleteSession,
			"message": "Error while deleting session data",
		})
		return
	}
	if errDeleteItem != nil {
		log.Println(errDeleteItem)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errDeleteItem,
			"message": "Error while deleting item data",
		})
		return
	}

	c.JSON(http.StatusOK, "Successfully delete bid session")
	return
}

//	@Tags BidSession Controller
//	@Summary Hiển thị lịch sử đấu giá của user.
//  @Description Hiển thị những hoạt động đấu giá của user trên các session.
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

//	@Tags BidSession Controller
//	@Summary Hiển thị lịch sử bán hàng của user
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
