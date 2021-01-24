package crypt

import (
	"fmt"
	"testing"
)

type data struct {
	pass         string
	passToVerify string
	result       bool
}

var dataset = []data{
	{pass: "", passToVerify: "", result: true},
	{pass: "senha", passToVerify: "s3nha", result: false},
	{pass: "pass12345!@#$%", passToVerify: "pass12345!@#$%", result: true},
}

func TestArgon2(t *testing.T) {
	for _, data := range dataset {
		hash, err := Generate(data.pass)
		fmt.Println(hash, err)
	}
}
