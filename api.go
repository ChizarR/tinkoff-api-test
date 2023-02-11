package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// SANDBOX
	baseURL = "https://sandbox-invest-public-api.tinkoff.ru/rest/tinkoff.public.invest.api.contract.v1.SandboxService/"

	// SANDBOX METHODS
	createSandboxAccount = "OpenSandboxAccount"
	getSandboxAccounts   = "GetSandboxAccounts"
	closeSandboxAccount  = "CloseSandboxAccount"
	getSandboxPositions  = "GetSandboxPositions"
	getSandboxPortfolio  = "GetSandboxPortfolio"
	payInSandbox         = "SandboxPayIn"
)

type TinkoffSandboxAPI interface {
	CreateSandboxAccount() (string, error)
	GetSandboxAccounts() ([]map[string]any, error)
	CloseSandboxAccount(account string) bool
}

type tinkoffAPI struct {
	client http.Client
	token  string
}

func NewTinkoffSandboxAPI(token string) TinkoffSandboxAPI {
	client := http.Client{Timeout: 15 * time.Second}
	return &tinkoffAPI{client: client, token: token}
}

func (t *tinkoffAPI) CreateSandboxAccount() (string, error) {
	resp, err := t.sendReq(http.MethodPost, createSandboxAccount, nil)
	if err != nil {
		return "", err
	}

	body := resp.Body
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	defer body.Close()

	var res map[string]any
	if err := json.Unmarshal(bytes, &res); err != nil {
		return "", err
	}

	fmt.Println(res)

	val, ok := res["accountId"]
	if !ok {
		return "", fmt.Errorf("No accountId")
	}
	return val.(string), nil
}

type Accounts []map[string]any

func (t *tinkoffAPI) GetSandboxAccounts() ([]map[string]any, error) {
	var accounts Accounts

	resp, err := t.sendReq(http.MethodPost, getSandboxAccounts, nil)
	if err != nil {
		return []map[string]any{}, err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []map[string]any{}, err
	}

	var respData map[string]Accounts
	if err := json.Unmarshal(bytes, &respData); err != nil {
		return []map[string]any{}, err
	}

	val, ok := respData["accounts"]
	if !ok {
		return accounts, fmt.Errorf("No accountId")
	}
	return val, nil
}

func (t *tinkoffAPI) CloseSandboxAccount(account string) bool {
	body := map[string]string{"accountId": account}
	resp, err := t.sendReq(http.MethodPost, closeSandboxAccount, body)
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

func (t *tinkoffAPI) sendReq(httpMethod, apiMethod string, body any) (*http.Response, error) {
	if body == nil {
		raw := map[string]any{}
		bt, err := json.Marshal(raw)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bt)
	}

	bt, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	b := bytes.NewReader(bt)

	url := baseURL + apiMethod
	req, err := http.NewRequest(httpMethod, url, b)
	if err != nil {
		return &http.Response{}, err
	}

	authString := "Bearer " + t.token
	req.Header.Add("Authorization", authString)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
