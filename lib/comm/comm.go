package comm

import (
	"fmt"
	"log"
	"time"

	"github.com/avanier/gorecptprint/lib/util"
	"github.com/spf13/viper"
	"go.bug.st/serial"
)

/*
When the printer transmits those characters, they are not counted in the
byte count of status messages. See page 29 for mildly confusing instructions.
When "Transparent Mode" is disabled, those bytes should not appear within
messages. For additional speed, one could go back to transparent mode, but
then you have to extract those bytes, which may be XORed within some of the
messages. No thank your siree.

See relevant documentation on pp. 29, 176.
*/

// DLE is some kind of prefix. I don't quite understand what it's for.
const DLE = 0x10

// XON announces the printer's buffer is below 20% capacity.
const XON = 0x11

// XOFF announces the printer's buffer is above 80% capacity.
const XOFF = 0x13

// Exposed the global instance of the serial port.
var Port serial.Port

// writeChunkSize is a tunable that allows to change the size of every write
// loop
var writeChunkSize uint

// Init prepares the com port for communication with the printer
func Init() {
	var err error
	var serial_port string
	writeChunkSize = viper.GetUint("write-chunk-size")

	mode := &serial.Mode{
		BaudRate: 19200,
		DataBits: 8,
		StopBits: serial.OneStopBit,
		Parity:   serial.NoParity,
	}

	serial_port = viper.GetString("port")
	fmt.Printf("using port %s\n", serial_port)

	Port, err = serial.Open(serial_port, mode)
	if err != nil {
		log.Fatal(err)
	}

	startFlowReader()
	resetPrinter()
	disableTransparentXONXOFF()
}

func displayResponse(resp []byte) {
	Port.Read(resp)

	fmt.Printf("response:\t %v\n", resp)
	fmt.Printf("response:\t %X\n", resp)

	for i, c := range resp {
		fmt.Printf("byte %d: %2X %v\n", i, c, util.Bytes2Bits([]byte{c}))
	}

	fmt.Println(resp)
}

// Check page 175

// ExecuteCommand executes a command with an arbitrary byte payload.
func ExecuteCommand(cmd []byte) {
	fmt.Printf("executing:\t 0x%X\n", cmd)
	writeToPort(cmd)
}

// UnsafeExecuteCommand sends a command directly to the printer without checking
// the buffer status. This is useful for executing async commands (see p. xxx) or
// for configuring the printer before the flow control system is enabled.
func UnsafeExecuteCommand(cmd []byte) {
	fmt.Printf("unsafely executing:\t 0x%X\n", cmd)
	Port.Write(cmd)
}

// chunkData chunks data in groups sized by `writeChunkSize`.
// Shamelessly lifted from SliceTricks.
// https://github.com/golang/go/wiki/SliceTricks
func chunkData(unchunked []byte) [][]byte {
	dataChunks := make([][]byte, 0, (len(unchunked)+int(writeChunkSize)-1)/int(writeChunkSize))

	for int(writeChunkSize) < len(unchunked) {
		unchunked, dataChunks = unchunked[writeChunkSize:], append(dataChunks, unchunked[0:writeChunkSize:writeChunkSize])
	}
	dataChunks = append(dataChunks, unchunked)
	return dataChunks
}

func writeToPort(cmd []byte) {
	// take the payload, split it in chunks of writeChunkSize
	// for every chunk check if we can currently write, if the buffer is full wait
	// for the channel to continue

	dataChunks := chunkData(cmd)

	for i := 0; i <= len(dataChunks)-1; i++ {
		if !recipientReady {
			_ = <-readyToWrite
		}
		Port.Write(dataChunks[i])
	}
}

// disableTransparentXONXOFF disables transparent flow control with an unsafe
// write.
func disableTransparentXONXOFF() {
	UnsafeExecuteCommand([]byte{0x10, 0x05, 0x43})
}

// resetPrinter performs a soft power-cycle of the printer and returns its to the
// default parameters. This is an async command and is done with an unsafe write.
func resetPrinter() {
	UnsafeExecuteCommand([]byte{0x10, 0x05, 0x40})
	time.Sleep(time.Second * 3)
}

// WatchPrinterOutput is a debug convenience function that read the data sent from
// the printer byte by byte.
func WatchPrinterOutput() {
	var err error
	var i int
	fmt.Println("watching printer output forever, 1 byte at a time")

	buf := make([]byte, 1)
	c := -1

	for {
		c++
		i, err = Port.Read(buf)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Printf("[%6d]read %d byte: 0x%X\n", c, i, buf)
	}

	// scanner := bufio.NewScanner(port)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text()) // Println will add back the final '\n'
	// }
	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }
}

func exerciseProgram() {
	ExecuteCommand([]byte{0x1b, 0x78})
}

func requestPrinterID() {
	var buf = make([]byte, 8+5)
	ExecuteCommand([]byte{0x1d, 0x49, 0x01}) // Request printer ID
	displayResponse(buf)
}

// verifyPreviousCommand unsafely sends a status query.
func verifyPreviousCommand() {
	var buf = make([]byte, 8)
	UnsafeExecuteCommand([]byte{0x1b, 0x00, 0x80, 0x00}) // Verify previous command completed
	displayResponse(buf)
}

// readyToWrite exposes a channel signalling that the printer's buffer is ready
// to accept data again after being declared full.
var readyToWrite chan bool

// recipientReady tracks whether the printer's buffer is ready for data or not.
var recipientReady bool

// startFlowReader sets up the required machinery for managing XON/XOFF flow
// control.
func startFlowReader() {
	go initReaderLoop()
}

func emitReadyToWrite() {
	// Here we do a non-blocking unbuffered channel write to announce we are ready
	// to work again.
	select {
	case readyToWrite <- true:
		// If there is a receiver available for that message, yay.
	default:
		// Otherwise carry on.
	}
}

func initReaderLoop() {
	var err error
	var buf = make([]byte, 1)
	recipientReady = false

	for {
		_, err = Port.Read(buf)
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, b := range buf {
			switch b {
			case XON:
				if !recipientReady {
					recipientReady = true
					go emitReadyToWrite()
				}
			case XOFF:
				recipientReady = false
			}
		}
	}
}

/*
Pseudocode

when writing
	if recipientReady
		then write
	else
		await chan buffer not full anymore

continuously
	if receive XON
		if recipientReady true
			continue
		else
			set recipientReady true
			go emitReadyToWrite or timeout after 3s
	if receive XOFF
		set recipientReady false

When printing graphics, be aware the print head prints 8 dots high at a time at
200 dpi. See p.157 for double density graphics options.

The buffer size on 4610-TF6 is 64KiB
See p.158
*/
