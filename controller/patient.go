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

type PatientController struct {
	svc svc.PatientSvc
}

func NewPatientController(svc svc.PatientSvc) *PatientController {
	return &PatientController{svc: svc}
}

func (c *PatientController) RegisterPatient(ctx echo.Context) error {
	var newPatient entities.PatientRegistrationPayload
	if err := ctx.Bind(&newPatient); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	err := c.svc.RegisterPatient(ctx.Request().Context(), newPatient)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusCreated)
}

func (c *PatientController) GetPatient(ctx echo.Context) error {
	var patient entities.GetPatientQueries
	if err := ctx.Bind(&patient); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if patient.Limit == 0 {
		patient.Limit = 5
	}

	if patient.Limit < 0 || patient.Offset < 0 {
		return responses.NewBadRequestError("invalid query param")
	}

	requestCtx, cancel := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
	defer cancel()

	resp, err := c.svc.GetPatient(requestCtx, patient)
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
