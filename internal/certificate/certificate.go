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

func writeToFile(filename, buf string) error{
	file, err := os.Create(filename)
	if (err != nil){
		return err
	}
	defer file.Close()
	file.WriteString(buf)
	
	return nil
}

func generateKeyAndCertificate(keyfile, certfile string) error{

	privatekey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	bufkey, _ := x509.MarshalECPrivateKey(privatekey)
	out := &bytes.Buffer{}
	pem.Encode(out,  &pem.Block{Type: "EC PRIVATE KEY", Bytes: bufkey})
	if err := writeToFile(keyfile, out.String()); err != nil{
		return err
	}

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
		return err
	}

	out.Reset()
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certificate})
	if err := writeToFile(certfile, out.String()); err != nil{
		return err
	}

	return nil
}

func tryToOpenFile(f string) bool{
	_, err := os.Open(f)
	if (err != nil){
		return false
	}
	return true
}

func SetupKeyAndCertificate(c Config) error{

	tryopen := tryToOpenFile(c.Certfile) && tryToOpenFile(c.Keyfile)
	if (!tryopen){
		err := generateKeyAndCertificate(c.Keyfile, c.Certfile)
		if err != nil{
			return err
		}
	}

	return nil
}