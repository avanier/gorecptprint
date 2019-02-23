package extras

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jacobsa/go-serial/serial"
	"itkettle.org/avanier/gorecptprint/lib/tf6"

	"github.com/Showmax/go-fqdn"
	"github.com/denisbrodbeck/machineid"
)

func PrintDummyGraphic(options serial.OpenOptions) {
	var printDummyGraph = []byte{
		0x1b, 0x2a, // Select and print Graphic

		0x02, // Dot density

		0x04, // N x 8 dots wide
		0x02, // N x 8 dots high

		0xff, 0x00, 0xff, 0x00,
		0xff, 0x00, 0xff, 0x00,
		0xff, 0x00, 0xff, 0x00,
		0xff, 0x00, 0xff, 0x00,
		0x00, 0xff, 0x00, 0xff,
		0x00, 0xff, 0x00, 0xff,
		0x00, 0xff, 0x00, 0xff,
		0x00, 0xff, 0x00, 0xff,

		0xff, 0x00, 0xff, 0x00,
		0xff, 0x00, 0xff, 0x00,
		0xff, 0x00, 0xff, 0x00,
		0xff, 0x00, 0xff, 0x00,
		0x00, 0xff, 0x00, 0xff,
		0x00, 0xff, 0x00, 0xff,
		0x00, 0xff, 0x00, 0xff,
		0x00, 0xff, 0x00, 0xff,
	}
	tf6.ExecuteHex(printDummyGraph, options)
}

func ByeTune(options serial.OpenOptions) {
	var readyTune = []byte{
		0x1b, 0x07,
		0x02,
		0x9a,
		0x1b, 0x07,
		0x01,
		0x99,
		0x1b, 0x07,
		0x01,
		0x95,
	}
	tf6.ExecuteHex(readyTune, options)
}

func ReadyTune(options serial.OpenOptions) {
	// Plays some beeps to signal end of initialization
	// See <p.144>
	var readyTune = []byte{
		0x1b, 0x07, // Start the sequence
		0x02, // Set the duration from 01 - FF times 0.1 seconds
		0x90, // Binary conversion of 10010000 - (10)<soft>(01)<octave 2>(0000)<note c>
		0x1b, 0x07,
		0x01,
		0x95,
		0x1b, 0x07,
		0x01,
		0x99,
	}
	tf6.ExecuteHex(readyTune, options)
}

func SplitString(longString string, maxLen int) []string {
	splits := []string{}

	var l, r int
	for l, r = 0, maxLen; r < len(longString); l, r = r, r+maxLen {
		for !utf8.RuneStart(longString[r]) {
			r--
		}
		splits = append(splits, longString[l:r])
	}
	splits = append(splits, longString[l:])
	return splits
}

// KeyValuePair is for poor people who can't affort tuples
type KeyValuePair struct {
	Key, Value string
}

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
		output = append(output, KeyValuePair{"DNS SANs: ", strings.Join(cert.DNSNames, "\n")})
	}
	if len(cert.IPAddresses) > 0 {
		var ipAddresses []string
		for i := 0; i < len(cert.IPAddresses); i++ {
			ipAddresses = append(ipAddresses, []string{cert.IPAddresses[i].String()}...)
		}
		output = append(output, KeyValuePair{"IP SANs: ", strings.Join(ipAddresses, "\n")})
	}

	return output, err
}
