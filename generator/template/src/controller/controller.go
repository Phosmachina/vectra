package controller

import (
	"Vectra/src/model/i18n"
	. "Vectra/src/model/service"
	. "Vectra/src/model/storage"
	"bytes"
	"context"
	. "github.com/Phosmachina/FluentKV/reldb"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	. "github.com/gofiber/fiber/v2/middleware/session"
	"io"
	"strings"
)

var (
	conform  = modifiers.New()
	validate = validator.New()
)

type Controller struct {
	router fiber.Router
	store  *Store
}

func NewController(router fiber.Router, store *Store) Controller {
	return Controller{router: router, store: store}
}

// HandleView is a function handling view logic for a certain page.
// It first retrieves the user session from the store, using the given context.
// If session retrieval is successful, it gets the API instance and extracts the user ID from the session.
// It then checks if the user has access to the requested path and if not, it returns a 403 Forbidden error.
// If a user has access, it creates a new bytes.Buffer and calls the provided writer function.
// The writer function is supposed to write the required data into the provided buffer.
// If writing is successful, it sets the context's content type to 'text/html; charset=UTF-8'
// and sends the buffer's bytes as the response. If writing fails, it returns the error.
//
// Parameters:
//
// - ctx: It's a fiber.Ctx pointer, the context of the request.
//
// - controller: It's an instance of ViewController having the store to get the session.
//
// - writer: It's a custom function to write into the provided io.Writer. It should handle the writing logic based on the provided user ID.
//
// Returns:
//
// - error: It returns an error in case something goes wrong. It would be 'nil' for successful execution.
func HandleView(ctx *fiber.Ctx, controller ViewController,
	writer func(buf io.Writer, userId string) error) error {

	sess, err := controller.store.Get(ctx)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	userId := sess.Get(SessionKeyForUserId).(string)
	if !GetApiV1().GetAccessManager().CheckAccessForRoute(userId,
		ctx.Path()) {
		return fiber.ErrForbidden
	}

	var buf = new(bytes.Buffer)
	err = writer(buf, userId)
	if err != nil {
		return err
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

	return ctx.Send(buf.Bytes())
}

// HandleRequest is a generic function designed to handle HTTP requests within a fiber context.
// This function is flexible and allows you to provide specific logic via higher-order functions
// and handle different types of objects and errors.
//
// This function takes four parameters:
//
// ctx: A pointer to a fiber context which holds all request information and methods to control the flow of a HTTP request.
//
// useInfo: A function that takes an object of given type T and returns an error and a pointer to an ObjWrapper of type K.
// This function encapsulates how we want to use the provided information, specifying both the potential error that could arise and the wrapped object we may get.
//
// errToReason: A mapping from error messages to string reasons.
// This map provides a clean way to handle possible errors and respond with appropriate reasons.
//
// onSuccess: a function that takes a pointer to an ObjWrapper of type K.
// This function executes if there are no errors from useInfo function.
// The ObjWrapper that is passed to this function is the one that is returned from the useInfo function.
//
// This function will return an error if something goes wrong during the processing of the HTTP request, otherwise, it will return nil.
func HandleRequest[T any, K IObject](ctx *fiber.Ctx, useInfo func(T) (error,
	*ObjWrapper[K]), onSuccess func(*ObjWrapper[K])) error {

	_i18n := i18n.GetInstance()

	var data T
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ReasonExch{
			Reason: i18n.Error.InvalidRequestStructure(),
		})
	}

	conform.Struct(context.Background(), &data)
	if validate.Struct(data) != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ReasonExch{
			Reason: i18n.Error.InvalidDataStructure(),
		})
	}

	err, k := useInfo(data)
	r := ReasonExch{}

	if err != nil {
		errName, _ := strings.CutPrefix(err.Error(), "Error")
		r.Reason = _i18n.Get("error." + errName)
		return ctx.Status(fiber.StatusBadRequest).JSON(r)
	} else {
		if onSuccess != nil {
			onSuccess(k)
		}
		return ctx.JSON(r)
	}
}

func checkAccess[T IObject](ctx *fiber.Ctx, sess *Session, idOfT string) error {

	userId := sess.Get(SessionKeyForUserId).(string)

	if !CheckAccessForTable[T](userId, idOfT) {
		return ctx.Status(fiber.StatusForbidden).JSON(ReasonExch{
			Reason: i18n.Error.InsufficientRoleLevel(),
		})
	}

	return nil
}
