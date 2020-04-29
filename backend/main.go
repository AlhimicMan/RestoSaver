package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var pghost = flag.String("pghost", "127.0.0.1", "PostgreSQL Host")
var pgDB = flag.String("pgdb", "", "PostgreSQL DB name")
var pglogin = flag.String("pglogin", "postgres", "PostgreSQL DB login")
var pgpass = flag.String("pgpass", "", "PostgreSQL DB Password")

func main() {
	flag.Parse()
	if (*pgDB == "") || (*pgpass == "") {
		fmt.Println("Please set PG Databse name and PG password")
		return
	}
	fmt.Println("Server starting...")
	srvHandler, err := NewAPIHandler()
	if err != nil {
		fmt.Println(err)
		return
	}

	//err = ParseJSON(srvHandler.db)
	if err != nil {
		log.Println("Error parsing json", err)
	}
	//return
	fmt.Println("Restaurants created")
	defer srvHandler.db.Close()
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/api/rests/all", srvHandler.GetAllPlaces)
	http.HandleFunc("/api/rests/add", srvHandler.AddNewPlace)
	http.HandleFunc("/api/card/buy", srvHandler.BuyNewCard)
	http.HandleFunc("/api/mode/print", srvHandler.PrintDBPlaces)
	http.HandleFunc("/api/mode/getimg", srvHandler.GetRestsImages)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error runnibng server", err)
	}
}
