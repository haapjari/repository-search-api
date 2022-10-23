package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/haapjari/glass/pkg/models"
	"github.com/haapjari/glass/pkg/utils"
	"gorm.io/gorm"
)

type Handler struct {
	Context  *gin.Context
	Database *gorm.DB
	Client   *http.Client
}

func NewHandler(c *gin.Context) *Handler {
	h := new(Handler)

	h.Context = c
	h.Database = c.MustGet("db").(*gorm.DB)
	h.Client = &http.Client{}

	return h
}

func (h *Handler) HandleGetRepositories() {
	var e []models.Repository

	h.Database.Find(&e)

	h.Context.JSON(http.StatusOK, gin.H{"data": e})
}

func (h *Handler) HandleGetRepositoryById() {
	var e models.Repository

	if err := h.Database.Where("id = ?", h.Context.Param("id")).First(&e).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.Context.JSON(http.StatusOK, gin.H{"data": e})
}

func (h *Handler) HandleCreateRepository() {
	var i models.Repository

	if err := h.Context.ShouldBindJSON(&i); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	r := models.Repository{RepositoryName: i.RepositoryName, RepositoryUrl: i.RepositoryUrl, OpenIssueCount: i.OpenIssueCount, ClosedIssueCount: i.ClosedIssueCount, OriginalCodebaseSize: i.OriginalCodebaseSize, LibraryCodebaseSize: i.LibraryCodebaseSize, RepositoryType: i.RepositoryType, PrimaryLanguage: i.PrimaryLanguage}

	h.Database.Create(&r)

	h.Context.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *Handler) HandleDeleteRepositoryById() {
	var r models.Repository

	if err := h.Database.Where("id = ?", h.Context.Param("id")).First(&r).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.Database.Delete(&r)

	h.Context.JSON(http.StatusOK, gin.H{"succeed": r})
}

func (h *Handler) HandleUpdateRepositoryById() {
	var r models.Repository

	if err := h.Database.Where("id = ?", h.Context.Param("id")).First(&r).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var i models.Repository

	if err := h.Context.ShouldBindJSON(&i); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.Database.Model(&r).Updates(i)

	h.Context.JSON(http.StatusOK, gin.H{"data": r})
}

// TODO
func (h *Handler) HandleFetchRepositories() {
	var (
		repositoryCount, err = strconv.Atoi(h.Context.Query("count"))
		language             = h.Context.Query("type")
		requestUrl           = utils.GetSourceGraphGraphQlApiBaseurl()
		queryString          = "{search(query:\"lang:" + language + " AND select:repo AND repohasfile:go.mod AND count:" + strconv.Itoa(repositoryCount) + "\", version:V2){results{repositories{name}}}}"
	)

	if err != nil {
		log.Fatalln(err)
	}

	rawQuery := map[string]string{
		"query": queryString,
	}

	queryAsJson, err := json.Marshal(rawQuery)
	if err != nil {
		log.Fatalln(err)
	}

	requestBody := bytes.NewBuffer(queryAsJson)

	request, err := http.NewRequest("POST", requestUrl, requestBody)
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalln(err)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Parsing Unstructured Data from SourceGraph Response

	resMap := h.ParseSourceGraphResponse(string(data))

	// Print Values of the Response

	// TODO: How-To: "Has Next" - Functionality

	for i := 0; i < repositoryCount; i++ {
		if resMap["repositories"].([]interface{})[i].(map[string]interface{}) != nil {
			for _, value := range resMap["repositories"].([]interface{})[i].(map[string]interface{}) {
				h.SaveToRepositoryTable(fmt.Sprintf("%v", value))
			}
		}
	}
}

func (h *Handler) ParseSourceGraphResponse(data string) map[string]interface{} {
	var responseAsJsonMap map[string]interface{}
	err := json.Unmarshal([]byte(string(data)), &responseAsJsonMap)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: This might return an error, or not succeed.

	dataArray := responseAsJsonMap["data"]
	dataMap := dataArray.(map[string]interface{})

	searchArray := dataMap["search"]
	searchMap := searchArray.(map[string]interface{})

	resultsArray := searchMap["results"]
	resultsMap := resultsArray.(map[string]interface{})

	return resultsMap
}

func (h *Handler) SaveToRepositoryTable(s string) {
	fmt.Println(s)
}
