package middlewares

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"gitlab.com/Fratbe/addglee/src/server/controller"
	"gitlab.com/Fratbe/addglee/src/server/logger"
)

const (
	FN_REGISTER_HANDLER_FROM_ENDPOINT = "RegisterHandlersFromEndpoint"
	FN_HTTP_REQUEST_HANDLER           = "HttpRequestHandler"
)

func RegisterHandlersFromEndpoint(ctx context.Context, mux *runtime.ServeMux, grpcIpPort string) error {
	err := controller.RegisterUserAccountServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}
	err = controller.RegisterDisciplineServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}
	err = controller.RegisterLevelServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}
	err = controller.RegisterTechniqueServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}
	err = controller.RegisterVideoInstanceServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}
	err = controller.RegisterVideoSegmentServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}
	err = controller.RegisterCurriculumServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}
	err = controller.RegisterCurriculumContentServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}
	err = controller.RegisterTechVideoSegmentServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}
	err = controller.RegisterVersionServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}
	err = controller.RegisterVideoPracServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024 * 1024 * 30))})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}

	err = controller.RegisterUserDashboardServiceHandlerFromEndpoint(ctx, mux, grpcIpPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REGISTER_HANDLER_FROM_ENDPOINT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return err
	}

	return nil
}

func RegisterServiceServer(s *grpc.Server, h *controller.Handler) {
	controller.RegisterUserAccountServiceServer(s, h)
	controller.RegisterDisciplineServiceServer(s, h)
	controller.RegisterLevelServiceServer(s, h)
	controller.RegisterTechniqueServiceServer(s, h)
	controller.RegisterVideoInstanceServiceServer(s, h)
	controller.RegisterVideoSegmentServiceServer(s, h)
	controller.RegisterCurriculumServiceServer(s, h)
	controller.RegisterCurriculumContentServiceServer(s, h)
	controller.RegisterTechVideoSegmentServiceServer(s, h)
	controller.RegisterVersionServiceServer(s, h)
	controller.RegisterVideoPracServiceServer(s, h)
	controller.RegisterUserDashboardServiceServer(s, h)
}
