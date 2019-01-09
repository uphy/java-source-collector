package java

import (
	"io"
	"text/scanner"
)

func DetectPackage(r io.Reader) (string, error) {
	var s scanner.Scanner
	s.Init(r)
	isPackage := false
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		txt := s.TokenText()
		if isPackage {
			return txt, nil
		}
		if txt == "package" {
			isPackage = true
		}
	}
	return "", nil
}
