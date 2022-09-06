package masker

import (
	"math"
	"strings"
)

// Masker masks the data
type Masker struct{}

// New is the constructor for Masker
func New() *Masker {
	return &Masker{}
}

// https://github.com/ggwhite/go-masker
func (m *Masker) overlay(str string, overlay string, start int, end int) (res string) {
	r := []rune(str)
	l := len(r)

	if l == 0 {
		return ""
	}

	if start < 0 {
		start = 0
	}
	if start > l {
		start = l
	}
	if end < 0 {
		end = 0
	}
	if end > l {
		end = l
	}
	if start > end {
		tmp := start
		start = end
		end = tmp
	}

	res = ""
	res += string(r[:start])
	res += overlay
	res += string(r[end:])
	return res
}

// Name masks the name
func (m *Masker) Name(i string) string {
	ri := []rune(i)
	l := len(ri)
	if l == 0 {
		return ""
	}

	ov := strings.Repeat("*", len(ri[1:]))
	return m.overlay(i, ov, 1, l)
}

// Email masks the email
func (m *Masker) Email(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}

	tmp := strings.Split(i, "@")
	addr := tmp[0]
	domain := tmp[1]

	start := 3
	if l < 3 {
		start = l
	}

	addr = m.overlay(addr, "******", start, l)

	return addr + "@" + domain
}

// Mobile masks the mobile
func (m *Masker) Mobile(i string) string {
	l := float64(len([]rune(i)))
	if l == 0 {
		return ""
	}

	mid := math.Floor(l / 2)

	start := int(mid - 2)
	if start < 0 {
		start = 0
	}

	end := int(mid + 2)
	if end > int(l) {
		end = int(l)
	}

	return m.overlay(i, "****", start, end)
}
