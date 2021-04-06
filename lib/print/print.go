package print

import (
	"log"
	"strings"

	"github.com/avanier/gorecptprint/lib/comm"
	"github.com/avanier/gorecptprint/lib/util"
)

func PrintTest() {
	PlayReadyTune()

	PrintDummyGraphic()
	PrintTitle("This is a test title\n\n")
	PrintParagraph(strings.Repeat("this is a test ", 30) + "\n\n")
	PrintParagraph(loremGibson)

	FeedPaper()
	CutPaper()

	PlayByeTune()
}

// PlayReadyTune plays a melody indicating the printer is ready to work
func PlayReadyTune() {
	comm.ExecuteCommand(readyTune)
}

// PlayReadyTune plays a melody indicating the printer is done working
func PlayByeTune() {
	comm.ExecuteCommand(byeTune)
}

func PrintDummyGraphic() {
	props := GraphicProps{
		D: 2, // print double-density
		W: 4,
		H: 2,
	}

	PrintGraphic(props, dummyGraphic)
}

// PrintString prints string on paper.
func PrintString(inputString string) {
	if util.IsASCII(inputString) {
		comm.ExecuteCommand([]byte(inputString))
	} else {
		log.Fatalf("string `%s` contains non-ASCII characters, bailing out", inputString)
	}
}

// PrintGraphic prints a black and white graphic of known dimensions.
func PrintGraphic(props GraphicProps, data []byte) {
	comm.ExecuteCommand([]byte{0x1b, 0x2F, 0x03}) // set the print speed to 15 lps for better graphics

	comm.ExecuteCommand(alignCenter)

	commandPrefix := []byte{0x1b, 0x2a}                             // Select and print Graphic
	wholeCommand := append(commandPrefix, []byte{byte(props.D)}...) // add doubleprint
	wholeCommand = append(wholeCommand, []byte{byte(props.W)}...)   // add width
	wholeCommand = append(wholeCommand, []byte{byte(props.H)}...)   // add height
	wholeCommand = append(wholeCommand, data...)                    // and the data
	comm.ExecuteCommand(wholeCommand)

	comm.ExecuteCommand([]byte{0x1b, 0x2F, 0x00}) // set the print speed back to default 52 lps
	comm.ExecuteCommand(alignLeft)
}

// PrintTitle prints a center aligned title
func PrintTitle(title string) {
	comm.ExecuteCommand(alignCenter)
	comm.ExecuteCommand(setTitleFont)
	PrintString(title)
}

// PrintParagraph prints a center aligned title, followed by a left aligned paragraph
func PrintParagraph(values string) {

	comm.ExecuteCommand(alignLeft)
	comm.ExecuteCommand(setParaFont)
	PrintString(values)
}

// PrintTitleValues print a title with a bunch of values
func PrintTitleValues(title string, values string) {
	PrintTitle(title)
	PrintParagraph(values + "\n")
}

// CutsPaper cuts the receipt paper. See p.181.
func CutPaper() {
	comm.ExecuteCommand([]byte{0x0c})
}

// FeedPaper finishes printing the data in the buffer and feeds paper using
// minimum units of line feed steps. 2mm in the case of the TF6. See p.156.
func FeedPaper() {
	comm.ExecuteCommand([]byte{0x1b, 0x4a, 0x20})
}
