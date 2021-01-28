package handler

import (
	"context"
	"net/http"
	"time"

	"bitbucket.org/evaly/go-boilerplate/logger"
	"bitbucket.org/evaly/go-boilerplate/rpcrestproxy"
	"bitbucket.org/evaly/go-boilerplate/rpcs/pb"
	"bitbucket.org/evaly/go-boilerplate/utils"
	"github.com/go-chi/chi"
	"google.golang.org/grpc/status"
)

// BrandHandler ...
type BrandHandler struct {
	rpcClients *rpcrestproxy.GoBoilerplateClients
	lgr        logger.StructLogger
}

// NewBrandHandler  ...
func NewBrandHandler(rpcClients *rpcrestproxy.GoBoilerplateClients, lgr logger.StructLogger) *BrandHandler {
	return &BrandHandler{
		rpcClients: rpcClients,
		lgr:        lgr,
	}
}

// GetRouter ...
func (bh *BrandHandler) GetRouter() *chi.Mux {
	r := chi.NewRouter()
	// r.Use(middleware.AuthenticatedOnly)
	r.Get("/", bh.ListBrand)

	return r
}

func (bh *BrandHandler) ListBrand(w http.ResponseWriter, r *http.Request) {
	tid := utils.GetTracingID(r.Context())

	req := &pb.ListBrandReq{}
	if err := ParseToProto(r, req); err != nil {
		bh.lgr.Errorln("listBrands", tid, err.Error())
		ServeJSON(w, "", http.StatusInternalServerError, "unable to parse body", nil, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	res, err := bh.rpcClients.BrandServiceClient.ListBrand(ctx, req)
	if err != nil {
		bh.lgr.Errorln("listBrands", tid, err.Error())
		e, ok := status.FromError(err)
		if !ok {
			ServeJSON(w, "E_INTERNAL", http.StatusInternalServerError, "", nil, nil)
			return
		}

		ServeJSON(w, Code(e.Code().String()), rpcrestproxy.HTTPStatusFromCode(e.Code()), e.Message(), nil, e.Details())
		return
	}

	bh.lgr.Println("listBrands", tid, "sending response")
	ServeJSON(w, "", http.StatusOK, "Successful", res, nil)
	return
}
