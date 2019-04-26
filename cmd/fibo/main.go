package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/frnksgr/fibo/pkg/fibo"
)

const defaultPort = 8080

const help = `Calculate Fibonacci numbers < 90
Usage:
  /                           this help message
  /fibo?n=<num>               same as /fibo/loop/?n=<num>
  /fibo/loop?n=<num>          calculate fibonacci(n) using a loop
  /fibo/recursive?n=<num>     calculate fibonacci(n) using recursion
  /fibo/srecursive?n=<num>    calculate fibonacci(n) using sequential recursion
  /fibo/microservice?n=<num>  calculate fibonacci(n) using microservices

  /load?s=<size>[k|m]		  download data of size <size> byte | k kilobyte | m megabyte
  /load                       upload data e.g. curl -X PUT --data-binary @<(dd if=/dev/urandom bs=1M count=10)
`

func getEnv(name string, fallback string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		value = fallback
	}
	return value
}

func getParameterNum(r *http.Request) uint {
	value := r.URL.Query().Get("n")
	if len(value) == 0 {
		log.Println("Paramter 'n' missing or has no value")
		return 0
	}
	n, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		log.Println("Parameter 'n' has wrong format")
		log.Println(err)
		return 0
	}
	if n >= 90 {
		log.Println("Number too high")
		n = 0
	}
	return uint(n)
}

func handleLoop(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "%d\n", fibo.Loop(getParameterNum(r)))
}

func handleRecursive(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "%d\n", fibo.Recursive(getParameterNum(r)))
}

func handleRecursiveSequential(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "%d\n", fibo.RecursiveSequential(getParameterNum(r)))
}

func handleMicroservice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "%d\n", fibo.Microservice(getParameterNum(r)))
}

// middleware doing request logging
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Got request: %s %s %s \n", r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func main() {
	address := fmt.Sprintf("0.0.0.0:%s", getEnv("PORT", "8080"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}	
		fmt.Fprint(w, help)
	})

	http.HandleFunc("/fibo", handleLoop)
	http.HandleFunc("/fibo/loop", handleLoop)
	http.HandleFunc("/fibo/srecursive", handleRecursiveSequential)
	http.HandleFunc("/fibo/recursive", handleRecursive)
	http.HandleFunc("/fibo/microservice", handleMicroservice)

	http.HandleFunc("/load", handleLoad)

	fmt.Printf("Starting server on %s\n", address)
	log.Fatal(http.ListenAndServe(address, requestLogger(http.DefaultServeMux)))
}
