package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	randomWordURL = "https://random-indonesian-word.p.rapidapi.com/words/random?limit=5"
	rapidAPIHost  = "random-indonesian-word.p.rapidapi.com"
	rapidAPIKey   = "e9b079b792msh3e1aa8608df520dp15d272jsn0831149d10d6"
	serverPort    = ":8080"
)

type randomWordResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type newResp struct {
	Words []string `json:"words"`
}

func getRandomWord() ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, randomWordURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("X-RapidAPI-Host", rapidAPIHost)
	req.Header.Set("X-RapidAPI-Key", rapidAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTTP response body: %w", err)
	}

	var data randomWordResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling HTTP response body: %w", err)
	}

	return data.Data, nil
}

func randWords(w http.ResponseWriter, r *http.Request) {
	words, err := getRandomWord()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := newResp{Words: words}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/words", randWords)

	server := &http.Server{
		Addr:    serverPort,
		Handler: mux,
	}

	fmt.Println("Listening on http://localhost" + serverPort)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
