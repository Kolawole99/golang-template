package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword uses bcrypt to HashAndSalt a plain string password
func HashPassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	return string(hash)
}

// ComparePassword uses bcrypt to hash a plain password and compares it to the stored hash to determine if they are equal
func ComparePassword(hashedPassword string, plainPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), plainPassword)

	return err == nil
}

// func randSeq(n int) string {
// 	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
// 	b := make([]rune, n)

// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}

// 	return string(b)
// }

func GenerateToken() string {
	// rand.Seed(time.Now().UnixNano())

	// fmt.Println(randSeq(10))
	return "!12345678g"
}
