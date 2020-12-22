# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

You can see a live version of the app at
https://pulley-shakesearch.herokuapp.com/. Try searching for "Hamlet" to display
a set of results.

In it's current state, however, the app is just a rough prototype. The search is
case sensitive, the results are difficult to read, and the search is limited to
exact matches.

## Your Mission

Improve the search backend. Think about the problem from the user's perspective
and prioritize your changes according to what you think is most useful.

## Submission

1. Fork this repository and send us a link to your fork after pushing your changes.
2. Heroku hosting - The project includes a Heroku Procfile and, in its
   current state, can be deployed easily on Heroku's free tier.
3. In your submission, share with us what changes you made and how you would prioritize changes if you had more time.

## Erick's Notes

### Goals of this exercise

- Demonstrate that I can write code
- Demonstrate I can design features with a user in mind
- Demonstrate that I can document what I'm working on (albeit in a .md instead of Jira/Notion/etc.)

### Define the user(s)

- A student doing a research project who needs to find which play/act to cite for a quote
- A Shakespeare fan who wants to find specific across different works

### Brainstorming improvements

- Performance: query run times, server startup times, concurrency, pagination, alternative data structures
- UX: autosuggest, live(ish) results, case-insensitive, wrapping lines, fuzzy match, show where it is, combine nearby results, index by play + act + scene instead of lines/character numbers
- UI: make it prettier (not focus of this exercise), highlight search terms
- Bugs: panic when finding something at the beginning or end (250 characters both ends), check for encoding bugs (looks like there are some fun utf-8 characters in file)
- Tests: no tests currently

```
http: panic serving [::1]:57139: runtime error: slice bounds out of range [-217:]
goroutine 26 [running]:
net/http.(*conn).serve.func1(0xc0000a8280)
        /usr/local/go/src/net/http/server.go:1801 +0x147
panic(0x12b5a60, 0xc000016340)
        /usr/local/go/src/runtime/panic.go:975 +0x3e9
main.(*Searcher).Search(0xc00000e080, 0xc0000162ee, 0x3, 0x1, 0xc000134088, 0xc000132001)
        /Users/efriis/projects/shakesearch/main.go:79 +0x1fc
main.handleSearch.func1(0x132a040, 0xc0001260e0, 0xc000112100)
        /Users/efriis/projects/shakesearch/main.go:51 +0x175
net/http.HandlerFunc.ServeHTTP(0xc000012730, 0x132a040, 0xc0001260e0, 0xc000112100)
        /usr/local/go/src/net/http/server.go:2042 +0x44
net/http.(*ServeMux).ServeHTTP(0x14a6c40, 0x132a040, 0xc0001260e0, 0xc000112100)
        /usr/local/go/src/net/http/server.go:2417 +0x1ad
net/http.serverHandler.ServeHTTP(0xc000126000, 0x132a040, 0xc0001260e0, 0xc000112100)
        /usr/local/go/src/net/http/server.go:2843 +0xa3
net/http.(*conn).serve(0xc0000a8280, 0x132a500, 0xc00007e040)
        /usr/local/go/src/net/http/server.go:1925 +0x8ad
created by net/http.(*Server).Serve
        /usr/local/go/src/net/http/server.go:2969 +0x36c
```

### Weighing improvements

- Bug fix first
- Performance improvements are lower priority. Server startup takes 298ms, and queries are sub-millisecond. If load becomes an issue, this server would be easy to replicate as-is because it's not dependent on a mutable data store.
- Case insensitive because it was suggested :)
- Highlighting search terms next. Can "hack" it in by simply adding `<strong>` tags to the search term
- Fuzzy whitespace matching -- can either do with regex (more flexible) or SuffixArray, where we equivalently save spaces and newlines (higher performance)
- Regex takes ~1ms for smaller queries vs microseconds with the SuffixArray (pretty big performance improvement but not a big deal with this sized dataset)
- Realistically, if we wanted full fuzzy matching, would probably index it with something like elasticsearch (not within my scope for coding challenge)
- Autosuggest could be fun -- would be a lookahead based on word frequency. Not super important user-wise though
- Pagination not important with this dataset
- Indexing by play at least (not necessarily act) is pretty important if you're searching shakespeare
- Testing should be done in future, but it's not as interesting for a recruiting project
- Combining nearby results is important. Searching for something like "Denmark" is a mess at the moment (lots of the same-ish section of text)

## Improvements to implement

[x] Check for encoding bugs + read up on runes vs. chars in golang --> looks like it's actually ok as-is
[x] Fix OOB bug resulting in server panic
[x] Case-insensitive queries (index lowercase, search lowercase)
[x] Highlighting search terms with `<strong>`
[ ] Prevent slicing in middle of runes/words
[ ] Combine nearby results
[ ] Index by work (parsing script to split up by work + show results by work)

## Improvements to propose for future

- Real levenshtein-distance matching with an index + unicode similarities like "ae" <> "æ"
- Tests -- make it tolerant to future changes
- Show the act/scene/line citation information too. Requires expanding on the parser. Important for student doing research
