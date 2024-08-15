package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/assiljaby/trafic-toll-calculator/types"
)

func main() {
	listenPort := flag.String("listenPort", ":3000", "listening port")
	flag.Parse()
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)

	http.HandleFunc("/aggregate", handleAggragate(svc))
	http.ListenAndServe(*listenPort, nil)
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