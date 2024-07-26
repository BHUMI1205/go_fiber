package model

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Id        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" validate:"required,max=20"`
	Email     string `json:"email" validate:"required,email" gorm:"unique"`
	Password  string `json:"password" validate:"required,min=6,max=6"`
	Carts     []Cart `gorm:"foreignKey:UserId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category struct {
	gorm.Model
	Id        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" validate:"required"`
	Products  []Product `gorm:"foreignKey:CategoryId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	gorm.Model
	Id         uint     `json:"id" gorm:"primaryKey"`
	Name       string   `json:"name" validate:"required"`
	CategoryId uint     `json:"categoryId" gorm:"not null" validate:"required"`
	Category   Category `gorm:"foreignKey:CategoryId"`
	Image      string   `json:"image" validate:"required"`
	Carts      []Cart   `gorm:"foreignKey:ProductId"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Cart struct {
	gorm.Model
	Id        uint    `json:"id" gorm:"primaryKey"`
	ProductId uint    `json:"productId" gorm:"not null" validate:"required"`
	Product   Product `gorm:"foreignKey:ProductId"`
	UserId    uint    `json:"userId" gorm:"not null" validate:"required"`
	User      User    `gorm:"foreignKey:UserId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Token struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}
