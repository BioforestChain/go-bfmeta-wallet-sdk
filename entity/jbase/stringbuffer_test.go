package jbase_test

import (
	"encoding/json"
	"testing"

	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"
	"github.com/stretchr/testify/assert"
)

func TestJson(t *testing.T) {
	type QAQ struct {
		Value jbase.StringBuffer `json:"qaq"`
	}
	{
		hex := jbase.NewUtf8StringBuffer("ff12345\"6")
		hexJson, _ := json.Marshal(QAQ{Value: hex.StringBuffer})

		assert.Equal(t, `{"qaq":"ff12345\"6"}`, string(hexJson))
	}
	type QAQ2 struct {
		Value *jbase.StringBuffer `json:"qaqx"`
	}
	{
		hex := jbase.NewUtf8StringBuffer("ff12345\"6")
		hexJson, _ := json.Marshal(QAQ2{Value: &hex.StringBuffer})

		assert.Equal(t, `{"qaqx":"ff12345\"6"}`, string(hexJson))
	}
}
