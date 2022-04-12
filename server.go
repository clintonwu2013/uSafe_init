package main

import (
	api "changingtec/usafe/IDEAPIFunction"
	idecrypt "changingtec/usafe/IDECrypt"
	models "changingtec/usafe/Models"
	IDET "changingtec/usafe/Utility"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

const settingServerFile = "./settingServer.json"

func main() {

	settingServer, err := getServerSetting()
	if err != nil {
		IDET.DebuggerPrintf("getServerSetting() error, err=" + err.Error())
		panic(err.Error())
	}

	PORT := strconv.Itoa(settingServer.Port)
	CertFilePath := settingServer.CertFilePath
	CertFilePassword := settingServer.CertFilePassword

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

	tlsCert, err := getTLSCert(CertFilePath, CertFilePassword)
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

func getTLSCert(CertFilePath, CertFilePassword string) (tlsCert tls.Certificate, err error) {
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

func getServerSetting() (settingServer models.SettingServer, err error) {
	fileName := settingServerFile
	jsonByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonByte, &settingServer)
	if err != nil {
		return
	}
	return
}
