package testHelper

import (
	"regexp"
	"strings"
)

//NormalizeSpace is a function that takes whitespace out of the equation for comparing strings
func NormalizeSpace(str string) string {
	reSpace := regexp.MustCompile("\\s+")
	normalizedSpace := string(reSpace.ReplaceAllString(str, ` `))
	reBrackets := regexp.MustCompile("(\\}|\\{|\"|,|:|\\])\\s+(,|\\[|\\}|\")")
	adjustBrackets := string(reBrackets.ReplaceAllString(normalizedSpace, `$1$2`))
	reProp := regexp.MustCompile("(\\}|,|:)\\s+([^\\s])")
	adjustProperties := string(reProp.ReplaceAllString(adjustBrackets, `$1$2`))
	return strings.TrimSpace(adjustProperties)
}
