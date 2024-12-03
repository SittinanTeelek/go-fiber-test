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

type Dog struct {
	gorm.Model
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
}

type NewUser struct {
	Email        string `json:"email" validate:"required,email"`
	UserName     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	LineID       string `json:"lineid,omitempty" validate:"required"`
	PhoneNumber  string `json:"phonenumber" validate:"required"`
	BusinessType string `json:"businesstype" validate:"required"`
	WebSite      string `json:"website" validate:"required"`
}

// {
//     "email":"sads@addfas.com",
// 	"username":"asdf-",
// 	"password":"123313",
// 	"lineid":"qwewqewq",
// 	"phonenumber":"213123",
// 	"businesstype":"341232",
// 	"website":"21321 "
// }

type DogsRes struct {
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
	Type  string `json:"type"`
}

type ResultDogData struct {
	Data       []DogsRes `json:"data"`
	Name       string    `json:"name"`
	Count      int       `json:"count"`
	SumRed     int       `json:"sum_red"`
	SumGreen   int       `json:"sum_green"`
	SumPink    int       `json:"sum_pink"`
	SumNoColor int       `json:"sum_nocolor"`
}

type Companies struct {
	gorm.Model
	CompanyID    int    `json:"company_id"`
	Name         string `json:"name"`
	BusinessType string `json:"businesstype"`
}

// "company_id":001,
// "name":"sdaas",
// "businesstype":"asdas"

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

// "employee_id":"",
// "name":"",
// "lastname":"",
// "birthday":"",
// "age":"",
// "email":"",
// "tel":"",

type ResultProfileData struct {
	Data         []Profile `json:"data"`
	Count        int       `json:"count"`
	GenZ         int       `json:"genz"`
	GenY         int       `json:"geny"`
	GenX         int       `json:"genx"`
	BabyBoomer   int       `json:"babyboomer"`
	GIGeneration int       `json:"gi_generation"`
}
