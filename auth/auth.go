// Package dio for jwt management
package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"gitlab.com/Fratbe/addglee/src/server/i18n"
	"gitlab.com/Fratbe/addglee/src/server/logger"
)

const (
	FN_PARSE_JWT_TOKEN = "dio.ParseJwtToken"
)

type Claims struct {
	UserAccount string `json:""`
	jwt.StandardClaims
}

func JwtHandler(r *http.Request, jwtKey string) (int, error) {
	if r.URL.String() != "/" && r.URL.String() != "/" && r.URL.String() != "/" && r.URL.String() != "/" {
		reqToken := r.Header.Get("Authorization")
		_, err := ParseJwtToken(jwtKey, reqToken)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: "", logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
			if err.Error() == i18n.T("jwt_parsing_problem") {
				return http.StatusBadRequest, err
			}
			return http.StatusUnauthorized, err
		}
	}
	return 0, nil
}

func ParseJwtToken(jwtKey, token string) (*Claims, error) {
	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) != 2 {
		return nil, errors.New(i18n.T("no_bearer"))
	}
	reqToken := strings.TrimSpace(splitToken[1])

	claims := &Claims{}

	jwtToken, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_PARSE_JWT_TOKEN, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return claims, err
	}
	if !jwtToken.Valid {
		return nil, errors.New(i18n.T("jwt_parsing_problem"))
	}
	return claims, nil
}
