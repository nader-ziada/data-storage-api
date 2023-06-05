package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

/*
	This is a rather naive function for figuring out the
	repository and objectID from the URL.
*/
func urlMatch(url string) (repository string, objectID string) {
	fragments := strings.SplitN(url, "/", -1)
	repository = fragments[2]
	objectID = ""
	if len(fragments) > 3 {
		objectID = fragments[3]
	}
	return repository, objectID
}

func main() {
	// We'll store the data in memory in a map.
	storage := make(map[string][]byte)

	handler := func(w http.ResponseWriter, req *http.Request) {
		repository, objectID := urlMatch(req.URL.Path)
		switch req.Method {
		/*
			Download an Object

			GET /data/{repository}/{objectID}

			Response

			Status: 200 OK
			{object data}
		*/
		case "GET":
			/*
				This implementation of GET is incomplete at this time and won't
				pass the tests, please improve it.
			*/
			_, err := w.Write(storage[objectID])
			if err != nil {
				log.Printf("Error sending object to client: %v\n", err)
				http.Error(w, "Error sending object to client", http.StatusInternalServerError)
				return
			}
		}
		fmt.Println(req.Method + " repository: " + repository + " objectID: " + objectID)
	}

	http.HandleFunc("/data/", handler)
	log.Fatal(http.ListenAndServe(":8282", nil))
}
