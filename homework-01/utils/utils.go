package utils

import (
	"fmt"
	"strings"
	"time"
	"github.com/prometheus/client_golang/prometheus"
)

const base = uint32(62)
const strLen = 7
const mapString = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func uintToByte(i uint32) (byte, error) {
	if int(i) >= len(mapString) {
		return 'a', fmt.Errorf("too big uint")
	}

	return mapString[i], nil
}

func byteToUint(u byte) (uint32, error) {
	i := strings.Index(mapString, string(u))

	if i == -1 {
		return 0, fmt.Errorf("bad byte on input")
	}

	return uint32(i), nil
}

func UintToString(i uint32) (string, error) {
	s := ""

	for j := 0; j < strLen; j++ {
		mod := i % base
		b, err := uintToByte(mod)

		if err != nil {
			return "", err

		}

		s += string(b)
		i /= base
	}

	return s, nil
}

func StringToUint(s string) (uint32, error) {
	if len(s) != strLen {
		return 0, fmt.Errorf("wrong string len")
	}

	i := uint32(0)
	sLen := len(s)

	for j := sLen - 1; j >= 0; j-- {
		b := s[j]
		u, err := byteToUint(b)

		if err != nil {
			return 0, err
		}

		i *= base
		i += u
	}

	return i, nil
}

func MeasureTime(action func(), counter prometheus.Counter, summary prometheus.Summary) {
	start := time.Now()
	action()
	elapsed := time.Since(start).Microseconds()

	counter.Inc()
	summary.Observe(float64(elapsed))
}
