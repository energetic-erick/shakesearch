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
	"sort"
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

func groupIdxs(idxs []int, mergeTolerance int) [][]int {
	rtn := [][]int{}
	sort.Ints(idxs)
	lastIdx := -1 - mergeTolerance
	curr := []int{}
	for _, idx := range idxs {
		if lastIdx+mergeTolerance < idx {
			// start anew
			if len(curr) > 0 {
				rtn = append(rtn, curr)
				curr = []int{}
			}
		}
		curr = append(curr, idx)
		lastIdx = idx
	}
	if len(curr) > 0 {
		rtn = append(rtn, curr)
	}
	return rtn
}

func (s *Searcher) formatGroup(group []int, queryLen int, lookaround int) string {
	if len(group) == 0 {
		return ""
	}

	rtn := ""
	lastIndex := max(0, group[0]-lookaround)
	for _, idx := range group {
		lookbehind := s.CompleteWorks[lastIndex:idx]
		word := s.CompleteWorks[idx : idx+queryLen]
		rtn += fmt.Sprintf("%s<strong>%s</strong>", lookbehind, word)
		lastIndex = idx + queryLen
	}
	lookaheadLimit := min(len(s.CompleteWorks), lastIndex+lookaround)
	rtn += s.CompleteWorks[lastIndex:lookaheadLimit]

	// Note: there are no lines longer than 100 characters in the corpus, so the 250 character slice will work fine here (will always have some leading text)
	trim, _ := regexp.Compile(`^[^\n]*\n*|\n*[^\n]*$`)
	rtn = trim.ReplaceAllString(rtn, "")
	return rtn
}

// Search the complete works for a given query
func (s *Searcher) Search(query string) []string {
	queryBytes := []byte(strings.ToLower(query))
	idxs := s.SuffixArray.Lookup(queryBytes, -1)
	results := []string{}

	lookaround := 250
	mergeTolerance := lookaround*2 + len(query)
	idxsGrouped := groupIdxs(idxs, mergeTolerance)

	fmt.Println(idxsGrouped)

	for _, group := range idxsGrouped {
		curr := s.formatGroup(group, len(query), lookaround)
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
