package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"gitlab.com/Fratbe/addglee/src/server/logger"
)

const (
	MN_GET_ACCOUNT            = ""
	MN_CREATE_ACCOUNT         = ""
	MN_SET_ACCOUNT_FIELD_DATA = ""
	FN_GET_ACCOUNT_BY_PSEUDO  = ""

	TB_ACCOUNT = ""
)

func (db DB) CreateUserAccount(userAccount *UserAccount) error {
	if userAccount.Modified == nil {
		now := time.Now()
		userAccount.Modified = &now
	}
	err := db.Database.Table(TB_ACCOUNT).Create(userAccount).Error
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_CREATE_ACCOUNT, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return err
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_CREATE_ACCOUNT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("userAccount:%+v\n", userAccount)
	return nil
}

func (db DB) GetUserAccount(condition, val string) *UserAccount {
	userAccount := &UserAccount{}
	if db.Database.Table(TB_ACCOUNT).Where(condition, val).First(&userAccount).RecordNotFound() {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_GET_ACCOUNT, logger.SEVERITY: logger.WARN_LEVEL}).Warnln("no userAccount with value:", val)
		return nil
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_GET_ACCOUNT, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("userAccount:%+v\n", userAccount)
	return userAccount
}

func GetUserAccountByPseudo(tx *gorm.DB, pseudo string) *UserAccount {
	userAccount := &UserAccount{}
	if tx.Table(TB_ACCOUNT).Where("", pseudo).First(&userAccount).RecordNotFound() {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_GET_ACCOUNT_BY_PSEUDO, logger.SEVERITY: logger.WARN_LEVEL}).Warnln("no userAccount with pseudo:", pseudo)
		return nil
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_GET_ACCOUNT_BY_PSEUDO, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("userAccount:%+v\n", userAccount)
	return userAccount
}

func (db DB) SetUserAccountFieldData(pseudonym, field string, val interface{}) error {
	var err error
	switch val.(type) {
	case string:
		err = db.Database.Table(TB_ACCOUNT).Where("", pseudonym).Update(field, val.(string)).Error
	case time.Time:
		err = db.Database.Table(TB_ACCOUNT).Where("", pseudonym).Update(field, val.(time.Time)).Error
	}
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_SET_ACCOUNT_FIELD_DATA, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return err
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_SET_ACCOUNT_FIELD_DATA, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugln("update. field:", field, "value:", val)
	return nil
}
