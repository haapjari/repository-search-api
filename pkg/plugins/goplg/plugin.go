package goplg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/haapjari/glass/pkg/utils"
)

var (
	GITHUB_API_TOKEN string = fmt.Sprintf("%v", utils.GetGithubApiToken())
	GITHUB_USERNAME  string = fmt.Sprintf("%v", utils.GetGithubUsername())
)

type GoPlugin struct {
	GitHubApiToken string
	GitHubUsername string
	HttpClient     *http.Client
	Parser         *Parser
}

func NewGoPlugin() *GoPlugin {
	g := new(GoPlugin)

	g.HttpClient = &http.Client{}
	g.Parser = NewParser()

	return g
}

func (g *GoPlugin) GetRepositoryMetadata(count int) {
	g.getBaseMetadataFromSourceGraph(count)
	g.enrichMetadataWithGitHubApi() // TODO
}

// PRIVATE METHODS

func (g *GoPlugin) getBaseMetadataFromSourceGraph(count int) {
	url := utils.GetSourceGraphGraphQlApiBaseurl()
	queryStr := "{search(query:\"lang:go +  AND select:repo AND repohasfile:go.mod AND count:" + strconv.Itoa(count) + "\", version:V2){results{repositories{name}}}}"

	rawReqBody := map[string]string{
		"query": queryStr,
	}

	jsonReqBody, err := json.Marshal(rawReqBody)
	utils.LogErr(err)

	bytesReqBody := bytes.NewBuffer(jsonReqBody)

	request, err := http.NewRequest("POST", url, bytesReqBody)
	request.Header.Set("Content-Type", "application/json")
	utils.LogErr(err)

	res, err := g.HttpClient.Do(request)
	utils.LogErr(err)

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	utils.LogErr(err)

	resMap, err := g.Parser.ParseSourceGraphResponse(string(resBody))
	utils.LogErr(err)

	if resMap != nil {
		for i := 0; i < count; i++ {
			if resMap["repositories"].([]interface{})[i].(map[string]interface{}) != nil {
				for _, value := range resMap["repositories"].([]interface{})[i].(map[string]interface{}) {
					mapBody := map[string]string{
						"repository_name": fmt.Sprintf("%v", value),
						"repository_url":  fmt.Sprintf("%v", value),
					}

					jsonBody, err := json.Marshal(mapBody)
					utils.LogErr(err)

					reqBody := bytes.NewBuffer(jsonBody)

					req, err := http.NewRequest("POST", utils.GetBaseurl()+"/"+"api/glass/v1/repository", reqBody)
					req.Header.Set("Content-Type", "application/json")
					utils.LogErr(err)

					_, err = g.HttpClient.Do(req)
					utils.LogErr(err)
				}
			}
		}
	}
}

func (g *GoPlugin) enrichMetadataWithGitHubApi() {
	fmt.Println("TODO")
}

// --- --- --- --- --- --- //
// --- --- --- --- --- --- //

// 		// Construct a Request Object
// 		req, err := http.NewRequest("GET", baseString+queryString, nil)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		// Execute Request Object
// 		req.Header.Set("Accept", "application/vnd.github.v3+json")
// 		src := oauth2.StaticTokenSource(
// 			&oauth2.Token{AccessToken: GITHUB_API_TOKEN},
// 		)

// 		client := oauth2.NewClient(context.Background(), src)

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		// Read Response Body (Will be read as ByteStream)
// 		body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		// Parse Response Body to JSON
// 		var repositoryResponse models.RepositoryResponse
// 		json.Unmarshal([]byte(body), &repositoryResponse)

// 		// Read Batches Data to Slice and Write to Database
// 		for _, item := range repositoryResponse.Items {
// 			var repository models.Repository = models.Repository{Name: item.Name, Uri: item.Url}
// 			db.Create(&repository)
// 		}
// 	}

// 	return c
// }

// // Counts lines of code in project, that is listed in repositories table and counts lines
// //
// //	library code in that project. Writes individual library code in libraries table and
// //	analyzed projects code in data_entties table.
// func countLines(c *gin.Context) *gin.Context {

// 	db := c.MustGet("db").(*gorm.DB)

// 	var r models.Repository
// 	var count int64

// 	db.Model(&r).Count(&count)
// 	countAsInteger := int(count)

// 	var repositories []models.Repository
// 	db.Find(&repositories)

// 	// loop through repositories
// 	for i := 0; i < countAsInteger; i++ {

// 		// parse uri
// 		var repositoryUri string = repositories[i].Uri + ".git"
// 		repositoryUri = strings.ReplaceAll(repositoryUri, "api.", "")
// 		repositoryUri = strings.ReplaceAll(repositoryUri, "repos/", "")

// 		// parse name
// 		var repositoryName string = repositories[i].Name

// 		var flag bool
// 		c, flag = utils.CheckIfAnalyzed(repositoryName, "main", c)
// 		fmt.Println("Check if the Repository is already analyzed: ", flag)

// 		var output []byte
// 		if !flag {
// 			output = utils.ExecuteInShell("git", "clone", repositoryUri)
// 			fmt.Println("Clone the Repository: ", output)
// 		}

// 		if !flag {
// 			utils.RunMainAnalysis(repositoryName, repositoryUri, db)
// 		}

// 		path := repositoryName + "/" + "package.json"

// 		// Search for: "package.json"
// 		if _, err := os.Stat(path); err == nil {

// 			packageJson, err := ioutil.ReadFile(path)
// 			utils.CheckError(err)

// 			packageJsonAsString := string(packageJson)

// 			if strings.Contains(packageJsonAsString, "dependencies") {
// 				utils.RunLibraryAnalysis(packageJsonAsString, db, repositoryName, `"dependencies": {`)
// 			}

// 			if strings.Contains(packageJsonAsString, "devDependencies") {
// 				utils.RunLibraryAnalysis(packageJsonAsString, db, repositoryName, `"devDependencies": {`)
// 			}

// 			if strings.Contains(packageJsonAsString, "peerDependencies") {
// 				utils.RunLibraryAnalysis(packageJsonAsString, db, repositoryName, `"peerDependencies": {`)
// 			}

// 			if strings.Contains(packageJsonAsString, "optionalDependencies") {
// 				utils.RunLibraryAnalysis(packageJsonAsString, db, repositoryName, `"optionalDependencies": {`)
// 			}

// 			if strings.Contains(packageJsonAsString, "bundledDependencies") {
// 				utils.RunLibraryAnalysis(packageJsonAsString, db, repositoryName, `"bundledDependencies": {`)
// 			}

// 		} else if errors.Is(err, os.ErrNotExist) {
// 			fmt.Println("Unable to locate: package.json")
// 		} else {
// 			fmt.Println("Unable to locate: package.json")
// 		}

// 		output = utils.ExecuteInShell("rm", "-rf", repositoryName)
// 		fmt.Println("Delete the Repository: ", output)
// 	}

// 	return c
// }

// // Fetches Issue count for every repository, that is in data_entities table
// func fetchIssues(c *gin.Context) *gin.Context {

// 	// Initialize GORM
// 	db := c.MustGet("db").(*gorm.DB)

// 	// Fetch all the repositories to slice (Similar to GET all)
// 	var dataEntities []models.Entity
// 	db.Find(&dataEntities)

// 	length := len(dataEntities)
// 	for i := 0; i < length; i++ {
// 		// If URI is not "" -> Parse

// 		var parsedDataEntityUri string
// 		var sliceOfDataEntityUri []string
// 		var repositoryOwner string
// 		var repositoryName string

// 		if dataEntities[0].Uri != "" {
// 			parsedDataEntityUri = strings.ReplaceAll(dataEntities[i].Uri, "https://github.com/", "")
// 			parsedDataEntityUri = strings.ReplaceAll(parsedDataEntityUri, ".git", "")
// 			sliceOfDataEntityUri = strings.Split(parsedDataEntityUri, "/")
// 			repositoryOwner = string(sliceOfDataEntityUri[0])
// 			repositoryName = string(sliceOfDataEntityUri[1])
// 		}

// 		// GraphQL Query
// 		jsonData := map[string]string{
// 			"query": `
//             {
//   				repository(owner:"` + repositoryOwner + `" name:"` + repositoryName + `") {
//     				issues {
//       					totalCount
//     				}
//   				}
//             }
//         `,
// 		}

// 		jsonValue, _ := json.Marshal(jsonData)

// 		// Construct a Request Object
// 		request, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonValue))
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		request.Header.Set("Accept", "application/vnd.github.v3+json")

// 		src := oauth2.StaticTokenSource(
// 			&oauth2.Token{AccessToken: GITHUB_API_TOKEN},
// 		)

// 		// If Issue Count is 0 -> Update it.
// 		if dataEntities[i].Issct == 0 {
// 			client := oauth2.NewClient(context.Background(), src)

// 			response, err := client.Do(request)
// 			utils.CheckError(err)
// 			defer response.Body.Close()

// 			body, err := ioutil.ReadAll(response.Body)
// 			if err != nil {
// 				log.Fatalln(err)
// 			}

// 			issueCount := strings.ReplaceAll(string(body), `{"data":{"repository":{"issues":{"totalCount":`, "")
// 			issueCount = strings.ReplaceAll(issueCount, `}}}}`, "")
// 			issueCountAsInteger, err := strconv.Atoi(issueCount)
// 			utils.CheckError(err)

// 			var d models.Entity
// 			db.First(&d, "name LIKE ?", dataEntities[i].Name)
// 			db.Model(&d).Update("Issct", issueCountAsInteger)
// 		}
// 	}

// 	return c
// }

// // Removes empty entries from data_entities table
// func removeEmptyEntriesFromDataEntities(c *gin.Context) *gin.Context {
// 	// Initialize GORM
// 	db := c.MustGet("db").(*gorm.DB)

// 	// fetch all the repositories to array (Similar to GET all)
// 	var dataEntities []models.Entity
// 	db.Find(&dataEntities)

// 	length := len(dataEntities)

// 	for i := 0; i < length; i++ {
// 		if dataEntities[i].Issct == 0 {
// 			var dataEntity models.Entity
// 			db.Where("id = ?", dataEntities[i].Id).First(&dataEntity)
// 			db.Delete(&dataEntity)
// 		}
// 	}

// 	for i := 0; i < length; i++ {
// 		if dataEntities[i].Liblin == 0 {
// 			var dataEntity models.Entity
// 			db.Where("id = ?", dataEntities[i].Id).First(&dataEntity)
// 			db.Delete(&dataEntity)
// 		}
// 	}

// 	return c
// }

// // Removes empty entries from libraries table
// func removeEmptyEntriesFromLibraries(c *gin.Context) *gin.Context {
// 	// Initialize GORM
// 	db := c.MustGet("db").(*gorm.DB)

// 	// fetch all the repositories to array (Similar to GET all)
// 	var libraries []models.Library
// 	db.Find(&libraries)

// 	length := len(libraries)

// 	for i := 0; i < length; i++ {
// 		if libraries[i].Liblin == 0 {
// 			var library models.Library
// 			db.Where("id = ?", libraries[i].Id).First(&library)
// 			db.Delete(&library)
// 		}

// 	}

// 	return c
// }

// // Fetches and Writes the Repositories to <repositories> table
// func FetchRepositories(c *gin.Context) {
// 	c = fetchMetadata(10, 100, c)

// 	c.IndentedJSON(http.StatusOK, "Metadata Succesfully Fetched to Database!")
// }

// // Fetches data from data_entities table and convert that to data.csv / data/data.csv
// func GenerateCsvFile(c *gin.Context) {

// 	// Initialize GORM
// 	db := c.MustGet("db").(*gorm.DB)

// 	// fetch all the repositories to array (Similar to GET all)
// 	var dataEntities []models.Entity
// 	db.Find(&dataEntities)

// 	sourceData := [][]string{
// 		{"name", "uri", "original_code_lines", "library_code_lines", "issue_count"},
// 	}

// 	length := len(dataEntities)

// 	for i := 0; i < length; i++ {
// 		var stringSlice []string
// 		stringSlice = append(stringSlice, dataEntities[i].Name)
// 		stringSlice = append(stringSlice, dataEntities[i].Uri)
// 		stringSlice = append(stringSlice, strconv.Itoa(dataEntities[i].Lin))
// 		stringSlice = append(stringSlice, strconv.Itoa(dataEntities[i].Liblin))
// 		stringSlice = append(stringSlice, strconv.Itoa(dataEntities[i].Issct))
// 		sourceData = append(sourceData, stringSlice[:])
// 	}

// 	csvFile, err := os.Create("data/data.csv")

// 	if err != nil {
// 		log.Fatalf("failed creating file: %s", err)
// 	}

// 	csvWriter := csv.NewWriter(csvFile)

// 	for _, row := range sourceData {
// 		_ = csvWriter.Write(row)
// 	}

// 	csvWriter.Flush()
// 	csvFile.Close()

// 	c.IndentedJSON(http.StatusOK, "Converted Generated Data to CSV!")
// }

// // Enriches the data with lines of code and library lines.
// func EnrichMetadata(c *gin.Context) {
// 	c = countLines(c)
// 	c = fetchIssues(c)

// 	c.IndentedJSON(http.StatusOK, "Metadata Succesfully Enriched!")
// }
