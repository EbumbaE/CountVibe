package server

import (
	"CountVibe/internal/authorization"
	"CountVibe/internal/log"

	"bytes"
	"os"
	"time"
	"net/http"
	"net"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
)

type Server struct {
	Port string

	homepage string 
	loginpage string 
	authpage string 
	refreshpage string
	diarypage string
	registrationpage string

	Certfile string
	Keyfile string

	Logger log.Logger
}

func CreateServer(c Config, logger log.Logger) *Server{
	return &Server{
		Port: c.Port,

		homepage: c.Homepage,
		loginpage: c.Loginpage,
		authpage: c.Authpage,
		refreshpage: c.Refreshpage,
		diarypage: c.Diarypage,
		registrationpage: c.Registrationpage,

		Certfile: c.Certfile,
		Keyfile: c.Keyfile,
		Logger: logger,
	}
}

func (s *Server) Run(){

	s.setupServerHandlers()
	s.setupAuthHandlers()
	s.setupKeyAndCertificate()
	s.Logger.Error(http.ListenAndServeTLS(s.Port, s.Certfile, s.Keyfile, nil))

}

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

func (s *Server)setupKeyAndCertificate(){

	tryopen := tryToOpenFile(s.Certfile) && tryToOpenFile(s.Keyfile)
	if (!tryopen){
		generateKeyAndCertificate(s.Keyfile, s.Certfile)
	}

}

func handler(w http.ResponseWriter, r *http.Request){
	http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
}

func homehandler(w http.ResponseWriter, r *http.Request){
    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, "../../static/html/home.html")

    }
}

func (s *Server) setupServerHandlers(){
	http.HandleFunc("/", handler)
	http.HandleFunc(s.homepage, homehandler)
}


func (s *Server) setupAuthHandlers(){

    http.HandleFunc(s.authpage, authorization.AuthHandler)
    http.HandleFunc(s.loginpage, authorization.LoginHandler)
    http.HandleFunc(s.refreshpage, authorization.RefreshHandler)
    http.HandleFunc(s.registrationpage, authorization.RegistrationHandler)
             
}