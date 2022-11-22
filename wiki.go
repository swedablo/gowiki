package main

// Follow these instructions: https://go.dev/doc/articles/wiki/

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
	filename := p.Title + ".txt"
	// 0600 gives read-write premissions for current user only
	return os.WriteFile(filename, p.Body, 0600)
}

/*
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression
}
*/

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	fmt.Println("*** LoadPage:::")
	fmt.Printf("Filename : %s, body=%s\n", filename, string(body))
	if err != nil {
		fmt.Printf("Returning error in LoadPage: %v\n", err)
		return nil, err
	}
	return &Page{
			Title: title,
			Body:  body,
		},
		nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	//fmt.Printf("renderTemplate, tmpl=%v\n", tmpl)
	/*
		t, err := template.ParseFiles(tmpl + ".html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	//err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//fmt.Printf("renderTemplate DONE!!!!!=%v\n", tmpl)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	/*
		title := r.URL.Path[len("/view/"):]
		fmt.Printf("ViewHandler: title=%v\n", title)
	*/
	/*
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
	*/
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	//fmt.Printf("viewHandler: p.Title=%v, p.Body=%v\n", p.Title, string(p.Body))
	renderTemplate(w, "view", p)

	//t, _ := template.ParseFiles("view.html")
	//t.Execute(w, p)
	/*
		OLD SOLUTION:
		fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	*/
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := r.URL.Path[len("/edit/"):]
	/*
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
	*/
	p, err := loadPage(title)
	//fmt.Printf("editHandler: title=%v\n", title)
	if err != nil {
		p = &Page{Title: title}
	}
	//fmt.Printf("editHandler: p.Title=%v, p.Body=%v\n", p.Title, string(p.Body))
	renderTemplate(w, "edit", p)
	//t, _ := template.ParseFiles("edit.html")
	//t.Execute(w, p)
	/*
		OLD SOLUTION:
		fmt.Fprintf(w, "<h1>Editing %s</h1>"+
			"<form action=\"/save/%s\" method=\"POST\">"+
			"<textarea name=\"body\">%s</textarea><br>"+
			"<input type=\"submit\" value=\"Save\">"+
			"<\form>",
			p.Title, p.Title, p.Body)
	*/
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := r.URL.Path[len("/save/"):]
	/*
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
	*/
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))

	/*
		p1 := &Page{
			Title: "TestPage",
			Body:  []byte("This is a simple Page."),
		}
		p1.save()
		p2, err := loadPage("TestPage")
		if err != nil {
			fmt.Printf("Error, and error occured loading the page: %v\n", err)
		}
		fmt.Println(string(p2.Body))
	*/
}
