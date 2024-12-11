package helper

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateNanoId() (string, error) {
	nanoId, err := gonanoid.Generate("abcdefg1234567890", 10)
	return nanoId, err
}
