package helper

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomString() string {
	return fmt.Sprint(time.Now().Unix() + int64(rand.Int()))
}
