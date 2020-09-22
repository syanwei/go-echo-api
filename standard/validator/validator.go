package main

import (
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Name  string `validate:"required"`
	Age   uint   `validate:"gte=1,lte=130"`
	Email string `validate:"required,email"`
}

var (
	name  string
	age   uint
	email string
)

func init() {
	flag.StringVar(&name, "name", "", "输入名字")
	flag.UintVar(&age, "age", 0, "输入年龄")
	flag.StringVar(&email, "email", "", "输入邮箱")
}

func main() {
	flag.Parse()

	user := &User{
		Name:  name,
		Age:   age,
		Email: email,
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		fmt.Println(err)
	}
}