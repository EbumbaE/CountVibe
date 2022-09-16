package certificate

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

type CertAuhtory struct {
	cert     *x509.Certificate
	buffCert *bytes.Buffer

	key     *ecdsa.PrivateKey
	buffKey *bytes.Buffer
}

type Certificate struct {
	cert     *x509.Certificate
	buffCert *bytes.Buffer

	key     *ecdsa.PrivateKey
	buffKey *bytes.Buffer
}

func writeToFile(filename, buf string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(buf)

	return nil
}

func newSubject() pkix.Name {
	return pkix.Name{
		Organization: []string{"Ebumba_E"},
		Country:      []string{"RU"},
		Locality:     []string{"Moscow"},
	}
}

func newCACert(subject pkix.Name) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber:          big.NewInt(2022),
		Subject:               subject,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

func newCert(subject pkix.Name) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber:          big.NewInt(2022),
		Subject:               subject,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		KeyUsage:              x509.KeyUsageDigitalSignature,
	}
}

func newCertAuhtory(subject pkix.Name) (*CertAuhtory, error) {

	cert := newCACert(subject)

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	bytesCert, err := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}

	buffCert := &bytes.Buffer{}
	pem.Encode(buffCert, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: bytesCert,
	})

	out := &bytes.Buffer{}
	buffKey, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}
	pem.Encode(out, &pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: buffKey,
	})

	return &CertAuhtory{
		cert:     cert,
		buffCert: buffCert,
		key:      key,
		buffKey:  out,
	}, nil
}

func newCertificate(ca *CertAuhtory, subject pkix.Name) (*Certificate, error) {

	cert := newCert(subject)

	certKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	bytesCert, err := x509.CreateCertificate(rand.Reader, cert, ca.cert, &certKey.PublicKey, ca.key)
	if err != nil {
		return nil, err
	}

	buffCert := &bytes.Buffer{}
	pem.Encode(buffCert, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: bytesCert,
	})

	out := &bytes.Buffer{}
	buffKey, err := x509.MarshalECPrivateKey(certKey)
	if err != nil {
		return nil, err
	}
	pem.Encode(out, &pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: buffKey,
	})

	return &Certificate{
		cert:     cert,
		buffCert: buffCert,
		key:      certKey,
		buffKey:  out,
	}, nil
}

func generateKeyAndCertificate(keyPath, certPath string) error {

	subject := newSubject()

	ca, err := newCertAuhtory(subject)
	if err != nil {
		return err
	}

	cert, err := newCertificate(ca, subject)
	if err != nil {
		return err
	}

	if err := writeToFile(keyPath, cert.buffKey.String()); err != nil {
		return err
	}
	if err := writeToFile(certPath, cert.buffCert.String()); err != nil {
		return err
	}

	return nil
}

func tryToOpenFile(f string) bool {
	_, err := os.Open(f)
	if err != nil {
		return false
	}
	return true
}

func SetupKeyAndCertificate(c Config) error {

	tryopen := tryToOpenFile(c.CertPath) && tryToOpenFile(c.KeyPath)
	if !tryopen {
		err := generateKeyAndCertificate(c.KeyPath, c.CertPath)
		if err != nil {
			return err
		}
	}

	return nil
}
