package main

import (
	"fmt"
	Router "golang-server/router"
	"log"
	"net/http"
)

func main() {
    r := Router.Router()
    log.Fatal(http.ListenAndServe(":4000", r))
    fmt.Println("Server is started at port 4000.")
    
}