package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil {
		p := &HomePage{Name: "Caveman"}
		t, e := template.ParseFiles("./templates/home.html")
		if e != nil {
			log.Printf("parsing templates home error: %v", e)
		}

		t.Execute(w, p)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}

}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fname := r.FormValue("username")
	var p *UserPage
	if len(cname.Value) != 0 {
		p = &UserPage{cname.Value}
	} else if len(fname) != 0 {
		p = &UserPage{fname}
	}

	t, e := template.ParseFiles("/templates/userhome.html")
	if e != nil {
		log.Printf("Parsing userhome template error: %s", e)
		return
	}

	t.Execute(w, p)

}
