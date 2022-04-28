package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func makeHandleVideoUrlPost(videoUrls chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read request body")
			return
		}

		url := string(bytes)
		videoUrls <- url
		w.WriteHeader(http.StatusOK)
	}
}

func listenHttp(port int, videoUrls chan string) {
	handleVideoUrlPost := makeHandleVideoUrlPost(videoUrls)

	r := mux.NewRouter()
	r.HandleFunc("/download", handleVideoUrlPost).Methods("POST")
	http.Handle("/", r)

	addr := getListenAddr(port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Msgf("http server listening at %s", addr)
	if err := http.Serve(l, nil); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func getListenAddr(port int) string {
	if isDevelopment() {
		return fmt.Sprintf("localhost:%d", port)
	} else {
		return fmt.Sprintf(":%d", port)
	}
}
