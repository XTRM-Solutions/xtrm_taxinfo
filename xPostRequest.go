package main

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

var transport *http.Transport = nil
var sepString string = "/**********************************************/"
var client *http.Client = nil

func XPostRequest(url string, payload io.Reader) (resp []byte, err error) {
	var requestHeaders, responseHeaders, sb strings.Builder

	if nil == client {
		xConfigHttp()
	}

	req, err := http.NewRequest("POST", url, payload)
	if nil != err {
		return resp, err
	}

	WriteSB(&sb, xData["TokenType"], " ", xData["AccessToken"])
	if GetFlagBool("debug") {
		xLog.Println("DGB authorization: " + sb.String())
	}
	req.Header.Add("Authorization", sb.String())
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "*/*")
	req.Header.Del("Transfer-Encoding")
	req.Header.Del("Keep-Alive")

	if FlagDebug {
		headers := req.Header
		for key, values := range headers {
			for ix, val := range values {
				WriteSB(&requestHeaders, "\theader: ", key, " [", strconv.Itoa(ix), "] : ", val, "\n")
			}
		}
	}

	rsp, err := client.Do(req)

	if nil == err {
		defer DeferError(rsp.Body.Close)
		resp, err = io.ReadAll(rsp.Body)
		setTimeoutExpires()
	}

	if FlagDebug {
		headers := rsp.Header
		for key, values := range headers {
			for ix, val := range values {
				WriteSB(&responseHeaders, "\theader: ", key, " ", strconv.Itoa(ix), " : ", val, "\n")
			}
		}
		sb.Reset()
		xLog.Println(sepString, "\nDBG Request/Response:\n",
			sepString, "\nRequest Headers: \n", requestHeaders.String(),
			sepString, "\nResponse Headers:\n", responseHeaders.String(),
			sepString, "\n", string(resp),
			sepString)
	}
	return resp, err
}

func xConfigHttp() {
	if nil == transport {
		transport = &http.Transport{
			Proxy:                  nil,
			DialContext:            nil,
			Dial:                   nil,
			DialTLSContext:         nil,
			DialTLS:                nil,
			TLSClientConfig:        nil,
			TLSHandshakeTimeout:    0,
			DisableKeepAlives:      false,
			DisableCompression:     false,
			MaxIdleConns:           0,
			MaxIdleConnsPerHost:    0,
			MaxConnsPerHost:        0,
			IdleConnTimeout:        0,
			ResponseHeaderTimeout:  0,
			ExpectContinueTimeout:  0,
			TLSNextProto:           nil,
			ProxyConnectHeader:     nil,
			GetProxyConnectHeader:  nil,
			MaxResponseHeaderBytes: 0,
			WriteBufferSize:        0,
			ReadBufferSize:         0,
			ForceAttemptHTTP2:      false,
		}
	}

	if nil == client {
		client = &http.Client{
			Transport:     transport,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}
	}
}
