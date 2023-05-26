package kucoin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

type KucoinSigner struct {
	key []byte
	apiKey        string
	apiSecret     string
	apiPassPhrase string
	apiKeyVersion string
}

func (ks *KucoinSigner) Sign(plain []byte) []byte {
	hm := hmac.New(sha256.New, ks.key)
	hm.Write(plain)
	return []byte(base64.StdEncoding.EncodeToString(hm.Sum(nil)))
}

func NewKcSigner(key, secret, passPhrase string) *KucoinSigner {
	ks := &KucoinSigner{
		apiKey:        key,
		apiSecret:     secret,
		apiPassPhrase: passPhraseEncrypt([]byte(secret), []byte(passPhrase)),
		apiKeyVersion: "2",
	}
	ks.key = []byte(secret)
	return ks
}

func (ks *KucoinSigner) Headers(plain string) map[string]string {
	
	t := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
	p := []byte(t + plain)
	s := string(ks.Sign(p))
	ksHeaders := map[string]string{
		"KC-API-KEY":        ks.apiKey,
		"KC-API-PASSPHRASE": ks.apiPassPhrase,
		"KC-API-TIMESTAMP":  t,
		"KC-API-SIGN":       s,
		"KC-API-KEY-VERSION": ks.apiKeyVersion,
	}

	return ksHeaders
}

func passPhraseEncrypt(key, plain []byte) string {
	hm := hmac.New(sha256.New, key)
	hm.Write(plain)
	return base64.StdEncoding.EncodeToString(hm.Sum(nil))
}
