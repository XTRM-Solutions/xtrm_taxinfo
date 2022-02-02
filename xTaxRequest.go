package main

import (
	"bufio"
	"bytes"
	"encoding/json"
)

type TransferTaxInformationRequest struct {
	TransferTaxInformation struct {
		Request struct {
			IssuerAccountNumber      string `json:"IssuerAccountNumber"`
			BeneficiaryAccountNumber string `json:"BeneficiaryAccountNumber"`
			TaxYear                  string `json:"TaxYear"`
			EmailBody                string `json:"EmailBody"`
			TaxYearStartMonth        string `json:"TaxYearStartMonth"`
		} `json:"Request"`
	} `json:"TransferTaxInformation"`
}

type TransferTaxInformationResponse struct {
	TransferTaxInformationResponse struct {
		TransferTaxInformationResult struct {
			Message         string `json:"Message"`
			OperationStatus struct {
				Errors struct {
				} `json:"Errors"`
				Success bool `json:"Success"`
			} `json:"OperationStatus"`
		} `json:"TransferTaxInformationResult"`
	} `json:"TransferTaxInformationResponse"`
}

func PostTaxRequest(spn string) (err error) {
	var jsonData []byte
	var xBody []byte
	var payload *bufio.Reader
	var ttir TransferTaxInformationRequest
	var rsp TransferTaxInformationResponse
	var url = xData["xUrl"] + "Report/TransferTaxInformation"

	err = nil

	req := &ttir.TransferTaxInformation.Request
	req.IssuerAccountNumber = xData["SPN"]
	req.BeneficiaryAccountNumber = FlagRecipient
	req.TaxYear = FlagYear
	req.TaxYearStartMonth = "1" // January
	req.EmailBody = "Your Tax Information"

	jsonData, err = json.Marshal(ttir)
	if nil != err {
		return err
	}
	if FlagDebug {
		xLog.Printf("DBG TransferTaxInformationRequest:\n%s\n%s\n",
			url,
			jsonData)
	}

	payload = bufio.NewReader(bytes.NewReader(jsonData))
	xBody, err = XPostRequest(url, payload)
	if nil != err {
		return err
	}

	err = json.Unmarshal(xBody, &rsp)

	return err
}
