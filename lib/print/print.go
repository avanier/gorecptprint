package print

import (
	"github.com/avanier/gorecptprint/lib/comm"
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
	playReadyTune()
	// verifyPreviousCommand()
	// requestPrinterID()
	// verifyPreviousCommand()
	// exerciseProgram()
	// verifyPreviousCommand()
	// comm.WatchPrinterOutput()
}

func playReadyTune() {
	comm.ExecuteCommand(readyTune)
}

// var options = serial.OpenOptions{
// 	PortName:        "/dev/ttyS0",
// 	BaudRate:        19200,
// 	DataBits:        8,
// 	StopBits:        1,
// 	MinimumReadSize: 4,
// }
