package luhn

import(
	//"fmt"
)

func is_odd(n int) bool {
	return n % 2 == 1
}

func Luhn(target string) bool {
	sum := 0
	table := []int {1,2,4,6,8,0,3,5,7,9}
	if 1 >= len(target){
		return false
	}
	for i,c := range target[0:len(target)]{
		dig := int(c) - int('0')
		if is_odd(len(target)) == is_odd(i){
			sum += table[dig]
			//fmt.Println("summing", dig, i)
		}else{
			sum += dig
		}
	}
	//fmt.Println("resut", target, sum, sum * 9)
	return (sum % 10) == 0
}