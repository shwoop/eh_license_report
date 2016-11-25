package main

import "testing"

func TestValidateUuid(t *testing.T) {
	validUuids := []string{
		"af236ef8-dd2a-4eb9-9564-91983a2bb800",
		"e694b11f-de56-4376-8bea-a449cb2f9368",
		"ab4290e5-af26-4422-89dc-4a6a1f6212b5",
	}
	badUuids := []string{
		"",
		"words",
		"9",
		"ab4290e5af26442289dc4a6a1f6212b5",
		"ab4290e5-af26:4422-89dc-4a6a1f6212b5",
	}
	for _, val := range validUuids {
		if result := ValidateUuid(val); !result {
			t.Error("Valid UUID reporting as invalid: ", val)
		}
	}
	for _, val := range badUuids {
		if result := ValidateUuid(val); result {
			t.Error("Invalid UUID reporting as valid: ", val)
		}
	}
}

var result bool

func benchmarkValidateUuid(s string, b *testing.B) {
	var r bool
	for i := 0; i < b.N; i++ {
		r = ValidateUuid(s)
	}
	result = r
}

func BenchmarkValidateUuid(b *testing.B) {
	validUuids := []string{
		"af236ef8-dd2a-4eb9-9564-91983a2bb800",
		"e694b11f-de56-4376-8bea-a449cb2f9368",
		"ab4290e5-af26-4422-89dc-4a6a1f6212b5",
	}
	badUuids := []string{
		"",
		"words",
		"9",
		"ab4290e5af26442289dc4a6a1f6212b5",
		"ab4290e5-af26:4422-89dc-4a6a1f6212b5",
	}
	for _, val := range validUuids {
		b.Run("ok "+val, func(b *testing.B) {
			benchmarkValidateUuid(val, b)
		})
	}
	for _, val := range badUuids {
		b.Run("bad "+val, func(b *testing.B) {
			benchmarkValidateUuid(val, b)
		})
	}

	// b.Run("ok 1", func(b *testing.B) {
	// 	benchmarkValidateUuid("af236ef8-dd2a-4eb9-9564-91983a2bb800", b)
	// })
	// b.Run("ok 2", func(b *testing.B) {
	// 	benchmarkValidateUuid("e694b11f-de56-4376-8bea-a449cb2f9368", b)
	// })
	// b.Run("ok 3", func(b *testing.B) {
	// 	benchmarkValidateUuid("ab4290e5-af26-4422-89dc-4a6a1f6212b5", b)
	// })
	// b.Run("bad 1", func(b *testing.B) {
	// 	benchmarkValidateUuid("bad", b)
	// })
	// b.Run("bad 2", func(b *testing.B) {
	// 	benchmarkValidateUuid("bad", b)
	// })
	// b.Run("bad 3", func(b *testing.B) {
	// 	benchmarkValidateUuid("ab4290e5af26442289dc4a6a1f6212b5", b)
	// })
	// b.Run("bad 4", func(b *testing.B) {
	// 	benchmarkValidateUuid("ab4290e5-af26:4422-89dc-4a6a1f6212b5", b)
	// })
}
