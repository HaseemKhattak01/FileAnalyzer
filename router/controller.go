package router

import (
	"FileReader/database"
	"FileReader/models"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Results struct {
	Vowels         int
	Spaces         int
	Capitalletters int
	Smallletters   int
	Words          int
}

type Identity struct {
	User     string
	Email    string
	Password string
}

func FileReader(g *gin.Context) {
	defer timer("main")()
	numstring := g.Query("num")
	if numstring == "" {
		g.JSON(http.StatusBadRequest, "Missing 'num' query parameter")
		return
	}

	num, err := strconv.Atoi(numstring)
	if err != nil {
		g.JSON(http.StatusBadRequest, err)
	}
	file, _, err := g.Request.FormFile("file")
	if err != nil {
		g.JSON(http.StatusBadRequest, err)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		g.JSON(http.StatusBadRequest, err)
	}
	text := string(content)
	size := len(text)
	chunkSize := size / num
	channel := make(chan Results)
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == num-1 {
			end = size
		}
		chunk := content[start:end]

		wg.Add(1)
		go func(chunk []byte) {
			defer wg.Done()
			Count(string(chunk), channel)
		}(chunk)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	finalResults := Results{}
	for res := range channel {
		finalResults.Vowels += res.Vowels
		finalResults.Spaces += res.Spaces
		finalResults.Capitalletters += res.Capitalletters
		finalResults.Smallletters += res.Smallletters
		finalResults.Words += res.Words
	}
	err = database.CreateUser(finalResults.Vowels, finalResults.Spaces, finalResults.Capitalletters, finalResults.Smallletters, finalResults.Words)
	if err != nil {
		g.JSON(http.StatusBadRequest, err)
	}
	g.JSON(http.StatusOK, finalResults)
}

func Count(str string, channel chan Results) {
	VowelCount := 0
	WordsCount := 0
	SpaceCount := 0
	CapitalCount := 0
	SmallCount := 0

	for _, char := range str {
		switch {
		case char == ' ':
			SpaceCount++
		case char >= 'A' && char <= 'Z':
			CapitalCount++
			if char == 'A' || char == 'E' || char == 'I' || char == 'O' || char == 'U' {
				VowelCount++
			}
		case char >= 'a' && char <= 'z':
			SmallCount++
			if char == 'a' || char == 'e' || char == 'i' || char == 'o' || char == 'u' {
				VowelCount++
			}
		}
		WordsCount = SpaceCount + 1
	}
	res := Results{
		VowelCount, SpaceCount, CapitalCount, SmallCount, WordsCount,
	}
	channel <- res

}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func Getall(g *gin.Context) {
	data, err := database.Getdata()
	if err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(data)
	g.JSON(http.StatusOK, data)
}

func Deletedata(g *gin.Context) {
	idString := g.Param("id")
	fmt.Println(idString)
	id, err := strconv.Atoi(idString)
	if err != nil {
		g.JSON(http.StatusBadRequest, err)
	}
	del, err := database.DeleteRecords(id)
	if err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(del)
	g.JSON(http.StatusOK, del)
}

func Updatedata(g *gin.Context) {
	up := models.UpdateField{}
	if err := g.ShouldBindJSON(&up); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	database.Update(up)

	g.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
}

func Signup(g *gin.Context) {
	var input models.Identity
	if err := g.BindJSON(&input); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	exists, err := database.SiginingUp(input)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Failed to check user existence"})
		return
	}
	if exists {
		g.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}
	g.JSON(http.StatusOK, gin.H{"message": "User signed up successfully"})
}

func Login(g *gin.Context) {
	var input models.Identify
	if err := g.ShouldBindJSON(&input); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	isvalid, err := database.LoggingIn(input)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Failed to authenticate user"})
	}
	if !isvalid {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	g.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
