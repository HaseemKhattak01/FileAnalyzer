package controllers

import (
	"FileReader/Jwt"
	"FileReader/database"
	"FileReader/models"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func FileReader(g *gin.Context) {

	defer timer("main")()
	numstring := g.Query("num")
	if numstring == "" {
		g.JSON(http.StatusBadRequest, "Missing 'num' query parameter")
		return
	}

	num, err := strconv.Atoi(numstring)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	file, _, err := g.Request.FormFile("file")
	if err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	content, err := io.ReadAll(file)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	text := string(content)
	size := len(text)
	chunkSize := size / num
	channel := make(chan models.Results)
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

	finalResults := models.Results{}
	for res := range channel {
		finalResults.Vowels += res.Vowels
		finalResults.Spaces += res.Spaces
		finalResults.Capitalletters += res.Capitalletters
		finalResults.Smallletters += res.Smallletters
		finalResults.Words += res.Words
	}

	err = database.InsertData(finalResults.Vowels, finalResults.Spaces, finalResults.Capitalletters, finalResults.Smallletters, finalResults.Words)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	g.JSON(http.StatusOK, models.Response{
		Message: "file read successfully",
		Data:    finalResults,
		Status:  http.StatusOK,
	})
}

func Count(str string, channel chan models.Results) {

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

	res := models.Results{
		Vowels: VowelCount, Spaces: SpaceCount, Capitalletters: CapitalCount, Smallletters: SmallCount, Words: WordsCount,
	}

	channel <- res

}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func GetAll(g *gin.Context) {

	data, err := database.Getdata()
	if err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	g.JSON(http.StatusOK, models.Response{
		Message: "Get all data successfully",
		Data:    data,
		Status:  http.StatusOK,
	})
}

func DeleteData(g *gin.Context) {

	idString := g.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}
	del, err := database.DeleteRecords(id)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	g.JSON(http.StatusOK, models.Response{
		Message: "Deleted data successfully",
		Data:    del,
		Status:  http.StatusOK,
	})
}

func UpdateData(g *gin.Context) {

	up := models.UpdateField{}
	if err := g.ShouldBindJSON(&up); err != nil {
		g.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	database.Update(up)

	g.JSON(http.StatusOK, models.Response{
		Message: "Updates data successfully",
		Status:  http.StatusOK,
	})
}

func SignUp(g *gin.Context) {

	var input models.Identity
	fmt.Println(input)
	if err := g.BindJSON(&input); err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	exists, err := database.SignUp_db(input)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if exists {
		g.JSON(http.StatusOK, models.Response{
			Message: "User signed up successfully",
			Status:  http.StatusConflict,
		})
		return
	}
}

func LogIn(g *gin.Context) {

	var input models.Identify
	if err := g.ShouldBindJSON(&input); err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	isvalid, err := database.LogIn_db(input)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	if !isvalid {
		g.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
		return
	}

	refreshToken, err := Jwt.CreateRefreshToken(input.Username)
	if err != nil {
		fmt.Println("Error creating access token:", models.ErrorResponse{Error: err.Error()})
		return
	}

	response := models.Response{
		Data: gin.H{
			"refresh_token": refreshToken,
		},
		Message: "Authentication successful",
		Status:  http.StatusOK,
	}

	g.JSON(http.StatusCreated, response)

}

func Refresh(g *gin.Context) {

	refreshToken := g.Request.Header.Get("Authorization")
	_, err := Jwt.RefreshTokenValidity(refreshToken)
	if err != nil {
		g.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
		return
	}

	accessToken, err := Jwt.CreateAccessToken(models.RefreshData.RefreshToken)
	if err != nil {
		g.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	response := models.Response{
		Data: gin.H{
			"access_token": accessToken,
		},
		Message: "Authentication successful",
		Status:  http.StatusOK,
	}

	g.JSON(http.StatusCreated, response)

}

func JokeHandler(g *gin.Context) {

	joke, err := database.GetJoke()
	if err != nil {
		g.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"joke": joke})

}

func HealthHandler(g *gin.Context) {
	fmt.Println("I am here in the health function!")
	g.Status(http.StatusNoContent)
}

func ReadinessHandler(g *gin.Context) {
	fmt.Println("i am here in readiness function!")
	g.Status(http.StatusOK)
}

func DBReadinessHandler(g *gin.Context) {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		g.String(http.StatusInternalServerError, fmt.Sprintf("Error connecting to database: %v", models.ErrorResponse{Error: err.Error()}))
	}

	defer dbConn.Close()

	err = dbConn.Ping()
	if err != nil {
		g.String(http.StatusInternalServerError, fmt.Sprintf("Database is not ready: %v", models.ErrorResponse{Error: err.Error()}))
	}

	g.String(http.StatusOK, "Database is ready")

}
