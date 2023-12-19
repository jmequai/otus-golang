package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

const (
	isEscaped stateType = iota
	isCaptured
	isMultiplied
)

type stateType int8

type stateRune struct {
	t stateType
	r rune
}

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b strings.Builder
	var err error
	var st *stateRune

	for _, r := range s {
		st, err = next(st, r, &b)

		if err != nil {
			return "", err
		}
	}

	if st != nil {
		if st.t == 0 || st.t == 2 {
			return "", ErrInvalidString
		}

		b.WriteRune(st.r)
	}

	return b.String(), nil
}

func next(prev *stateRune, rn rune, b *strings.Builder) (*stateRune, error) {
	var stateTp stateType

	switch {
	case rn == '\\':
		stateTp = isEscaped
	case rn >= '0' && rn <= '9':
		stateTp = isMultiplied
	default:
		stateTp = isCaptured
	}

	state := &stateRune{stateTp, rn}

	if prev == nil {
		return state, nil
	}

	if prev.t == 0 {
		if state.t == 1 {
			return nil, ErrInvalidString
		}

		state.t = 1

		return state, nil
	}

	if prev.t == 1 {
		if state.t == 2 {
			m, e := strconv.Atoi(string(state.r))

			if e != nil {
				return nil, ErrInvalidString
			}

			s := strings.Repeat(string(prev.r), m)

			b.WriteString(s)

			return nil, nil
		}

		b.WriteRune(prev.r)

		return state, nil
	}

	return nil, ErrInvalidString
}
