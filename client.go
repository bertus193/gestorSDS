package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const CERTIFICATE_PATH = "cert.pem"

func buildClient() *http.Client {
	// Set up our own certificate pool
	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Load our trusted certificate path
	pemData, err := ioutil.ReadFile(CERTIFICATE_PATH)
	if err != nil {
		panic(err)
	}
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(pemData)
	if !ok {
		panic("Couldn't load PEM data")
	}

	return client
}

func main() {
	client := buildClient()
	// two different ways to execute an http.GET
	response, err := client.Get("https://127.0.0.1:8021/hello")
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(contents))
	}
}
