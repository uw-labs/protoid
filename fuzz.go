// +build gofuzz

package protoid

func Fuzz(data []byte) int {
	dec, err := Decode(data)
	if err != nil {
		if dec != nil {
			panic("decoded despite error")
		}
		return 0
	} else {
		return 1
	}
}
