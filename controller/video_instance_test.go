package controller

import (
	//	"context"
	"encoding/json"
	"fmt"
	//	"io/ioutil"
	//"os"
	//	"reflect"
	//	"testing"

	//	"gitlab.com/Fratbe/addglee/src/server/i18n"
	"gitlab.com/Fratbe/addglee/src/server/models"
	youtube "google.golang.org/api/youtube/v3"
)

var vmd = `{
 "kind": "youtube#videoListResponse",
 "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/f5LSa2dKT3ZZuSThDmELgNrO1_o\"",
 "pageInfo": {
  "totalResults": 1,
  "resultsPerPage": 1
 },
 "items": [
  {
   "kind": "youtube#video",
   "etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/sMlN93o69Sk6wZ4bLUvquI1sfng\"",
   "id": "R0WcFxtKFj8",
   "snippet": {
    "publishedAt": "2013-04-28T14:06:34.000Z",
    "channelId": "UCV5hEOemT-y3o-5pd0QmauQ",
    "title": "How to Do a Roundhouse Kick | Kickboxing Lessons",
    "description": "Side Kicks for Your Kickboxing Lifestyle\r\nSubliminal Training - Kickboxing CD (Ultrasonic Martial Arts Series):http://amzn.to/1iN9SLd\r\nTiger Claw Feiyue Martial Arts Shoes: http://amzn.to/1OV6AQj   \r\nShock Doctor Gel Max Mouth Guard: http://amzn.to/1Y9lUih  \r\nKettlebell Kickboxing: Scorcher Video Series: http://amzn.to/1OcVbOD \r\nTaekwondo Durable Kick Pad for Training: http://amzn.to/1NBwoUq  \r\nMuay Thai Punching - Resin Statue: http://amzn.to/1KQFqwd \n\nWatch more How to Do Kickboxing videos: http://www.howcast.com/videos/510172-How-to-Do-a-Roundhouse-Kick-Kickboxing-Lessons\n\n\n\nUnable to read transcription file",
    "thumbnails": {
     "default": {
      "url": "https://i.ytimg.com/vi/R0WcFxtKFj8/default.jpg",
      "width": 120,
      "height": 90
     },
     "medium": {
      "url": "https://i.ytimg.com/vi/R0WcFxtKFj8/mqdefault.jpg",
      "width": 320,
      "height": 180
     },
     "high": {
      "url": "https://i.ytimg.com/vi/R0WcFxtKFj8/hqdefault.jpg",
      "width": 480,
      "height": 360
     },
     "standard": {
      "url": "https://i.ytimg.com/vi/R0WcFxtKFj8/sddefault.jpg",
      "width": 640,
      "height": 480
     },
     "maxres": {
      "url": "https://i.ytimg.com/vi/R0WcFxtKFj8/maxresdefault.jpg",
      "width": 1280,
      "height": 720
     }
    },
    "channelTitle": "HowcastSportsFitness",
    "tags": [
     "how to do a roundhouse kick",
     "kickboxing",
     "kickboxing training",
     "roundhouse kick",
     "howcast",
     "kickboxer",
     "kickboxing workout",
     "cardio kickboxing",
     "kick boxing",
     "kickboxing techniques",
     "kickboxing class",
     "kickboxing moves",
     "kickboxing lessons",
     "kickboxing drills"
    ],
    "categoryId": "17",
    "liveBroadcastContent": "none",
    "localized": {
     "title": "How to Do a Roundhouse Kick | Kickboxing Lessons",
     "description": "Side Kicks for Your Kickboxing Lifestyle\r\nSubliminal Training - Kickboxing CD (Ultrasonic Martial Arts Series):http://amzn.to/1iN9SLd\r\nTiger Claw Feiyue Martial Arts Shoes: http://amzn.to/1OV6AQj   \r\nShock Doctor Gel Max Mouth Guard: http://amzn.to/1Y9lUih  \r\nKettlebell Kickboxing: Scorcher Video Series: http://amzn.to/1OcVbOD \r\nTaekwondo Durable Kick Pad for Training: http://amzn.to/1NBwoUq  \r\nMuay Thai Punching - Resin Statue: http://amzn.to/1KQFqwd \n\nWatch more How to Do Kickboxing videos: http://www.howcast.com/videos/510172-How-to-Do-a-Roundhouse-Kick-Kickboxing-Lessons\n\n\n\nUnable to read transcription file"
    }
   },
   "contentDetails": {
    "duration": "PT5M36S",
    "dimension": "2d",
    "definition": "hd",
    "caption": "false",
    "licensedContent": true,
    "projection": "rectangular"
   },
   "status": {
    "uploadStatus": "processed",
    "privacyStatus": "public",
    "license": "youtube",
    "embeddable": true,
    "publicStatsViewable": true
   },
   "statistics": {
    "viewCount": "1467573",
    "likeCount": "10766",
    "dislikeCount": "556",
    "favoriteCount": "0",
    "commentCount": "770"
   },
   "player": {
    "embedHtml": "\u003ciframe width=\"480\" height=\"270\" src=\"//www.youtube.com/embed/R0WcFxtKFj8\" frameborder=\"0\" allow=\"accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture\" allowfullscreen\u003e\u003c/iframe\u003e"
   },
   "recordingDetails": {
    "recordingDate": "2013-01-24T00:00:00.000Z"
   }
  }
 ]
}`

func (testDB *TestDB) SetId(key string, val uint32) {
	key += "_id"
	(*testDB)[key] = val
}

func (testDB *TestDB) GetId(key string) uint32 {
	key += "_id"
	return (*testDB)[key].(uint32)
}

func (testDB *TestDB) CreateYoutubeVideoInstance(youtubeVideoInstance *models.YoutubeVideoInstance) error {
	if _, ok := (*testDB)["CreateYoutubeVideoInstanceFailure_id"]; ok {
		return fmt.Errorf("error")
	}

	youtubeVideoInstance.Id = (*testDB).GetId("YoutubeVideoInstance")
	return nil
}

func (testDB *TestDB) UpdateYoutubeVideoInstance(youtubeVideoInstance *models.YoutubeVideoInstance, val map[string]interface{}) error {
	if _, ok := (*testDB)["UpdateYoutubeVideoInstanceFailure_id"]; ok {
		return fmt.Errorf("error")
	}

	var value string
	for _, v := range val {
		switch sv := v.(type) {
		case uint32:
			value += fmt.Sprint(sv)
		case string:
			value += v.(string)
		}
	}
	(*testDB)[value] = youtubeVideoInstance
	return nil
}

func (testDB *TestDB) GetYoutubeVideoInstanceByValue(val map[string]interface{}) *models.YoutubeVideoInstance {
	if val, ok := (*testDB)["YoutubeVideoInstanceTest_id"]; ok {
		ymd := []byte(vmd)
		youtubeVideoInstance, _ := createVideoInstance(ymd)
		youtubeVideoInstance.Id = val.(uint32)
		return youtubeVideoInstance
	}
	if _, ok := (*testDB)["UpdateYoutubeVideoInstanceFailure_id"]; ok {
		ymd := []byte(vmd)
		youtubeVideoInstance, _ := createVideoInstance(ymd)
		return youtubeVideoInstance
	}

	var num, value string
	for _, v := range val {
		switch sv := v.(type) {
		case uint32:
			num = fmt.Sprintf("%d", sv)
		case string:
			value += v.(string)
		}
	}
	value += num

	if videoInstance, ok := (*testDB)[value]; ok {
		return videoInstance.(*models.YoutubeVideoInstance)
	}

	return nil
}

func (testDB *TestDB) CreateVideoSource(videoSource *models.VideoSource) error {
	(*testDB)[videoSource.Name] = videoSource
	return nil
}

func (testDB *TestDB) GetVideoSource(name string) *models.VideoSource {
	if videoSource, ok := (*testDB)[name]; ok {
		return videoSource.(*models.VideoSource)
	}
	return &models.VideoSource{
		Id:       1,
		Name:     "youtube.com",
		UrlStart: "https://www.youtube.com",
	}
}

func (testDB *TestDB) GetLanguageByLetterCode(code string) *models.Language {
	return nil
}

/*
func TestCreateVideoInstance(t *testing.T) {
	h := &Handler{mockDB}

	fmt.Println("------- TestSignUpUserAccount Success Case -------")
	msg := fmt.Sprintf(i18n.T("video_instance_creation"), 1)
	tests := []struct {
		userAccountId uint32
		urls          string
		want          string
	}{
		{
			userAccountId: 1,
			urls:          "https://test/watch?v=VLYKx-Fwhxg https://test/watch?v=VLYKx-Fwhxg2",
			want:          msg + " / " + msg,
		},
	}

	for _, tt := range tests {
		req := &VideoInstanceRequest{
			UserAccountId: tt.userAccountId,
			Urls:          tt.urls,
		}
		resp, err := h.CreateVideoInstance(context.Background(), req)
		if err != nil {
			t.Errorf("CreateVideoInstance got unexpected error: %s", err)
		}
		if resp.Msg != tt.want {
			t.Errorf("CreateVideoInstance(%v)=\"%v\", wanted: \"%v\"", req, resp.Msg, tt.want)
		}
	}

	fmt.Println("------- TestSignUpUserAccount Failure Case -------")
	testFailure := []struct {
		userAccountId uint32
		urls          string
		want          string
	}{
		{
			userAccountId: 1,
			urls:          "",
			want:          i18n.T("no_video_instance_creation"),
		},
	}

	for _, tt := range testFailure {
		mockDB.SetId("YoutubeVideoInstance", 1)
		req := &VideoInstanceRequest{
			UserAccountId: tt.userAccountId,
			Urls:          tt.urls,
		}
		resp, _ := h.CreateVideoInstance(context.Background(), req)
		if resp.Msg != tt.want {
			t.Errorf("CreateVideoInstance(%v)=\"%v\", wanted: \"%v\"", req, resp.Msg, tt.want)
		}
	}

}

func TestLurkVideo(t *testing.T) {
	LoadAppConfig("../config/config.json")

	fmt.Println("------- TestLurkVideo Success Case -------")
	testSuccess := []struct {
		testDB        TestDB
		idForSetting  string
		url           string
		userAccountId uint32
		want          uint32
	}{
		{
			testDB:        mockDB,
			idForSetting:  "YoutubeVideoInstance",
			url:           "https://www.youtube.com/watch?v=e3cdiSWpxC0",
			userAccountId: 1,
			want:          1,
		},
		{
			testDB:        mockDB,
			idForSetting:  "YoutubeVideoInstance",
			url:           "https://www.youtube.com/watch?v=VLYKx-Fwhxg",
			userAccountId: 2,
			want:          2,
		},
		{
			testDB:        mockDB,
			idForSetting:  "YoutubeVideoInstanceTest",
			url:           "https://www.youtube.com/watch?v=VLYKx-Fwhxg",
			userAccountId: 2,
			want:          2,
		},
	}

	for _, tt := range testSuccess {
		yvi := YoutubeVideoInstance{}
		tt.testDB.SetId(tt.idForSetting, tt.want)
		h := &Handler{tt.testDB}
		videoInstanceId, _, err := yvi.lurk_video(h, tt.url, tt.userAccountId)
		if err != nil {
			t.Errorf("lurk_video got unexpected error: %s", err)
		}
		if videoInstanceId != tt.want {
			t.Errorf("lurk_video(%v, %s, %d)=%v, wanted %v", h, tt.url, tt.userAccountId, videoInstanceId, tt.want)
		}
	}

	fmt.Println("------- TestLurkVideo Failure Case -------")
	createTestConfigFile("../config/test_config.json")
	testFailure := []struct {
		configPath    string
		testDB        TestDB
		idForSetting  string
		url           string
		userAccountId uint32
		want          uint32
	}{
		{
			configPath:    "../config/config.json",
			testDB:        mockDB,
			idForSetting:  "YoutubeVideoInstance",
			url:           "https://www.google.com/watch?v=e3cdiS",
			userAccountId: 1,
			want:          0,
		},
		{
			configPath:    "../config/config.json",
			testDB:        mockDB,
			idForSetting:  "YoutubeVideoInstance",
			url:           "https://www.google.com/watch?v=e3cdiSWpxC0",
			userAccountId: 1,
			want:          0,
		},
		{
			configPath:    "../config/test_config.json",
			testDB:        mockDB,
			idForSetting:  "YoutubeVideoInstance",
			url:           "https://www.youtube.com/watch?v=VLYKx-Fwhxg",
			userAccountId: 2,
			want:          0,
		},
		{
			configPath:    "../config/config.json",
			testDB:        mockDB,
			idForSetting:  "UpdateYoutubeVideoInstanceFailure",
			url:           "https://www.youtube.com/watch?v=VLYKx-Fwhxg",
			userAccountId: 2,
			want:          0,
		},
		{
			configPath:    "../config/config.json",
			testDB:        mockDB,
			idForSetting:  "CreateYoutubeVideoInstanceFailure",
			url:           "https://www.youtube.com/watch?v=VLYKx-Fwhxg",
			userAccountId: 2,
			want:          0,
		},
	}

	for _, tt := range testFailure {
		LoadAppConfig(tt.configPath)
		yvi := YoutubeVideoInstance{}
		tt.testDB.SetId(tt.idForSetting, tt.want)
		h := &Handler{tt.testDB}
		videoInstanceId, _, err := yvi.lurk_video(h, tt.url, tt.userAccountId)
		if err == nil {
			t.Errorf("lurk_video got nil error")
		}
		if videoInstanceId != tt.want {
			t.Errorf("lurk_video(%v, %s, %d)=\"%v\", wanted \"%v\"", h, tt.url, tt.userAccountId, videoInstanceId, tt.want)
		}
	}

	_ = os.Remove("../config/test_config.json")
}

type TestYoutubeConfig struct {
	DeveloperKey string `json:"developer_key"`
}

func createTestConfigFile(path string) {
	test := TestYoutubeConfig{
		DeveloperKey: "test",
	}

	file, _ := json.Marshal(test)

	_ = ioutil.WriteFile(path, file, 0644)
}

func TestParseYoutubeVideoInstance(t *testing.T) {
	// Success Case
	fmt.Println("------- TestParseYoutubeVideoInstance Success Case -------")
	// set up test cases
	tests := []struct {
		testDB        TestDB
		vSource       *models.VideoSource
		videoMetaData string
		urlEnd        string
		userAccountId uint32
		want          *models.YoutubeVideoInstance
	}{
		{
			videoMetaData: vmd,
		},
	}

	for _, tt := range tests {
		ymd := []byte(tt.videoMetaData)
		youtubeVideoInstance, err := createVideoInstance(ymd)
		if err != nil {
			t.Errorf("parseYoutubeVideoInstance got unexpected error: %s", err)
		}
		if reflect.TypeOf(youtubeVideoInstance) != reflect.TypeOf(&models.YoutubeVideoInstance{}) {
			t.Errorf("parseYoutubeVideoInstance returns wrong data: %v", youtubeVideoInstance)
		}
	}

	delete(mockDB, "YoutubeVideoInstanceTest_id")
	delete(mockDB, "UpdateYoutubeVideoInstanceFailure_id")
}
*/

func createVideoInstance(ymd []byte) (*models.YoutubeVideoInstance, error) {
	metaData := &youtube.VideoListResponse{}
	err := json.Unmarshal(ymd, metaData)
	youtubeVideoInstance, _, err := parseYoutubeVideoInstance(metaData)
	return youtubeVideoInstance, err
}
