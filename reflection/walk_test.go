package reflection

import (
	"reflect"
	"testing"
)

func TestWalk(t *testing.T) {
	type Profile struct {
		Age int
		City string
	}
	type Person struct {
		Name string
		Profile Profile
	}

	cases := []struct{
		Name string
		Input interface{}
		ExpectedCalls []string
	}{
		{
			"struct with one string field",
			struct{
				Name string
			}{"Chris"},
			[]string{"Chris"},
		},
		{
			"struct with two string fields",
			struct{
				Name string
				City string
			}{
				"Chris",
				"London",
			},
			[]string{"Chris", "London"},
		},
		{
			"struct with non string field",
			struct{
				Name string
				Age int
			}{
				"Chris",
				37,
			},
			[]string{"Chris"},
		},
		{
			"nested fields",
			Person{
				"Chris",
				Profile{
					33, 
					"London",
				},
			},
			[]string{"Chris", "London"},
		},
		{
			"pointers to things",
			&Person{
				"Chris",
				Profile{
					33, 
					"London",
				},
			},
			[]string{"Chris", "London"},
		},
		{
			"slices",
			[]Profile{
				{
					33, 
					"London",
				},
				{
					44,
					"Pta",
				},
			},
			[]string{"London", "Pta"},
		},
		{
			"arrays",
			[2]Profile{
				{ 33, "London", },
				{ 44, "Pta", },
			},
			[]string{"London", "Pta"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %q, want %q", got, test.ExpectedCalls)
			}
		})
	}
	t.Run("with maps", func(t *testing.T) {
		var got []string
		m := map[string]string{
				"33": "London",
				"44": "Pta",
			}
		want := []string{"London", "Pta"}
		walk(m, func(input string) {
			got = append(got, input)
		})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("map got %q, want %q", got, want)
		}
	})
	t.Run("with channels", func(t *testing.T) {
		var got []string
		m := make(chan Profile)
		go func() {
			m <- Profile{33, "London"}
			m <- Profile{33, "Pta"}
			close(m)
		}()
		want := []string{"London", "Pta"}
		walk(m, func(input string) {
			got = append(got, input)
		})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("map got %q, want %q", got, want)
		}
	})
	t.Run("with functions", func(t *testing.T) {
		var got []string
		aFunc := func() (Profile, Profile) {
			return Profile{33, "London"}, Profile{44, "Pta", }
		}
		want := []string{"London", "Pta"}
		walk(aFunc, func(input string) {
			got = append(got, input)
		})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("map got %q, want %q", got, want)
		}
	})
}