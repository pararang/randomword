package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	host          = "https://random-indonesian-word.p.rapidapi.com"
	randomWordURL = host + "/words/random"
	rapidAPIHost  = "random-indonesian-word.p.rapidapi.com"
	rapidAPIKey   = "e9b079b792msh3e1aa8608df520dp15d272jsn0831149d10d6"
	serverPort    = ":8080"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type RandomWordResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type WordResponse struct {
	Words []string `json:"words"`
}

func getRandomWords(limit int) ([]string, error) {
	url := randomWordURL
	if limit > 0 {
		url += "?limit=" + strconv.Itoa(limit)
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Host", rapidAPIHost)
	req.Header.Add("X-RapidAPI-Key", rapidAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	var randomWordResp RandomWordResponse
	err = json.NewDecoder(resp.Body).Decode(&randomWordResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing data: %v", err)
	}

	if randomWordResp.Status != "success" {
		return nil, fmt.Errorf("error: %s", randomWordResp.Message)
	}

	return randomWordResp.Data, nil
}

func handleWords(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "invalid limit parameter", http.StatusBadRequest)
		return
	}

	words, err := getRandomWords(limit)
	if err != nil {
		errResp := ErrorResponse{Message: err.Error()}
		jsonData, _ := json.Marshal(errResp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonData)
		return
	}

	wordResp := WordResponse{Words: words}
	jsonData, _ := json.Marshal(wordResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/words", handleWords)

	fmt.Println("Listening on http://localhost" + serverPort)
	err := http.ListenAndServe(serverPort, nil)
	if err != nil {
		panic(err)
	}
}
