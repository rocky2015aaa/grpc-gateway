package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/joshbetz/config"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"

	"gitlab.com/Fratbe/addglee/src/server/controller"
	"gitlab.com/Fratbe/addglee/src/server/logger"
	"gitlab.com/Fratbe/addglee/src/server/middlewares"
	"gitlab.com/Fratbe/addglee/src/server/models"
	"gitlab.com/Fratbe/addglee/src/server/responses"
	"gitlab.com/Fratbe/addglee/src/server/util"
)

var (
	DB              models.DB
	applicationHost string
)

const (
	INIT                     = ""
	FN_MAIN                  = ""
	DATABASE_CONNECTION_INFO = ""

	APPLICATION_CONFIG_KEY_LIST          = ""
	APPLICATON_OPERATION_CONFIG_KEY_LIST = ""
)

func init() {
	var conf string
	util.SetConfig(&conf)
	c := config.New(conf)

	applicationConfigKeyList := strings.Split(APPLICATION_CONFIG_KEY_LIST, ",")
	applicationConfigValueList := make([]string, len(applicationConfigKeyList))
	for i, applicationConfigKey := range applicationConfigKeyList {
		c.Get(applicationConfigKey, &applicationConfigValueList[i])
		if applicationConfigValueList[i] == "" {
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: INIT, logger.SEVERITY: logger.ERROR_LEVEL}).Errorf("Config file: the %s parameter is missing", applicationConfigKey)
			os.Exit(0)
		}
	}

	applicationOperationConfigKeyList := strings.Split(APPLICATON_OPERATION_CONFIG_KEY_LIST, ",")
	applicationOperationConfigValueList := make([]string, len(applicationOperationConfigKeyList))
	for i, applicationOperationConfigKey := range applicationOperationConfigKeyList {
		c.Get(applicationOperationConfigKey, &applicationOperationConfigValueList[i])
		if applicationOperationConfigValueList[i] == "" {
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: INIT, logger.SEVERITY: logger.ERROR_LEVEL}).Errorf("Config file: the %s parameter is missing", applicationOperationConfigKey)
			os.Exit(0)
		}
	}

	applicationHost = applicationOperationConfigValueList[0]
	initDBConfig := fmt.Sprintf(DATABASE_CONNECTION_INFO, applicationOperationConfigValueList[1], applicationOperationConfigValueList[2], applicationOperationConfigValueList[3], applicationOperationConfigValueList[4], applicationOperationConfigValueList[5], applicationOperationConfigValueList[6])

	DB = models.DB{
		Database: models.InitDB(initDBConfig),
	}

	go util.WatchingSignal()
}

func main() {
	grpcIpPort := applicationHost + ":8991"
	go func() {
		ctx := context.Background()
		mux := http.NewServeMux()
		mux.HandleFunc("", controller.UploadFile)
		m := new(util.JSONPb)
		gwmux := runtime.NewServeMux(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, m),
			runtime.WithForwardResponseOption(responses.HttpResponseModifier),
		)
		err := middlewares.RegisterHandlersFromEndpoint(ctx, gwmux, grpcIpPort)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_MAIN, logger.SEVERITY: logger.WARN_LEVEL}).Warnf("failed to start HTTP gateway: %v\n", err)
			return
		}
		mux.Handle("/", gwmux)

		srv := &http.Server{
			Addr:    applicationHost + ":",
			Handler: responses.ResponseHandler(mux, controller.AppConf.JwtKey),
		}

		srv.ListenAndServe()
	}()

	conn, err := net.Listen("", grpcIpPort)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_MAIN, logger.SEVERITY: logger.WARN_LEVEL}).Warnf("%v\n", err)
	}
	controller.H = controller.Handler{DB}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(), // server runtime error(panic) recovery
		)),
	)
	middlewares.RegisterServiceServer(s, &controller.H)

	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_MAIN, logger.SEVERITY: logger.INFO_LEVEL}).Infof("started")
	if err := s.Serve(conn); err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_MAIN, logger.SEVERITY: logger.WARN_LEVEL}).Warnf("%v\n", err)
	}
}
