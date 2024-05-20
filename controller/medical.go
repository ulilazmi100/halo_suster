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

type MedicalController struct {
	svc svc.MedicalSvc
}

func NewMedicalController(svc svc.MedicalSvc) *MedicalController {
	return &MedicalController{svc: svc}
}

func (c *MedicalController) RegisterPatient(ctx echo.Context) error {
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

func (c *MedicalController) GetPatient(ctx echo.Context) error {
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

func (c *MedicalController) RegisterRecord(ctx echo.Context) error {
	var newRecord entities.RecordRegistrationPayload
	if err := ctx.Bind(&newRecord); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	// Retrieve the values from the context
	userID := ctx.Get("user_id").(string)
	nip := ctx.Get("nip").(string)
	name := ctx.Get("name").(string)

	// Create the CreatedByDetail struct with the retrieved values
	nipInt, err := entities.StringToInt64(nip)
	if err != nil {
		return err
	}
	createdByDetail := entities.CreatedByDetail{
		UserId: userID,
		Nip:    nipInt,
		Name:   name,
	}

	err = c.svc.RegisterRecord(ctx.Request().Context(), newRecord, createdByDetail)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusCreated)
}

func (c *MedicalController) GetRecord(ctx echo.Context) error {
	var record entities.GetRecordQueries
	if err := ctx.Bind(&record); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if record.Limit == 0 {
		record.Limit = 5
	}

	if record.Limit < 0 || record.Offset < 0 {
		return responses.NewBadRequestError("invalid query param")
	}

	requestCtx, cancel := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
	defer cancel()

	resp, err := c.svc.GetRecord(requestCtx, record)
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
