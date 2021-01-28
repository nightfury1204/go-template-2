package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ParseToProto ...
func ParseToProto(r *http.Request, m proto.Message) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return protojson.Unmarshal(b, m)
}

// Code represents response code
type Code string

// Response reponse serializer util
type Response struct {
	Code    Code        `json:"code,omitempty"`
	Status  int         `json:"-"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ServeJSON serves json to http client
func (r *Response) ServeJSON(w http.ResponseWriter) error {
	resp := &Response{
		Code:    r.Code,
		Status:  r.Status,
		Message: r.Message,
		Data:    r.Data,
		Errors:  r.Errors,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

type protoJSON struct {
	Message protoreflect.ProtoMessage
}

// ResponseProtoJSON reponse serializer util
type ResponseProtoJSON struct {
	Code    Code        `json:"code,omitempty"`
	Status  int         `json:"-"`
	Message string      `json:"message,omitempty"`
	Data    protoJSON   `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func (p protoJSON) MarshalJSON() ([]byte, error) {
	po := protojson.MarshalOptions{
		EmitUnpopulated: false,
		UseProtoNames:   true,
	}
	jsnRes, err := po.Marshal(p.Message)
	if err != nil {
		return nil, err
	}

	return jsnRes, nil
}

// ServeJSON a utility func which serves json to http client
func ServeJSON(w http.ResponseWriter, code Code, status int, message string, data interface{}, errors interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var resp interface{}

	p, ok := data.(protoreflect.ProtoMessage)
	if ok {
		resp = &ResponseProtoJSON{
			Code:    code,
			Status:  status,
			Message: message,
			Data:    protoJSON{Message: p},
			Errors:  errors,
		}
	} else {
		resp = &Response{
			Code:    code,
			Status:  status,
			Message: message,
			Data:    data,
			Errors:  errors,
		}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}
