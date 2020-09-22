package main

import (
	"errors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"

	zhTranslate "github.com/go-playground/validator/v10/translations/zh"
)

type UserEcho struct {
	Name  string `json:"name" validate:"required"`
	Age   uint   `json:"age" validate:"gte=1,lte=130"`
	Email string `json:"email" validate:"required,email"`
}

type Error struct {
	Error string `json:"error"`
}

type Validator struct {
	trans     ut.Translator
	validator *validator.Validate
}

// Validate do validation for request value.
func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	msg := ""
	for _, v := range errs.Translate(v.trans) {
		if msg != "" {
			msg += ", "
		}
		msg += v
	}
	return errors.New(msg)
}

func main() {
	e := echo.New()
	e.Debug = true

	// setup japanese translation
	vvv := validator.New()

	english := en.New()
	uniTrans := ut.New(english, english, zh.New(), zh_Hant_TW.New())
	translator, _ := uniTrans.GetTranslator("zh")
	zhTranslate.RegisterDefaultTranslations(vvv, translator)

	e.Validator = &Validator{validator: vvv, trans: translator}

	e.Any("/", func(ctx echo.Context) error {
		user := new(UserEcho)
		if err := ctx.Bind(user); err != nil {
			return err
		}

		if err := ctx.Validate(user); err != nil {
			return ctx.JSON(http.StatusBadRequest, &Error{Error: err.Error()})
		}

		return ctx.JSON(http.StatusOK, user)
	})

	err := e.Start(":8888")

	if err != nil {
		log.Fatalln("run error :", err)
	}
}