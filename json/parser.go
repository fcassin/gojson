package gojson

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func parseObject(l lexer, o reflect.Value) {
	var state ParserState = SeekingName
	var name string

	for {
		token, err := l.next()
		if err != nil {
			fmt.Printf("err: %s\n", err)
			break
		}

		switch state {
		case SeekingName:
			if token.kind != String {
				fmt.Printf("expected to find a name! found %v\n", token)
			} else {
				name = token.token
				state = SeekingNameValueSeparator
			}
		case SeekingNameValueSeparator:
			if token.kind != Colon {
				fmt.Println("expected to find a colon!")
			} else {
				state = SeekingValue
			}
		case SeekingValue:
			// TODO: Handle all case
			switch token.kind {
			case ArrayOpen:
				var structIndex []int
				var structField reflect.StructField

				for i := 0; i < o.Type().NumField(); i++ {
					if o.Type().Field(i).Tag.Get("json") == name {
						structIndex = append(structIndex, i)
						structField = o.Type().Field(i)
						break
					}

					// TODO: Error if the field could not be found
				}

				workingSlice := reflect.New(structField.Type).Elem()
				parseArray(l, workingSlice)
				o.FieldByIndex(structIndex).Set(workingSlice)

				state = SeekingSeparator
			case Number:
				// TODO: Refactor into populate
				var structIndex []int

				for i := 0; i < o.Type().Elem().NumField(); i++ {
					if o.Type().Elem().Field(i).Tag.Get("json") == name {
						structIndex = append(structIndex, i)
						break
					}

					// TODO: Error if the field could not be found
				}

				value, err := strconv.ParseFloat(token.token, 64)
				if err != nil {
					fmt.Printf("err: %s\n", err)
				}

				o.Elem().FieldByIndex(structIndex).Set(reflect.ValueOf(value))

				state = SeekingSeparator
			}
		case SeekingSeparator:
			switch token.kind {
			case ObjectClose:
				return
			case Comma:
				state = SeekingName
			}
		}
	}
}

func parseArray(l lexer, o reflect.Value) {
	var state ParserState = SeekingElements

	for {
		token, err := l.next()
		if err != nil {
			fmt.Printf("err: %s\n", err)
			break
		}

		switch state {
		case SeekingElements:
			switch token.kind {
			case ObjectOpen:
				workingElem := reflect.New(o.Type().Elem())
				parseObject(l, workingElem)
				o.Set(reflect.Append(o, workingElem.Elem()))
			case ArrayClose:
				return
			case Comma:
				continue
			default:
				fmt.Printf("missing handling for %d\n", token.kind)
			}
		}
	}
}

func Unmarshall(rawBytes []byte, element any) {
	var rootValue reflect.Value = reflect.ValueOf(element).Elem()

	if !rootValue.CanSet() {
		// TODO: Return an error here
		fmt.Println("cannot set root value!!!")
		fmt.Println("cannot set root value!!!")
		fmt.Println("cannot set root value!!!")
	}

	var jsonLexer lexer = NewLexer(rawBytes)

	for {
		t, err := jsonLexer.next()

		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println(err)
			break
		}

		if t.kind == ObjectOpen {
			parseObject(jsonLexer, rootValue)
		}
	}
}
