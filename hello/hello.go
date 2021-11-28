package main

// import "fmt"

const englishHelloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "
const frenchHelloPrefix = "Bonjour, "
const dutchHelloPrefix = "Hallo, "
const spanish = "Spanish"
const french = "French"
const dutch = "Dutch"
const world = "World"

func greetingPrefix(language string) (prefix string) {
	switch language {
	case french: 
		prefix = frenchHelloPrefix 
	case spanish:
		prefix = spanishHelloPrefix
	case dutch:
		prefix = dutchHelloPrefix 
	default:
		prefix = englishHelloPrefix
	}
	return
}

func Hello(name string, language string) string {
	if name == "" {
		name = world
	}
	prefix := greetingPrefix(language)
	return prefix + name
}

func main() {
}