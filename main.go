// Copyright © 2017 Solf1re2 <jy1v07@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/solf1re2/config"
	"github.com/solf1re2/gosol/cmd"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	DBHost  = "127.0.0.1"
	DBPort  = ":3306"
	DBUser  = "root"
	DBPass  = "password"
	DBDbase = "cms"
)

var database *sql.DB

type Page struct {
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	GUID       string
}

func main() {
	// Load the config file - PORT,
	cfg := config.LoadConfig("./config.json")

	dbConn := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", DBUser, DBPass, DBHost, DBPort, DBDbase)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("Couldn't connect!")
		log.Println(err.Error)
	}
	database = db

	rtr := mux.NewRouter()
	//rtr.HandleFunc("/pages/{id:[0-9]+}", pageHandler)
	//rtr.HandleFunc("/homepage", pageHandler)
	//rtr.HandleFunc("/contact", pageHandler)
	//rtr.HandleFunc("/page/{id:[0-9]+}", ServePage)

	rtr.HandleFunc("/api/pages", APIPage).
		Methods("GET").
		Schemes("https")
	rtr.HandleFunc("/api/pages/{guid:[0-9a-zA\\-]+}", APIPage).
		Methods("GET").
		Schemes("https")

	rtr.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
	rtr.HandleFunc("/", RedirIndex)
	rtr.HandleFunc("/home", ServeIndex)

	fmt.Printf("Server port :%v\n", cfg.Server.Port)

	http.Handle("/", rtr)

	http.ListenAndServe(":"+cfg.Server.Port, nil)
	cmd.Execute()
}

func RedirIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", 301)
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	var Pages = []Page{}
	pages, err := database.Query("SELECT page_title, page_content,page_date,page_guid FROM pages ORDER BY ? DESC", "page_date")
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	defer pages.Close()
	for pages.Next() {
		thisPage := Page{}
		pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.GUID)
		thisPage.Content = template.HTML(thisPage.RawContent)
		Pages = append(Pages, thisPage)
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, Pages)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageID := vars["id"]
	fileName := "files/" + pageID + ".html"
	_, err := os.Stat(fileName)
	if err != nil {
		fileName = "files/404.html"
	}

	http.ServeFile(w, r, fileName)
}

func APIPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	fmt.Println(pageGUID)
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT page_title, page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get the page: " + pageGUID)
		log.Println(err.Error)
		return
	}
	//html := `<html><head><title>` + thisPage.Title + `</title></head><body><h1>` + thisPage.Title + `</h1><div>` + thisPage.Content + `</div></body></html>`
	//fmt.Fprintln(w, html)
	t, _ := template.ParseFiles("templates/blog.html")
	t.Execute(w, thisPage)
}

func (p Page) TruncatedText() string {
	chars := 0
	for i, _ := range p.RawContent {
		chars++
		if chars > 150 {
			return p.RawContent[:i] + `...`
		}
	}
	return p.RawContent
}
