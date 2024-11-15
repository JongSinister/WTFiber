package models

import (
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" validate:"required"`
	Tel       string             `bson:"tel" validate:"required"`
	Email     string             `bson:"email" valodate:"required"`
	Role      string             `bson:"role" valodate:"required"`
	Password  string             `bson:"password" valodate:"required"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty"`
}

// Check Email Validation
func (user *User) ValidateEmail() bool {
	regex := `^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	re := regexp.MustCompile(regex)
	return re.MatchString(user.Email)
}

// HashPassword hashes the user's password using bcrypt
func (user *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}

// CheckPasword compares a hashed password with the provided password
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// GenerateAuthToken generates a new JWT token for the user
func (user *User) GenerateToken(secret string) (string, error) {
	// 1) Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	// 2) Sign the token with the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
