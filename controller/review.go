package controller

import (
	"hellogorm/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//  @Description Show review of User by User ID, return a JSON form
//	@Param id path string true "User id"
//  @Success 200 {object} model.UserReview
//	@Failure 500 {body} string "Error message"
//  @Router /review/:id [GET]
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
			"error":   errgetReview,
			"message": "Error while fetching review data",
		})
		return
	}
	c.JSON(200, review)
	return
}
