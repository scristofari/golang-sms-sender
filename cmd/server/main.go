package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/scristofari/sms-sender/twilio"
)

func main() {
	http.HandleFunc("/send-sms", sendHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("Method %s not allowed !", r.Method), http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	if err := twilio.SendSMS(query.Get("to"), query.Get("content")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Done !!")
}
