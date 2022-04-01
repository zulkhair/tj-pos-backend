package stringutil

func LeftPad(text string, length int, char byte) string {
	textLen := len(text)

	textToAdd := ""
	for i := textLen; i < length; i++ {
		textToAdd = string(char) + textToAdd
	}

	return textToAdd + text
}

func RightPad(text string, length int, char byte) string {
	textLen := len(text)

	textToAdd := ""
	for i := textLen; i < length; i++ {
		textToAdd = textToAdd + string(char)
	}

	return text + textToAdd
}
