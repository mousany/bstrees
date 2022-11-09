package console

import (
	"bstrees/pkg/trait/number"
	"bufio"
)

func Read[T number.Integer](istream *bufio.Reader) (T, error) {
	res, sign := T(0), 1
	readed := false
	c, err := istream.ReadByte()
	for ; err == nil && (c < '0' || c > '9'); c, err = istream.ReadByte() {
		if c == '-' {
			sign = -1
		}
	}
	for ; err == nil && c >= '0' && c <= '9'; c, err = istream.ReadByte() {
		readed = true
		res = res<<3 + res<<1 + T(c-'0')
	}
	if !readed {
		return 0, err
	}
	return res * T(sign), err
}
