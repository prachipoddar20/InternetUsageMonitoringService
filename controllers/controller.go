package controllers

import (
	"IUMS/database"
	"IUMS/models"
	"IUMS/service"
	"strconv"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Dependency Injection
type iumsRepo struct {
	Db *gorm.DB
}

func New() *iumsRepo {
	db := database.InitDb()

	return &iumsRepo{Db: db}
}

func (repository *iumsRepo) InsertDataIntoDb() {
	database.InsertDatasetIntoDB(repository.Db)
}

/*
 *  /analytics?date=<param>&limit=<param>&page=<param>
 *  sample: http://localhost:8000/analytics?date=10122022&limit=100&page=2
 */
func (repository *iumsRepo) GetUsageDetails(c *gin.Context) {
	dateProvided, isDateProvided := c.GetQuery("date")
	limitProvided := c.DefaultQuery("limit", "100")
	pageProvided := c.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitProvided)
	page, _ := strconv.Atoi(pageProvided)
	var date time.Time

	var errorMessage = make(map[string]string)

	//handling edge case
	if isDateProvided {
		//converting string to time in date format
		date, _ = time.Parse("02-01-2006", dateProvided[:2]+"-"+dateProvided[2:4]+"-"+dateProvided[4:])

		if date.After(time.Now()) {
			handleError(errorMessage, "invalid date", http.StatusUnprocessableEntity, c)
			return
		}
	} else {
		handleError(errorMessage, "please enter date", http.StatusBadRequest, c)
		return
	}

	defer c.Request.Body.Close()

	//calling service
	response, err := service.GetUsageDetails(repository.Db, date, limit, page)

	if err != nil {
		handleError(errorMessage, err.Error(), http.StatusInternalServerError, c)
		return
	}
	//return response as json
	c.JSON(http.StatusOK, successResponseForUsageDetails{Ok: true, Response: *response})
}

/*
 *  /search/user?username=<param>
 *  sample: http://localhost:8000/user/search?username=obsessedHare3
 */
func (repository *iumsRepo) GetUserDetails(c *gin.Context) {

	username, isUsernameProvided := c.GetQuery("username")

	var errorMessage = make(map[string]string)

	//handling edge case
	if isUsernameProvided == false || username == "" {
		errorMessage["message"] = errorMessage["message"] + ", " + "please enter username"
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"ok": false, "error": errorMessage})
	}

	//calling service
	userDetails, err := service.GetSearchUserDetails(repository.Db, username)

	if err != nil {
		handleError(errorMessage, err.Error(), http.StatusInternalServerError, c)
		return
	}

	//return response as json
	c.JSON(http.StatusOK, successResponseForSearchUser{Ok: true, Response: *userDetails})
}

// private structures
type successResponseForUsageDetails struct {
	Ok       bool                  `json:"ok"`
	Response []models.UsageDetails `json:"data"`
}

type successResponseForSearchUser struct {
	Ok       bool               `json:"ok"`
	Response models.UserDetails `json:"data"`
}

// private functions
func handleError(errorMessage map[string]string, message string, httpstatus int, c *gin.Context) {
	errorMessage["message"] = message
	c.AbortWithStatusJSON(httpstatus, gin.H{
		"ok":    false,
		"error": errorMessage,
	})
}
