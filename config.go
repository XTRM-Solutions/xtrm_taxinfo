package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

const ProgramName = "XtrmTaxinfo"
const LogName = ProgramName + ".log"
const RegExpAccountPattern = "^(PAT|SPN)[0-9]+$"
const RegExpTaxYearPattern = "^2[012][0-9]$"

var xLogFile *os.File
var xLogBuffer *bufio.Writer

var xLog log.Logger

var FlagManager string
var FlagRecipient string
var FlagDebug bool
var FlagVerbose bool
var FlagYear string

/* var FlagPretty bool  */

var nFlags *pflag.FlagSet

func ConfigLog() {

	var err error
	var logWriters []io.Writer

	xLogFile, err = os.OpenFile(ProgramName+".log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		xLog.Fatalf("error opening log file %s: %v", LogName, err)
	}
	xLogBuffer = bufio.NewWriter(xLogFile)
	logWriters = append(logWriters, os.Stderr)
	logWriters = append(logWriters, xLogBuffer)

	xLog.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	xLog.SetOutput(io.MultiWriter(logWriters...))
}

func ConfigFlags() {
	var rPattern *regexp.Regexp = nil

	nFlags = pflag.NewFlagSet("default", pflag.ContinueOnError)

	nFlags.StringP("taxyear", "t", "2020",
		"Tax year of report")
	nFlags.BoolP("debug", "d", false,
		"Enable debug logging and messages")
	nFlags.BoolP("verbose", "v", false,
		"Enable verbose output")
	nFlags.StringP("manager", "m",
		// "",
		"",
		"Company requesting tax information to send")
	nFlags.StringP("recipient", "r",
		"",
		"Beneficiary to receive tax information")

	err := nFlags.Parse(os.Args[1:])
	if nil != err {
		xLog.Fatalf("\nerror parsing flags: %s\n%s %s\n%s\n\t%v\n",
			err.Error(),
			"common issue: 2 hyphens for long-form arguments,",
			"1 hyphen for short-form argument",
			"Program arguments are:",
			os.Args)
	}

	FlagDebug = GetFlagBool("debug")
	FlagVerbose = GetFlagBool("verbose")
	if FlagVerbose {
		xLog.Print("Verbose mode engaged ... ")
	}
	FlagManager = strings.ToUpper(GetFlagString("manager"))
	if !IsStringSet(&FlagManager) {
		defer UsageMessage()
		xLog.Fatal("Oops! --manager is mandatory.")
	}
	rPattern = regexp.MustCompile(RegExpAccountPattern)
	if FlagDebug || FlagVerbose {
		xLog.Printf("compiled pattern: %s", rPattern.String())
	}
	if !rPattern.MatchString(FlagManager) {
		defer UsageMessage()
		xLog.Fatalf("%s\n%s\n ==> %s <== does not match pattern ==> %s <== \n",
			"Managing account must start with SPN followed by one or more digits",
			"PROBLEM: Sending account invalid format",
			FlagManager, rPattern.String())
	}

	FlagRecipient = GetFlagString("recipient")
	if !IsStringSet(&FlagRecipient) {
		defer UsageMessage()
		xLog.Fatal("Oops! --recipient is mandatory.")
	}
	if !rPattern.MatchString(FlagRecipient) {
		defer UsageMessage()
		xLog.Fatalf("%s\n%s\n ==> %s <== does not match pattern ==> %s <== \n",
			"Account must start with PAT followed by one or more digits",
			"PROBLEM: Sending account invalid format",
			FlagManager, rPattern.String())
	}
	FlagYear = GetFlagString("taxyear")
	if !IsStringSet(&FlagYear) {
		rPattern = regexp.MustCompile(RegExpTaxYearPattern)
		if !rPattern.MatchString("FlagYear") {
			xLog.Fatalf("%s\n==> %s <== does not match pattern ==> %s <==\n",
				"TaxYear must be year 2000 or later",
				FlagYear, rPattern.String())
		}
	}
}

// GetFlagBool fetch the bool for a boolean flag
func GetFlagBool(key string) (value bool) {
	var err error
	value, err = nFlags.GetBool(key)
	if nil != err {
		xLog.Fatalf("error fetching value for boolean flag [ %s ]: %s \n", key, err.Error())
		return false
	}
	return value
}

// GetFlagString fetch the string associated with a CLI arg
func GetFlagString(key string) (value string) {
	var err error
	value, err = nFlags.GetString(key)
	if nil != err {
		xLog.Fatalf("error fetching value for string flag [ %s ]: %s \n", key, err.Error())
		return ""
	}
	return value
}

// GetFlagInt fetch the value of integer flag
func GetFlagInt(key string) (value int) {
	var err error
	value, err = nFlags.GetInt(key)
	if nil != err {
		xLog.Fatalf("%s [ %s ]: %s \n",
			"error fetching value for integer flag",
			key, err.Error())
		return 0
	}
	return value
}

// UsageMessage /* UsageMessage
func UsageMessage() {
	_, _ = fmt.Printf("%s\n\t%s\n",
		nFlags.FlagUsages(), "\tHuh? --account is required but not present")
}
