package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Searcher is a struct to store supporting data for the server
type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

// Load the dataset from a given file at server setup
func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New([]byte(strings.ToLower(s.CompleteWorks)))
	return nil
}

// Search the complete works for a given query
func (s *Searcher) Search(query string) []string {
	// Note: there are no lines longer than 100 characters in the corpus, so the 250 character slice will work fine here (will always have some leading text)
	leadingTrim, _ := regexp.Compile(`^[^\n]*\n*`)
	trailingTrim, _ := regexp.Compile(`\n*[^\n]*$`)
	queryBytes := []byte(strings.ToLower(query))
	idxs := s.SuffixArray.Lookup(queryBytes, -1)
	results := []string{}
	for _, idx := range idxs {
		lowerLimit := max(0, idx-250)
		upperLimit := min(idx+250, len(s.CompleteWorks))

		lookbehindRaw := s.CompleteWorks[lowerLimit:idx]
		lookbehind := leadingTrim.ReplaceAllString(lookbehindRaw, "")

		word := s.CompleteWorks[idx : idx+len(query)]

		lookaheadRaw := s.CompleteWorks[idx+len(query) : upperLimit]
		lookahead := trailingTrim.ReplaceAllString(lookaheadRaw, "")

		curr := fmt.Sprintf("%s<strong>%s</strong>%s", lookbehind, word, lookahead)
		results = append(results, curr)
	}
	return results
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
