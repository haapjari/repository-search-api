package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/haapjari/glass/pkg/models"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func RemoveCharFromString(str string, chr rune) string {
	var builder strings.Builder

	for _, ch := range str {
		if ch != chr {
			builder.WriteRune(ch)
		}
	}

	return builder.String()
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

func CopyFile(src string, dst string) error {
	out, err := exec.Command("cp", src, dst).CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		return err
	}

	return nil
}

func RemoveFile(path string) error {
	out, err := exec.Command("rm", path).CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		return err
	}

	return nil
}
