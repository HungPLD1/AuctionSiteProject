package controller

import (
	"AuctionSiteProject/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//	@Tags Review Controller
//	@Summary Lấy review bằng userid
//  @Description Lấy review bằng userid dưới dạng path, trả về form JSON.
//	@Param id path string true "User id"
//  @Success 200 {object} model.UserReview
//	@Failure 500 {body} string "Error message"
//  @Router /review/:userid [GET]
func ShowReview(c *gin.Context) {
	db := GetDBInstance().Db
	userid := c.Param("userid")
	var review []model.UserReview

	errgetReview := db.Table("user_review").
		Where("user_target = ?", userid).
		Select("*").
		Scan(&review).
		Error
	if errgetReview != nil {
		log.Println(errgetReview)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errgetReview,
			"message": "Error while fetching review data",
		})
		return
	}
	c.JSON(200, review)
	return
}

//	@Tags Review Controller
//	@Summary Tạo review mới
//  @Description Tạo review từ user đến user khác. Review phải được đính với id session đấu già mà user đó thắng. Mỗi id session chỉ có duy nhất 1 review. Một user không thể review chính bản thân mình.
//  @Param Authorization header string true "Session token"
//  @Param NewSessionInfo body model.NewReview true "Information to be provided"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /review [POST]
func CreateReview(c *gin.Context) {
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
	var reviewinfo model.NewReview
	errJSON := c.BindJSON(&reviewinfo)
	if errJSON != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errJSON,
			"message": "Not a valid JSON!",
		})
		return
	}
	//check if the user have the right to write review
	db := GetDBInstance().Db
	var sessionstatus []string
	var sellerID []string
	var winnerID []string
	db.Table("bid_session").
		Where("session_id = ?", reviewinfo.SessionID).
		Pluck("session_status", &sessionstatus).
		Pluck("winner_id", &winnerID).
		Pluck("seller_id", &sellerID)

	if (sessionstatus[0] != FINISHED) && (sessionstatus[0] != CANCELLED) {
		c.JSON(http.StatusUnauthorized, "Unauthorized: Không thể viết review khi phiên đấu giá đang chạy")
		return
	}
	if userID != winnerID[0] {
		c.JSON(http.StatusUnauthorized, "Unauthorized: Không phải là user thắng phiên đấu giá!")
		return
	}
	if userID == sellerID[0] {
		c.JSON(http.StatusUnauthorized, "Unauthorized: Không thể review phiên đấu giá của chính mình!")
		return
	}
	//Add review info to database
	newreview := model.UserReview{
		UserWriter:    userID,
		UserTarget:    reviewinfo.UserTarget,
		SessionID:     reviewinfo.SessionID,
		ReviewContent: reviewinfo.ReviewContent,
		ReviewScore:   reviewinfo.ReviewScore,
	}
	if err := db.Table("user_review").Create(newreview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Cannot add review to database",
		})
		return
	}
	c.JSON(http.StatusOK, "Successfully add review to database!")
	return
}

//	@Tags Review Controller
//	@Summary Xóa review của user dựa theo session id (administrator only)
//  @Param Authorization header string true "Session token"
//  @Param sessionID path string true "ID session"
//  @Success 200 {body} string "Success message"
//	@Failure 400 {body} string "Error message"
//	@Failure 401 {body} string "Error message"
//	@Failure 500 {body} string "Error message"
//  @Router /review/:sessionid [DELETE]
func DeleteReview(c *gin.Context) {
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

	//delete review
	sessionid := c.Param("sessionid")
	db := GetDBInstance().Db
	errDeleteReview := db.Table("user_review").Where("session_id = ?", sessionid).Delete(&model.UserReview{}).Error
	if errDeleteReview != nil {
		log.Println(errDeleteReview)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errDeleteReview,
			"message": "Error while deleting review",
		})
		return
	}
	c.JSON(http.StatusOK, "Successfully delete review!")
	return
}
