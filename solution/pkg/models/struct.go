package models

import (

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID               uint   `json:"id" db:"id"`
	Username         string `json:"login" db:"login"`
	Email            string `json:"email" db:"email"`
	Password         string `json:"password" db:"password"`
	CountryCode          string `json:"countryCode" db:"countrycode"`
	PublicProfile    bool   `json:"isPublic" db:"ispublic"`
	PhoneNumber      string `json:"phone" db:"phone"`
	Image            string `json:"image" db:"image"`
}
type Profile struct{
    Login string    `json:"login"`
    Email string `json:"email"`
    CountryCode string `json:"countryCode"`
    IsPublic bool `json:"isPublic"`
    Phone string `json:"phone"`
	Image string `json:"image"`
}
type Countries struct {
	ID     uint   `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Alpha2 string `json:"alpha2" db:"alpha2"`
	Alpha3 string `json:"alpha3" db:"alpha3"`
	Region string `json:"region" db:"region"`
}
type CountryResponse struct {
    Name   string `json:"name"`
    Alpha2 string `json:"alpha2"`
    Alpha3 string `json:"alpha3"`
    Region string `json:"region"`
}
type Claims struct {
	User_id uint `json:"user_id"`
	jwt.StandardClaims
}
type EditParameters struct{
	CountryCode string `json:"countryCode"`
	IsPublic bool	`json:"isPublic"`
	Phone string `json:"phone"`
	Image string  `json:"image"`
}