package main

import (
	"bytes"
	"encoding/json"
)

type GetPaymentMenthodsResponseHeader struct {
	GetPaymentMethodsResponse struct {
		PaymentMethodResult struct {
			OperationStatus struct {
				Errors struct {
				} `json:"Errors"`
				Success bool `json:"Success"`
			} `json:"OperationStatus"`
			PaymentMethods struct {
				PaymentMethodDetails []struct {
					PaymentMethodID   string `json:"PaymentMethodID"`
					PaymentMethodName string `json:"PaymentMethodName"`
				} `json:"PaymentMethodDetails"`
			} `json:"PaymentMethods"`
		} `json:"PaymentMethodResult"`
	} `json:"GetPaymentMethodsResponse"`
}

func xGetPaymentMethods() (GetPaymentMenthodsResponseHeader, error) {
	var resp GetPaymentMenthodsResponseHeader
	var url = xData["xUrl"] +
		"Payment/GetBeneficiaryCompanyPaymentMethods"

	var jsonBytes = []byte("{}")

	payload := bytes.NewReader(jsonBytes)
	xBody, err := XPostRequest(url, payload)
	if nil != err {
		return resp, err
	}
	return resp, json.Unmarshal(xBody, &resp)

}
