package json

import (
	"fmt"
	"regexp"
)

type parserState int

const (
	readingName parserState = iota
	readingNumber
	seekingBaseStructure
	seekingName
	seekingNameValueSeparator
	seekingValue
	seekingSeparator
	// object handling
	seekingMembers
	// TODO: Refactor, refactor, refactor
	seekingMember
	// array handling
	seekingElements
	seekingElement
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
	zeroPattern     string = "0"
	oneNinePattern  string = "[1-9]"
	minusPattern    string = "-"
	signPattern     string = "[+\\-]"
	fractionPattern string = "\\."
	exponentPattern string = "[eE+\\-]"
)

var (
	integerPattern *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("%s|%s|%s",
		zeroPattern, oneNinePattern, minusPattern))
	numberPattern *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s",
		zeroPattern, oneNinePattern, minusPattern, fractionPattern, exponentPattern))
)
