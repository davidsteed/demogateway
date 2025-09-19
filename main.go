package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tpl = template.Must(template.ParseFiles("form.tmpl"))

type FormData struct {
	Action  string  // where the form POSTs to
	Message string  // optional flash/message
	Step    string  // e.g. "1" or "0.01"
	Min     string  // e.g. "0"
	Default string  // prefill the input
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := FormData{
			Action:  "/submit",
			Message: "",
			Step:    "1",
			Min:     "0",
			Default: "",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.ExecuteTemplate(w, "form", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Parse the numeric value from the form.
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		raw := r.Form.Get("value")
		num, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			http.Error(w, "value must be a number", http.StatusBadRequest)
			return
		}

		// TODO: do something with `num` (store, call a service, etc.)
		log.Printf("received: %f", num)

		// Re-render the form with a confirmation message and prefill.
		data := FormData{
			Action:  "/submit",
			Message: "Thanks! Received value: " + raw,
			Step:    "1",
			Min:     "0",
			Default: raw,
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.ExecuteTemplate(w, "form", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
