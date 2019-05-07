package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func parseSize(parameter string) (uint64, error) {
	dim := uint64(1)
	switch parameter[len(parameter)-1:] {
	case "k", "K":
		dim = 1024
		parameter = parameter[:len(parameter)-1]
	case "m", "M":
		dim = 1024 * 1024
		parameter = parameter[:len(parameter)-1]
	}
	size, err := strconv.ParseUint(parameter, 10, 32)
	if err != nil {
		return 0, err
	}
	return dim * size, nil
}

func writeBody(size uint64, w http.ResponseWriter) error {
	in, err := os.Open("/dev/urandom")
	if err != nil {
		return err
	}
	defer in.Close()
	buf := make([]byte, 1024)
	for size > 1024 {
		size -= 1024
		if _, err = w.Write(buf); err != nil {
			return err
		}
	}
	_, err = w.Write(buf[:size])
	return err
}

func handleLoad(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		value := r.URL.Query().Get("s")
		if len(value) == 0 {
			return
		}
		size, err := parseSize(value)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Add("Content-Type", "application/binary")
		if err := writeBody(size, w); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return

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
