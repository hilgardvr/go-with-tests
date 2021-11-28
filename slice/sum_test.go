package slice

import (
	"reflect"
	"testing"
)


func TestSum(t *testing.T) {
	t.Run("collection of any size", func(t *testing.T) {

		numbers := []int{1,2,3}

		want := 6
		got := Sum(numbers)

		if got != want {
			t.Errorf("got %d wanted %d, given %v", got, want, numbers)
		}
	})
}

func TestSumAll(t *testing.T) {
	got := SumAll([]int{1,2}, []int{3,4,5})
	want := []int{3, 12}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestSumTails(t *testing.T) {

	checkSums := func(t testing.TB, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("make sum of slices", func(t *testing.T) {
		got := SumTails([]int{1,2}, []int{3,4,5}, []int{1})
		want := []int{2,9,0}
		checkSums(t, got, want)
	})

	t.Run("handle empty slices", func(t *testing.T) {
		got := SumTails([]int{1,2}, []int{3,4,5}, []int{})
		want := []int{2,9,0}
		checkSums(t, got, want)
	})
}