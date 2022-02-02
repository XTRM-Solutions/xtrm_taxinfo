package main

// ConfigSecrets
// Get the login information
// Hardcoded for now, but could get information from a db
// or some other source. For testing, this is sufficient.
func ConfigSecrets() {
	xData["xAuthorizeUrl"] = "https://xapisandbox.xtrm.com/oAuth/token"
	xData["xUrl"] = "https://xapisandbox.xtrm.com/API/V4/"
	xData["xClient"] = "DONOTUSE_API_User"
	xData["xSecret"] = "NOTREALLYASECRETUSEYOURAPISECRETHERE"
	xData["SPN"] = "SPN20136817"
	xAuthorize(
		"POST",
		xData["xAuthorizeUrl"],
		xData["xClient"],
		xData["xSecret"],
	)
}
