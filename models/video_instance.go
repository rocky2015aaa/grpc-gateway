package models

import (
	"reflect"

	"github.com/sirupsen/logrus"

	"gitlab.com/Fratbe/addglee/src/server/logger"
)

const (
	MN_CREATE_VIDEO_INSTANCE       = ""
	MN_UPDATE_VIDEO_INSTANCE       = ""
	MN_GET_VIDEO_INSTANCE_BY_VALUE = ""

	TB_VIDEO_INSTANCE = ""
)

func (db DB) CreateYoutubeVideoInstance(youtubeVideoInstance *YoutubeVideoInstance) error {
	err := db.Database.Table(TB_VIDEO_INSTANCE).Create(youtubeVideoInstance).Error
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_CREATE_VIDEO_INSTANCE, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return err
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_CREATE_VIDEO_INSTANCE, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("videoInstance:%+v\n", youtubeVideoInstance)
	return nil
}

func (db DB) UpdateYoutubeVideoInstance(youtubeVideoInstance *YoutubeVideoInstance, val map[string]interface{}) error {
	err := db.Database.Table(TB_VIDEO_INSTANCE).Where(val).Update(youtubeVideoInstance).Error
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_UPDATE_VIDEO_INSTANCE, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return err
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_UPDATE_VIDEO_INSTANCE, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("videoInstance:%+v\n", youtubeVideoInstance)
	return nil
}

func (db DB) GetYoutubeVideoInstanceByValue(val map[string]interface{}) *YoutubeVideoInstance {
	videoInstance := &YoutubeVideoInstance{}
	db.Database.Table(TB_VIDEO_INSTANCE).Where(val).First(videoInstance)
	if reflect.DeepEqual(YoutubeVideoInstance{}, *videoInstance) {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_GET_VIDEO_INSTANCE_BY_VALUE, logger.SEVERITY: logger.WARN_LEVEL}).Warnf("no video instance with %v\n", val)
		return nil
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_GET_VIDEO_INSTANCE_BY_VALUE, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("video instance: %+v\n", videoInstance)
	return videoInstance
}
