package main

import (
	"RestoSave/backend/models"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type APIResponse struct {
	Text   string
	Status int `json:"_"`
}

type APIHandler struct {
	db *gorm.DB
}

func NewAPIHandler() (*APIHandler, error) {
	handler := &APIHandler{}
	DBConnArgs := fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable",
		*pghost,
		*pglogin,
		*pgDB,
		*pgpass)
	db, err := gorm.Open("postgres", DBConnArgs)
	if err != nil {
		fmt.Println("Error connecting to database")
		return nil, err
	}
	handler.db = db
	db.CreateTable(&models.Restaurant{})
	db.CreateTable(&models.GiftCard{})
	return handler, nil
}

func (h *APIHandler) SendResponseState(resp *APIResponse, w http.ResponseWriter) {
	result, _ := json.Marshal(resp)
	w.WriteHeader(resp.Status)
	w.Write(result)
}

// apigen:api {"url": "/rests/all", "auth": false}
func (h *APIHandler) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	var DBRests []models.Restaurant
	h.db.Where("moderated = ?", true).Find(&DBRests)
	results, err := json.Marshal(&DBRests)
	if err != nil {
		NErr := &APIResponse{
			Text:   "Internal error",
			Status: http.StatusInternalServerError,
		}
		h.SendResponseState(NErr, w)
		log.Println(err)
		return
	}
	w.Write(results)
}

// apigen:api {"url": "/rests/add", "auth": false}
func (h *APIHandler) AddNewPlace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		NErr := &APIResponse{
			Text:   "Only POST method allowed",
			Status: http.StatusBadRequest,
		}
		h.SendResponseState(NErr, w)
		return
	}
	NewRestaurant := &models.Restaurant{}
	err := json.NewDecoder(r.Body).Decode(NewRestaurant)
	if err != nil {
		log.Println("/rest/add. Error getting body", err)
	}
	h.db.Create(NewRestaurant)
	if len(h.db.GetErrors()) == 0 {
		result := &APIResponse{
			Text:   "OK",
			Status: http.StatusOK,
		}
		h.SendResponseState(result, w)
	}
}

func (h *APIHandler) BuyNewCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		NErr := &APIResponse{
			Text:   "Only POST method allowed",
			Status: http.StatusBadRequest,
		}
		h.SendResponseState(NErr, w)
		return
	}
	NewCard := &models.GiftCard{}
	err := json.NewDecoder(r.Body).Decode(NewCard)
	if err != nil {
		log.Println("/rest/add. Error getting body", err)
		NErr := &APIResponse{
			Text:   "Error parsing request",
			Status: http.StatusBadRequest,
		}
		h.SendResponseState(NErr, w)
		return
	}
	h.db.Create(NewCard)
	if len(h.db.GetErrors()) == 0 {
		result := &APIResponse{
			Text:   "OK",
			Status: http.StatusOK,
		}
		h.SendResponseState(result, w)
	}
}

func (h *APIHandler) PrintDBPlaces(w http.ResponseWriter, r *http.Request) {
	var DBRests []models.Restaurant
	h.db.Find(&DBRests)
	_, err := json.Marshal(&DBRests)
	if err != nil {
		NErr := &APIResponse{
			Text:   "Internal error",
			Status: http.StatusInternalServerError,
		}
		h.SendResponseState(NErr, w)
		log.Println(err)
		return
	}
	//w.Write(results)

	for _, rest := range DBRests {
		resStr := fmt.Sprintf("%s - %s\n", rest.Name, rest.ImageURL)
		w.Write([]byte(resStr))
	}
}

func (h *APIHandler) GetRestsImages(w http.ResponseWriter, r *http.Request) {

	fullUrlFile := "https://cloud.mail.ru/public/LDfH/8aUrjwFtf/P3170008.jpg"

	// Build fileName from fullPath
	fileName := buildFileName(fullUrlFile)

	// Create blank file
	file := createFile(fileName)

	// Put content on file
	putFile(file, httpClient(), fullUrlFile, fileName)

}

func putFile(file *os.File, client *http.Client, fullUrlFile string, fileName string) {
	resp, err := client.Get(fullUrlFile)

	checkError(err)

	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	checkError(err)

	fmt.Printf("Just Downloaded a file %s with size %d\n", fileName, size)
}

func buildFileName(fullUrlFile string) string {
	fileUrl, err := url.Parse(fullUrlFile)
	checkError(err)

	path := fileUrl.Path
	segments := strings.Split(path, "/")

	return segments[len(segments)-1]
}

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func createFile(fileName string) *os.File {
	file, err := os.Create(fileName)

	checkError(err)
	return file
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
