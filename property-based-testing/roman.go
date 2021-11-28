package propertybasedtesting

import (
	"errors"
	"fmt"
	"strings"
)

type RomanNumeral struct {
	Value  uint16
	Symbol string
}

var allNumerals = []RomanNumeral{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func ConvertToRoman(arabic uint16) string {
	var result strings.Builder
	for _, numeral := range allNumerals {
		for arabic >= numeral.Value {
			result.WriteString(numeral.Symbol)
			arabic -= numeral.Value
		}
	}
	return result.String()
}

func ConvertToArabic(roman string) uint16 {
	var total uint16

	for i := 0; i < len(roman); i++ {
		var symbol string
		if i + 1 < len(roman) {
			symbol = string(roman[i:i+2])
			v, err := symbolLookup(symbol)
			if err == nil {
				total += v
				i++
				continue
			}
		}
		symbol = string(roman[i])
		v, err := symbolLookup(symbol)
		if err != nil {
			panic(fmt.Sprintf("unknown symbol %s", symbol))
		}
		total += v
	}
	return total
}

//attempt at improved performance - slower on benchmark tests though
func ConvertToArabic2(roman string) uint16 {
	var total uint16
	numeralMap := make(map[string]uint16)
	for _, v := range allNumerals {
		numeralMap[v.Symbol] = v.Value
	}


	for i := 0; i < len(roman); i++ {
		var symbol string
		if i + 1 < len(roman) {
			symbol = string(roman[i:i+2])
			if v, ok := numeralMap[symbol]; ok {
				total += v
				i++
				continue
			}
		}
		symbol = string(roman[i])
		v, err := symbolLookup(symbol)
		if err != nil {
			panic(fmt.Sprintf("unknown symbol %s", symbol))
		}
		total += v
	}
	return total
}

func symbolLookup(symbol string) (uint16, error) {
	for _, v := range allNumerals {
		if v.Symbol == symbol {
			return v.Value, nil
		}
	}
	return 0, errors.New("not found")
}
