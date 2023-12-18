package helpers

import "testing"

type urlTest struct {
	arg, expected string
}

var urlTests = []urlTest{
	{"http://www.youtube.com", "http://youtube.com"},
	{"http://www.youtube.com/", "http://youtube.com"},
	{"https://www.youtube.com/", "http://youtube.com"},
	{"https://www.youtube.com/", "http://youtube.com"},
	{"www.youtube.com", "http://youtube.com"},
	{"www.youtube.com/", "http://youtube.com"},
	{"youtube.com", "http://youtube.com"},
	{"youtube.com/", "http://youtube.com"},
}

func TestGeneralizeURL(t *testing.T) {
	for _, test := range urlTests {
		if output := GeneralizeURL(test.arg); output != test.expected {
			t.Errorf("Output %s not equal to expected %s", output, test.expected)
		}
	}
}