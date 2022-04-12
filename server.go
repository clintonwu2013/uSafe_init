package main

import (
	api "changingtec/usafe/IDEAPIFunction"
	idecrypt "changingtec/usafe/IDECrypt"
	IDET "changingtec/usafe/Utility"
	"crypto/tls"
	"encoding/pem"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var (
	PORT             = os.Getenv("PORT")
	CertFilePath     = os.Getenv("CERT_FILE_PATH")
	CertFilePassword = os.Getenv("CERT_FILE_PASSWORD")
	LogDebug         *log.Logger
)

func main() {

	engine := gin.New()

	normalRouter := engine.Group("/")

	APIList := api.GetAPIs()

	for _, api := range APIList {
		switch strings.ToUpper(api.Method) {
		case "GET":
			normalRouter.GET(api.URLPath, api.HandlerFunc)
		case "POST":
			normalRouter.POST(api.URLPath, api.HandlerFunc)
		case "PUT":
			normalRouter.PUT(api.URLPath, api.HandlerFunc)
		case "DELETE":
			normalRouter.DELETE(api.URLPath, api.HandlerFunc)
		}
	}

	tlsCert, err := getTLSCert()
	if err != nil {
		IDET.DebuggerPrintf("getTLSCert() error, err=" + err.Error())
		panic(err.Error())
	}

	srv := &http.Server{
		Addr:    ":" + PORT,
		Handler: engine,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			MinVersion:   tls.VersionTLS10,
			MaxVersion:   tls.VersionTLS12,
		},
	}

	IDET.DebuggerPrintf("starting uSafe server at :" + PORT)
	log.Fatal(srv.ListenAndServeTLS("", ""))

}

func getTLSCert() (tlsCert tls.Certificate, err error) {
	cert, privk, err := idecrypt.LoadPKCS12FromFile(CertFilePath, CertFilePassword)
	if err != nil {
		IDET.DebuggerPrintf("idecrypt.LoadPKCS12FromFile error, err=" + err.Error())
		panic(err.Error())
	}

	bytes, err := idecrypt.MarshalPKCS8PrivateKey(privk)
	if err != nil {
		IDET.DebuggerPrintf("idecrypt.MarshalPKCS8PrivateKey error, err=" + err.Error())
		panic(err.Error())
	}

	certPEMByte := pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw},
	)
	keyPEMByte := pem.EncodeToMemory(
		&pem.Block{Type: "PRIVATE KEY", Bytes: bytes},
	)
	tlsCert, err = tls.X509KeyPair(certPEMByte, keyPEMByte)
	return
}
