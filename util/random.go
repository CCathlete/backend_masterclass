package u

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

/*
We don't need the random numbers to be cryptographically safe so it's
enough to use math/rand.
*/
var (
	seed    = time.Now().UnixNano()
	source  = rand.NewSource(seed)
	randGen = rand.New(source)

	mailProviders = []string{"gmail", "yahoo", "outlook"}
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

// Returns a random int64 between min and max.
func RandomInt(min, max int64) int64 {
	// We add 1 because Int63n panics if its argument <= 0
	// and we need to cover the case of min = max.
	return min + randGen.Int63n(max-min+1)
}

func RandomStr(n int) string {
	var sb strings.Builder
	k := len(alphabet) // Number of possible options for a random character.

	for i := 0; i < n; i++ {
		// We choose an index randomly and alphabet[index].
		randChar := alphabet[rand.Intn(k)]
		sb.WriteByte(randChar)
	}

	return sb.String()
}

// Returns a random account owner name.
func RandomOwner() string {
	return RandomStr(6)
}

// Returns a random amount of money.
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandCurrency() string {
	currencies := []string{"EUR", "USD", "ILS"}
	l := len(currencies)
	return currencies[rand.Intn(l)]
}

func RandomUsername() string {
	return RandomStr(6)
}

func RandomEmail() (email string) {
	mailProvider := mailProviders[RandomInt(0, 2)]
	email = RandomStr(6) + "@" + mailProvider + ".com"
	return
}

func RandomFullName() string {
	return RandomStr(6) + " " + RandomStr(6)
}

func RandomPassword() string {
	return RandomStr(8)
}

func HashPassword(password string) (hash string) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	hash = string(bcryptPassword)
	return
}
