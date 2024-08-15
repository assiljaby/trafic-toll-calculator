package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/assiljaby/trafic-toll-calculator/types"
)

func main() {
	listenPort := flag.String("listenPort", ":3000", "listening port")
	flag.Parse()
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)

	makeHttpTransport(*listenPort, svc)
}

func makeHttpTransport(listenPort string, svc Aggregator) {
	fmt.Println("HTTP transport running on port:", listenPort)
	http.HandleFunc("/aggregate", handleAggragate(svc))
	http.ListenAndServe(listenPort, nil)
}

func handleAggragate(svc Aggregator) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}