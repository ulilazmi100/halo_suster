package controller

import (
	"halo_suster/db/entities"
	"halo_suster/responses"
	"halo_suster/svc"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	svc svc.UserSvc
}

func NewUserController(svc svc.UserSvc) *UserController {
	return &UserController{svc: svc}
}

type registerResponse struct {
	UserId      string `json:"userId"`
	Nip         int64  `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken,omitempty"`
}

type loginResponse struct {
	UserId      string `json:"userId"`
	Nip         int64  `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type simpleResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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

	respData := registerResponse{
		UserId:      userId,
		Nip:         newUser.Nip,
		Name:        newUser.Name,
		AccessToken: accessToken,
	}

	return ctx.JSON(http.StatusCreated, simpleResponse{
		Message: "User registered successfully",
		Data:    respData,
	})
}

func (c *UserController) Login(ctx echo.Context) error {
	var user entities.RegistrationPayload
	if err := ctx.Bind(&user); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	loginPayload := entities.Credential{
		Nip:      entities.Int64ToString(user.Nip),
		Password: user.Password,
	}

	userId, name, accessToken, err := c.svc.Login(ctx.Request().Context(), loginPayload)
	if err != nil {
		return err
	}

	respData := loginResponse{
		UserId:      userId,
		Nip:         user.Nip,
		Name:        name,
		AccessToken: accessToken,
	}

	return ctx.JSON(http.StatusOK, simpleResponse{
		Message: "User logged in successfully",
		Data:    respData,
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

	respData := registerResponse{
		UserId: userId,
		Nip:    newUser.Nip,
		Name:   newUser.Name,
	}

	return ctx.JSON(http.StatusCreated, simpleResponse{
		Message: "User registered successfully",
		Data:    respData,
	})
}

func (c *UserController) NurseLogin(ctx echo.Context) error {
	var user entities.RegistrationPayload
	if err := ctx.Bind(&user); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	loginPayload := entities.Credential{
		Nip:      entities.Int64ToString(user.Nip),
		Password: user.Password,
	}

	userId, name, accessToken, err := c.svc.NurseLogin(ctx.Request().Context(), loginPayload)
	if err != nil {
		return err
	}

	respData := loginResponse{
		UserId:      userId,
		Nip:         user.Nip,
		Name:        name,
		AccessToken: accessToken,
	}

	return ctx.JSON(http.StatusOK, simpleResponse{
		Message: "User logged in successfully",
		Data:    respData,
	})
}

func (c *UserController) NurseUpdate(ctx echo.Context) error {
	id := ctx.Param("userId")
	var updatePayload entities.NurseUpdatePayload
	if err := ctx.Bind(&updatePayload); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if err := c.svc.UpdateNurse(ctx.Request().Context(), id, updatePayload); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (c *UserController) NurseDelete(ctx echo.Context) error {
	id := ctx.Param("userId")

	if err := c.svc.DeleteNurse(ctx.Request().Context(), id); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (c *UserController) NurseAccess(ctx echo.Context) error {
	id := ctx.Param("userId")
	var accessPayload entities.NurseAccessPayload
	if err := ctx.Bind(&accessPayload); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if err := c.svc.AccessNurse(ctx.Request().Context(), id, accessPayload); err != nil {
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

	resp, err := c.svc.GetUser(ctx.Request().Context(), user)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return ctx.JSON(http.StatusOK, simpleResponse{
			Message: "success",
			Data:    []interface{}{},
		})
	}

	return ctx.JSON(http.StatusOK, simpleResponse{
		Message: "success",
		Data:    resp,
	})
}
