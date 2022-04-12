package idecrypt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"io/ioutil"

	p12chain "software.sslmate.com/src/go-pkcs12"
)

// LoadPKCS12FromFile : pfxFilePath as .pfx file path, return return a x509 certificate and rsa private key when error == nil
func LoadPKCS12FromFile(pfxFilePath, pass string) (*x509.Certificate, *rsa.PrivateKey, error) {
	pf, e := ioutil.ReadFile(pfxFilePath)
	if e != nil {
		fmt.Println("ReadFile:", pfxFilePath, e.Error())
		return nil, nil, e
	}

	return LoadPKCS12(pf, pass)
}

// LoadPKCS12 : load pkcs 12 from byte
func LoadPKCS12(contentByte []byte, pass string) (*x509.Certificate, *rsa.PrivateKey, error) {
	var prik *rsa.PrivateKey
	prikInter, cert, _, e := p12chain.DecodeChain(contentByte, pass)
	if e == nil {
		prik = prikInter.(*rsa.PrivateKey)
	}
	return cert, prik, e
}

type PKCS8Key struct {
	Version             int
	PrivateKeyAlgorithm []asn1.ObjectIdentifier
	PrivateKey          []byte
}

func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) ([]byte, error) {
	var pkey PKCS8Key
	pkey.Version = 0
	pkey.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	pkey.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	pkey.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	return asn1.Marshal(pkey)
}
