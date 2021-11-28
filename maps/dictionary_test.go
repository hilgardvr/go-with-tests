package maps

import "testing"

func TestDictionary(t *testing.T) {
	t.Run("known word", func(t *testing.T) {
		dictionary := Dictionary{"test": "this is a test"}
		word := "test"
		got, _ := dictionary.Search(word)
		want := "this is a test"
		assertStrings(t, got, want)
	})
	t.Run("unknown word", func(t *testing.T) {
		dictionary := Dictionary{"test": "this is a test"}
		word := "unknown"
		_, err := dictionary.Search(word)
		assertError(t, err, ErrNotFound)
	})
	t.Run("add word", func(t *testing.T) {
		dict := Dictionary{}
		word := "test"
		def := "this is a test"
		dict.Add(word, def)

		got, err := dict.Search(word)

		if err != nil {
			t.Fatalf("should find added word: %q", err)
		}
		assertStrings(t, got, def)
	})
	t.Run("add existing word", func(t *testing.T) {
		word := "test"
		def := "this is a test"
		dict := Dictionary{word: def}
		err := dict.Add(word, "new dup")
		if err == nil {
			t.Fatalf("expecting error on duplicate insert")
		}
		
		got, err := dict.Search(word)
		if err != nil {
			t.Fatalf("should find added word: %q", err)
		}
		assertStrings(t, got, def)
	})
	t.Run("update a word", func(t *testing.T) {
		dict := Dictionary{}
		word := "test"
		def := "this is a test"
		def2 := "this is a test2"
		dict.Add(word, def)

		got, err := dict.Search(word)
		if err != nil {
			t.Fatalf("should find added word: %q", err)
		}
		assertStrings(t, got, def)

		dict.Update(word, def2)

		got, err = dict.Search(word)

		if err != nil {
			t.Fatalf("should find added word: %q", err)
		}
		assertStrings(t, got, def2)
	})
	t.Run("update a word that doesn't exist", func(t *testing.T) {
		dict := Dictionary{}
		word := "test"
		def := "this is a test"

		err := dict.Update(word, def)

		assertError(t, err, ErrWordDoesNotExist)
	})
	t.Run("delete a word", func(t *testing.T) {
		dict := Dictionary{}
		word := "test"
		def := "this is a test"

		dict.Add(word, def)

		got, err := dict.Search(word)
		if err != nil {
			t.Fatalf("should find added word: %q", err)
		}

		assertStrings(t, got, def)

		dict.Delete(word)

		_, err = dict.Search(word)

		assertError(t, err, ErrNotFound)

	})
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got error %q want error %q", got, want)
	}
}