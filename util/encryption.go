// Package util works as useful utilities for the application.
// Encryption and so on.
package util

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/Fratbe/addglee/src/server/logger"
)

const (
	FN_HASH_GENERATE = "util.HashGenerate"
	FN_HASH_COMPARE  = "util.HashCompare"
)

func HashGenerate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_HASH_GENERATE, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return "", err
	}

	hash := string(hashedBytes[:])
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_HASH_GENERATE, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln("hash:", hash)
	return hash, nil
}

func HashCompare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	err := bcrypt.CompareHashAndPassword(existing, incoming)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_HASH_COMPARE, logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
	}
	return err
}
