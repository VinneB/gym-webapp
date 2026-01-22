package server

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

func StartServer() {
	fmt.Println("Entering start server")
	http.HandleFunc("/test", testHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("test/test.html")
	if err != nil {
		log.Println("Failed to open test.html")
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		fmt.Fprintf(w, line)
	}
	fmt.Fprintf(w, "I love %s\n", r.URL.Path[1:])

}
