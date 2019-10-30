package luhn

import "testing"

func Test_A_is_odd(t *testing.T) {
	cases := []struct {
		in int;
		want bool
	}{
		{3, true},
		{4, false},
	}
	for _, c := range cases {
		got := is_odd(c.in)
		if got != c.want {
			t.Errorf("strlen_is_odd(%q) == %t, want %t", c.in, got, c.want)
		}
	}
}

func TestOne(t *testing.T){
	cases := []struct {
		in string;
		want bool
	}{
		{"79927398710", false},
		{"79927398711", false},
		{"79927398712", false},
		{"79927398713", true},
		{"79927398714", false},
		{"79927398715", false},
		{"79927398716", false},
		{"79927398717", false},
		{"79927398718", false},
		{"79927398719", false},
	}
	for _, c := range cases {
		got := Luhn(c.in)
		if got != c.want {
			t.Errorf("Luhn(%q) == %t, want %t", c.in, got, c.want)
		}
	}
}

func Test_B_odd_and_even_digits(t *testing.T){
	cases := []struct {
		in string;
		want bool
	}{
		// 11 digits, should double 1,8,3,2,9
		{"79927398713", true},
		// 10 digits, should double 7,9,3,2,9
		{"79927399877", true},
	}
	for _, c := range cases {
		got := Luhn(c.in)
		if got != c.want {
			t.Errorf("Luhn(%q) == %t, want %t", c.in, got, c.want)
		}
	}
}


func Test_B_various_lengths(t *testing.T){
	cases := []struct {
		in string;
		want bool
	}{
		{"", false},  // length of 0
		{"0", false},  // length of 1
		{"75", true},
		{"794", true},
		{"7997", true},
		{"79921", true},
		{"799270", true},
		{"79927398", true},
		{"799273990", true},
	}
	for _, c := range cases {
		got := Luhn(c.in)
		if got != c.want {
			t.Errorf("Luhn(%q) == %t, want %t", c.in, got, c.want)
		}
	}
}

