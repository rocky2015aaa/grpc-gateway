// Package controller provides restful api handler functions
package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"

	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"gitlab.com/Fratbe/addglee/src/server/i18n"
	"gitlab.com/Fratbe/addglee/src/server/logger"
	"gitlab.com/Fratbe/addglee/src/server/models"
)

const (
	UNIT_TEST_DOMAIN  = ""
	YOUTUBE_REGEX     = ""
	YOUTUBE_DOMAIN    = ""
	YOUTUBE_ID_QUERY  = ""
	YOUTUBE_ID_LENGTH = 1

	MN_CREATE_VIDEO_INSTANCE        = "controller.CreateVideoInstance"
	MN_HANDLE_REQUEST               = "controller.YoutubeVideoInstance.createYoutubeVideoInstance"
	MN_YOUTUBE_LURK_VIDEO           = "controller.YoutubeVideoInstance.lurk_video"
	FN_GET_YOUTUBE_METADATA         = "controller.getYoutubeMetaData"
	FN_PARSE_YOUTUBE_VIDEO_INSTANCE = "controller.parseYoutubeVideoInstance"
)

var (
	videoDomains = map[string]VideoInstance{
		UNIT_TEST_DOMAIN: VideoInstanceTest{},
		YOUTUBE_DOMAIN:   YoutubeVideoInstance{},
	}
)

type YoutubeVideoInstance struct {
	YoutubeVideoSourceId uint32
}

type VideoInstance interface {
	lurk_video(h *Handler, url string, userAccountId uint32) (uint32, string, error)
}

// VideoInstanceRequest has the field ([]string)
func (h *Handler) CreateVideoInstance(ctx context.Context, videoInstance *VideoInstanceRequest) (*Response, error) {
	// ValidateStruct checks if essential JSON field is empty like "" or invalid value
	valid, err := govalidator.ValidateStruct(videoInstance)
	if InputValidation(valid, err) != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_CREATE_VIDEO_INSTANCE, logger.SEVERITY: logger.WARN_LEVEL}).Warnln(err)
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "400"))
		return &Response{Message: i18n.T("invalid_request")}, err
	}

	resp, err := videoInstance.HandleRequest(h)
	if err != nil {
		return nil, err
	}

	if resp.Message == i18n.T("invalid_request") {
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "400"))
	}
	if resp.Message == i18n.T("no_video_instance_creation") {
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
	}

	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_CREATE_VIDEO_INSTANCE, logger.SEVERITY: logger.INFO_LEVEL}).Infof("Video instance created\n")
	return resp, err
}

func (vir *VideoInstanceRequest) HandleRequest(h *Handler) (*Response, error) {
	respMsg := ""
	urls := strings.Fields(vir.Urls)
	if len(urls) == 0 {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_HANDLE_REQUEST, logger.SEVERITY: logger.INFO_LEVEL}).Infof("There is no video URL.")
		return &Response{Message: i18n.T("no_video_instance_creation")}, nil
	}
	for domain, instance := range videoDomains {
		for _, url := range urls {
			if strings.Contains(url, domain) {
				logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_HANDLE_REQUEST, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("url: %s, user account id: %d\n", url, vir.UserAccountId)
				videoId, msg, err := instance.lurk_video(h, url, vir.UserAccountId)
				if err != nil {
					if err.Error() == i18n.T("invalid_url_domain") || err.Error() == i18n.T("invalid_request") {
						return &Response{Message: i18n.T("invalid_request")}, nil
					}

					return nil, errors.New(i18n.T("video_instance_creation_failure"))
				}
				if respMsg == "" {
					respMsg += fmt.Sprintf(i18n.T(msg), videoId)
				} else {
					respMsg += " / " + fmt.Sprintf(i18n.T(msg), videoId)
				}
			}
		}
	}
	return &Response{Message: respMsg}, nil
}

func (yvi YoutubeVideoInstance) lurk_video(h *Handler, url string, userAccountId uint32) (uint32, string, error) {
	match, err := regexp.MatchString(YOUTUBE_REGEX, url)
	if !match {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_YOUTUBE_LURK_VIDEO, logger.SEVERITY: logger.DEBUG_LEVEL}).Debugf("Youtube regex: %s, url: %s\n", YOUTUBE_REGEX, url)
		return 0, "", errors.New(i18n.T("invalid_url_domain"))
	}
	urlStart := url[:strings.Index(url, YOUTUBE_DOMAIN)+len(YOUTUBE_DOMAIN)]
	urlEndStart := strings.Index(url, YOUTUBE_ID_QUERY) + 2
	if len(url[urlEndStart:]) < YOUTUBE_ID_LENGTH {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_YOUTUBE_LURK_VIDEO, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln("Youtube video id length is too short")
		return 0, "", errors.New(i18n.T("invalid_request"))
	}
	urlEnd := url[urlEndStart : urlEndStart+YOUTUBE_ID_LENGTH]
	logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_YOUTUBE_LURK_VIDEO, logger.SEVERITY: logger.INFO_LEVEL}).Infof("url_start: %s, url_end: %s", urlStart, urlEnd)

	val := map[string]interface{}{
		"": urlEnd,
		"": userAccountId,
	}
	client := &http.Client{
		Transport: &transport.APIKey{Key: AppConf.DeveloperKey},
	}
	service, err := youtube.New(client)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_YOUTUBE_LURK_VIDEO, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return 0, "", errors.New(i18n.T("youtube_api_client_creation_failure"))
	}
	// Make the API call to YouTube.
	call := service.Videos.List("").
		Id(urlEnd)

	previousVideoInstance := h.GetYoutubeVideoInstanceByValue(val)
	if previousVideoInstance != nil {
		_, err := call.IfNoneMatch(previousVideoInstance.Etag).Do()
		if err != nil && err.Error() == "googleapi: got HTTP response code 304 with body: " {
			return previousVideoInstance.Id, "no_video_instance_change", nil
		}
		videoInstance, err := getYoutubeMetaData(h, call)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_YOUTUBE_LURK_VIDEO, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
			return 0, "", err
		}
		videoInstance.UrlEnd = urlEnd
		videoInstance.UserAccountId = userAccountId
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_YOUTUBE_LURK_VIDEO, logger.SEVERITY: logger.DEBUG_LEVEL}).Infof("video for update: %+v\n", videoInstance)
		if err = h.UpdateYoutubeVideoInstance(videoInstance, val); err != nil {
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_YOUTUBE_LURK_VIDEO, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
			return 0, "", err
		}
		return previousVideoInstance.Id, "video_instance_update", nil
	}

	videoInstance, err := getYoutubeMetaData(h, call)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_YOUTUBE_LURK_VIDEO, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return 0, "", err
	}

	if yvi.YoutubeVideoSourceId == 0 {
		videoSource := h.GetVideoSource(YOUTUBE_DOMAIN)
		if videoSource == nil {
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_YOUTUBE_LURK_VIDEO, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
			return 0, "", err
		}
		yvi.YoutubeVideoSourceId = videoSource.Id
		videoInstance.VideoSourceId = videoSource.Id
	} else {
		videoInstance.VideoSourceId = yvi.YoutubeVideoSourceId
	}

	videoInstance.UrlEnd = urlEnd
	videoInstance.UserAccountId = userAccountId
	if err = h.CreateYoutubeVideoInstance(videoInstance); err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: MN_YOUTUBE_LURK_VIDEO, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return 0, "", err
	}

	return videoInstance.Id, "video_instance_creation", nil
}

func getYoutubeMetaData(h *Handler, call *youtube.VideosListCall) (*models.YoutubeVideoInstance, error) {
	response, err := call.Do()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_GET_YOUTUBE_METADATA, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return nil, err
	}
	videoInstance, defaultLanguage, err := parseYoutubeVideoInstance(response)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_GET_YOUTUBE_METADATA, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return nil, err
	}

	if language := h.GetLanguageByLetterCode(defaultLanguage); language != nil {
		videoInstance.Defaultlanguage = language.Id
	} else {
		videoInstance.Defaultlanguage = uint32(1)
	}

	return videoInstance, nil
}

func parseYoutubeVideoInstance(metaData *youtube.VideoListResponse) (*models.YoutubeVideoInstance, string, error) {
	item := metaData.Items[0]
	publishedat, err := time.Parse("", item.Snippet.PublishedAt)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_PARSE_YOUTUBE_VIDEO_INSTANCE, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return nil, "", err
	}

	categoryId, err := strconv.ParseUint(item.Snippet.CategoryId, 10, 64)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_PARSE_YOUTUBE_VIDEO_INSTANCE, logger.SEVERITY: logger.ERROR_LEVEL}).Errorln(err)
		return nil, "", err
	}

	return &models.YoutubeVideoInstance{
		Etag:                     metaData.Etag,
		Title:                    item.Snippet.Title,
		Description:              item.Snippet.Description,
		Category:                 item.Kind,
		Publishedat:              &publishedat,
		SnippetThumbnailUrl:      item.Snippet.Thumbnails.Default.Url,
		Categoryid:               uint32(categoryId),
		ContentdetailsDuration:   item.ContentDetails.Duration,
		ContentdetailsDimension:  []byte(item.ContentDetails.Dimension),
		ContentdetailsDefinition: []byte(item.ContentDetails.Definition),
		ContentdetailsProjection: []byte(item.ContentDetails.Projection),
		StatusEmbeddable:         item.Status.Embeddable,
		StatisticsViewcount:      uint32(item.Statistics.ViewCount),
		StatisticsLikecount:      uint32(item.Statistics.LikeCount),
		StatisticsDislikecount:   uint32(item.Statistics.DislikeCount),
		StatisticsCommentcount:   uint32(item.Statistics.CommentCount),
		PlayerEmbedhtml:          item.Player.EmbedHtml,
		PlayerEmbedheight:        uint32(item.Player.EmbedHeight),
		PlayerEmbedwidth:         uint32(item.Player.EmbedWidth),
		SnippetChannelid:         item.Snippet.ChannelId,
	}, item.Snippet.DefaultLanguage, nil
}
