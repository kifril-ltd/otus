package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const (
	EscapeSymbol = rune('\\')
	ZeroRune     = rune(0)
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(in string) (string, error) {
	if in == "" {
		return "", nil
	}

	var (
		res      strings.Builder
		escape   bool
		prevRune rune
	)

	for _, curRune := range in {
		switch {
		case escape:
			if !unicode.IsDigit(curRune) && curRune != EscapeSymbol {
				return "", ErrInvalidString
			}
			prevRune = curRune
			escape = false
		case curRune == EscapeSymbol:
			if prevRune != ZeroRune {
				res.WriteRune(prevRune)
			}
			prevRune = ZeroRune
			escape = true
		case unicode.IsDigit(curRune):
			if prevRune == ZeroRune {
				return "", ErrInvalidString
			}

			count, err := strconv.Atoi(string(curRune))
			if err != nil {
				return "", err
			}

			res.WriteString(strings.Repeat(string(prevRune), count))
			prevRune = ZeroRune
		default:
			if prevRune != ZeroRune {
				res.WriteRune(prevRune)
			}
			prevRune = curRune
		}
	}

	if escape {
		return "", ErrInvalidString
	}

	if prevRune != ZeroRune {
		res.WriteRune(prevRune)
	}

	return res.String(), nil
}
