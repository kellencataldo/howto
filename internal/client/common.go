package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

type ConnectionInfo struct {
	ServerAddress string
	ServerPort    int
	ClientCert    *os.File
	ClientKey     *os.File
	ServerCert    *os.File
}

var (
	client        *http.Client
	serverAddress string
)

func InitializeConnectionInfo(conn ConnectionInfo) error {

	defer conn.ClientCert.Close()
	defer conn.ClientKey.Close()
	defer conn.ServerCert.Close()

	//@TODO: validate port range here, return error
	serverAddress = fmt.Sprintf("https://%s:%d", conn.ServerAddress, conn.ServerPort)

	// read the server cert and add it to the cert pool
	fileInfo, err := conn.ServerCert.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	serverCertBuffer := make([]byte, fileSize)
	if _, err := conn.ServerCert.Read(serverCertBuffer); err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(serverCertBuffer)

	// read the client cert
	fileInfo, err = conn.ClientCert.Stat()
	if err != nil {
		return err
	}

	fileSize = fileInfo.Size()
	clientCertBuffer := make([]byte, fileSize)
	if _, err := conn.ClientCert.Read(clientCertBuffer); err != nil {
		return err
	}

	// read the client private key
	fileInfo, err = conn.ClientKey.Stat()
	if err != nil {
		return err
	}

	fileSize = fileInfo.Size()
	clientKeyBuffer := make([]byte, fileSize)
	if _, err := conn.ClientKey.Read(clientKeyBuffer); err != nil {
		return err
	}

	certKeyPair, err := tls.X509KeyPair(clientCertBuffer, clientKeyBuffer)
	if err != nil {
		return err
	}

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certKeyPair},
			},
		},
	}

	return nil
}
