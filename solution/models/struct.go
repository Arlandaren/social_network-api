package models

type User struct {
    ID              int    `json:"id" db:"id"`
    Username        string `json:"username" db:"username"`
    Email           string `json:"email" db:"email"`
    Password        string `json:"password" db:"password"`
    RegistrationDate string `json:"registration_date" db:"created_at"`
    Country         string `json:"country" db:"country"`
    PublicProfile   bool   `json:"public_profile" db:"public_profile"`
    PhoneNumber     string `json:"phone_number" db:"phone_number"`
    Image           string `json:"image" db:"image"`
}

type Countries struct {
    ID     uint   `json:"id" db:"id"`
    Name   string `json:"name" db:"name"`
    Alpha2 string `json:"alpha2" db:"alpha2"`
    Alpha3 string `json:"alpha3" db:"alpha3"`
    Region string `json:"region" db:"region"`
}