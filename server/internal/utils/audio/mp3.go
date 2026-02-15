package audio

import (
	"encoding/binary"
	"errors"
	"os"
)

// ParseMp3Props reads audio properties from an MP3 file
// Supports CBR and VBR (Xing/Info header)
func ParseMp3Props(path string) (*AudioProps, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 1. Skip ID3v2 tag
	header := make([]byte, 10)
	if _, err := f.Read(header); err != nil {
		return nil, err
	}
	
	startOffset := int64(0)
	if string(header[0:3]) == "ID3" {
		// Calculate tag size
		// The ID3v2 tag size is encoded with 4 bytes where the most significant bit (bit 7) is set to zero in every byte
		size := (int(header[6]&0x7F) << 21) | (int(header[7]&0x7F) << 14) | (int(header[8]&0x7F) << 7) | int(header[9]&0x7F)
		startOffset = int64(size + 10)
		f.Seek(startOffset, 0)
	} else {
		f.Seek(0, 0)
	}

	// 2. Find first MP3 Frame
	// Scan for 0xFF 0xE0 (Sync word, 11 bits set to 1)
	// Increase buffer size to 16KB to handle larger gaps/garbage
	buf := make([]byte, 16384)
	n, err := f.Read(buf)
	if err != nil {
		return nil, err
	}

	var frameHeader uint32
	found := false
	frameOffset := 0

	for i := 0; i < n-4; i++ {
		// Sync word: 11 bits set to 1 (0xFF, and top 3 bits of next byte)
		if buf[i] == 0xFF && (buf[i+1]&0xE0) == 0xE0 {
			h := binary.BigEndian.Uint32(buf[i : i+4])
			version := (h >> 19) & 3
			layer := (h >> 17) & 3
			bitrateIdx := (h >> 12) & 0xF
			sampleRateIdx := (h >> 10) & 3
			
			// Basic validity check
			if version != 1 && layer != 0 && bitrateIdx != 0 && bitrateIdx != 15 && sampleRateIdx != 3 {
				frameHeader = h
				found = true
				frameOffset = i
				break
			}
		}
	}

	if !found {
		return nil, errors.New("MP3 frame sync not found")
	}

	// Parse Header
	version := (frameHeader >> 19) & 3 // 3=MPEG1, 2=MPEG2, 0=MPEG2.5
	layer := (frameHeader >> 17) & 3   // 3=Layer1, 2=Layer2, 1=Layer3
	bitrateIdx := (frameHeader >> 12) & 0xF
	sampleRateIdx := (frameHeader >> 10) & 3
	channelMode := (frameHeader >> 6) & 3
	// padding := (frameHeader >> 9) & 1 // Unused for now

	// Tables
	// Bitrates (kbps)
	var kbps int
	if version == 3 && layer == 1 { // MPEG1 Layer3
		bitrates := []int{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0}
		kbps = bitrates[bitrateIdx]
	} else if version == 3 && layer == 2 { // MPEG1 Layer2
		bitrates := []int{0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 0}
		kbps = bitrates[bitrateIdx]
	} else if version == 3 && layer == 3 { // MPEG1 Layer1
		bitrates := []int{0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 0}
		kbps = bitrates[bitrateIdx]
	} else { // MPEG2/2.5
		// Layer3
		bitrates := []int{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0}
		kbps = bitrates[bitrateIdx]
	}
	
	// Sample Rates
	var sampleRate int
	if version == 3 { // MPEG1
		sampleRates := []int{44100, 48000, 32000, 0}
		sampleRate = sampleRates[sampleRateIdx]
	} else if version == 2 { // MPEG2
		sampleRates := []int{22050, 24000, 16000, 0}
		sampleRate = sampleRates[sampleRateIdx]
	} else if version == 0 { // MPEG2.5
		sampleRates := []int{11025, 12000, 8000, 0}
		sampleRate = sampleRates[sampleRateIdx]
	}

	channels := 2
	if channelMode == 3 { // Mono
		channels = 1
	}

	// Check for Xing/Info header (VBR)
	// Side info size calculation
	sideInfoSize := 0
	if version == 3 { // MPEG1
		if channelMode == 3 { // Mono
			sideInfoSize = 17
		} else { // Stereo/Joint/Dual
			sideInfoSize = 32
		}
	} else { // MPEG2/2.5
		if channelMode == 3 { // Mono
			sideInfoSize = 9
		} else { // Stereo/Joint/Dual
			sideInfoSize = 17
		}
	}
	
	// Xing/Info header offset relative to frame start
	xingOffset := frameOffset + 4 + sideInfoSize
	
	if xingOffset+12 <= n {
		tag := string(buf[xingOffset : xingOffset+4])
		if tag == "Xing" || tag == "Info" {
			// Read flags
			flags := binary.BigEndian.Uint32(buf[xingOffset+4 : xingOffset+8])
			// If FRAMES flag is set (0x0001)
			if flags&1 != 0 {
				// Read frames count
				frames := binary.BigEndian.Uint32(buf[xingOffset+8 : xingOffset+12])
				
				// Samples per frame
				samplesPerFrame := 1152 // Layer 2/3
				if version == 3 && layer == 3 { // Layer 1
					samplesPerFrame = 384
				} else if version != 3 && layer == 1 { // MPEG2 Layer 3
					samplesPerFrame = 576
				}

				if sampleRate > 0 {
					duration := float64(frames) * float64(samplesPerFrame) / float64(sampleRate)
					
					// Calculate average bitrate
					fileInfo, _ := f.Stat()
					fileSize := fileInfo.Size()
					avgBitrate := 0
					if duration > 0 {
						avgBitrate = int(float64(fileSize-startOffset) * 8 / duration / 1000)
					}
					
					return &AudioProps{
						Duration:   duration,
						SampleRate: sampleRate,
						Channels:   channels,
						BitRate:    avgBitrate,
						BitDepth:   16, // MP3 decodes to 16-bit usually
					}, nil
				}
			}
		}
	}

	// Fallback to CBR
	fileInfo, _ := f.Stat()
	fileSize := fileInfo.Size()
	var duration float64
	if kbps > 0 {
		// Calculate average frame size including padding for more accuracy?
		// FrameSize = (SamplesPerFrame / 8 * BitRate) / SampleRate + Padding
		// But for total duration, file size / bitrate is good enough for CBR
		duration = float64(fileSize-startOffset) * 8 / float64(kbps*1000)
	}

	return &AudioProps{
		Duration:   duration,
		SampleRate: sampleRate,
		Channels:   channels,
		BitRate:    kbps,
		BitDepth:   16,
	}, nil
}
