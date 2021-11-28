package propertybasedtesting

import (
	"fmt"
	"log"
	"testing"
	"testing/quick"
)

var cases = []struct {
	Arabic uint16
	Roman  string
}{
	{1, "I"},
	{2, "II"},
	{3, "III"},
	{4, "IV"},
	{5, "V"},
	{6, "VI"},
	{7, "VII"},
	{8, "VIII"},
	{9, "IX"},
	{10, "X"},
	{10, "X"},
	{14, "XIV"},
	{18, "XVIII"},
	{20, "XX"},
	{39, "XXXIX"},
	{40, "XL"},
	{47, "XLVII"},
	{49, "XLIX"},
	{50, "L"},
	{100, "C"},
	{90, "XC"},
	{200, "CC"},
	{400, "CD"},
	{1984, "MCMLXXXIV"},
	{3999, "MMMCMXCIX"},
	{2014, "MMXIV"},
	{1006, "MVI"},
	{798, "DCCXCVIII"},
}

func TestRomanNumerals(t *testing.T) {
	for _, test := range cases {
		desc := fmt.Sprintf("%d gets converted to %s", test.Arabic, test.Roman)
		t.Run(desc, func(t *testing.T) {
			got := ConvertToRoman(test.Arabic)

			if got != test.Roman {
				t.Errorf("got %q, want %q", got, test.Roman)
			}
		})
	}
}

func TestConvertToArabic(t *testing.T) {
	for _, test := range cases {
		desc := fmt.Sprintf("%s gets converted to %d", test.Roman, test.Arabic)
		t.Run(desc, func(t *testing.T) {
			got := ConvertToArabic(test.Roman)
			if got != test.Arabic {
				t.Errorf("got %d, want %d", got, test.Arabic)
			}
		})
		t.Run(desc, func(t *testing.T) {
			got := ConvertToArabic2(test.Roman)
			if got != test.Arabic {
				t.Errorf("got %d, want %d", got, test.Arabic)
			}
		})
	}
}

func BenchmarkConvertToArabic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range cases {
			desc := fmt.Sprintf("%s gets converted to %d", test.Roman, test.Arabic)
			b.Run(desc, func(t *testing.B) {
				got := ConvertToArabic(test.Roman)
				if got != test.Arabic {
					t.Errorf("got %d, want %d", got, test.Arabic)
				}
			})
			desc = fmt.Sprintf("%s gets converted to %d 2", test.Roman, test.Arabic)
			b.Run(desc, func(t *testing.B) {
				got := ConvertToArabic2(test.Roman)
				if got != test.Arabic {
					t.Errorf("got %d, want %d", got, test.Arabic)
				}
			})
		}
	}
}

func TestPropertiesOfConversion(t *testing.T) {
	assersion := func(arabic uint16) bool {
		log.Println(arabic)
		roman := ConvertToRoman(arabic)
		fromRoman := ConvertToArabic(roman)
		return fromRoman == arabic
	}

	if err := quick.Check(assersion, nil); err != nil {
		t.Error("failed checks", err)
	}
}