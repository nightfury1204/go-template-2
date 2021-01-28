package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"bitbucket.org/evaly/go-boilerplate/api/response"
	"bitbucket.org/evaly/go-boilerplate/infra"
	"bitbucket.org/evaly/go-boilerplate/infra/rabbitmq"
)

type SystemController struct {
	db    infra.DB
	kv    infra.KV
	queue rabbitmq.Rabbit
}

func NewSystemController(db infra.DB, kv infra.KV) *SystemController {
	return &SystemController{
		db: db,
		kv: kv,
	}
}

func (s *SystemController) systemCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.connCheck(); err != nil {
		_ = response.ServeJSON(w, http.StatusInternalServerError, nil, nil, err.Error(), nil)
		return
	}
	response.ServeJSONData(w, "ok", http.StatusOK)
	return
}

func (s *SystemController) workerCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.connCheck(); err != nil {
		_ = response.ServeJSON(w, http.StatusInternalServerError, nil, nil, err.Error(), nil)
		return
	}
	if err := s.queue.Ping(); err != nil {
		_ = response.ServeJSON(w, http.StatusInternalServerError, nil, nil, err.Error(), nil)
		return
	}
	response.ServeJSONData(w, "ok", http.StatusOK)
	return
}

func (s *SystemController) apiCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("apiCheck")
	if err := s.connCheck(); err != nil {
		_ = response.ServeJSON(w, http.StatusInternalServerError, nil, nil, err.Error(), nil)
		return
	}
	response.ServeJSONData(w, "ok", http.StatusOK)
	return
}

func (s *SystemController) connCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	log.Println("db ping")
	if err := s.db.Ping(ctx); err != nil {
		return fmt.Errorf("mongo conn error: %v", err)
	}
	log.Println("kv ping")
	if err := s.kv.Ping(ctx); err != nil {
		return fmt.Errorf("redis conn error: %v", err)
	}
	return nil
}
