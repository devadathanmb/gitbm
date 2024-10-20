package utils

import (
	"time"

	"github.com/goombaio/namegenerator"
)

// Gets random name as the title says
func GetRandomName() (name string) {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name = nameGenerator.Generate()

	return name
}
