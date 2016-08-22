package main

import (
	"net/http"
	"os"
	"text/template"
)

// Render html to the ResponseWriter
func renderTemplate(w http.ResponseWriter, tmpl string, name string, data interface{}) {

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(tmpl)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, nil)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, nil)
		return
	}

	// Parse the template file
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute
	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
