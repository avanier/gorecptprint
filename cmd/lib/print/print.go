package print

import (
	"fmt"
	"log"

	"go.bug.st/serial"
)

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

func Print() {
	// ports, err := serial.GetPortsList()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(ports)

	mode := &serial.Mode{
		BaudRate: 19200,
		DataBits: 8,
		StopBits: serial.OneStopBit,
		Parity:   serial.NoParity,
	}
	port, err := serial.Open("/dev/ttyS0", mode)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(port.GetModemStatusBits())

	defer port.Close()

	playReadyTune(port)
	requestPrinterID(port)
}

func playReadyTune(port serial.Port) {
	_, err := port.Write(readyTune)
	if err != nil {
		log.Fatal(err)
	}
}

func requestPrinterID(port serial.Port) {
	var buf = make([]byte, 8+7)

	port.Write([]byte{0x1d, 0x49, 0x01}) // Request printer ID
	// port.Write([]byte{0x1b, 0x76})

	// for {
	// 	// Reads up to 100 bytes
	// 	n, err := port.Read(buf)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	if n == 0 {
	// 		fmt.Println("\nEOF")
	// 		break
	// 	}

	// 	fmt.Printf("%d: %x\n", n, buf[:n])

	// 	// If we receive a newline stop reading
	// 	if strings.Contains(string(buf[:n]), "\n") {
	// 		break
	// 	}
	// }

	port.Read(buf)

	fmt.Printf("response:\t %v\n", buf)
	fmt.Printf("response:\t %x\n", buf)

	for i, c := range buf {
		fmt.Printf("byte %d: %2x %v\n", i, c, bytes2Bits([]byte{c}))
	}

	fmt.Print(buf)
}

func bytes2Bits(data []byte) []int {
	dst := make([]int, 0)
	for _, v := range data {
		for i := 0; i < 8; i++ {
			move := uint(7 - i)
			dst = append(dst, int((v>>move)&1))
		}
	}
	// fmt.Println(len(dst))
	return dst
}

// var options = serial.OpenOptions{
// 	PortName:        "/dev/ttyS0",
// 	BaudRate:        19200,
// 	DataBits:        8,
// 	StopBits:        1,
// 	MinimumReadSize: 4,
// }
