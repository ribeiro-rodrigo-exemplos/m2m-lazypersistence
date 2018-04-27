package util

import "encoding/json"

//CheckDecimal - verifica se um número possui casas decimais
func CheckDecimal(number float64) bool {

	result := number / float64(int(number))

	if result != 1 {
		return true
	}

	return false
}

func CheckDecimalInterface(numberInterface interface{}) bool {
	numberFloat, ok := numberInterface.(float64)

	if !ok {
		return false
	}

	return CheckDecimal(numberFloat)
}

func CheckDecimalNumber(number json.Number) bool {

	_, err := number.Int64()

	if err != nil {
		return true
	}

	return false
}
