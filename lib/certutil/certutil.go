package certutil

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Showmax/go-fqdn"
	"github.com/denisbrodbeck/machineid"
)

// KeyValuePair is for poor people who can't affort tuples
type KeyValuePair struct {
	Key, Value string
}

// CertToKVPairs parses the certificate data to extract some interesting metadata
// to attach to our printout. There is probably some smarter way to do this with
// reflection or something.
func CertToKVPairs(certData []byte) ([]KeyValuePair, error) {
	var output = []KeyValuePair{}
	var err error

	output = append(output, KeyValuePair{"PrintTime", time.Now().Format(time.RFC3339)})

	output = append(output, KeyValuePair{"Hostname", fqdn.Get()})

	machineID, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}
	output = append(output, KeyValuePair{"MachineID", machineID})

	pemCert, _ := pem.Decode(certData)

	cert, err := x509.ParseCertificate(pemCert.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	rawIssuer := string(cert.Issuer.String())
	issuer := strings.Split(rawIssuer, ",")
	output = append(output, KeyValuePair{"Issuer", strings.Join(issuer, "\n")})

	rawSubject := string(cert.Subject.String())
	subject := strings.Split(rawSubject, ",")
	output = append(output, KeyValuePair{"Subject", strings.Join(subject, "\n")})

	output = append(output, KeyValuePair{"Version", fmt.Sprintf("%X", cert.Version)})
	output = append(output, KeyValuePair{"SerialNumber", cert.SerialNumber.String()})
	output = append(output, KeyValuePair{"NotBefore", cert.NotBefore.Format(time.RFC3339)})
	output = append(output, KeyValuePair{"NotAfter", cert.NotAfter.Format(time.RFC3339)})
	validityLength := cert.NotAfter.Sub(cert.NotBefore).Hours() / 24
	output = append(output, KeyValuePair{"ValidFor", strconv.FormatFloat(validityLength, 'f', 0, 64) + " days"})
	output = append(output, KeyValuePair{"SigAlgo", cert.SignatureAlgorithm.String()})
	output = append(output, KeyValuePair{"PubKeyAlgo", cert.PublicKeyAlgorithm.String()})
	output = append(output, KeyValuePair{"IsCA", strconv.FormatBool(cert.IsCA)})

	if cert.IsCA {
		output = append(output, KeyValuePair{"MaxPathLen", fmt.Sprintf("%X", cert.MaxPathLen)})
	}
	if len(cert.DNSNames) > 0 {
		output = append(output, KeyValuePair{"DNS SANs", strings.Join(cert.DNSNames, "\n")})
	}
	if len(cert.IPAddresses) > 0 {
		var ipAddresses []string
		for i := 0; i < len(cert.IPAddresses); i++ {
			ipAddresses = append(ipAddresses, []string{cert.IPAddresses[i].String()}...)
		}
		output = append(output, KeyValuePair{"IP SANs", strings.Join(ipAddresses, "\n")})
	}

	return output, err
}

func PrintCertificate(certData []byte) {

}
