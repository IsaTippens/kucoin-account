package kucoin

import (
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type TransferRequest struct {
	Oid      string  `json:"clientOid"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	From     string  `json:"from"`
	To       string  `json:"to"`
}

type TransferModel struct {
	OrderId string `json:"orderId"`
}

func (c Client) Transfer(tr TransferRequest) (*KucoinResponse, error) {
	// set oid to uuid
	if tr.Oid == "" {
		tr.Oid = uuid.New().String()
	}
	req := NewRequest(fasthttp.MethodPost, "/api/v2/accounts/inner-transfer", nil, tr)
	return c.Do(&req)
}

type AccountsRequest struct {
	Currency string `json:"currency;omitempty"`
	Type     string `json:"type;omitempty"`
}

type AccountModel struct {
	Currency  string `json:"currency"`
	Type      string `json:"type"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
}

// An AccountsModel is the set of *AccountModel.
type AccountsModel []*AccountModel

func (c Client) GetAccounts() (*KucoinResponse, error) {
	req := NewRequest(fasthttp.MethodGet, "/api/v1/accounts", nil, nil)
	return c.Do(&req)
}

type Order struct {
	Id 		  string `json:"id"`
	Symbol    string `json:"symbol"`
	OpType    string `json:"opType"`
	Type      string `json:"type"`
	Side      string `json:"side"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Funds	 string `json:"funds"`
	DealFunds string `json:"dealFunds"`
	DealSize  string `json:"dealSize"`
	Fee       string `json:"fee"`
	FeeCurrency string `json:"feeCurrency"`
	CreatedAt int64  `json:"createdAt"`
}

type FilledHFOrder struct {
	LastId int64 `json:"lastId"`
	Items []Order `json:"items"`
}

type ActiveHFOrders []Order

func (c Client) GetFilledHFOrders(symbol string) (*KucoinResponse, error) {
	params := map[string]string{
		"symbol": symbol,
	}
	req := NewRequest(fasthttp.MethodGet, "/api/v1/hf/orders/done", params, nil)
	return c.Do(&req)
}

func (c Client) GetActiveHFOrders(symbol string) (*KucoinResponse, error) {
	params := map[string]string{
		"symbol": symbol,
	}
	req := NewRequest(fasthttp.MethodGet, "/api/v1/hf/orders/active", params, nil)
	return c.Do(&req)
}
