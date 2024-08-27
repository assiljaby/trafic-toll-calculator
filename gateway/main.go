package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/assiljaby/trafic-toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenPort := flag.String("listenPort", ":8080", "Listening port")
	flag.Parse()
	httpClient := client.NewHTTPClient("http://localhost:3000")
	invoiceHandler := NewInvoiceHandler(httpClient)
	http.HandleFunc("/invoice", makeApiFunc(invoiceHandler.handleGetInvoice))
	logrus.Infof("Listening no port %s", *listenPort)
	log.Fatal(http.ListenAndServe(*listenPort, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	invoice, err := h.client.GetInvoice(context.Background(), 2040347648)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, invoice)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeApiFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri:": r.RequestURI,
			}).Info("REQ :: ")
		}(time.Now())
		if err := fn(w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"err": err.Error()})
		}
	}
}
