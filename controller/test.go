package controller

import (
	"sort"

	"github.com/jinzhu/gorm"

	"gitlab.com/Fratbe/addglee/src/server/models"
)

type TestDB map[string]interface{}

var (
	mockDB = &TestDB{
		"king": &models.UserAccount{
			Id:          1,
			LastName:    "Snow",
			FirstName:   "Jon",
			Pseudo:      "king",
			Password:    "$2a$10$aY4BQ4iw6pbXHKMs3bi.yOClCN79NP4E288GJxGo7GLJDnJON97UW",
			Activated:   true,
			FakeAccount: false,
		},
		"str_100": &models.Str{
			Id:         100,
			LanguageId: 4,
			Str:        "white belt",
			Comment:    "level name",
		},
		"level_100": &models.Level{
			Id:             100,
			Level:          1,
			DisciplineId:   1,
			Name:           100,
			Active:         true,
			InactiveReason: "",
		},
		"str_200": &models.Str{
			Id:         200,
			LanguageId: 4,
			Str:        "throw body",
			Comment:    "technique name",
		},
		"technique_200": &models.Technique{
			Id:               200,
			FormIsometricp:   true,
			DisciplineOrigin: 1,
			Plyometric:       true,
			Name:             200,
			UserAccountId:    1,
			Active:           true,
		},
		"videoprac_1": &models.VideoPrac{
			Id:             1,
			Trainee:        1,
			MainDiscipline: 1,
			VideoTypeId:    1,
			Title:          "test",
			Language:       1,
		},
		"videoprac_2": &models.VideoPrac{
			Id:             2,
			Trainee:        2,
			MainDiscipline: 2,
			VideoTypeId:    2,
			Title:          "test2",
			Language:       2,
		},
		"videoprac_proposition_1": &struct {
			Id                 uint32
			VideopracTechnique uint32
			SourceUser         uint32
			TargetUser         uint32
		}{
			1,
			1,
			1,
			2,
		},
		"videoprac_technique_1": &struct {
			Id          uint32
			VideopracId uint32
		}{
			1,
			3,
		},
		"videoprac_3": &models.VideoPrac{
			Id:             3,
			Trainee:        1,
			MainDiscipline: 3,
			VideoTypeId:    3,
			Title:          "test3",
			Language:       3,
		},

		/*
				"piano": &models.Discipline{
					Id:   1,
					Name: "piano",
					//Description: "piano1",
				},
				"disciplineId/10000": &models.Discipline{
					Id:   10000,
					Name: "football",
					//Description: "football10000",
				},
				"disciplineId/5": &models.Discipline{
					Id:   5,
					Name: "football",
					//Description: "football5",
				},
				"disciplineId/10": &models.Discipline{
					Id:   10,
					Name: "football",
					//Description: "football10",
				},
				"disciplineId/20": &models.Discipline{
					Id:   20,
					Name: "baseball",
					//Description: "baseball20",
				},
				"disciplineId/30": &models.Discipline{
					Id:   30,
					Name: "basketball",
					//Description: "basketball30",
				},
			"disciplineIdForLevelorsub/5": &models.Levelorsub{
				Id:             5,
				NameEn:         "basic",
				Level:          5,
				Sublevel:       1,
				DisciplineId:   5,
				Active:         true,
				InactiveReason: "",
			},
			"disciplineIdForLevelorsub/10000": &models.Levelorsub{
				Id:             7,
				NameEn:         "basic",
				Level:          5,
				Sublevel:       1,
				DisciplineId:   5,
				Active:         true,
				InactiveReason: "",
			},
			"1": &models.Technique{
				Id:               1,
				FormIsometricp:   true,
				DisciplineOrigin: 1,
				NameEn_1:         "type",
				Plyometric:       true,
				OldNewNoise:      true,
			},
			"1/fire": &models.Technique{
				Id:               1,
				FormIsometricp:   true,
				DisciplineOrigin: 1,
				NameEn_1:         "fire",
				Plyometric:       true,
				OldNewNoise:      true,
				Active:           false,
				InactiveReason:   "",
			},
			"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/XKqBUSyPcMUfxS027TmNHgWsM8g": &models.YoutubeVideoInstance{
				Id: 1,
			},
			"ABC_SUCCESS10": &models.YoutubeVideoInstance{
				ContentdetailsDuration: "PT1H10M2S",
			},
			"ABCDE_ERROR10": &models.YoutubeVideoInstance{
				ContentdetailsDuration: "PT1T10M2S",
			},
			"idForVideoSegment/1": &models.VideoSegment{
				Id:               1,
				VideoInstanceId:  1,
				StartCentisecond: 140,
				StopCentisecond:  200,
			},
			"idForTechnique/2": &models.Technique{
				Id:               2,
				FormIsometricp:   true,
				DisciplineOrigin: 1,
				NameEn_1:         "type",
				Plyometric:       true,
				OldNewNoise:      true,
			},
			"idForTechnique/12": &models.Technique{
				Id:               12,
				FormIsometricp:   true,
				DisciplineOrigin: 1,
				NameEn_1:         "type2",
				Plyometric:       true,
				OldNewNoise:      true,
			},
			"idForTechnique/13": &models.Technique{
				Id:               13,
				FormIsometricp:   true,
				DisciplineOrigin: 1,
				NameEn_1:         "type2",
				Plyometric:       true,
				OldNewNoise:      true,
			},
			"idForTechnique/10000": &models.Technique{
				Id: 10000,
			},
			"idForDiscipline/1/idForLevel/1": &models.Curriculum{
				Id:           1,
				Active:       true,
				LevelId:      1,
				DisciplineId: 1,
			},
			"idForDiscipline/5/idForLevel/5": &models.Curriculum{
				Id:           5,
				Active:       true,
				LevelId:      5,
				DisciplineId: 5,
			},
			"idForDiscipline/12/idForLevel/12": &models.Curriculum{
				Id:           12,
				Active:       true,
				LevelId:      12,
				DisciplineId: 12,
			},
			"idForDiscipline/13/idForLevel/13": &models.Curriculum{
				Id:           13,
				Active:       true,
				LevelId:      13,
				DisciplineId: 13,
			},
			"idForDiscipline/14/idForLevel/14": &models.Curriculum{
				Id:           14,
				Active:       true,
				LevelId:      14,
				DisciplineId: 14,
			},
			"idForCurriculum/12/idForTechnique/12": &models.CurriculumContent{
				CurriculumId: 12,
				TechniqueId:  12,
			},
			"idForCurriculum/12/idForTechnique/13": &models.CurriculumContent{
				CurriculumId: 12,
				TechniqueId:  13,
			},
			"idForCurriculum/14/idForTechnique/14": &models.CurriculumContent{
				CurriculumId: 14,
				TechniqueId:  14,
			},
			"idForTForVs/111/idForVideoSegment/111": &models.TechVideoSegment{
				TechniqueId:    111,
				VideoSegmentId: 111,
			},
			"idForTForVs/111/idForVideoSegment/12": &models.TechVideoSegment{
				TechniqueId:    111,
				VideoSegmentId: 122,
			},
			"idForVideoSegment/111": &models.VideoSegment{
				Id:               111,
				VideoInstanceId:  111,
				StartCentisecond: 140,
				StopCentisecond:  200,
			},
			"idForVideoSegment/122": &models.VideoSegment{
				Id:               122,
				VideoInstanceId:  122,
				StartCentisecond: 140,
				StopCentisecond:  200,
			},
			"idForTForVs/99/idForVideoSegment/12": &models.TechVideoSegment{
				TechniqueId:    99,
				VideoSegmentId: 100,
			},
		*/
	}
)

type VideoInstanceTest struct{}

func (vit VideoInstanceTest) lurk_video(h *Handler, url string, userAccountId uint32) (uint32, string, error) {
	return 1, "video_instance_creation", nil
}

func getDBIndexKeys(testDB *TestDB) []string {
	keys := []string{}
	for k, _ := range *testDB {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (testDB *TestDB) GetOrAddStringId(tx *gorm.DB, languageId uint32, str ...string) (stringId uint32, err error) {
	return stringId, err
}

func (testDB *TestDB) GetStringFromId(stringId uint32) (languageId uint32, str string, err error) {
	return languageId, str, err
}
