package command

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/12yanogden/errors"
	"github.com/12yanogden/str"
	"github.com/12yanogden/strslices"
)

type ParsedInput struct {
	Keys			[]string
	ValueType		int
	Value 			interface{}
}

const (
	IS_OPTION = iota
	VALUE_TYPE
	IS_OPTIONAL
	DEFAULT
	DEFAULT_TYPE
)

const (
	BOOLEAN_TYPE = iota
	STRING_TYPE
	ARRAY_TYPE
)

func ParseInput(rules []Rule) ([]ParsedInput, error) {
	inputs := os.Args[1:]
	parsedRules := parseRules(rules)
	parsedInputs := []ParsedInput{}

	for len(inputs) > 0 {
		// Get key from in
		key := getKeyFromInput(&inputs)

		// Get parsed rule by key
		rule := getParsedRuleByKey(&key, &parsedRules)

		// Parse and append input
		append(parsedInputs, parseInput(&inputs, &rule))
	}

	return parsedInputs, nil
}


func isOption(def *string) bool {
	isOption := string((*def)[0]) == "-"

	if isOption {
		*def = (*def)[1:]
	}

	return isOption
}

func match(bytes *[]byte, pattern *regexp.Regexp) []byte {
	match := pattern.Find(*bytes)

	if match != nil {
		trimBytesLeft(bytes, len(match))
	}

	return match
}

func trimBytesLeft(bytes *[]byte, len int) {
	*bytes = (*bytes)[len:]
}

func trimLeftAlphaNums(str *string) string {
	trimmed := ""
	i := 0

	// Collect substring to be trimmed
	for i := range *str {
		char := string((*str)[i])

		if str.IsAlphaNum(char) {
			trimmed += char
		}
	}

	// Trim string given
	*str = (*str)[i:]

	// Return trimmed substring
	return trimmed
}

func getKeyFromInput(inputs *[]string) string {
	key := ""

	for i, input := range *inputs {
		j := 0

		for j := range input {
			char := string(input[j])

			if (i == 0 && char == "-") {
				if ((j + 1) < len(input) && string(input[j + 1]) == "-") {
					j++
				}

				continue

			} else if str.IsAlphaNum(char) {
				key += char

			} else {
				break
			}
		}

		(*inputs)[i] = input[j:]
	}

	return key
}

func getParsedRuleByKey(key *string, rules *[]ParsedRule) ParsedRule {
	for _, rule := range *rules {
		for _, ruleKey := range rule.Keys {
			if *key == ruleKey {
				return rule
			}
		}
	}

	errors.Scream("no rule found for key: " + *key)
}

func isOption(arg string) (bool, error) {
	isOption, err := regexp.MatchString("--[a-zA-Z0-9=\?\*]+", arg)
	peek(err)

	return isOption, nil
}

