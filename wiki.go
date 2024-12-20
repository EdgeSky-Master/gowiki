package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)
type Page struct {
	Title string
	Body []byte
}
var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
// func getTitle (w http.ResponseWriter, r *http.Request) (string,error){
// 	m := validPath.FindStringSubmatch(r.URL.Path)
// 	if m == nil {
// 		http.NotFound(w,r)
// 		return "", errors.New("invalid Page Title")
// 	}
// 	return m[2], nil 
// }
func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}
func loadPage(title string) (*Page, error){
	filename:= "data/" + title + ".txt"
	body,err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body},nil
}
// func handler(w http.ResponseWriter, r *http.Request){
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }
func convertLinks (title string) string {
	temp := regexp.MustCompile(`\[(\w+)\]`)
	return temp.ReplaceAllStringFunc(title, func(match string) string{
		pageName := match[1 : len(match)-1]
		return fmt.Sprintf(`<a href="/view/%s">%s</a>`,pageName,pageName)
	})
}
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	// t, err := template.ParseFiles(tmpl + ".html")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// err = t.Execute(w,p)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	p.Body = []byte(convertLinks(string(p.Body)))
	err := templates.ExecuteTemplate(w, tmpl+".html",p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func viewHandler (w http.ResponseWriter, r *http.Request, title string){
	p,err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request, title string){
	p,err := loadPage(title)
	if err != nil {
		p= &Page{Title: title}
	} 
	renderTemplate(w, "edit", p)
}
func saveHandler(w http.ResponseWriter, r *http.Request, title string){
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w,r,"/view/"+title, http.StatusFound)

}
func rootHandler(w http.ResponseWriter, r *http.Request){
	http.Redirect(w,r,"/view/FrontPage", http.StatusFound)
}
func makeHandler (fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		// Here we will extract the page title from the Request,
        // and call the provided handler 'fn'
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w,r)
			return 
		}
		fn (w,r,m[2]) 
	}
}
func main () {
	// p1 := &Page {Title: "TestPage", Body: []byte("This is a sample Page.")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string (p2.Body))
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8081", nil))
}