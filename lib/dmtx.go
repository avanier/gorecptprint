package dmtx

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func GenDMXT() (error error, barcode barcode.Barcode) {
	// Create the barcode
	qrCode, _ := qr.Encode("Hello World ougabouga", qr.M, qr.Auto)

	// Scale the barcode to 200x200 pixels
	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	return nil, qrCode
}
