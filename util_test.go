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
