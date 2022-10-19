package metadata

import (
	"net/http"
	"sync"

	"github.com/haapjari/glass/pkg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRepositories(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var repositories []models.Repository
	db.Find(&repositories)

	c.JSON(http.StatusOK, gin.H{"data": repositories})
}

func CreateRepository(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate Input
	var input models.CreateRepositoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Model
	repository := models.Repository{Name: input.Name, Uri: input.Uri}
	db.Create(&repository)
	c.JSON(http.StatusOK, gin.H{"data": repository})
}

func GetRepositoryById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Attempt to Fetch Models
	var repository models.Repository
	if err := db.Where("id = ?", c.Param("id")).First(&repository).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": repository})
}

func DeleteRepositoryById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Attempt to Fetch Models
	var repository models.Repository
	if err := db.Where("id = ?", c.Param("id")).First(&repository).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&repository)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func UpdateRepositoryById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Attempt to Fetch Models
	var repository models.Repository
	if err := db.Where("id = ?", c.Param("id")).First(&repository).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate Input
	var input models.UpdateRepositoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&repository).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": repository})
}

func CreateRepositoriesFromFile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Check the ContentLenght, if it's 0, there is no File.
	if c.Request.ContentLength == 0 {
		c.IndentedJSON(http.StatusBadRequest, "Include File in Request")
	}

	// Initialize Object
	var obj []models.Repository

	// Bind JSON Data to Object
	c.BindJSON(&obj)

	// Begin Multi-Thread For-Loop

	var objectLength int = len(obj)

	// Initialize Wait Group
	var waitGroup sync.WaitGroup
	waitGroup.Add(objectLength)

	for i := 0; i < objectLength; i++ {

		// Spawn a Thread for invidual Iteration
		// Pass 'i' into the GoRoutines function
		//		in order to make sure each goroutine
		//		uses different value for 'i'
		go func(i int) {
			// At the end of the GoRoutine, tell the WaitGroup
			//		that another thread has completed.
			defer waitGroup.Done()
			repository := models.Repository{Name: obj[i].Name, Uri: obj[i].Uri}
			db.Create(&repository)
		}(i)
	}

	// Wait should be called as many times as Add.
	waitGroup.Wait()

	c.IndentedJSON(http.StatusOK, "File Data Succesfully Loaded to Database.")
}
