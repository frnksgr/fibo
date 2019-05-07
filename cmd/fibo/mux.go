package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/frnksgr/fibo/pkg/fibo"
)

var mainServeMux = initServeMux()

func initServeMux() *http.ServeMux {
	http.HandleFunc("/fibo", handleLoop)
	http.HandleFunc("/fibo/loop", handleLoop)
	http.HandleFunc("/fibo/srecursive", handleRecursiveSequential)
	http.HandleFunc("/fibo/recursive", handleRecursive)
	http.HandleFunc("/fibo/microservice", handleMicroservice)

	http.HandleFunc("/load", handleLoad)

	return http.DefaultServeMux
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
