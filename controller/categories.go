package controller

import (
	"hellogorm/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//  @Description Search categories by id, return all by default, return a JSON form
//  @Param id query string true "id of categories, if empty then return all"
//  @Success 200 {object} model.Categories
//	@Failure 500 {body} string "Error message"
//  @Router /categories [GET]
func SearchCategories(c *gin.Context) {
	db := GetDBInstance().Db
	var categories []model.Categories
	id := c.DefaultQuery("id", "all")

	errGetCategories := db.Table("categories").
		Where("categories_id = ? OR 'all' = ?", id, id).
		Select("*").
		Scan(&categories).
		Error
	if errGetCategories != nil {
		log.Println(errGetCategories)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errGetCategories,
			"message": "Error while fetching categories data",
		})
		return
	}
	c.JSON(200, categories)
	return
}
