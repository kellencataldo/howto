package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		log.Println("get called, cool")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "get called"}`))
	case "PUT":
		log.Println("put called, cool")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "put called"}`))
	case "DELETE":
		log.Println("delete called, cool")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "delete called"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

func main() {

	logfile := flag.String("logfile", "", "Logfile for debugging (default outputs to stdout")
	clientCert := flag.String("clientcert", "", "client's cert file to be used for secure connections")
	serverCert := flag.String("servercert", "", "server's cert file to be use for secure connections")
	serverKey := flag.String("serverkey", "", "server's private key file")
	serverPort := flag.Int("port", 8080, "port the server will listen for connections on")
	flag.Parse()

	if "" != *logfile {
		if f, err := os.OpenFile(*logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); nil != err {
			log.Fatalf("Unable to open log file, error recieved: %s\n", err)
		} else {
			log.SetOutput(f)
			defer f.Close()
		}
	}

	if "" == *clientCert {
		log.Fatalln("No client certificate specified, unable to secure connections")
	} else if "" == *serverCert {
		log.Fatalln("No server certificate specified, unable to secure connections")
	} else if "" == *serverKey {
		log.Fatalln("No server private key specified, unable to secure connections")
	} else if *serverPort > 65535 || *serverPort < 0 { // kellen: think of a better check here, type assert?
		log.Fatalf("Invalid port specified: %d, does not fall within valid port range", *serverPort)
	}

	rootCAs := x509.NewCertPool()
	if certs, err := ioutil.ReadFile(*clientCert); err != nil {
		log.Fatalf("Failed read client certfile %s, error: %w", clientCert, err)
	} else if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Fatalf("Client cert file did not contain a certificate")
	}

	address := fmt.Sprintf(":%d", *serverPort)
	tlsConfig := &tls.Config{
		RootCAs: rootCAs,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", httpHandler)

	srv := &http.Server{
		Addr:      address,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	log.Fatal(srv.ListenAndServeTLS(*serverCert, *serverKey))
}
