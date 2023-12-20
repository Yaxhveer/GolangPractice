package main

import (
	"fmt"
	"io"
	"net/http"
)

const URL string = "https://lco.dev";

func main () {
	fmt.Println("Yooo");
	resonse, err := http.Get(URL);
	if (err != nil) {
		panic(err);
	}
	fmt.Println(resonse.Body);

	defer resonse.Body.Close()  // caller must close the connection

	// Reading the body
	databytes, err := io.ReadAll(resonse.Body);
	if (err != nil) {
		panic(err);
	}
	data := string(databytes);
	fmt.Println(data);
}
