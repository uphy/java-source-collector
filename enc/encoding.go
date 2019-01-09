package enc

import (
	"errors"
	"io"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

var (
	ErrNotFound = errors.New("encoding not found")
)

func NewReader(reader io.Reader, encodingName string) (io.Reader, error) {
	encoding, err := FindEncoding(encodingName)
	if err != nil {
		return nil, err
	}
	if encoding == nil {
		return reader, nil
	}
	return encoding.NewDecoder().Reader(reader), nil
}

func NewWriter(writer io.Writer, encodingName string) (io.Writer, error) {
	encoding, err := FindEncoding(encodingName)
	if err != nil {
		return nil, err
	}
	if encoding == nil {
		return writer, nil
	}
	return encoding.NewEncoder().Writer(writer), nil
}

func FindEncoding(name string) (encoding.Encoding, error) {
	switch strings.ToLower(name) {
	case "sjis":
		fallthrough
	case "shift_jis":
		return japanese.ShiftJIS, nil
	case "eucjp":
		return japanese.EUCJP, nil
	case "utf8":
		fallthrough
	case "utf-8":
		return nil, nil
	}
	return nil, ErrNotFound
}
