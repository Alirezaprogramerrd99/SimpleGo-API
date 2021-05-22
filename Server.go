package main


import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

type BodyRequest struct {
	Task    string `json:"task"`
	Numbers []int  `json:"numbers"`
}
type POSTBodyResponse struct {
	Task              string  `json:"task"`
	Numbers           []int   `json:"numbers"`
	Answer            float64 `json:"answer ,omitempty"`
	SortedSliceAnswer []int   `json:"answer,omitempty"`
	StatusCode        int     `json:"code"`
	Massage           string  `json:"massage"`
}

type HistoryRecord struct {
	Task              string  `json:"task"`
	Numbers           []int   `json:"numbers"`
	Answer            float64 `json:"answer ,omitempty"`
	SortedSliceAnswer []int   `json:"answer,omitempty"`
}

var requestRecords []HistoryRecord

type GETResponse struct { // we will send this struct in json format to client.
	Size       int             `json:"size"`
	History    []HistoryRecord `json:"history"`
	StatusCode int             `json:"code"`
	Massage    string          `json:"massage"`
}

func meanCalculate(numbersArr []int) float64 {
	sum := 0
	for i := 0; i < len(numbersArr); i++ {
		sum += numbersArr[i]
	}
	return float64(sum) / float64(len(numbersArr))
}

func historyRequestHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"code": "405"`))
		w.Write([]byte(`
		"message": "Method is not supported."}`))
		return
	}

	w.Header().Set("Content-Type", "application/json") // setting header.
	var byteSliceJSON []byte
	w.WriteHeader(http.StatusOK)

	byteSliceJSON, _ = json.Marshal(GETResponse{
		Size:       len(requestRecords),
		History:    requestRecords,
		StatusCode: 200,
		Massage:    "History sent successfully!",
	})
	w.Write(byteSliceJSON)
}

func calculatorRequestHandler(w http.ResponseWriter, r *http.Request) {

	status := 200
	msg := "Task done successfully!"

	if r.Method != "POST" {

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"code": "405"`))
		w.Write([]byte(`
	"message": "Method is not supported."}`))
		return
	}

	var byteSliceJSON []byte
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	POSTBody, err := ioutil.ReadAll(r.Body) // pass IO closer object for reading the requets data.

	if err != nil {
		panic(err)
	}

	postReqBody := string(POSTBody)
	splitBody := strings.Split(postReqBody, ":")

	if len(splitBody) > 3 {

		w.WriteHeader(http.StatusNotAcceptable)
		status = 406
		msg = "Request not acceptable."
	}

	fmt.Println("from client: " + postReqBody) // print the body of client post req.
	reqBody := new(BodyRequest)
	err2 := json.Unmarshal(POSTBody, reqBody) // set json fields to struct fields.

	if err2 != nil {

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code": "404"`))
		w.Write([]byte(`
	"message": "Invalid Types."}`))
		return
	}
	// if request was valid must check before calling this function.
	var Avg float64
	var sortedAns []int // sorted list answer from server.

	if status == 200 {

		switch reqBody.Task {
		case "mean":
			Avg = meanCalculate(reqBody.Numbers)
		case "sort":
			var tempSl = make([]int, len(reqBody.Numbers))
			copy(tempSl, reqBody.Numbers)
			sort.Ints(tempSl)
			sortedAns = tempSl

		default:
			msg = "Task not supported."
			status = http.StatusNotFound
		}
		// store the request in history Records.
		if status == 200 {

			requestRecords = append(requestRecords, HistoryRecord{
				Task:              reqBody.Task,
				Numbers:           reqBody.Numbers,
				Answer:            Avg,
				SortedSliceAnswer: sortedAns,
			})
		}
	}

	byteSliceJSON, _ = json.Marshal(POSTBodyResponse{
		Task:              reqBody.Task,
		Numbers:           reqBody.Numbers,
		Answer:            Avg,
		SortedSliceAnswer: sortedAns,
		StatusCode:        status,
		Massage:           msg,
	})

	w.Write(byteSliceJSON)
}

func ErrorRequestHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"code": "404"`))
	w.Write([]byte(`
"message": "Not found."}`))
}

func main() {

	const PORT string = ":8080"
	requestRecords = []HistoryRecord{}
	log.Printf("server is runinng on port " + PORT)
	http.HandleFunc("/calculator", calculatorRequestHandler)
	http.HandleFunc("/history", historyRequestHandler)
	http.HandleFunc("/", ErrorRequestHandler)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
