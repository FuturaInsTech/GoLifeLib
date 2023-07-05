package utilities

import (
	"crypto/rand"
	"fmt"
)

func GetRandom() {
	RandomCrypto, _ := rand.Prime(rand.Reader, 128)
	fmt.Println("Random")
	fmt.Println(RandomCrypto)
	//return (RandomCrypto)
}
