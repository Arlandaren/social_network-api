package models

import (
	// "time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID               uint   `json:"id" db:"id"`
	Username         string `json:"login" db:"login"`
	Email            string `json:"email" db:"email"`
	Password         string `json:"password" db:"password"`
	// RegistrationDate time.Time `json:"registration_date" db:"created_at"`
	Country          string `json:"countryCode" db:"countryCode"`
	PublicProfile    bool   `json:"isPublic" db:"isPublic"`
	PhoneNumber      string `json:"phone" db:"phone"`
	Image            string `json:"image" db:"image"`
}

type Countries struct {
	ID     uint   `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Alpha2 string `json:"alpha2" db:"alpha2"`
	Alpha3 string `json:"alpha3" db:"alpha3"`
	Region string `json:"region" db:"region"`
}
type Claims struct {
	User_id uint `json:"user_id"`
	jwt.StandardClaims
}
