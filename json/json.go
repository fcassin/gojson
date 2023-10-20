package gojson

import (
	"fmt"
	"regexp"
)

type ParserState int

const (
	ReadingName ParserState = iota
	ReadingNumber
	SeekingBaseStructure
	SeekingName
	SeekingNameValueSeparator
	SeekingValue
	SeekingSeparator
	// object handling
	SeekingMembers
	// TODO: Refactor, refactor, refactor
	SeekingMember
	// array handling
	SeekingElements
	SeekingElement
)

/* type JSONValue int

const (
	Object JSONValue = iota
	Array
	String
	Number
	True
	False
	Null
) */

const (
	ZeroPattern     string = "0"
	OneNinePattern  string = "[1-9]"
	MinusPattern    string = "-"
	SignPattern     string = "[+\\-]"
	FractionPattern string = "\\."
	ExponentPattern string = "[eE+\\-]"
)

var (
	IntegerPattern *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("%s|%s|%s",
		ZeroPattern, OneNinePattern, MinusPattern))
	NumberPattern *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s",
		ZeroPattern, OneNinePattern, MinusPattern, FractionPattern, ExponentPattern))
)
