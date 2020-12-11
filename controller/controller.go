package controller

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joshbetz/config"
	"github.com/sirupsen/logrus"

	"gitlab.com/Fratbe/addglee/src/server/auth"
	"gitlab.com/Fratbe/addglee/src/server/i18n"
	"gitlab.com/Fratbe/addglee/src/server/logger"
	"gitlab.com/Fratbe/addglee/src/server/models"
)

var H Handler

const (
	FN_CONTROLLER_INIT = "controller.LoadAppConfig"
	AUTHORIZATION      = "authorization"
)

// Package controller provides restful api handler functions
type Handler struct {
	models.DataHandler
}

type AppConfig struct {
	Version                  string `json:""`
	DeveloperKey             string `json:""`
	JwtKey                   string `json:""`
	JwtExpirationTimeMinutes int    `json:""`
}

var AppConf *AppConfig

func init() {
	LoadAppConfig("")
}

func LoadAppConfig(configFilePath string) {
	c := config.New(configFilePath)
	AppConf = &AppConfig{}
	c.Get("", &AppConf.Version)
	c.Get("", &AppConf.DeveloperKey)
	c.Get("", &AppConf.JwtKey)

	var jwtExpirationTimeMinutes string
	c.Get("", &jwtExpirationTimeMinutes)
	AppConf.JwtExpirationTimeMinutes, _ = strconv.Atoi(jwtExpirationTimeMinutes)
	if AppConf.JwtExpirationTimeMinutes < 1 || AppConf.JwtExpirationTimeMinutes > 1200 {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_CONTROLLER_INIT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln("Config file: check the jwt_expiration_time_minutes value")
	}
}

func InputValidation(valid bool, err error) error {
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(i18n.T("invalid_request"))
	}
	return nil
}

func IssueToken(claims *auth.Claims) (string, error) {
	expirationTime := time.Now().Add(time.Duration(AppConf.JwtExpirationTimeMinutes) * time.Minute)
	claims.StandardClaims.ExpiresAt = expirationTime.Unix()

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(AppConf.JwtKey))
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: "", logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln(err)
		return "", err
	}
	return tokenString, nil
}
