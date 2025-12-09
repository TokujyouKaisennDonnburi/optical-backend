package calendar

import (
	"strings"
	"unicode/utf8"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
)

type Color string

var NilColor Color = ""

func NewColor(color string) (Color, error) {
	// '#' + 6文字
	length := 7
	if utf8.RuneCountInString(color) != length {
		return "", apperr.ValidationError("invalid color length")
	}
	if color[0] != '#' {
		return "", apperr.ValidationError("invalid color format")
	}
	// 16進数チェック
	for i, ch := range color {
		if i == 0 {
			continue
		}
		if !strings.ContainsRune("0123456789abcdefABCDEF", ch) {
			return "", apperr.ValidationError("invalid color format")
		}
	}
	return Color(color), nil
}
