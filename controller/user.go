package controller

import (
	"context"
	"halo_suster/db/entities"
	"halo_suster/responses"
	"halo_suster/svc"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	svc svc.UserSvc
}

func NewUserController(svc svc.UserSvc) *UserController {
	return &UserController{svc: svc}
}

func (c *UserController) Register(ctx echo.Context) error {
	var newUser entities.RegistrationPayload
	if err := ctx.Bind(&newUser); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	userId, accessToken, err := c.svc.Register(ctx.Request().Context(), newUser)
	if err != nil {
		return err
	}

	respData := map[string]interface{}{
		"userId":      userId,
		"nip":         newUser.Nip,
		"name":        newUser.Name,
		"accessToken": accessToken,
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"data":    respData,
	})
}

func (c *UserController) Login(ctx echo.Context) error {
	var user entities.RegistrationPayload
	if err := ctx.Bind(&user); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	// Set a timeout for the context
	requestCtx, cancel := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
	defer cancel()

	loginPayload := entities.Credential{
		Nip:      user.Nip,
		Password: user.Password,
	}

	userId, name, accessToken, err := c.svc.Login(requestCtx, loginPayload)
	if err != nil {
		return err
	}

	respData := map[string]interface{}{
		"userId":      userId,
		"nip":         user.Nip,
		"name":        name,
		"accessToken": accessToken,
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "User logged in successfully",
		"data":    respData,
	})
}

func (c *UserController) NurseRegister(ctx echo.Context) error {
	var newUser entities.NurseRegistrationPayload
	if err := ctx.Bind(&newUser); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	userId, err := c.svc.NurseRegister(ctx.Request().Context(), newUser)
	if err != nil {
		return err
	}

	respData := map[string]interface{}{
		"userId": userId,
		"nip":    newUser.Nip,
		"name":   newUser.Name,
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"data":    respData,
	})
}

func (c *UserController) NurseLogin(ctx echo.Context) error {
	var user entities.RegistrationPayload
	if err := ctx.Bind(&user); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	// Set a timeout for the context
	requestCtx, cancel := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
	defer cancel()

	loginPayload := entities.Credential{
		Nip:      user.Nip,
		Password: user.Password,
	}

	userId, name, accessToken, err := c.svc.NurseLogin(requestCtx, loginPayload)
	if err != nil {
		return err
	}

	respData := map[string]interface{}{
		"userId":      userId,
		"nip":         user.Nip,
		"name":        name,
		"accessToken": accessToken,
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "User logged in successfully",
		"data":    respData,
	})
}

func (c *UserController) NurseUpdate(ctx echo.Context) error {
	id := ctx.Param("userId")
	var updatePayload entities.NurseUpdatePayload
	if err := ctx.Bind(&updatePayload); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	// Set a timeout for the context
	requestCtx, cancel := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
	defer cancel()

	err := c.svc.UpdateNurse(requestCtx, id, updatePayload)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (c *UserController) NurseDelete(ctx echo.Context) error {
	id := ctx.Param("userId")

	// Set a timeout for the context
	requestCtx, cancel := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
	defer cancel()

	err := c.svc.DeleteNurse(requestCtx, id)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (c *UserController) NurseAccess(ctx echo.Context) error {
	id := ctx.Param("userId")
	var acessPayload entities.NurseAccessPayload
	if err := ctx.Bind(&acessPayload); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	// Set a timeout for the context
	requestCtx, cancel := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
	defer cancel()

	err := c.svc.AccessNurse(requestCtx, id, acessPayload)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (c *UserController) GetUser(ctx echo.Context) error {
	var user entities.GetUserQueries
	if err := ctx.Bind(&user); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if user.Limit == 0 {
		user.Limit = 5
	}

	if user.Limit < 0 || user.Offset < 0 {
		return responses.NewBadRequestError("invalid query param")
	}

	// Set a timeout for the context
	requestCtx, cancel := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
	defer cancel()

	resp, err := c.svc.GetUser(requestCtx, user)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "success",
			"data":    []interface{}{},
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    resp,
	})
}
