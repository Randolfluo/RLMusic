package audio

import (
	"errors"
	"os"
)

type AudioProps struct {
	Duration   float64
	SampleRate int
	BitDepth   int
	Channels   int
	BitRate    int // kbps (approximate or extracted)
}

// ParseFlacProps reads audio properties from a FLAC file
func ParseFlacProps(path string) (*AudioProps, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 1. Check "fLaC" signature
	marker := make([]byte, 4)
	if _, err := f.Read(marker); err != nil {
		return nil, err
	}
	if string(marker) != "fLaC" {
		return nil, errors.New("not a valid FLAC file")
	}

	// 2. Scan metadata blocks
	for {
		header := make([]byte, 4)
		if _, err := f.Read(header); err != nil {
			return nil, err
		}

		isLast := (header[0] & 0x80) != 0
		blockType := header[0] & 0x7F
		length := int(uint32(header[1])<<16 | uint32(header[2])<<8 | uint32(header[3]))

		if blockType == 0 { // STREAMINFO
			data := make([]byte, length)
			if _, err := f.Read(data); err != nil {
				return nil, err
			}
			if length < 34 { // Min length for STREAMINFO is 34 bytes
				return nil, errors.New("invalid STREAMINFO block length")
			}
			return parseStreamInfo(data, path)
		} else {
			// Skip other blocks
			if _, err := f.Seek(int64(length), 1); err != nil {
				return nil, err
			}
		}

		if isLast {
			break
		}
	}

	return nil, errors.New("STREAMINFO block not found")
}

func parseStreamInfo(data []byte, path string) (*AudioProps, error) {
	// Offset 10: Sample Rate (20 bits), Channels (3 bits), Bits Per Sample (5 bits), Total Samples (36 bits)
	// data[10-12]:
	// 10: RRRR RRRR
	// 11: RRRR RRRR
	// 12: RRRR CCBB
	// 13: B BBTT TTTT ...

	// Let's implement correctly
	// Bytes 0-9: Min/Max Block/Frame Size. Skip.
	// Start at byte 10.

	// Combine bytes 10-17 into a uint64 to easily shift
	// But we only need a few bytes.
	// Sample Rate: 20 bits.
	// Channels: 3 bits.
	// Bits Per Sample: 5 bits.
	// Total Samples: 36 bits.

	// Use binary.BigEndian for standard access? No, fields are not byte aligned.

	// Read relevant bytes manually
	b10 := uint64(data[10])
	b11 := uint64(data[11])
	b12 := uint64(data[12])
	b13 := uint64(data[13])

	// Sample Rate (20 bits) from b10, b11, and high 4 bits of b12
	sampleRate := int((b10 << 12) | (b11 << 4) | ((b12 & 0xF0) >> 4))

	// Channels (3 bits) from b12 [3:1] (bits 1-3). Wait, let's trace bit by bit.
	// b12: RRRR C CC B
	// bits 7-4: SampleRate end
	// bits 3-1: Channels
	// bit 0: BitsPerSample start
	channels := int(((b12 & 0x0E) >> 1) + 1)

	// Bits Per Sample (5 bits)
	// b12: ... B
	// b13: BBBB ...
	// bit 0 of b12 is MSB. bits 7-4 of b13 are LSBs.
	bitDepth := int(((b12 & 0x01) << 4) | ((b13 & 0xF0) >> 4) + 1)

	// Total Samples (36 bits)
	// b13: ... T TTT (low 4 bits)
	// b14: TTTTTTTT
	// b15, b16, b17
	totalSamples := ((uint64(data[13]) & 0x0F) << 32) |
		(uint64(data[14]) << 24) |
		(uint64(data[15]) << 16) |
		(uint64(data[16]) << 8) |
		uint64(data[17])

	var duration float64
	if sampleRate > 0 {
		duration = float64(totalSamples) / float64(sampleRate)
	}

	// Calculate BitRate approximatly from file size
	fileInfo, err := os.Stat(path)
	var bitRate int
	if err == nil && duration > 0 {
		bitRate = int(float64(fileInfo.Size()*8) / duration / 1000)
	}

	return &AudioProps{
		Duration:   duration,
		SampleRate: sampleRate,
		BitDepth:   bitDepth,
		Channels:   channels,
		BitRate:    bitRate,
	}, nil
}
