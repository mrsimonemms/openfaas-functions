package function

import (
	"fmt"
	"net/http"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	// var input []byte

	// if r.Body != nil {
	// 	defer r.Body.Close()

	// 	body, _ := io.ReadAll(r.Body)

	// 	input = body
	// }

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	data, err := GetFromFeed(client)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Body: %+v", data)))
}
