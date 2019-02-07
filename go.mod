module itkettle.org/avanier/gorecptprint

require (
	github.com/boombuler/barcode v1.0.0
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	github.com/mkideal/pkg v0.0.0-20170503154153-3e188c9e7ecc
	golang.org/x/image v0.0.0-20190118043309-183bebdce1b2
	golang.org/x/sys v0.0.0-20190201152629-afcc84fd7533 // indirect
	golang.org/x/text v0.3.0 // indirect
)

replace itkettle.org/avanier/gorecptprint/lib/dmtx => ./dmtx

replace itkettle.org/avanier/gorecptprint/lib/extras => ./extras

replace itkettle.org/avanier/gorecptprint/lib/tf6 => ./tf6
