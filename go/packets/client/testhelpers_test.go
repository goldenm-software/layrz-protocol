package client_test

import "time"

var fixedTime = time.Unix(1700000000, 0)

func stringPtr(s string) *string  { return &s }
func floatPtr(f float64) *float64 { return &f }
