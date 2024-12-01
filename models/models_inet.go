package models

import "gorm.io/gorm"

type Person struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type User struct {
	Name     string `json:"name" validate:"required,min=3,max=32"`
	IsActive *bool  `json:"isactive" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required,email,min=3,max=32"`
}

type Dogs struct {
	gorm.Model
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
}

type NewUser struct {
	Email        string `json:"email,omitempty" validate:"required,email,min=3,max=32"`
	UserName     string `json:"username" validate:"required,min=6,max=20"`
	Password     string `json:"password" validate:"required,min=3,max=32"`
	LineID       string `json:"lineid,omitempty" validate:"required,min=3,max=32"`
	PhoneNumber  string `json:"phonenumber" validate:"required,min=3,max=32"`
	BusinessType string `json:"businesstype" validate:"required,min=3,max=32"`
	WebSite      string `json:"website" validate:"required,min=2,max=30"`
}

type DogsRes struct {
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
	Type  string `json:"type"`
}

type Companys struct {
	gorm.Model
	CompanyID    int    `json:"company_id"`
	Name         string `json:"name"`
	BusinessType string `json:"businesstype"`
}

type Profile struct {
	gorm.Model
	EmployeeID int    `json:"employee_id"`
	Name       string `json:"name"`
	LastName   string `json:"lastname"`
	BirthDay   int    `json:"birthday"`
	Age        int    `json:"age"`
	Email      string `json:"email"`
	Tel        int    `json:"tel"`
}
