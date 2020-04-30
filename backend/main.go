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

func AccessMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		next.ServeHTTP(w, r)
	})

}
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

	if err != nil {
		log.Println("Error parsing json", err)
	}
	//return
	fmt.Println("Restaurants created")
	defer srvHandler.db.Close()
	//http.Handle("/", http.FileServer(http.Dir("./frontend")))

	portalMux := http.NewServeMux()
	portalMux.HandleFunc("/api/rests/all", srvHandler.GetAllPlaces)
	portalMux.HandleFunc("/api/rests/add", srvHandler.AddNewPlace)
	portalMux.HandleFunc("/api/card/buy", srvHandler.BuyNewCard)
	portalMux.HandleFunc("/api/mode/print", srvHandler.PrintDBPlaces)
	portalMux.HandleFunc("/api/mode/getimg", srvHandler.GetRestsImages)
	portalMW := AccessMW(portalMux)
	err = http.ListenAndServe(":8080", portalMW)
	if err != nil {
		log.Println("Error runnibng server", err)
	}
}
