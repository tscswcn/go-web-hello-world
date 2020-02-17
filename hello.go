package main

import (
"fmt"
"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Web Hello World!")
		//fmt.print("Go Web Hello World!")
		//fmt.Print("Go Web Hello World!")
	})

	http.ListenAndServe(":8081", nil)
}

