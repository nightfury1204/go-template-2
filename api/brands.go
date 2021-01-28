package api

import (
	"net/http"

	"bitbucket.org/evaly/go-boilerplate/api/response"
	"bitbucket.org/evaly/go-boilerplate/logger"
	"bitbucket.org/evaly/go-boilerplate/service"
	"bitbucket.org/evaly/go-boilerplate/utils"
)

// BrandsController ...
type BrandsController struct {
	svc service.BrandService
	lgr logger.StructLogger
}

// NewBrandsController ...
func NewBrandsController(svc service.BrandService) *BrandsController {
	return &BrandsController{
		svc: svc,
	}
}

// SetLogger ...
func (cc *BrandsController) SetLogger(lgr logger.StructLogger) {
	cc.lgr = lgr
}

// ListBrand ...
func (cc *BrandsController) ListBrand(w http.ResponseWriter, r *http.Request) {
	tid := utils.GetTracingID(r.Context())
	pageQ, skipQ, limitQ, err := parseSkipLimit(r, 10, 100)
	if err != nil {
		cc.lgr.Errorln("listBrands", tid, err.Error())
		_ = response.ServeJSON(w, http.StatusBadRequest, nil, nil, err.Error(), nil)
		return
	}
	pager := utils.Pager{Skip: skipQ, Limit: limitQ}

	cc.lgr.Println("listBrands", tid, "getting brands")
	result, err := cc.svc.ListBrand(r.Context(), pager)
	if err != nil {
		cc.lgr.Errorln("listBrands", tid, err.Error())
		_ = response.ServeJSON(w, http.StatusInternalServerError, nil, nil, err.Error(), nil)
		return
	}

	cc.lgr.Println("listBrands", tid, "sending response")
	prev, next := getNextPreviousPager(r.URL.Path, pageQ, limitQ)
	_ = response.ServeJSON(w, http.StatusOK, prev, next, utils.SuccessMessage, result)
	return
}
