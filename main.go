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
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/solf1re2/config"
	"github.com/solf1re2/gosol/cmd"
)

func main() {
	// Load the config file - PORT,
	cfg := config.LoadConfig("./config.json")

	rtr := mux.NewRouter()
	rtr.HandleFunc("/pages/{id:[0-9]+}", pageHandler)
	//rtr.HandleFunc("/homepage", pageHandler)
	//rtr.HandleFunc("/contact", pageHandler)

	fmt.Printf("Server port :%v\n", cfg.Server.Port)

	http.Handle("/", rtr)

	http.ListenAndServe(":"+cfg.Server.Port, nil)
	cmd.Execute()
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
