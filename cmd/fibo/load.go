package main

import (
	"fmt"
	"net/http"
	"errors"
	"strconv"
	"io"
	"io/ioutil"
	"log"
)

func getParameterSize(r *http.Request) (uint64, error) {
	value := r.URL.Query().Get("s")
	if len(value) == 0 {
		return 0, errors.New("Paramter 's' missing or has no value")
	}
	dim := uint64(1)
	switch value[len(value)-1:] {
	case "k", "K":
		dim = 1024
		value = value[:len(value)-1]
	case "m", "M":
		dim = 1024 * 1024
		value = value[:len(value)-1]
	}
	s, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return dim * s, nil
}

func handleLoad(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// TODO: implement
		fmt.Fprintln(w, "not implemented yet")
	case "PUT":
		// read body
		if r.Body == nil {
			fmt.Fprintln(w, "Body empty, no data")
		} else {
			defer r.Body.Close()
			size, err := io.Copy(ioutil.Discard, r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err)
				return
			}
			fmt.Fprintf(w, "Fetched %d size byte\n", size)
		}
	default:		
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
