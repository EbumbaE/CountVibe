package certificate

import (
	"bytes"
	"os"
	"time"
	"net"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
)

func writeToFile(filename string, buf string){
	file, err := os.Create(filename)
	if (err != nil){
		panic(err)
	}
	file.WriteString(buf)
	file.Close()
}

func generateKeyAndCertificate(keyfile string, certfile string){

	privatekey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	bufkey, _ := x509.MarshalECPrivateKey(privatekey)
	out := &bytes.Buffer{}
	pem.Encode(out,  &pem.Block{Type: "EC PRIVATE KEY", Bytes: bufkey})
	writeToFile(keyfile, out.String())

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2022),
		Subject: pkix.Name{
			Organization:  []string{"Ebumba_E"},
			Country:       []string{"RU"},
			Locality:      []string{"Moscow"},
		},
		IPAddresses:  		   []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	certificate, err := x509.CreateCertificate(rand.Reader, ca, ca, &privatekey.PublicKey, privatekey)
	if err != nil {
		panic(err)
	}

	out.Reset()
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certificate})
	writeToFile(certfile, out.String())

}

func tryToOpenFile(f string) bool{
	_, err := os.Open(f)
	if (err != nil){
		return false
	}
	return true
}

func SetupKeyAndCertificate(c Config){

	tryopen := tryToOpenFile(c.Certfile) && tryToOpenFile(c.Keyfile)
	if (!tryopen){
		generateKeyAndCertificate(c.Keyfile, c.Certfile)
	}

}