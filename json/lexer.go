package json

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unicode/utf8"
)

type tokenKind int

const (
	objectOpen tokenKind = iota
	objectClose
	arrayOpen
	arrayClose

	stringToken
	number

	colon
	comma

	unknown
)

type token struct {
	kind  tokenKind
	token string
}

type lexer struct {
	// TODO(fca): Hide details into a jsonscanner
	reader io.RuneScanner
	token  rune
}

func newLexer(rawBytes []byte) lexer {
	return lexer{
		reader: bytes.NewReader(rawBytes),
	}
}

func (l *lexer) nextString() (t token, e error) {
	var r rune
	var value string
	for {
		r, _, e = l.reader.ReadRune()

		if e != nil {
			return
		}

		if r == '"' {
			return token{stringToken, value}, nil
		}

		value = value + string(r)
	}
}

func (l *lexer) nextNumber() (t token, e error) {
	var r rune
	var s int
	var value string

	for {
		r, s, e = l.reader.ReadRune()

		if e != nil {
			return
		}

		if isWhitespace(r) {
			return token{number, value}, nil
		}

		if r == ',' {
			l.reader.UnreadRune()
			return token{number, value}, nil
		}

		runeBytes := make([]byte, s)
		utf8.EncodeRune(runeBytes, r)

		if numberPattern.Match(runeBytes) {
			value = value + string(r)
		} else {
			return t, errors.New(fmt.Sprintf("unexpected rune while reading number: %c", r))
		}
	}
}

func (l *lexer) next() (t token, e error) {
	var r rune
	for {
		r, _, e = l.reader.ReadRune()

		if e != nil {
			return
		}

		// TODO(fca): Move isWhitespace here, package this all up
		if !isWhitespace(r) {
			break
		}
	}

	switch r {
	case '{':
		return token{objectOpen, "{"}, nil
	case '}':
		return token{objectClose, "}"}, nil
	case '[':
		return token{arrayOpen, "["}, nil
	case ']':
		return token{arrayClose, "]"}, nil
	case '"':
		return l.nextString()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.reader.UnreadRune()
		return l.nextNumber()
	// TODO(fca): Handle ':'
	case ':':
		return token{colon, ":"}, nil
	case ',':
		return token{comma, ","}, nil
	default:
		return token{unknown, string(r)}, nil
	}
}

func iterate(rawBytes []byte, element any) {
	var elementValue reflect.Value = reflect.ValueOf(&element).Elem()
	if !elementValue.CanSet() {
		// TODO: Return an error here
		fmt.Println("cannot set element value!!!")
		fmt.Println("cannot set element value!!!")
		fmt.Println("cannot set element value!!!")
	}

	var jsonlexer lexer = newLexer(rawBytes)

	for {
		t, err := jsonlexer.next()

		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println(err)
			break
		}

		if t.kind != unknown {
			fmt.Println(t)
		}
	}
}

func isWhitespace(r rune) bool {
	return r == '\n' || r == ' '
}
