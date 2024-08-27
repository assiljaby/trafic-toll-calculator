package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenPort := flag.String("listenPort", ":8080", "Listening port")
	flag.Parse()
	http.HandleFunc("/invoice", makeApiFunc(handleGetInvoice))
	logrus.Infof("Listening no port %s", *listenPort)
	log.Fatal(http.ListenAndServe(*listenPort, nil))
}

func handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, map[string]string{"Invoice": "invoices yet"})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeApiFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"err": err.Error()})
		}
	}
}
