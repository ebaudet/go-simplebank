package utils

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const passphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

// RandomString generates a random string of length n
func RandomString(length int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < length; i++ {
		sb.WriteByte(alphabet[rand.Intn(k)])
	}
	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 2500)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := SupportedCurrencies()
	return currencies[rand.Intn(len(currencies))]
}

// RandomEmail generates a random email address.
func RandomEmail() string {
	return RandomString(6) + "@" + RandomString(int(RandomInt(3, 5))) + ".com"
}

// Generate a random password with numbers, letters and special characters
func RandomPassword(length int) string {
	var sb strings.Builder
	k := len(passphabet)
	for i := 0; i < length; i++ {
		sb.WriteByte(passphabet[rand.Intn(k)])
	}
	return sb.String()
}

// RandomFullName returns a random full name
func RandomFullName() string {
	caser := cases.Title(language.English)
	return caser.String(RandomString(6) + " " + RandomString(6))
}
