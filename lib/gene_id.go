package lib

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func GeneID(personId string) string {
	unixNano := time.Now().Unix()
	randInt := rand.Intn(1000)
	return fmt.Sprintf("%s_%d_%03d", personId, unixNano, randInt)
}
