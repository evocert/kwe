package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"strings"
	"time"
)

type Certificate struct {
	certserial          int64
	Organization        string
	Country             string
	Province            string
	Locality            string
	StreetAddress       string
	PostalCode          string
	from                time.Time
	until               time.Time
	cert                *x509.Certificate
	certPrivateKey      *rsa.PrivateKey
	certPrivateKeyBytes []byte
	certPEM             *bytes.Buffer
	certPrivKeyPEM      *bytes.Buffer
	serverCert          tls.Certificate
	certPool            *x509.CertPool
}

func (cert *Certificate) ClientTLSConf() (clnttlsconf *tls.Config) {
	clnttlsconf = &tls.Config{
		RootCAs: cert.certPool,
	}
	return
}

func (cert *Certificate) ServerTLSConf() (svrtlsconf *tls.Config) {
	svrtlsconf = &tls.Config{
		Certificates: []tls.Certificate{cert.serverCert},
	}
	return
}

type CA struct {
	caserial          int64
	Organization      string
	Country           string
	Province          string
	Locality          string
	StreetAddress     string
	PostalCode        string
	from              time.Time
	until             time.Time
	ca                *x509.Certificate
	caPrivateKey      *rsa.PrivateKey
	caPrivateKeyBytes []byte
	caPEM             *bytes.Buffer
	caPrivKeyPEM      *bytes.Buffer
	certs             map[int64]*Certificate
}

type CAS struct {
	ca map[int64]*CA
}

func (ca *CA) Certificate(serial int64) (cert *Certificate) {
	cert, _ = ca.certs[serial]
	return
}

func (ca *CA) Register(serial int64) (err error) {
	if _, certok := ca.certs[serial]; !certok {
		cert := &x509.Certificate{
			SerialNumber: big.NewInt(serial),
			Subject: pkix.Name{
				Organization:  strings.Split(ca.Organization, ","),
				Country:       strings.Split(ca.Country, ","),
				Province:      strings.Split(ca.Province, ","),
				Locality:      strings.Split(ca.Locality, ","),
				StreetAddress: strings.Split(ca.StreetAddress, ","),
				PostalCode:    strings.Split(ca.PostalCode, ","),
			},
			IPAddresses:  []net.IP{net.IPv4zero, net.IPv6unspecified}, // IPv4(0, 0, 0, 0), net.IPv6loopback},
			NotBefore:    ca.from,
			NotAfter:     ca.until,
			SubjectKeyId: []byte{1, 2, 3, 4, 6},
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:     x509.KeyUsageDigitalSignature,
		}

		if certPrivKey, errcertPrivKey := rsa.GenerateKey(rand.Reader, 4096); errcertPrivKey == nil {
			if certBytes, errcertBytes := x509.CreateCertificate(rand.Reader, cert, ca.ca, &certPrivKey.PublicKey, ca.caPrivateKey); errcertBytes == nil {
				certPEM := new(bytes.Buffer)
				pem.Encode(certPEM, &pem.Block{
					Type:  "CERTIFICATE",
					Bytes: certBytes,
				})

				certPrivKeyPEM := new(bytes.Buffer)
				pem.Encode(certPrivKeyPEM, &pem.Block{
					Type:  "RSA PRIVATE KEY",
					Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
				})

				if serverCert, errserverCert := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes()); errserverCert == nil {
					certpool := x509.NewCertPool()
					certpool.AppendCertsFromPEM(ca.caPEM.Bytes())
					ca.certs[serial] = &Certificate{certserial: serial, cert: cert, from: ca.from, until: ca.until,
						Organization: ca.Organization, Country: ca.Country, Province: ca.PostalCode, Locality: ca.Locality, StreetAddress: ca.StreetAddress, PostalCode: ca.PostalCode,
						certPrivateKey: certPrivKey, certPrivateKeyBytes: certBytes, certPEM: certPEM, certPrivKeyPEM: certPrivKeyPEM,
						serverCert: serverCert}
				}
			}
		}
	}
	return
}

func (cas *CAS) CA(serial int64) (ca *CA) {
	ca, _ = cas.ca[serial]
	return ca
}

func (cas *CAS) Register(serial int64, stngs ...map[string]interface{}) (ca *CA, err error) {
	var Organization string = ""
	var Country string = ""
	var Province string = ""
	var Locality string = ""
	var StreetAddress string = ""
	var PostalCode string = ""
	var from time.Time = time.Now()
	var until time.Time = from.AddDate(10, 0, 0)
	if len(stngs) > 0 {
		for k, d := range stngs[0] {
			if v, _ := d.(string); v != "" {
				if strings.EqualFold(k, "orginization") {
					if Organization == "" {
						Organization = v
					}
				} else if strings.EqualFold(k, "country") {
					if Country == "" {
						Country = v
					}
				} else if strings.EqualFold(k, "province") {
					if Province == "" {
						Province = v
					}
				} else if strings.EqualFold(k, "locality") {
					if Locality == "" {
						Locality = v
					}
				} else if strings.EqualFold(k, "streetaddress") {
					if StreetAddress == "" {
						StreetAddress = v
					}
				} else if strings.EqualFold(k, "postalcode") {
					if PostalCode == "" {
						PostalCode = v
					}
				} else if strings.EqualFold(k, "from") {
					if t, terr := time.Parse(time.RFC3339, v); terr == nil {
						from = t
					}
				} else if strings.EqualFold(k, "until") {
					if t, terr := time.Parse(time.RFC3339, v); terr == nil {
						until = t
					}
				}
			}
		}
	}
	c := &x509.Certificate{
		SerialNumber: big.NewInt(serial),
		Subject: pkix.Name{
			Organization:  strings.Split(Organization, ","),  //[]string{"Company, INC."},
			Country:       strings.Split(Country, ","),       //[]string{"US"},
			Province:      strings.Split(Province, ","),      //[]string{""},
			Locality:      strings.Split(Locality, ","),      //[]string{"San Francisco"},
			StreetAddress: strings.Split(StreetAddress, ","), //[]string{"Golden Gate Bridge"},
			PostalCode:    strings.Split(PostalCode, ","),    //[]string{"94016"},
		},
		NotBefore:             from,
		NotAfter:              until,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	if caPrivKey, errcaPrivKey := rsa.GenerateKey(rand.Reader, 4096); errcaPrivKey == nil {
		if caBytes, errcaBytes := x509.CreateCertificate(rand.Reader, c, c, &caPrivKey.PublicKey, caPrivKey); errcaBytes == nil {
			caPEM := new(bytes.Buffer)
			pem.Encode(caPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: caBytes,
			})

			caPrivKeyPEM := new(bytes.Buffer)
			pem.Encode(caPrivKeyPEM, &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
			})
			ca = &CA{
				caserial:          serial,
				ca:                c,
				caPrivateKey:      caPrivKey,
				caPrivateKeyBytes: caBytes,
				caPEM:             caPEM,
				caPrivKeyPEM:      caPrivKeyPEM,
				Organization:      Organization,
				Country:           Country,
				Province:          Province,
				Locality:          Locality,
				StreetAddress:     StreetAddress,
				PostalCode:        PostalCode,
				from:              c.NotBefore, until: c.NotAfter, certs: map[int64]*Certificate{}}
			cas.ca[serial] = ca
		}
	}

	return
}

var gblcas *CAS = nil

func GLOBALCAS() *CAS {
	if gblcas == nil {
		gblcas = &CAS{ca: gblcas.ca}
	}
	return gblcas
}

func init() {
	if gblcas == nil {
		gblcas = &CAS{ca: map[int64]*CA{}}
	}
}
