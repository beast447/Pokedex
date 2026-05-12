package main

import "testing"


func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
		input: "   hello fella   ",
		expected: []string{"hello", "fella"},
		},
	
	{
		input: "HELLO FELLA",
		expected: []string{"hello", "fella"},
	},
	{
		input: "HeLLoFelLa",
		expected: []string{"hellofella"},
	},
	{
		input: "   hello fella",
		expected: []string{"hello", "fella"},
	},
}
	 
for i, c := range cases{
	actual := cleanInput(c.input)
	if len(actual) != len(c.expected){
		t.Errorf("Length not matching for test case: %v", i)
		t.Fail()
	}
	for i := range actual{
		word := actual[i]
		expectedWord := c.expected[i]
		if word != expectedWord{
			t.Errorf("Word doesnt match on test case: %v", i)
			t.Fail()
		}
	}
}
}
