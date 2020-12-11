// Package models handles database connection and management
package models

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"

	"gitlab.com/Fratbe/addglee/src/server/logger"
)

const (
	TB_STR          = ""
	FN_INIT_DB      = ""
	MN_GET_SUBQUERY = ""
)

type DataHandler interface {
	CreateUserAccount(userAccount *UserAccount) error
	CreateDiscipline(str *Str, discipline *Discipline, disciplineDescription *DisciplineDescription) error
	GetUserAccount(condition, val string) *UserAccount
	GetDisciplines() []*DisciplineDetail
	GetDisciplinesByUser(userAccountId uint32) []*DisciplineDetail
	GetDisciplineById(id uint32) *DisciplineDetail
	CreateLevel(str *Str, level *Level) error
	GetLevel(id uint32) *LevelWithName
	GetLevels(disciplineId uint32) []*LevelWithName
	UpdateLevelActivation(disciplineId uint32, name string, active bool, inactiveReason string) error
	GetOrAddStringId(tx *gorm.DB, languageId uint32, str ...string) (stringId uint32, err error)
	GetStringFromId(stringId uint32) (languageId uint32, str string, err error)
	CreateTechnique(str *Str, technique *Technique) error
	GetTechniques(disciplineId uint32) []*TechniqueWithName
	GetTechniqueById(techniqueId uint32) *TechniqueWithName
	UpdateTechniqueActivation(userAccountId uint32, nameEn1 string, active bool, inactiveReason string) error
	SetUserAccountFieldData(pseudonym, field string, val interface{}) error
	CreateCurriculum(str *Str, curriculum *Curriculum) error
	GetCurriculumByLevelId(levelId uint32) *Curriculum
	GetCurriculumById(id uint32) *Curriculum
	CreateCurriculumContent(cc *CurriculumContent) error
	GetCurriculumContentsByTechniqueId(techniqueId uint32) []*CurriculumContent
	GetCurriculumContentsByCurriculumId(curriculumId uint32) []*CurriculumContent
	GetVideoPracs(condition map[string]interface{}) []*VideoPrac
	GetVideoPracsByIds(ids []uint32) []*VideoPrac
	GetVideoPracDescriptionList(ids []uint32) []*VideoPrac
	UpdateVideoPrac(videoPracId uint32, updateValues map[string]interface{}) error
	GetVideoPracListPropositionUser(condition string, id uint32) []*VideoPrac
	GetVideoPracById(videoPracId uint32) *VideoPrac
	AddVideoPracProposition(videoPracProposition *VideoPracProposition) error
	GetVideoPracPropositionById(videoPracPropositionId uint32) *VideoPracProposition
	UpdateVideoPracProposition(videoPracProposition *VideoPracProposition) error
	GetVideoPracTechniqueById(videoPracTechniqueId uint32) *VideoPracTechnique
	UpdateVideoPracTechnique(videoPracTechniqueId uint32, updateValues map[string]interface{}) error
	GetVideoTypeById(videoTypeId uint32) *VideoType
	CreateYoutubeVideoInstance(youtubeVideoInstance *YoutubeVideoInstance) error
	GetYoutubeVideoInstanceByValue(val map[string]interface{}) *YoutubeVideoInstance
	UpdateYoutubeVideoInstance(youtubeVideoInstance *YoutubeVideoInstance, val map[string]interface{}) error
	CreateVideoSource(videoSource *VideoSource) error
	GetVideoSource(name string) *VideoSource
	GetLanguageByLetterCode(code string) *Language
	CreateNewVideoSegment(videoSegment *VideoSegment) error
	GetVideoSegmentById(videoSegmentId uint32) *VideoSegment
	CreateTechVideoSegment(tvs *TechVideoSegment) error
	GetTechVideoSegmentByTechniqueId(techniqueId uint32) []*TechVideoSegment
	CheckIfVideoPracTitleExistsForUser(title string) bool
	GetUserDashboard(userAccountId uint32) *UserDashboard
}

type DB struct {
	Database *gorm.DB
}

func InitDB(dataSourceName string) *gorm.DB {
	db, err := gorm.Open("", dataSourceName)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_INIT_DB, logger.SEVERITY: logger.PANIC_LEVEL}).Panicln(err)
		os.Exit(0)
	}
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_INIT_DB, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("database:%+v\n", db)
	return db
}

func (db DB) getSubQuery(selectTarget, table, whereCondition string, val interface{}) *gorm.SqlExpr {
	subQuery := db.Database.Select(selectTarget).Table(table).Where(whereCondition, val).SubQuery()
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_GET_SUBQUERY, logger.SEVERITY: logger.INFO_LEVEL}).Infoln(subQuery)
	return subQuery
}
