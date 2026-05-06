package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
)

// ParseWavProps reads audio properties from a WAV file.
func ParseWavProps(path string) (*AudioProps, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	header := make([]byte, 12)
	if _, err := io.ReadFull(f, header); err != nil {
		return nil, err
	}
	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return nil, errors.New("not a valid RIFF/WAVE file")
	}

	var channels, sampleRate, byteRate, bitsPerSample int
	var dataSize int64
	var headerSize int64 = 12

	for {
		chunkHeader := make([]byte, 8)
		_, err := io.ReadFull(f, chunkHeader)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		headerSize += 8

		chunkID := string(chunkHeader[0:4])
		chunkSize := binary.LittleEndian.Uint32(chunkHeader[4:8])
		padding := int64(chunkSize % 2)

		switch chunkID {
		case "fmt ":
			if chunkSize < 16 {
				return nil, errors.New("invalid fmt chunk size")
			}
			fmtData := make([]byte, chunkSize)
			if _, err := io.ReadFull(f, fmtData); err != nil {
				return nil, err
			}
			headerSize += int64(chunkSize)
			if padding > 0 {
				f.Seek(padding, io.SeekCurrent)
				headerSize += padding
			}
			channels = int(binary.LittleEndian.Uint16(fmtData[2:4]))
			sampleRate = int(binary.LittleEndian.Uint32(fmtData[4:8]))
			byteRate = int(binary.LittleEndian.Uint32(fmtData[8:12]))
			bitsPerSample = int(binary.LittleEndian.Uint16(fmtData[14:16]))

		case "data":
			dataSize = int64(chunkSize)
			f.Seek(int64(chunkSize)+padding, io.SeekCurrent)
			headerSize += int64(chunkSize) + padding

		default:
			f.Seek(int64(chunkSize)+padding, io.SeekCurrent)
			headerSize += int64(chunkSize) + padding
		}
	}

	if sampleRate == 0 || byteRate == 0 {
		return nil, errors.New("missing or invalid fmt chunk")
	}

	if dataSize == 0 {
		fileInfo, err := f.Stat()
		if err != nil {
			return nil, err
		}
		dataSize = fileInfo.Size() - headerSize
		if dataSize < 0 {
			dataSize = 0
		}
	}

	duration := float64(dataSize) / float64(byteRate)
	bitRate := byteRate * 8 / 1000

	return &AudioProps{
		Duration:   duration,
		SampleRate: sampleRate,
		BitDepth:   bitsPerSample,
		Channels:   channels,
		BitRate:    bitRate,
	}, nil
}

// WavMetadata implements tag.Metadata for WAV files.
type WavMetadata struct {
	title       string
	artist      string
	album       string
	albumArtist string
	composer    string
	year        int
	genre       string
	track       int
	trackTotal  int
	disc        int
	discTotal   int
	picture     *tag.Picture
	lyrics      string
	comment     string
	raw         map[string]interface{}
}

func (m *WavMetadata) Format() tag.Format     { return "WAV/RIFF" }
func (m *WavMetadata) FileType() tag.FileType { return "WAV" }
func (m *WavMetadata) Title() string          { return m.title }
func (m *WavMetadata) Album() string          { return m.album }
func (m *WavMetadata) Artist() string         { return m.artist }
func (m *WavMetadata) AlbumArtist() string    { return m.albumArtist }
func (m *WavMetadata) Composer() string       { return m.composer }
func (m *WavMetadata) Year() int              { return m.year }
func (m *WavMetadata) Genre() string          { return m.genre }
func (m *WavMetadata) Track() (int, int)      { return m.track, m.trackTotal }
func (m *WavMetadata) Disc() (int, int)       { return m.disc, m.discTotal }
func (m *WavMetadata) Picture() *tag.Picture  { return m.picture }
func (m *WavMetadata) Lyrics() string         { return m.lyrics }
func (m *WavMetadata) Comment() string        { return m.comment }
func (m *WavMetadata) Raw() map[string]interface{} {
	if m.raw == nil {
		return map[string]interface{}{}
	}
	return m.raw
}

func (m *WavMetadata) applyInfoTag(id, value string) {
	switch id {
	case "INAM":
		if m.title == "" {
			m.title = value
		}
	case "IART":
		if m.artist == "" {
			m.artist = value
		}
	case "IPRD":
		if m.album == "" {
			m.album = value
		}
	case "ICMT":
		if m.comment == "" {
			m.comment = value
		}
	case "ICRD":
		if m.year == 0 {
			if y, err := strconv.Atoi(value); err == nil {
				m.year = y
			}
		}
	case "IGNR":
		if m.genre == "" {
			m.genre = value
		}
	case "ITRK":
		if m.track == 0 {
			if t, err := strconv.Atoi(value); err == nil {
				m.track = t
			}
		}
	case "IPRT":
		if m.disc == 0 {
			if d, err := strconv.Atoi(value); err == nil {
				m.disc = d
			}
		}
	}
}

func mergeMetadata(dst *WavMetadata, src tag.Metadata) {
	if dst.title == "" && src.Title() != "" {
		dst.title = src.Title()
	}
	if dst.artist == "" && src.Artist() != "" {
		dst.artist = src.Artist()
	}
	if dst.album == "" && src.Album() != "" {
		dst.album = src.Album()
	}
	if dst.albumArtist == "" && src.AlbumArtist() != "" {
		dst.albumArtist = src.AlbumArtist()
	}
	if dst.composer == "" && src.Composer() != "" {
		dst.composer = src.Composer()
	}
	if dst.year == 0 && src.Year() > 0 {
		dst.year = src.Year()
	}
	if dst.genre == "" && src.Genre() != "" {
		dst.genre = src.Genre()
	}
	t, tt := src.Track()
	if dst.track == 0 && t > 0 {
		dst.track = t
		dst.trackTotal = tt
	}
	d, dt := src.Disc()
	if dst.disc == 0 && d > 0 {
		dst.disc = d
		dst.discTotal = dt
	}
	if dst.picture == nil && src.Picture() != nil {
		dst.picture = src.Picture()
	}
	if dst.lyrics == "" && src.Lyrics() != "" {
		dst.lyrics = src.Lyrics()
	}
	if dst.comment == "" && src.Comment() != "" {
		dst.comment = src.Comment()
	}
}

func isValidChunkID(id string) bool {
	if len(id) != 4 {
		return false
	}
	for _, b := range []byte(id) {
		if b < 0x20 || b > 0x7E {
			return false
		}
	}
	return true
}

// ParseWavMetadata parses WAV metadata from RIFF INFO chunks, id3 chunks, and trailing ID3 tags.
func ParseWavMetadata(path string) (tag.Metadata, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := info.Size()

	header := make([]byte, 12)
	if _, err := io.ReadFull(f, header); err != nil {
		return nil, err
	}
	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return nil, errors.New("not a valid RIFF/WAVE file")
	}

	meta := &WavMetadata{raw: make(map[string]interface{})}

	// Scan RIFF chunks
	for {
		chunkPos, _ := f.Seek(0, io.SeekCurrent)
		if chunkPos >= fileSize {
			break
		}

		chunkHeader := make([]byte, 8)
		if _, err := io.ReadFull(f, chunkHeader); err != nil {
			break
		}

		chunkID := string(chunkHeader[0:4])
		chunkSize := binary.LittleEndian.Uint32(chunkHeader[4:8])
		padding := int64(chunkSize % 2)
		chunkEnd := chunkPos + 8 + int64(chunkSize) + padding
		if chunkEnd > fileSize {
			chunkEnd = fileSize
		}

		// Safety: non-printable chunkID means we've hit trailing non-RIFF data (e.g., ID3)
		if !isValidChunkID(chunkID) {
			break
		}

		switch chunkID {
		case "LIST":
			listType := make([]byte, 4)
			if _, err := io.ReadFull(f, listType); err == nil && string(listType) == "INFO" {
				infoEnd := chunkPos + 8 + int64(chunkSize)
				for {
					infoPos, _ := f.Seek(0, io.SeekCurrent)
					if infoPos >= infoEnd-8 {
						break
					}
					infoID := make([]byte, 4)
					if _, err := io.ReadFull(f, infoID); err != nil {
						break
					}
					infoSizeBuf := make([]byte, 4)
					if _, err := io.ReadFull(f, infoSizeBuf); err != nil {
						break
					}
					infoSize := binary.LittleEndian.Uint32(infoSizeBuf)
					infoDataEnd := infoPos + 8 + int64(infoSize)
					if infoDataEnd > infoEnd {
						break
					}
					infoData := make([]byte, infoSize)
					if _, err := io.ReadFull(f, infoData); err != nil {
						break
					}
					value := strings.TrimRight(string(infoData), "\x00")
					meta.raw[string(infoID)] = value
					meta.applyInfoTag(string(infoID), value)
				}
			}
			f.Seek(chunkEnd, io.SeekStart)

		case "id3 ":
			id3Data := make([]byte, chunkSize)
			if _, err := io.ReadFull(f, id3Data); err == nil {
				if m, err := tag.ReadID3v2Tags(bytes.NewReader(id3Data)); err == nil {
					mergeMetadata(meta, m)
				}
			}
			f.Seek(chunkEnd, io.SeekStart)

		default:
			f.Seek(chunkEnd, io.SeekStart)
		}
	}

	// Trailing ID3v1 (always last 128 bytes)
	if fileSize >= 128 {
		buf := make([]byte, 128)
		if _, err := f.ReadAt(buf, fileSize-128); err == nil {
			if string(buf[0:3]) == "TAG" {
				if m, err := tag.ReadID3v1Tags(bytes.NewReader(buf)); err == nil {
					mergeMetadata(meta, m)
				}
			}
		}
	}

	// Trailing ID3v2 — search last 64KB for "ID3" marker
	const maxSearch = 64 * 1024
	searchStart := fileSize - maxSearch
	if searchStart < 0 {
		searchStart = 0
	}
	searchBuf := make([]byte, fileSize-searchStart)
	if _, err := f.ReadAt(searchBuf, searchStart); err == nil {
		if idx := bytes.Index(searchBuf, []byte("ID3")); idx >= 0 {
			id3Start := searchStart + int64(idx)
			id3Buf := make([]byte, fileSize-id3Start)
			if _, err := f.ReadAt(id3Buf, id3Start); err == nil {
				if m, err := tag.ReadID3v2Tags(bytes.NewReader(id3Buf)); err == nil {
					mergeMetadata(meta, m)
				}
			}
		}
	}

	return meta, nil
}
