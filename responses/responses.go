package responses

import (
	"context"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"

	"gitlab.com/Fratbe/addglee/src/server/auth"
	"gitlab.com/Fratbe/addglee/src/server/logger"
)

const (
	FN_RESPONSE_HANDLER = ""
)

func ResponseHandler(h http.Handler, jwtKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode, err := auth.JwtHandler(r, jwtKey)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_RESPONSE_HANDLER, logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
			http.Error(w, err.Error(), statusCode)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func HttpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get(""); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		w.WriteHeader(code)
		delete(md.HeaderMD, "")
		delete(w.Header(), "")
	}

	return nil
}
