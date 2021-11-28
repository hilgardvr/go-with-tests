package main

import "testing"

func TestHello(t *testing.T) {
	assertCorrectMessage := func(t testing.TB, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
	t.Run("sqying hello to people", func(t *testing.T) {
		got := Hello("Pete", "")
		want := "Hello, Pete"
		assertCorrectMessage(t, got, want)
	})
	t.Run("sqying hello to world when no string is supplied", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})
	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMessage(t, got, want)
	})
	t.Run("in French", func(t *testing.T) {
		got := Hello("Alizee", "French")
		want := "Bonjour, Alizee"
		assertCorrectMessage(t, got, want)
	})
	t.Run("in Dutch", func(t *testing.T) {
		got := Hello("Sam", "Dutch")
		want := "Hallo, Sam"
		assertCorrectMessage(t, got, want)
	})
}