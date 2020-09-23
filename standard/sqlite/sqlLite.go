package main

import (
	"database/sql"
	"errors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"os"
	"time"

	zhTranslate "github.com/go-playground/validator/v10/translations/zh"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"gopkg.in/gorp.v2"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/xeodou/go-sqlcipher"
)

var dbDriver = "sqlite3"

// Error indicate response error
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

// Comment is a struct to hold unit of request and response.
type Comment struct {
	Id      int64     `json:"id" db:"id,primarykey,autoincrement"`
	Name    string    `json:"name" form:"name" db:"name,notnull,size:200"`
	Text    string    `json:"text" form:"text" validate:"required,max=20" db:"text,notnull,size:399"`
	Created time.Time `json:"created" db:"created,notnull"`
	Updated time.Time `json:"updated" db:"updated,notnull"`
}

// PreInsert update fields Created and Updated.
func (c *Comment) PreInsert(s gorp.SqlExecutor) error {
	c.Created = time.Now()
	c.Updated = c.Created
	return nil
}

// PreInsert update field Updated.
func (c *Comment) PreUpdate(s gorp.SqlExecutor) error {
	c.Updated = time.Now()
	return nil
}

func setupDB() (*gorp.DbMap, error) {
	tmpStr := "file:test" + ".s3db?_auth&_auth_user=admin&_auth_pass=admin&_auth_crypt=sha256"
	db, err := sql.Open("sqlite3", tmpStr)
	//db, err := sql.Open(dbDriver, "test")
	if err != nil {
		return nil, err
	}
	dialect := gorp.SqliteDialect{}
	dbMap := &gorp.DbMap{Db: db, Dialect: dialect}
	dbMap.AddTableWithName(Comment{}, "comments").SetKeys(true, "id")
	err = dbMap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}
	return dbMap, nil
}

// Controller is a controller for this application.
type Controller struct {
	dbMap *gorp.DbMap
}

// InsertComment is GET handler to return record.
func (controller *Controller) GetComment(c echo.Context) error {
	var comment Comment
	// fetch record specified by parameter id
	err := controller.dbMap.SelectOne(&comment,
		"SELECT * FROM comments WHERE id = $1", c.Param("id"))
	if err != nil {
		if err != sql.ErrNoRows {
			c.Logger().Error("SelectOne: ", err)
			return c.String(http.StatusBadRequest, "SelectOne: "+err.Error())
		}
		return c.String(http.StatusNotFound, "Not Found")
	}
	return c.JSON(http.StatusOK, comment)
}

// InsertComment is GET handler to return records.
func (controller *Controller) ListComments(c echo.Context) error {
	var comments []Comment
	// fetch last 10 records
	_, err := controller.dbMap.Select(&comments,
		"SELECT * FROM comments ORDER BY created desc LIMIT 10")
	if err != nil {
		c.Logger().Error("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, comments)
}

// InsertComment is POST handler to insert record.
func (controller *Controller) InsertComment(c echo.Context) error {
	var comment Comment
	// bind request to comment struct
	if err := c.Bind(&comment); err != nil {
		c.Logger().Error("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}
	// validate request
	if err := c.Validate(&comment); err != nil {
		c.Logger().Error("Validate: ", err)
		return c.JSON(http.StatusBadRequest, &Error{Error: err.Error()})
	}
	// insert record
	if err := controller.dbMap.Insert(&comment); err != nil {
		c.Logger().Error("Insert: ", err)
		return c.String(http.StatusBadRequest, "Insert: "+err.Error())
	}
	c.Logger().Infof("inserted comment: %v", comment.Id)
	return c.NoContent(http.StatusCreated)
}

func main() {
	dbMap, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	controller := &Controller{dbMap: dbMap}

	// setup echo
	e := echo.New()
	e.Debug = true
	validate := validator.New()

	english := en.New()
	uniTrans := ut.New(english, zh.New())
	translator, _ := uniTrans.GetTranslator("zh")
	_ = zhTranslate.RegisterDefaultTranslations(validate, translator)

	e.Validator = &Validator{validator: validate, trans: translator}
	e.Logger.SetOutput(os.Stderr)

	e.GET("/api/comments/:id", controller.GetComment)
	e.GET("/api/comments", controller.ListComments)
	e.POST("/api/comments", controller.InsertComment)
	e.Logger.Fatal(e.Start(":8989"))
}
