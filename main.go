package main

import (
	middleware "UserMicroservice/middlewares"
	router "UserMicroservice/routers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello World")

	log.Fatal(http.ListenAndServe(":8000", middleware.RemoveTrailingSlash(router.Router())))

}
