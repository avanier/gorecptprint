module itkettle.org/avanier/gorecptprint

require (
	github.com/Showmax/go-fqdn v0.0.0-20180501083314-6f60894d629f
	github.com/boombuler/barcode v1.0.0
	github.com/denisbrodbeck/machineid v1.0.0
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	golang.org/x/sys v0.0.0-20190222072716-a9d3bda3a223 // indirect
)

replace itkettle.org/avanier/gorecptprint/lib/dmtx => ./dmtx

replace itkettle.org/avanier/gorecptprint/lib/extras => ./extras

replace itkettle.org/avanier/gorecptprint/lib/tf6 => ./tf6
