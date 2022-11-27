package internal

import (
	"log"
	"math/big"
	"strconv"
)

func hex2Int(hex string) int {
	res := hex[2:]
	num, err := strconv.ParseInt(res, 16, 64)
	if err != nil {
		log.Fatal(err)
	}
	return int(num)
}

func hex2BigInt(hexString string) big.Int {
	s := hexString[2:]
	bigInt := new(big.Int)
	bigInt.SetString(s, 16)
	return *bigInt
}