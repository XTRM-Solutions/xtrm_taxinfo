package main

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	Issued       string `json:".issued"`
	Expires      string `json:".expires"`
}

var xData = map[string]string{
	"": "default",
}

// XTRM time example: "Wed, 28 Oct 2020 20:15:16 GMT"
const xtrmTimeFormat string = "Mon, 02 Jan 2006 15:04:05 MST"

var xtrmTimeout, _ = time.ParseDuration("28m")

func setTimeoutExpires() {
	xData["InactiveTimeout"] =
		time.Now().Add(xtrmTimeout).Format(xtrmTimeFormat)
}

func checkTimeout() (isExpired bool) {
	ts := xData["InactiveTimeout"]
	if "" != ts {
		timeExpires, err := time.Parse(xtrmTimeFormat, ts)
		if nil != err || timeExpires.Before(time.Now()) {
			return true
		}
		return false
	}
	return true
}

func isTokenActive(duration time.Duration) (active bool) {

	// check inactive timeout
	if checkTimeout() {
		return false
	}

	// Do we already have an access token good for at least 2 hours?
	if "" != xData["AccessToken"] {
		// is it current?
		expires := xData["Expires"]
		if "" != expires {
			timeExpires, err := time.Parse(xtrmTimeFormat, expires)
			if nil != err {
				xLog.Fatal("Internal error: could not parse time [ " +
					expires + "] as format [ " +
					xtrmTimeFormat + " ]\n\tbecause\n" +
					err.Error())
			}
			if timeExpires.After(time.Now().Add(duration)) {
				return true
			}
		}
	}
	return false
}

func xAuthorize(xMethod, xUrl, xClient, xSecret string) (success bool) {

	if isTokenActive(1 * time.Hour) {
		return true
	}

	// otherwise need to authorize or reauthorize

	payload :=
		bufio.NewReader(strings.NewReader("grant_type=password" +
			"&client_id=" + xClient +
			"&client_secret=" + xSecret))

	req, err := http.NewRequest(xMethod, xUrl, payload)

	if err != nil {
		xLog.Fatal(err.Error())
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if nil != err {
		xLog.Fatal(err.Error())
	}
	defer DeferError(res.Body.Close)

	xBody, err := io.ReadAll(res.Body)
	if nil != err {
		xLog.Fatal(err.Error())
	}

	var tr tokenResponse
	err = json.Unmarshal(xBody, &tr)
	if nil != err {
		xLog.Fatal(err.Error())
	}
	/*
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		ClientID     string `json:"client_id"`
		Issued       string `json:".issued"`
		Expires      string `json:".expires"`
	*/

	xData["AccessToken"] = tr.AccessToken
	xData["TokenType"] = tr.TokenType
	xData["ExpiresIn"] = strconv.Itoa(tr.ExpiresIn)
	xData["RefreshToken"] = tr.RefreshToken
	xData["ClientID"] = tr.ClientID
	xData["Issued"] = tr.Issued
	xData["Expires"] = tr.Expires

	setTimeoutExpires()

	if FlagVerbose {
		xLog.Println("AuthorizeRequest Succeeded")
	}

	return len(tr.AccessToken) > 0
}
