package util

//CheckDecimal - verifica se um número possui casas decimais
func CheckDecimal(number float64) bool {

	result := number / float64(int(number))

	if result != 1 {
		return true
	}

	return false
}
