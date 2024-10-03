package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	swagger "oldcow/go"
	"os"
	"time"
)

const (
	RootCertificatePath string = "../minica.pem"
	ClientCertPath      string = "../client/cert.pem"
	ClientKeyPath       string = "../client/key.pem"
)

type Options struct {
	CT string
	CC int
	M  int
}

func main() {
	var opt Options
	flag.StringVar(&opt.CT, "k", "", "candy type")
	flag.IntVar(&opt.CC, "c", 0, "candy count")
	flag.IntVar(&opt.M, "m", 0, "money")
	flag.Parse()
	var answer swagger.BuyCandyBody
	var err error
	var r *http.Request
	var data string
	var buf bytes.Buffer

	answer.Money = int32(opt.M)
	answer.CandyCount = int32(opt.CC)
	answer.CandyType = opt.CT
	err = json.NewEncoder(&buf).Encode(answer)
	if err != nil {
		log.Fatal(err)
	}

	rootCaPool := x509.NewCertPool()
	rootCA, err := os.ReadFile(RootCertificatePath)
	if err != nil {
		log.Fatal(err)
	}
	rootCaPool.AppendCertsFromPEM(rootCA)
	log.Println("RootCA loaded")
	c := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{RootCAs: rootCaPool,
				GetClientCertificate: func(info *tls.CertificateRequestInfo) (certificate *tls.Certificate, e error) {
					c, err := tls.LoadX509KeyPair(ClientCertPath, ClientKeyPath)
					if err != nil {
						fmt.Printf("error loading client key pair %v\n", err)
						return nil, err
					}
					return &c, nil
				},
			},
		},
	}

	if r, err = http.NewRequest(http.MethodPost, "https://candy.tld:3333/buy_candy", &buf); err != nil {
		log.Fatal(err)
	}

	if data, err = callServer(c, r); err != nil {
		log.Fatal(err)
	}
	log.Println(data)
}

func callServer(c http.Client, r *http.Request) (string, error) {
	respons, err := c.Do(r)
	if err != nil {
		return "", err
	}
	defer respons.Body.Close()
	data, err := io.ReadAll(respons.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
