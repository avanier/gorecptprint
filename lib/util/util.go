package util

// Bytes2Bits converts bytes into a human readable string of 0s and 1s.
func Bytes2Bits(data []byte) []int {
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
