package main

import (
	"flag"
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
	zhTranslate "github.com/go-playground/validator/v10/translations/zh"
)

type User1 struct {
	Name  string `validate:"required"`
	Age   uint   `validate:"gte=1,lte=130"`
	Email string `validate:"required,email"`
}

var (
	name1  string
	age1   uint
	email1 string
)

func init() {
	flag.StringVar(&name1, "name", "", "输入名字")
	flag.UintVar(&age1, "age", 0, "输入年龄")
	flag.StringVar(&email1, "email", "", "输入邮箱")
}

func main() {
	flag.Parse()

	user := &User1{
		Name:  name1,
		Age:   age1,
		Email: email1,
	}

	validate := validator.New()

	e := en.New()
	uniTrans := ut.New(e, e, zh.New(), zh_Hant_TW.New())
	translator, _ := uniTrans.GetTranslator("zh")
	zhTranslate.RegisterDefaultTranslations(validate, translator)

	err := validate.Struct(user)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, err := range errs {
			fmt.Println(err.Translate(translator))
		}
	}
}