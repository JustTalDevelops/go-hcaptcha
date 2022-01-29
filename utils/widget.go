package utils

import "math/rand"

// widgetCharacters are the characters used in randomly generated widget IDs.
var widgetCharacters = []rune("abcdefghijkmnopqrstuvwxyz0123456789")

// WidgetID generates a new random widget ID.
func WidgetID() string {
	b := make([]rune, Between(10, 12))
	for i := range b {
		b[i] = widgetCharacters[rand.Intn(len(widgetCharacters))]
	}
	return string(b)
}
