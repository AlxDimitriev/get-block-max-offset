package test

import (
	"get_block_test/internal"
	"math/big"
	"testing"
)

func TestGetMaxBalanceOffset(t *testing.T) {
	client := internal.NewClient("qweqwe")
	val1 := new(big.Int)
	val1.SetString("-1234567890123456", 10)
	client.AddressesBalance["1"] = *val1

	val2 := new(big.Int)
	val2.SetString("-12345678901234567", 10)
	client.AddressesBalance["2"] = *val2

	val3 := new(big.Int)
	val3.SetString("-123456789012345678", 10)
	client.AddressesBalance["3"] = *val3

	val4 := new(big.Int)
	val4.SetString("12345678901234567", 10)
	client.AddressesBalance["4"] = *val4

	if addr, balance := client.GetMaxBalanceOffset(); addr != "3" || balance != "-123456789012345678" {
		t.Fatalf("wrong max balance offset. Addr: %s; Balance: %s", addr, balance)
	}

	val5 := new(big.Int)
	val5.SetString("12345678901234567890", 10)
	client.AddressesBalance["5"] = *val5

	if addr, balance := client.GetMaxBalanceOffset(); addr != "5" || balance != "12345678901234567890" {
		t.Fatalf("wrong max balance offset. Addr: %s; Balance: %s", addr, balance)
	}
}