package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	BetaClusterRedisURL = "beta-redis-master.redis.svc.cluster.local:6379"
	SuccessMessage      = "Successful"
)

var (
	RequiredFieldMessage = func(fields ...string) string {
		return fmt.Sprintf("%v required", fields)
	}
	reqid uint64
)

func BoolP(boolValue bool) *bool {
	return &boolValue
}

func CustomJsonMarshal(data interface{}, tag string) ([]byte, error) {
	var json = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 tag,
	}.Froze()

	return json.Marshal(data)
}

func GetTracingID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}

func SetTracingID(ctx context.Context) context.Context {
	uid := uuid.New().String()
	myid := atomic.AddUint64(&reqid, 1)
	requestID := fmt.Sprintf("%s-%06d", uid, myid)
	ctx = context.WithValue(ctx, middleware.RequestIDKey, requestID)
	return ctx
}

const myCharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func RandStr(length int) string {

	b := make([]byte, length)
	for i := range b {
		b[i] = myCharset[seededRand.Intn(len(myCharset))]
	}
	return string(b)
}

func DecodeInterface(input, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}

const float64EqualityThreshold = 1e-9

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
