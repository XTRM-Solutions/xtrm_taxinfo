package main

func main() {

	ConfigLog()                // must be first
	DeferError(xLogFile.Close) // close log file on exit
	ConfigFlags()
	ConfigSecrets()

	// xGetPaymentMethods() // make sure requests are working?

	err := PostTaxRequest(FlagManager)
	if nil != err {
		xLog.Fatalf("Oops! PostTaxRequest failed for account %s because: %s\n",
			FlagManager, err.Error())
	}

}
