package audio

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

// ParseWavProps reads audio properties from a WAV file
func ParseWavProps(path string) (*AudioProps, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read RIFF header
	header := make([]byte, 12)
	if _, err := f.Read(header); err != nil {
		return nil, err
	}
	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return nil, errors.New("not a valid WAV file")
	}

	var channels, sampleRate, byteRate, bitsPerSample int
	var dataSize int64

	// Read chunks
	for {
		chunkHeader := make([]byte, 8)
		if _, err := f.Read(chunkHeader); err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		chunkID := string(chunkHeader[0:4])
		chunkSize := binary.LittleEndian.Uint32(chunkHeader[4:8])

		// Handle padding byte if chunk size is odd
		padding := int64(chunkSize % 2)

		if chunkID == "fmt " {
			if chunkSize < 16 {
				return nil, errors.New("invalid fmt chunk")
			}
			fmtData := make([]byte, chunkSize)
			if _, err := f.Read(fmtData); err != nil {
				return nil, err
			}
			
			// Handle padding
			if padding > 0 {
				f.Seek(padding, 1)
			}

			// AudioFormat (2 bytes)
			// NumChannels (2 bytes)
			channels = int(binary.LittleEndian.Uint16(fmtData[2:4]))
			// SampleRate (4 bytes)
			sampleRate = int(binary.LittleEndian.Uint32(fmtData[4:8]))
			// ByteRate (4 bytes)
			byteRate = int(binary.LittleEndian.Uint32(fmtData[8:12]))
			// BlockAlign (2 bytes)
			// BitsPerSample (2 bytes)
			bitsPerSample = int(binary.LittleEndian.Uint16(fmtData[14:16]))
		} else if chunkID == "data" {
			dataSize = int64(chunkSize)
			// We found data. 
			// If we have fmt info, we can calculate duration.
			// But we should continue just in case (though unlikely to have fmt after data)
			// Skip data
			f.Seek(int64(chunkSize)+padding, 1)
		} else {
			// Skip chunk
			f.Seek(int64(chunkSize)+padding, 1)
		}
	}

	if sampleRate == 0 || byteRate == 0 {
		return nil, errors.New("invalid WAV format or missing fmt chunk")
	}

	if dataSize == 0 {
		// If data chunk was not found or size is 0 (unlikely for valid file)
		// Try to estimate from file size if possible, or return error
		fileInfo, _ := f.Stat()
		dataSize = fileInfo.Size() - 44 // Approximate
	}

	duration := float64(dataSize) / float64(byteRate)
	
	// BitRate = ByteRate * 8 / 1000 (kbps)
	bitRate := byteRate * 8 / 1000

	return &AudioProps{
		Duration:   duration,
		SampleRate: sampleRate,
		BitDepth:   bitsPerSample,
		Channels:   channels,
		BitRate:    bitRate,
	}, nil
}
