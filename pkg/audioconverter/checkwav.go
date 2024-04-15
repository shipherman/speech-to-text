package audioconverter

import "fmt"

// CheckWAVHeader checks if audiofile has apropriate wav header
func CheckWAVHeader(wav []byte) (bool, error) {
	// 8..11 bytes are contain format header
	// Approprite wave file header is WAVE
	//
	// TODO
	// Implement switch case to convert ogg to wav
	//
	waveHeader := string(wav[8:12])
	if waveHeader == "WAVE" {
		return true, nil
	}
	return false, fmt.Errorf("format header is not WAVE: %s", waveHeader)
}
