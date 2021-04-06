package cmd

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"strings"

	"github.com/avanier/gorecptprint/lib/certutil"
	"github.com/avanier/gorecptprint/lib/comm"
	"github.com/avanier/gorecptprint/lib/print"
	"github.com/avanier/gorecptprint/lib/util"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var printCmd = &cobra.Command{
	Use:   "print [certificate.pem ...]",
	Short: "reads certificates and outputs the data",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		comm.Init()

		print.PlayReadyTune()

		for certNum, cert := range args {
			var certData []byte
			var pairs []certutil.KeyValuePair

			if cert == "-" {
				cert = "/dev/stdin"
			}

			certData, err = ioutil.ReadFile(cert)
			if err != nil {
				log.Fatal(err)
			}

			pairs, err = certutil.CertToKVPairs(certData)
			if err != nil {
				log.Fatal(err)
			}

			for i := 0; i < len(pairs); i++ {
				print.PrintTitleValues(pairs[i].Key, pairs[i].Value)
			}

			print.FeedPaper()

			stringArray := util.SplitString(string(certData), 174)

			if len(stringArray[len(stringArray)-1]) < 174 {
				shortString := stringArray[len(stringArray)-1]
				paddedString := shortString + "\n" + strings.Repeat(string("\x00"), (173-len(shortString)))
				stringArray[len(stringArray)-1] = paddedString
			}

			scaleFactor := 2
			symbPerLine := 2
			offset := 16 // px @ 8px/mm

			// Print x DMTX per line

			var numGroups int
			if len(stringArray)%symbPerLine != 0 {
				numGroups = (len(stringArray) / symbPerLine) + 1
			} else {
				numGroups = (len(stringArray) / symbPerLine)
			}

			for grInd := 0; grInd < numGroups; grInd++ {
				stringBatch := make([]string, symbPerLine)
				var masterRectangle image.Rectangle
				var masterImage draw.Image
				var baseWidth int

				for sInd := 0; (sInd < symbPerLine) && ((grInd*symbPerLine)+sInd < len(stringArray)); sInd++ {
					stringBatch[sInd] = stringArray[(grInd*symbPerLine)+sInd]
					oneBarcode, _ := datamatrix.Encode(stringBatch[sInd])
					oneBarcode, _ = barcode.Scale(oneBarcode, oneBarcode.Bounds().Dx()*scaleFactor, oneBarcode.Bounds().Dy()*scaleFactor)
					if sInd == 0 {
						baseWidth = oneBarcode.Bounds().Dx()
					}
					sp := image.Point{(baseWidth * sInd) + (offset * sInd), 0}    // place image starting point with offset
					dr := image.Rectangle{sp, sp.Add(oneBarcode.Bounds().Size())} // make a drawing rectangle for holding them pixels
					if sInd == 0 {
						fullWidth := (oneBarcode.Bounds().Dx() * symbPerLine) + (offset * symbPerLine)
						masterRectangle = image.Rectangle{image.ZP, image.Point{fullWidth + (1 * scaleFactor), oneBarcode.Bounds().Dy() + (1 * scaleFactor)}}
						masterImage = image.NewRGBA(masterRectangle)
						white := color.RGBA{255, 255, 255, 255}
						draw.Draw(masterImage, masterImage.Bounds(), &image.Uniform{white}, image.ZP, draw.Src) // apply primer on the whole rectangle
					}
					draw.Draw(masterImage, dr, oneBarcode, image.ZP, draw.Src)
				}

				pixels, width, height := util.GetPixels(masterImage)

				dmtxProps := print.GraphicProps{D: 2, W: int16(width / 8), H: int16(height / 8)}
				print.PrintGraphic(dmtxProps, pixels)

				print.FeedPaper()
			}

			// Wrap up this one cert
			print.FeedPaper()
			// If the cert's not the last and we didn't set the `--no-cut-between` flag
			if certNum+1 != len(args) && !viper.GetBool("no-cut-between") {
				print.CutPaper()
			}
		}

		// Wrap up the whole operation
		if !viper.GetBool("no-cut-last") {
			print.CutPaper()
		}
		print.PlayByeTune()
	},
}

func init() {
	var defaultSerialPort = "/dev/ttyS0"
	var defaultWriteChunkSize = uint(192)

	printCmd.PersistentFlags().String("port", defaultSerialPort, "serial port to use for communication")
	printCmd.PersistentFlags().Uint("write-chunk-size", defaultWriteChunkSize, "the size in bytes of every chunk of data sent to the printer")
	printCmd.PersistentFlags().Bool("no-cut-between", false, "don't cut the receipt between every certificate")
	printCmd.PersistentFlags().Bool("no-cut-last", false, "don't cut the receipt after the last certificate")
	printCmd.PersistentFlags().Bool("no-ident", false, "don't print information that could identify the machine where this was printed")

	viper.BindPFlags(printCmd.PersistentFlags())
	viper.BindPFlags(printCmd.Flags())
}
