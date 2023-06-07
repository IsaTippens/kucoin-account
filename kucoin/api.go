package kucoin

import (
	"strconv"
	"time"

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
	Id          string `json:"id"`
	Symbol      string `json:"symbol"`
	OpType      string `json:"opType"`
	Type        string `json:"type"`
	Side        string `json:"side"`
	Price       string `json:"price"`
	Size        string `json:"size"`
	Funds       string `json:"funds"`
	DealFunds   string `json:"dealFunds"`
	DealSize    string `json:"dealSize"`
	Fee         string `json:"fee"`
	FeeCurrency string `json:"feeCurrency"`
	CreatedAt   int64  `json:"createdAt"`
}

type FilledHFOrder struct {
	LastId int64   `json:"lastId"`
	Items  []Order `json:"items"`
}

type ActiveHFOrders []Order

type OrdersRequest struct {
	Symbol	string
	StartDate	time.Time
	EndDate		time.Time
}

func convert_date_to_ms_str(date time.Time) (string) {
	// subtract 2 hours
	date = date.Add(-2 * time.Hour)
	return strconv.FormatInt(date.UnixMilli(), 10)
}

func (c Client) GetFilledHFOrders(orders OrdersRequest) (*KucoinResponse, error) {
	params := map[string]string{
		"symbol": orders.Symbol,
		"limit":  "100",
	}
	if !orders.StartDate.IsZero() {
		params["startAt"] = convert_date_to_ms_str(orders.StartDate)
	}
	if !orders.EndDate.IsZero() {
		params["endAt"] = convert_date_to_ms_str(orders.EndDate)
	}
	req := NewRequest(fasthttp.MethodGet, "/api/v1/hf/orders/done", params, nil)
	return c.Do(&req)
}

func (c Client) GetActiveHFOrders(orders OrdersRequest) (*KucoinResponse, error) {
	params := map[string]string{
		"symbol": orders.Symbol,
		"limit":  "100",
	}
	req := NewRequest(fasthttp.MethodGet, "/api/v1/hf/orders/active", params, nil)
	return c.Do(&req)
}

func (c Client) GetFills() (*KucoinResponse, error) {
	params := map[string]string{
		"tradeType": "TRADE_HF",
	}
	req := NewRequest(fasthttp.MethodGet, "/api/v1/fills", params, nil)
	return c.Do(&req)
}
