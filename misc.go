package main

import "strings"

// DeferError
// account for an at-close function that
// may return an error for its close
func DeferError(f func() error) {
	err := f()
	if nil != err {
		xLog.Printf("%s%s",
			"(may be harmless) error in deferred function: ",
			err.Error())
	}
}

// WriteSB Add a series of strings to a strings.Builder
func WriteSB(sb *strings.Builder, inputStrings ...string) {
	for _, val := range inputStrings {
		_, err := sb.WriteString(val)
		if nil != err {
			xLog.Print("strings.Builder failed to add " + val + " ??")
			xLog.Fatal("values: ", inputStrings)
		}
	}
}

func IsStringSet(s *string) (isSet bool) {
	if nil != s && "" != *s {
		return true
	}
	return false
}
