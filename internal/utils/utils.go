package utils

import (
	"math"
	"math/rand"
	"time"
)

func CatchPokemon(baseExp int) bool {
	maxBaseExperience := 600
	if maxBaseExperience > baseExp {
		maxBaseExperience = baseExp
	}

	//catchProb := 1.0 - (float64(baseExp) / float64(maxBaseExperience+1.0))
	catchProb := 1 / (1 + math.Exp(float64(baseExp-200)/50))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomValue := r.Float64()
	// fmt.Printf("randomValue:%v\n", randomValue)
	// fmt.Printf("catchProb:%v\n", catchProb)
	return randomValue < catchProb

}
