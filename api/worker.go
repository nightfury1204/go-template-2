package api

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/evaly/go-boilerplate/api/response"
	"bitbucket.org/evaly/go-boilerplate/logger"
	"bitbucket.org/evaly/go-boilerplate/model"
	"bitbucket.org/evaly/go-boilerplate/service"
	"bitbucket.org/evaly/go-boilerplate/utils"
)

// WorkerController ...
type WorkerController struct {
	svc service.BrandService
	lgr logger.StructLogger
}

// NewWorkerController...
func NewWorkerController(svc service.BrandService) *WorkerController {
	return &WorkerController{
		svc: svc,
	}
}

// SetLogger ...
func (wc *WorkerController) SetLogger(lgr logger.StructLogger) {
	wc.lgr = lgr
}

// AddNewBrand ...
func (c *WorkerController) AddNewBrand(w http.ResponseWriter, r *http.Request) {
	tid := utils.GetTracingID(r.Context())
	brand := &model.BrandInfo{}

	if err := json.NewDecoder(r.Body).Decode(&brand); err != nil {
		c.lgr.Errorf("addNewBrand", tid, "Failed to unmarshal, %v", err.Error())
	}
	if err := c.svc.InsertNewBrandFromEvent(r.Context(), brand); err != nil {
		_ = response.ServeJSON(w, http.StatusInternalServerError, nil, nil, err.Error(), nil)
		c.lgr.Errorf("addNewBrand", tid, "Failed to insert product, %v", err.Error())
	}
	_ = response.ServeJSON(w, http.StatusOK, nil, nil, utils.SuccessMessage, "")
	c.lgr.Println("addNewProduct", tid, "sending response")
}
