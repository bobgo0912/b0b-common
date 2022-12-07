package server

import (
	"context"
	"github.com/bobgo0912/b0b-common/internal/constant"
	"github.com/bobgo0912/b0b-common/internal/log"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"io"
	"net/http"
)

const ProtoMessageContextName = "protoMessage"

type Router struct {
	R *mux.Router
}

func NewRouter() *Router {
	return &Router{R: mux.NewRouter()}
}

func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	return r.R.HandleFunc(path, f)
}

type Fp func(req any, resp any) func(http.ResponseWriter, *http.Request)

func F(req proto.Message, f func(req proto.Message, w http.ResponseWriter)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		protoName := r.Header[constant.ProtoHeader]
		if len(protoName) < 1 {
			http.Error(w, "bad proto", http.StatusBadRequest)
			return
		}
		ip := RemoteIp(r)
		log.Info(ip)
		all, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		message := GetMsgV2ByFullName(protoName[0], all)
		r = r.WithContext(context.WithValue(r.Context(), ProtoMessageContextName, message))
		req = message
		f(message, w)
	}
}
func RemoteIp(req *http.Request) string {
	remoteAddr := req.Header.Get("Remote_addr")
	if remoteAddr == "" {
		if ip := req.Header.Get("ipv4"); ip != "" {
			remoteAddr = ip
		} else if ip = req.Header.Get("XForwardedFor"); ip != "" {
			remoteAddr = ip
		} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
			remoteAddr = ip
		} else {
			remoteAddr = req.Header.Get("X-Real-Ip")
		}
	}

	if remoteAddr == "::1" || remoteAddr == "" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
func (r *Router) HandleProtoFunc(path string, f func(req proto.Message, w http.ResponseWriter), req proto.Message) *mux.Route {
	return r.R.HandleFunc(path, F(req, f))
}
func F1(req any, f func(req any, w http.ResponseWriter)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		protoName := r.Header[constant.ProtoHeader]
		if len(protoName) < 1 {
			http.Error(w, "bad proto", http.StatusBadRequest)
			return
		}
		all, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		message := GetMsgV2ByFullName1(protoName[0], all)
		r = r.WithContext(context.WithValue(r.Context(), ProtoMessageContextName, message))
		req = message
		f(req, w)
	}
}

func (r *Router) HandleProtoFunc1(path string, f func(req any, w http.ResponseWriter), req proto.Message) *mux.Route {
	return r.R.HandleFunc(path, F1(&req, f))
}

func (r *Router) Use(mwf ...mux.MiddlewareFunc) {
	r.R.Use(mwf...)
}

func GrpcMide() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		protoName := r.Header[constant.ProtoHeader]
		if len(protoName) < 1 {
			http.Error(w, "bad proto", http.StatusBadRequest)
			return
		}
		all, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		message := GetMsgV2ByFullName(protoName[0], all)
		ret := r.WithContext(context.WithValue(r.Context(), ProtoMessageContextName, message))
		r = ret
	})
}

func GrpcMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		protoName := r.Header[constant.ProtoHeader]
		if len(protoName) < 1 {
			http.Error(w, "bad proto", http.StatusBadRequest)
			return
		}
		all, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		message := GetMsgV2ByFullName(protoName[0], all)
		ret := r.WithContext(context.WithValue(r.Context(), ProtoMessageContextName, message))
		next.ServeHTTP(w, ret)
	})
}
func GetMsgV2ByFullName(fullName string, data []byte) proto.Message {
	msgName := protoreflect.FullName(fullName)
	messageType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if err != nil {
		log.Warn("GMV2ByFN fail err=", err)
		return nil
	}
	message := messageType.New().Interface()
	err = proto.Unmarshal(data, message)
	if err != nil {
		log.Warn("GMV2ByFN Unmarshal fail err=", err)
		return nil
	}
	return message
}
func GetMsgV2ByFullName1(fullName string, data []byte) interface{} {
	msgName := protoreflect.FullName(fullName)
	messageType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if err != nil {
		log.Warn("GMV2ByFN fail err=", err)
		return nil
	}
	message := messageType.New().Interface()
	err = proto.Unmarshal(data, message)
	if err != nil {
		log.Warn("GMV2ByFN Unmarshal fail err=", err)
		return nil
	}
	return &message
}
