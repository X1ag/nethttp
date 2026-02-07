package utils

import "testing"

func TestAdd(t *testing.T) {
	expected := 5

	if Add(2, 3) != expected {
		t.Errorf("expected %d, got %d", expected, Add(2, 3))
	}
}

func TestEven(t *testing.T) {
	expected := true 

	if IsEven(2) != expected {
		t.Errorf("expected %v, got %v", expected, IsEven(2))
	}
}

