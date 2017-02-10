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

	"github.com/solf1re2/gosol/cmd"
	"github.com/solf1re2/gosol/config"
)

func main() {
	config := config.LoadConfig("./config.json")
	fmt.Printf("Server port :%v\n", config.Server.Port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+config.Server.Port, nil)
	cmd.Execute()
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
