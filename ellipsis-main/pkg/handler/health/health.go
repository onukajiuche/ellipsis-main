package health

import (
	"encoding/json"
	"fmt"
	"net/http"

	"brief/internal/model"
	"brief/service/ping"
	"brief/utility"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Validate      *validator.Validate
	Logger        *log.Logger
	HealthService ping.HealthService
}

func NewController(validate *validator.Validate, logger *log.Logger, hService ping.HealthService) *Controller {
	return &Controller{
		validate, logger, hService,
	}
}

// Post godoc
//
//	@Summary		check api health
//	@Description	check api health
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Param			ping	body		model.Ping	true	"Ping"
//	@Success		200		{object}	utility.Response
//	@Failure		400		{object}	utility.Response
//	@Router			/health [post]
func (base *Controller) Post(w http.ResponseWriter, r *http.Request) {
	var req model.Ping

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := base.Validate.Struct(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validate), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if !base.HealthService.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	base.Logger.Info("ping successfull")

	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successfull", req.Message)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// Get godoc
//
//	@Summary		check api health
//	@Description	check api health
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utility.Response
//	@Failure		400	{object}	utility.Response
//	@Router			/health [get]
func (base *Controller) Get(w http.ResponseWriter, r *http.Request) {
	if !base.HealthService.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	base.Logger.Info("ping successfull")
	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successfull", map[string]interface{}{"user": "user object"})
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
