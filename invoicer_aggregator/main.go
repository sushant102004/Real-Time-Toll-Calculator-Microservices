/*
	Purpose of this file:
		Start HTTP server and accept POST requests of calculated data.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sushant102004/Traffic-Toll-Microservice/types"
)

const basePrice = 0.12

func main() {
	httpListenAddr := flag.String("httpAddr", ":3000", "the listen address of http server")

	flag.Parse()

	svc := NewInvoiceAggregator()
	makeHTTP_Transport(*httpListenAddr, svc)
}

func makeHTTP_Transport(listenAddr string, svc *InvoiceAggregator) {
	fmt.Println("HTTP Transport Running on Port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/get-invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(svc *InvoiceAggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.CalculatedDistance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		fmt.Println(distance)

		// This will insert data into memory map.
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
	}
}

func handleGetInvoice(svc *InvoiceAggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusNotAcceptable, map[string]string{
				"error": "please provide a valid OBU ID",
			})
			return
		}

		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}
		_ = obuID
		// inv, err := svc.Get(obuID)
		// if err != nil {
		// 	writeJSON(w, http.StatusBadRequest, map[string]string{
		// 		"error": "OBU ID is invalid.",
		// 	})
		// 	return
		// }

		// inv := types.Invoice{
		// 	OBUID:         obuID,
		// 	TotalDistance: resp.TotalDistance,
		// 	TotalAmount:   resp.TotalAmount * basePrice,
		// }

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "this will work as invoice api",
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(body)
}