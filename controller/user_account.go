// Package controller provides restful api handler functions
package controller

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"gitlab.com/Fratbe/addglee/src/server/auth"
	"gitlab.com/Fratbe/addglee/src/server/i18n"
	"gitlab.com/Fratbe/addglee/src/server/logger"
	"gitlab.com/Fratbe/addglee/src/server/models"
	"gitlab.com/Fratbe/addglee/src/server/util"
)

var editableFields = map[string]struct{}{
	"": struct{}{},
	"": struct{}{},
	"": struct{}{},
	"": struct{}{},
}

const (
	FN_SIGNUP_ACCOUNT = "controller.SignUpAccount"
	FN_LOGIN_ACCOUNT  = "controller.LoginAccount"
	FN_UPDATE_ACCOUNT = "controller.UpdateAccount"
	FN_REFRESH_TOKEN  = "controller.RefreshToken"
)

func (h *Handler) SignUpAccount(ctx context.Context, userAccount *models.UserAccount) (*Response, error) {
	// ValidateStruct checks if essential JSON field is empty like "" or invalid value
	valid, err := govalidator.ValidateStruct(userAccount)
	if InputValidation(valid, err) != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_SIGNUP_ACCOUNT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "400"))
		return &Response{Message: i18n.T("invalid_request")}, nil
	}

	// HashGenerate generates bycrypt hash for password
	hashedPassword, err := util.HashGenerate(userAccount.Password)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_SIGNUP_ACCOUNT, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return nil, err
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_SIGNUP_ACCOUNT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("hashedPassword:%+v\n", hashedPassword)

	userAccount.Password = hashedPassword
	now := time.Now()
	userAccount.CreationDate = &now
	if err := h.CreateUserAccount(userAccount); err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_SIGNUP_ACCOUNT, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_SIGNUP_ACCOUNT, logger.SEVERITY: logger.INFO_LEVEL}).Infof("Signup done:%+v\n", userAccount)
	return &Response{Message: i18n.T("account_created")}, nil
}

func (h *Handler) LoginAccount(ctx context.Context, loginData *LoginInformation) (*LoginInformationResponse, error) {
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_LOGIN_ACCOUNT, logger.SEVERITY: logger.INFO_LEVEL}).Infof("test")
	// ValidateStruct checks if essential JSON field is empty like "" or invalid value
	valid, err := govalidator.ValidateStruct(loginData)
	if InputValidation(valid, err) != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_LOGIN_ACCOUNT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "400"))
		return &LoginInformationResponse{Message: i18n.T("invalid_request")}, nil
	}

	userAccount := h.GetUserAccount("", loginData.Pseudo)
	if userAccount == nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_LOGIN_ACCOUNT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln("no user account with pseudonym:", loginData.Pseudo)
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
		return &LoginInformationResponse{Message: i18n.T("no_account")}, nil
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_LOGIN_ACCOUNT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("user account:%+v\n", userAccount)

	// HashCompare validates user password with the password of account in database
	err = util.HashCompare(userAccount.Password, loginData.Password)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_LOGIN_ACCOUNT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_LOGIN_ACCOUNT, logger.SEVERITY: logger.INFO_LEVEL}).Infof("Login successful:%+v\n", userAccount)

	claims := &auth.Claims{
		UserAccount:    userAccount.Pseudo,
		StandardClaims: jwt.StandardClaims{},
	}

	tokenString, err := IssueToken(claims)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_LOGIN_ACCOUNT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
		return nil, err
	}
	// make response without password
	userAccount.Password = ""
	// add success msg
	return &LoginInformationResponse{Message: i18n.T("jwt_creation_success"), AccessToken: tokenString, UserAccount: userAccount}, nil
}

func (h *Handler) UpdateAccount(ctx context.Context, selfCareData *SelfCareInformation) (*Response, error) {
	// ValidateStruct checks if essential JSON field is empty like "" or invalid value
	valid, err := govalidator.ValidateStruct(selfCareData)
	if InputValidation(valid, err) != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_UPDATE_ACCOUNT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "400"))
		return &Response{Message: i18n.T("invalid_request")}, nil
	}
	userAccount := h.GetUserAccount("", selfCareData.Pseudo)
	if userAccount == nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_UPDATE_ACCOUNT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln("no user account with pseudonym:", selfCareData.Pseudo)
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
		return &Response{Message: i18n.T("no_account")}, nil
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_UPDATE_ACCOUNT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("user account:%+v\n", userAccount)

	err = util.HashCompare(userAccount.Password, selfCareData.Password)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_UPDATE_ACCOUNT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
		return nil, err
	}

	if _, ok := editableFields[selfCareData.Field]; !ok {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_UPDATE_ACCOUNT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln("no editable fields:", selfCareData.Field)
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
		return &Response{Message: i18n.T("")}, nil
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_UPDATE_ACCOUNT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("editable fields:%+v\n", editableFields)

	if selfCareData.Field == "password" {
		hashedPassword, err := util.HashGenerate(selfCareData.Value)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_UPDATE_ACCOUNT, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
			return nil, err
		}
		selfCareData.Value = hashedPassword
	}

	if err := h.SetUserAccountFieldData(userAccount.Pseudo, selfCareData.Field, selfCareData.Value); err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_UPDATE_ACCOUNT, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_UPDATE_ACCOUNT, logger.SEVERITY: logger.INFO_LEVEL}).Infoln("updated. field:", selfCareData.Field, "value:", selfCareData.Value)
	return &Response{Message: i18n.T("")}, nil
}

func (h *Handler) RefreshToken(ctx context.Context, refreshToken *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	mdata, _ := metadata.FromIncomingContext(ctx)
	claim, err := auth.ParseJwtToken(AppConf.JwtKey, mdata["authorization"][0])
	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors != jwt.ValidationErrorExpired {
			grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "400"))
			return &RefreshTokenResponse{Message: err.Error()}, nil
		}
	}
	tokenString, err := IssueToken(claim)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_REFRESH_TOKEN, logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
		return nil, err
	}
	return &RefreshTokenResponse{Message: i18n.T(""), AccessToken: tokenString}, nil
}
