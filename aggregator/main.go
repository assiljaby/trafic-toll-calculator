package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/assiljaby/trafic-toll-calculator/types"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		store          = makeStore()
		svc            = NewInvoiceAggregator(store)
		grpcListenPort = os.Getenv("AGG_GRPC_ENDPOINT")
		httpListenPort = os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLoggerMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(grpcListenPort, svc))
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
	log.Fatal(makeHttpTransport(httpListenPort, svc))
}

func makeHttpTransport(listenPort string, svc Aggregator) error {
	var (
		aggMetricHandler = newHTTPMetricsHandler("aggregate")
		invMetricHandler = newHTTPMetricsHandler("invoice")
		aggregateHandler = makeHTTPHandlerFunc(aggMetricHandler.instrument(handleAggregate(svc)))
		invoiceHandler   = makeHTTPHandlerFunc(invMetricHandler.instrument(handleGetInvoice(svc)))
	)
	http.HandleFunc("/invoice", invoiceHandler)
	http.HandleFunc("/aggregate", aggregateHandler)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("HTTP transport running on port ", listenPort)
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

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type given %s", storeType)
		return nil
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
