package main

import "testing"

func TestsqlRounding(t *testing.T) {
	tests := map[int]int{
		0:   0,
		1:   4,
		2:   4,
		3:   4,
		4:   4,
		5:   6,
		6:   6,
		7:   8,
		8:   8,
		12:  12,
		59:  60,
		900: 900,
		901: 902,
	}

	r := Report{}

	for val, ans := range tests {
		if res := r.sqlRounding(val); res != ans {
			t.Error(
				"Sql rouding error: %s returns %s when expecting %s",
				val,
				res,
				ans,
			)
		}
	}
}
