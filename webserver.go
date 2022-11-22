//go:build ignore

package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	fmt.Fprint(w, "Hi there, I love my angel Aimee!")
}

func main() {
	http.HandleFunc("/	", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
