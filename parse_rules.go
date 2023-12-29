package command

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/12yanogden/pepr/internal/errors"
	"github.com/12yanogden/pepr/internal/str"
	"github.com/12yanogden/pepr/internal/strslices"
)

type Rule struct {
	Def  string
	Desc string
}

type Key struct {
	Name		string
	IsGroupable	bool
}

type ParsedRule struct {
	Keys		[]Key
	Rule		*Rule
	Tokens 		map[int]string	// type => value
}

func parseRules(rules []Rule) []ParsedRule {
	for _, rule := range rules {
		defBytes := []byte(rule.Def)
		validateDef(defBytes)
		tokens := map[int]string{}
		parsedKeys := parseKeys(defBytes, &tokens)
		bToS := strconv.FormatBool	

		// Parse args
		if tokens[IS_OPTION] == "true" {
			if hasVariable := (match(&defBytes, regexp.MustCompile("^=")) != nil); hasVariable {
				if hasArray := (match(&defBytes, regexp.MustCompile("^[")) != nil); hasArray {
					defaultArray := match(&defBytes, regexp.MustCompile("[a-zA-Z0-9, ]+"))

					tokens[VALUE_TYPE] = "array"
					
					if defaultArray != nil {
						tokens[DEFAULT] = strslices.ToCSV()
					}

					
				} else {
					tokens[DEFAULT_TYPE] = "string"
					tokens[DEFAULT] = string(defBytes)
				}
			}
		}
		tokens[VALUE_TYPE] = 
	}

	// Validate rules


	return []ParsedRule{}
}

func validateDef(def []byte) {
	if len(def) == 0 {
		errors.Scream("rule definition cannot be empty")
	} else if def[0] == '-' {
		validateOption(def)
	} else {
		validateArg(def)
	}
}

func validateOption(def []byte) {
	if !(regexp.MustCompile(
		`^(
			(
				(
					\-[[:alnum:]]				// groupable option
				)|(
					\-\-[[:alnum:]]+			// non-groupable option
				)
			)\|									// 0 to many with '|' operator
		)*(
			(
				\-[[:alnum:]]					// groupable option
			)|(
				\-\-[[:alnum:]]+				// non-groupable option
			)
		)(										// at least 1 key required
			=(
				(
					'[[:alnum:]]'
				)|(
					\[(
						'(
							(
								[[:alnum:]]| 	// 0 to many defaults with ,
							)+', ?'
						)*(
							(
								[[:alnum:]]| 	// at least 1 default
							)+
						)?'
					)?\]						// ? allows for empty []
				)
			)?									// ? allows for having a '=' without defaults
		)?$`									// ? allows for boolean options
		// Option keys
		`^(` + 
			`(` +
				`(` +
					`\-[a-zA-Z0-9]` +				// groupable option
				`)` +
				`|` + 
				`(` +
					`\-\-[a-zA-Z0-9]+` +			// non-groupable option
				`)` +
			`)\|` +									// with bar
		`)*` +										// 0 to many times
		`(` +
			`(` +
				`\-[a-zA-Z0-9]` +					// groupable option
			`)` +
			`|` +
			`(` +
				`\-\-[a-zA-Z0-9]+` +				// non-groupable option
			`)` +
		`)` +										// w/out bar, 1 time

		// Variables and defaults
		`(` +
			`=` +
			'(' +
				`(` +
					`'[a-zA-Z0-9 ]*'` +				// allows for empty ''
				`)`	+
				`|`	+
				`(` +
					`\[` +
					`('` +
						`(` +
							`[a-zA-Z0-9 ]*', ?'` +	// default array value w/ comma, space delimiter is optional, allows for empty ''
						`)` +
						`*` +						// 0 to many times
						`(` +
							`[a-zA-Z0-9 ]*` +		// default array value w/out comma, allows for empty ''
						`)` +
					`')?` +							// default array values are optional, could be []
					`\]` +
				`)` +								
			`)?` +									// default value is optional, could be =
		`)?` +										// accepting a value is optional
		`$`,
	).Match(defBytes)) {
		Scream("invalid rule definition: " + string(def))
	}
}

func validateArg(def []byte) {
	if !(regexp.MustCompile(
		// Option keys
		`^(` + 
			`(` +
				`(` +
					`\-[a-zA-Z0-9]` +				// groupable option
				`)` +
				`|` + 
				`(` +
					`\-\-[a-zA-Z0-9]+` +			// non-groupable option
				`)` +
			`)\|` +									// with bar
		`)*` +										// 0 to many times
		`(` +
			`(` +
				`\-[a-zA-Z0-9]` +					// groupable option
			`)` +
			`|` +
			`(` +
				`\-\-[a-zA-Z0-9]+` +				// non-groupable option
			`)` +
		`)` +										// w/out bar, 1 time

		// Variables and defaults
		`(` +
			`=` +
			'(' +
				`(` +
					`'[a-zA-Z0-9 ]*'` +				// allows for empty ''
				`)`	+
				`|`	+
				`(` +
					`\[` +
					`('` +
						`(` +
							`[a-zA-Z0-9 ]*', ?'` +	// default array value w/ comma, space delimiter is optional, allows for empty ''
						`)` +
						`*` +						// 0 to many times
						`(` +
							`[a-zA-Z0-9 ]*` +		// default array value w/out comma, allows for empty ''
						`)` +
					`')?` +							// default array values are optional, could be []
					`\]` +
				`)` +								
			`)?` +									// default value is optional, could be =
		`)?` +										// accepting a value is optional
		`$`,
	).Match(defBytes)) {
		Scream("invalid rule definition: " + string(def))
	}
}

func parseKeys(def []byte, tokens *[]string) {
	keys := bytes.Explode(
		match(
			&def,
			regexp.MustCompile("^[-a-zA-Z0-9|]+"),
		),
		'|',
	)

	for _, key := range keys {
		isGroupable := false

		if _, isSet := (*tokens)[IS_OPTION]; (!isSet || (*tokens)[IS_OPTION] == "false") {
			(*tokens)[IS_OPTION] = bToS(match(&key, regexp.MustCompile("^-")) != nil)
		}

		if (*tokens)[IS_OPTION] == "true" {
			isGroupable = match(&key, regexp.MustCompile("^-")) == nil
		}

		append(parsedKeys, Key{
			Name: 			string(key),
			IsGroupable:	isGroupable,
		})
	}
}