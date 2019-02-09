package tf6

import (
	"log"

	"github.com/jacobsa/go-serial/serial"
)

type GraphicProps struct {
	D, W, H int16
}

func ExecuteHex(b []byte, options serial.OpenOptions) {
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	defer port.Close()

	_, err = port.Write(b)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
}

func PrintString(inputString string, options serial.OpenOptions) {
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	defer port.Close()

	_, err = port.Write([]byte(inputString))
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
}

func PrintGraphic() {

}
