package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Page is one of the wiki pages. A wiki consists of a series of interconnected
// pages, eacho of which has a title and a body.
type Page struct {
	Title string
	Body  []byte
}

// save will save the `Page`'s `Body` to a text file.
// `save` takes as its receiver `p`, a pointer to `Page`. It takes no parameters,
// and returns a value of type `error`.
// `save` returns returns an `error` value because that is the return type of `WriteFile`.
// `save` returns the error value, to let the application handle it should
// anything go wrong while writing the file.
// If all goes well, `Page.save()` will return `nil` (the zero-value for
// pointers, interfaces, and some other types).
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// loadPage constructs the filename from `title`, reads the files contents into
// `body`, and returns a pointer to a `Page` literal constructed with the
// proper title and body values.
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, _ := ioutil.ReadFile(filename)
	return &Page{Title: title, Body: body}, nil
}

// viewHandler extracts the page from `r.URL.Path`, loads the page data into
// HTML, and writes it to `w`, the `http.ResponseWriter`
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
