// Package xmath has integer helpers missing from the standard library.
package xmath

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func Abs[T Signed](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

func Sign[T Signed](n T) T {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	}
	return 0
}

// Mod returns n mod m in the range [0, m), unlike Go's % which can be negative.
func Mod[T Integer](n, m T) T {
	return (n%m + m) % m
}

func GCD[T Integer](a, b T) T {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func LCM[T Integer](a, b T) T {
	return a / GCD(a, b) * b
}

// Sum adds up a slice of numbers.
func Sum[T Integer](xs []T) T {
	var s T
	for _, x := range xs {
		s += x
	}
	return s
}
