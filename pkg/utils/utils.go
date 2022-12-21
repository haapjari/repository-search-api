package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/haapjari/glass/pkg/models"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func RemoveElement(s []models.Repository, index int) []models.Repository {
	ret := make([]models.Repository, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

// Execute in bash shell
func ExecuteInShell(app string, arg string, target string) []byte {
	cmd := exec.Command(app, arg, target)
	out, err := cmd.Output()
	CheckError(err)
	return out
}

func WriteToFile(data any, fileName string) {
	file, err := json.MarshalIndent(data, "", " ")
	CheckErr(err)

	err = ioutil.WriteFile(fileName+".json", file, 0644)
	CheckErr(err)
}

// Analyze libraries of packageJsonAsString and write them into repositoryNames element
//
//	in data_entities and libraries tables.
// func RunLibraryAnalysis(packageJsonAsString string, db *gorm.DB, repositoryName string, startOfObject string) {

// 	// Split the string to Slice from "devDependencies".
// 	strSlice := strings.Split(packageJsonAsString, startOfObject)

// 	// Split the right half of the Slice to new slice from every new line character.

// 	if len(strSlice) > 1 {
// 		strSlice = strings.Split(strSlice[1], "\n")
// 	}

// 	// do we have right slice here
// 	fmt.Println(strSlice)

// 	totalLines := 0

// 	// Loop through the right half of the slice
// 	for i := 0; i < len(strSlice); i++ {

// 		// remove whitespace
// 		packageJsonRow := strings.TrimSpace(strSlice[i])

// 		// if the element equals "}," - break out from the loop
// 		if packageJsonRow == "}," || packageJsonRow == "}" {
// 			break
// 		}

// 		// parse string into format of library@version
// 		packageJsonRow = strings.ReplaceAll(packageJsonRow, `"`, "")
// 		packageJsonRow = strings.ReplaceAll(packageJsonRow, `:`, "")
// 		packageJsonRow = strings.ReplaceAll(packageJsonRow, `^`, "")
// 		packageJsonRow = strings.ReplaceAll(packageJsonRow, `,`, "")

// 		libraryName := strings.ReplaceAll(packageJsonRow, `.`, "")
// 		libraryName = strings.ReplaceAll(libraryName, " ", "")
// 		libraryName = strings.ReplaceAll(libraryName, "0", "")
// 		libraryName = strings.ReplaceAll(libraryName, "1", "")
// 		libraryName = strings.ReplaceAll(libraryName, "2", "")
// 		libraryName = strings.ReplaceAll(libraryName, "3", "")
// 		libraryName = strings.ReplaceAll(libraryName, "4", "")
// 		libraryName = strings.ReplaceAll(libraryName, "5", "")
// 		libraryName = strings.ReplaceAll(libraryName, "6", "")
// 		libraryName = strings.ReplaceAll(libraryName, "7", "")
// 		libraryName = strings.ReplaceAll(libraryName, "8", "")
// 		libraryName = strings.ReplaceAll(libraryName, "9", "")

// 		// libraryName variable contains only the library name

// 		// slice might contain spaces, if element is not whitespace go into if clause
// 		if libraryName != "" {

// 			// write libraries into database

// 			// search, if library is already on the database, otherwise create it
// 			var lib models.Library
// 			db.First(&lib, "name LIKE ?", libraryName)

// 			if lib.Name == "" {

// 				// Create Models
// 				var library models.Library = models.Library{Name: libraryName}

// 				// Write to Databases
// 				db.Create(&library)
// 			}
// 		}

// 		packageJsonRow = strings.ReplaceAll(packageJsonRow, " ", "@")
// 		packageJsonRow = strings.ReplaceAll(packageJsonRow, " ", "")

// 		// Enter Analysis Phase

// 		if packageJsonRow != "" {

// 			var lib models.Library
// 			db.First(&lib, "name LIKE ?", libraryName)

// 			// if the library is not analyzed, gocloc it
// 			if lib.Liblin == 0 {
// 				ExecuteInShell("npm", "install", packageJsonRow)

// 				// output is a bytestream variable
// 				output := ExecuteInShell("gocloc", "--output-type=json", "node_modules/")

// 				// initialize empty struct
// 				var element models.CodeAnalysisStruct

// 				// parse JSON into empty struct
// 				json.Unmarshal(output, &element)

// 				// We don't create new model, we need to patch old model
// 				var libEntity models.Library
// 				db.First(&libEntity, "name LIKE ?", libraryName)
// 				db.Model(&libEntity).Update("Liblin", element.Total.Code)

// 				output = ExecuteInShell("rm", "-rf", "node_modules")
// 				output = ExecuteInShell("rm", "-rf", "package.json")
// 				output = ExecuteInShell("rm", "-rf", "package-lock.json")
// 			}
// 		}

// 		// Calculate Total

// 		// Fetch Library Line Data from Database
// 		var libraryStruct models.Library
// 		db.First(&libraryStruct, "name LIKE ?", libraryName)

// 		var dataStruct models.Entity
// 		db.First(&dataStruct, "name LIKE ?", repositoryName)

// 		totalLines = libraryStruct.Liblin + dataStruct.Liblin

// 		var entity models.Entity
// 		db.First(&entity, "name LIKE ?", repositoryName)

// 		db.Model(&entity).Update("Liblin", totalLines)
// 	}
// }

// // Analyze repository and write them into repositoryNames element in data_entities table
// func RunMainAnalysis(repositoryName string, repositoryUri string, db *gorm.DB) {

// 	// Output is a ByteStream value
// 	output := ExecuteInShell("gocloc", "--output-type=json", repositoryName)

// 	// Parse JSON into empty struct
// 	var element models.CodeAnalysisStruct
// 	json.Unmarshal(output, &element)

// 	// WRITE TO DATABASE: `name` and `liblin` to `libraries` - table.
// 	// WRITE TO DATABASE: `name` and `lin` to `data` - table

// 	// Create Models
// 	var library models.Library = models.Library{Name: repositoryName, Liblin: element.Total.Code, Uri: repositoryUri}
// 	var dataEntity models.Entity = models.Entity{Name: repositoryName, Lin: element.Total.Code, Uri: repositoryUri}

// 	// Write to Databases
// 	db.Create(&dataEntity)
// 	db.Create(&library)
// }

// // Check if library is already analyzed
// func CheckIfAnalyzed(repositoryName string, packageType string, c *gin.Context) (*gin.Context, bool) {

// 	db := c.MustGet("db").(*gorm.DB)

// 	if packageType == "library" {
// 		fmt.Println("Library")
// 	}

// 	if packageType == "main" {

// 		// query database for name field
// 		var dataEntity models.Entity
// 		db.First(&dataEntity, "name LIKE ?", repositoryName)

// 		// if the entity has value in field "Lin", return true
// 		if dataEntity.Lin != 0 {
// 			return c, true
// 		}
// 	}

// 	return c, false
// }

func Remove(slice []string, i int) []string {

	// Remove the element at index i from a.
	copy(slice[i:], slice[i+1:]) // Shift a[i+1:] left one index.
	slice[len(slice)-1] = ""     // Erase last element (write zero value).
	slice = slice[:len(slice)-1] // Truncate slice.

	return slice
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
