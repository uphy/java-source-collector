package java

import (
	"io"
	"text/scanner"
	"unicode"
)

func DetectPackage(r io.Reader) (string, error) {
	var s scanner.Scanner
	s.IsIdentRune = func(ch rune, i int) bool {
		return ch == '.' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
	}
	s.Init(r)
	isPackage := false
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		txt := s.TokenText()
		if isPackage {
			s.Scan()
			return txt, nil
		}
		if txt == "package" {
			isPackage = true
		}
	}
	return "", nil
}
