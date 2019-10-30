package luhn

import(
	/// "fmt"
	"strings"
	"unicode"
)

func is_odd(n int) bool {
	return n % 2 == 1
}

func non_digits(str string) string {
    return strings.Map(
    	func(r rune) rune {
	        if unicode.IsDigit(r) {
	            return -1
	        }
	        return r
    	},
    	str)
}

func Valid(target string) bool {
	sum := 0
	table := []int {0,2,4,6,8,1,3,5,7,9}
	spaceless := strings.Join(strings.Fields(target), "")
	// fmt.Println("testing with", s, "now", spaceless)
	if len(non_digits(spaceless)) > 0 {
		return false
	}
	if 1 >= len(spaceless){
		return false
	}
	for i,c := range spaceless{
		dig := int(c) - int('0')
		if is_odd(len(spaceless)) == is_odd(i){
			sum += table[dig]
			// fmt.Println("summing", dig, i)
		}else{
			sum += dig
		}
	}
	// fmt.Println("resut", spaceless, sum, sum * 9)
	return (sum % 10) == 0
}

