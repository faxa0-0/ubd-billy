package utils

import (
	"billy/models"
	"math/rand/v2"
)

func GenerateUsage(user_id int) models.Usage {
	return models.Usage{
		UserID:  user_id,
		Youtube: rand.Uint64(),
		Spotify: rand.Uint64(),
		Netflix: rand.Uint64(),
		Basic:   rand.Uint64(),
	}
}
