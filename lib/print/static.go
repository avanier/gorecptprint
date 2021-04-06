package print

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

var byeTune = []byte{
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

// A multi-line test graphic to check print head alignment, and byte count
// alignment.
var dummyGraphic = []byte{
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

var alignLeft = []byte{0x1b, 0x61, 0x00}
var alignCenter = []byte{0x1b, 0x61, 0x01}
var setTitleFont = []byte{0x1d, 0x21, 0x11} // make font twice width, twice as high
var setParaFont = []byte{0x1d, 0x21, 0x01}  // make font default size

// GraphicProps actually captures the properties found at pp.157-158 of the
// programmer's manual.
type GraphicProps struct {
	// Whether the graphic is to be printed double-size or not. Valid values are 0..2.
	D int16
	// The width in dots
	W int16
	// The height in dots.
	H int16
}

var loremGibson = "Images formed and reformed: a flickering montage of the Sprawl's towers and ragged Fuller domes, dim figures moving toward him in the coffin for Armitage's call. The alarm still oscillated, louder here, the rear wall dulling the roar of the spherical chamber. He'd waited in the human system. The alarm still oscillated, louder here, the rear wall dulling the roar of the Villa bespeak a turning in, a denial of the bright void beyond the hull. The semiotics of the spherical chamber. Then a mist closed over the black water and the amplified breathing of the blowers and the amplified breathing of the fighters. He woke and found her stretched beside him in the dark, curled in his sleep, and wake alone in the tunnel's ceiling. He'd waited in the puppet place had been a subunit of Freeside's security system. Case felt the edge of the previous century. Case felt the edge of the Villa bespeak a turning in, a denial of the bright void beyond the hull. The Tessier-Ashpool ice shattered, peeling away from the Chinese program's thrust, a worrying impression of solid fluidity, as though the shards of a broken mirror bent and elongated as they rotated, but it never told the correct time."
