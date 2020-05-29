package main

import (
    "net/http"
	"log"
	"os"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"fmt"
	"database/sql"
    "encoding/json"
	_ "github.com/go-sql-driver/mysql"
)

type Call struct {
  from string `json:"from"`
  to string `json:"to"`
  status string `json:"status"`
  direction string `json:"direction"`
  duration string `json:"duration"`
  user_id string `json:"user_id"`
  workspace_id string `json:"workspace_id"`
  started string `json:"started"`
}


var db* sql.DB;

func CreateCall(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  var call Call
  _ = json.NewDecoder(r.Body).Decode(call)
  // perform a db.Query insert
	stmt, err := db.Prepare("INSERT INTO calls (`from`, `to`, `status`, `direction`, `user_id`, `workspace_id`, `started`) VALUES ( ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		fmt.Printf("CreateCall Could not execute query..");
		return 
	}
	defer stmt.Close()
	res, err := stmt.Exec(call.from, call.to, call.status, call.duration, call.user_id, call.workspace_id, call.started)
	if err != nil {
		fmt.Printf("CreateCall Could not execute query..");
		return
	}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              

	callId, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("CreateCall Could not execute query..");
		return
	}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              
    w.Header().Set("X-Call-ID", strconv.FormatInt(callId, 10))
  	json.NewEncoder(w).Encode(&call)
}

func main() {
	fmt.Print("starting Lineblocs API server..");
    r := mux.NewRouter()
    // Routes consist of a path and a handler function.
    r.HandleFunc("/call/createCall", CreateCall);

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	var err error
	//db, err = sql.Open("mysql", "lineblocs:lineblocs@lineblocs.ckehyurhpc6m.ca-central-1.rds.amazonaws.com/lineblocs")
	db, err = sql.Open("mysql", "lineblocs:lineblocs@45.76.62.46/lineblocs")
	if err != nil {
		panic("Could not connect to MySQL");
		return
	}
    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(":8000", loggedRouter))
}