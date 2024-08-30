package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	httpListenPort := flag.String("httpListenPort", ":3000", "listening port for HTTP")
	grpcListenPort := flag.String("grpcListenPort", ":3001", "listening portfor GRPC")
	flag.Parse()
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLoggerMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(*grpcListenPort, svc))
	}()
	// time.Sleep(time.Second * 2)
	// c, err := client.NewGRPCClient(*grpcListenPort)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err = c.Aggregate(context.Background(), &types.AggregateRequest{
	// 	ObuID: 1,
	// 	Value: 15.15,
	// 	Unix:  time.Now().UnixNano(),
	// }); err != nil {
	// 	log.Fatal(err)
	// }
	log.Fatal(makeHttpTransport(*httpListenPort, svc))
}

func makeHttpTransport(listenPort string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port:", listenPort)
	http.HandleFunc("/aggregate", handleAggragate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(listenPort, nil)
}

func makeGRPCTransport(listenPort string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port:", listenPort)

	// Make TCP Listener
	ln, err := net.Listen("tcp", listenPort)
	if err != nil {
		return err
	}
	defer ln.Close()

	// Make GRPC server
	server := grpc.NewServer([]grpc.ServerOption{}...)

	// Registering grpc to the grpc package
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"err": "No obu in request"})
			return
		}

		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"err": "malformatted obu id"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"err": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggragate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"err": err.Error()})
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"err": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
