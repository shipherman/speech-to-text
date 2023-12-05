package audioconvertor

import "fmt"

func CheckWAVHeader(wav []byte) (bool, error) {
	waveHeader := string(wav[8:11])
	if waveHeader == "WAVE" {
		fmt.Println("zbs")
		return true, nil
	}
	return false, fmt.Errorf("format header is not WAVE: %s", waveHeader)
}
