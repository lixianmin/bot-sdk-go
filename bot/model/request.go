package model

import (
	"github.com/lixianmin/got/convert"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/lixianmin/bot-sdk-go/bot/data"
)

const (
	INTENT_REQUEST        = "IntentRequest"
	LAUNCH_REQUEST        = "LaunchRequest"
	SESSION_ENDED_REQUEST = "SessionEndedRequest"

	AUDIO_PLAYER_PLAYBACK_STARTED                 = "AudioPlayer.PlaybackStarted"
	AUDIO_PLAYER_PLAYBACK_STOPPED                 = "AudioPlayer.PlaybackStopped"
	AUDIO_PLAYER_PLAYBACK_FINISHED                = "AudioPlayer.PlaybackFinished"
	AUDIO_PLAYER_PLAYBACK_NEARLY_FINISHED         = "AudioPlayer.PlaybackNearlyFinished"
	AUDIO_PLAYER_PROGRESS_REPORT_INTERVAL_ELAPSED = "AudioPlayer.ProgressReportIntervalElapsed"

	VIDEO_PLAYER_PLAYBACK_STARTED                 = "VideoPlayer.PlaybackStarted"
	VIDEO_PLAYER_PLAYBACK_STOPPED                 = "VideoPlayer.PlaybackStopped"
	VIDEO_PLAYER_PLAYBACK_FINISHED                = "VideoPlayer.PlaybackFinished"
	VIDEO_PLAYER_PLAYBACK_NEARLY_FINISHED         = "VideoPlayer.PlaybackNearlyFinished"
	VIDEO_PLAYER_PLAYBACK_SCHEDULED_STOP_REACHED  = "VideoPlayer.PlaybackScheduledStopReached"
	VIDEO_PLAYER_PROGRESS_REPORT_INTERVAL_ELAPSED = "VideoPlayer.ProgressReportIntervalElapsed"
)

type Request struct {
	Type   string
	Common data.RequestPart
}

type IntentRequest struct {
	Data   data.IntentRequest
	Dialog *Dialog
	Request
}

type LaunchRequest struct {
	Data data.LaunchRequest
	Request
}

type SessionEndedRequest struct {
	Data data.SessionEndedRequest
	Request
}

type EventRequest struct {
	Data data.EventRequest
	Request
}

type AudioPlayerEventRequest struct {
	Data data.AudioPlayerEventRequest
	EventRequest
}

type VideoPlayerEventRequest struct {
	Data data.VideoPlayerEventRequest
	EventRequest
}

func (this *EventRequest) GetUrl() string {
	return this.Data.Request.Url
}

func (this *EventRequest) GetName() string {
	return this.Data.Request.Name
}

func (this *AudioPlayerEventRequest) GetOffsetInMilliseconds() int32 {
	return this.Data.Request.OffsetInMilliseconds
}

func (this *VideoPlayerEventRequest) GetOffsetInMilliseconds() int32 {
	return this.Data.Request.OffsetInMilliseconds
}

// 获取意图名
func (this *IntentRequest) GetIntentName() (string, bool) {
	return this.Dialog.GetIntentName()
}

// 槽位填充是否完成
func (this *IntentRequest) IsDialogStateCompleted() bool {
	return this.Dialog.DialogState == "COMPLETED"
}

// 获取用户请求query
func (this *IntentRequest) GetQuery() string {
	query, _ := this.Dialog.GetQuery()
	return query
}

// 获取用户id
func (my *Request) GetUserId() string {
	return my.Common.Context.System.User.UserId
}

// 获取设备id
func (my *Request) GetDeviceId() string {
	return my.Common.Context.System.Device.DeviceId
}

// 这个就是cuid/sn
func (my *Request) GetOriginalDeviceId() string {
	return my.Common.Context.System.Device.OriginalDeviceId
}

func (my *Request) GetUserDeviceId() string {
	return my.Common.Context.System.Device.UserDeviceId
}

func (my *Request) GetDeviceIPAddress() string {
	return my.Common.Context.System.Device.DeviceIPAddress
}

func (my *Request) GetBaiduUid() string {
	return my.Common.Context.System.User.UserInfo.Account.Baidu.BaiduUid
}

func (my *Request) GetCity() string {
	return my.Common.Context.System.User.UserInfo.Location.City
}

// 获取音频播放上下文
func (my *Request) GetAudioPlayerContext() data.AudioPlayerContext {
	return my.Common.Context.AudioPlayer
}

// 获取视频播放上下文
func (my *Request) GetVideoPlayerContext() data.VideoPlayerContext {
	return my.Common.Context.VideoPlayer
}

// 获取access token
func (my *Request) GetAccessToken() string {
	return my.Common.Context.System.User.AccessToken
}

func (my *Request) GetApiAccessToken() string {
	return my.Common.Context.System.ApiAccessToken
}

// 获取请求的时间戳
func (my *Request) GetTimestamp() int {
	i, err := strconv.Atoi(my.Common.Request.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// 获取请求id
func (my *Request) GetRequestId() string {
	return my.Common.Request.RequestId
}

// 获取技能id
func (my *Request) GetBotId() string {
	return my.Common.Context.System.Application.ApplicationId
}

// 验证请求时间戳合法性
func (my *Request) VerifyTimestamp() bool {

	if my.GetTimestamp()+180 > int(time.Now().Unix()) {
		return true
	}

	return false
}

// 获取设备支持的接口类型
func (my *Request) GetSupportedInterfaces() map[string]interface{} {
	return my.Common.Context.System.Device.SupportedInterfaces
}

func (my *Request) isSupportInterface(support string) bool {
	supportedInterfaces := my.GetSupportedInterfaces()
	_, ok := supportedInterfaces[support]

	if ok {
		return true
	}
	return false
}

// 检查是否支持展现
func (my *Request) IsSupportDisplay() bool {
	return my.isSupportInterface("Display")
}

// 检查是否支持音频播放
func (my *Request) IsSupportAudio() bool {
	return my.isSupportInterface("AudioPlayer")
}

// 检查是否支持视频播放
func (my *Request) IsSupportVideo() bool {
	return my.isSupportInterface("VideoPlayer")
}

// 验证技能id合法性
func (my *Request) VerifyBotID(myBotID string) bool {
	if my.GetBotId() == myBotID {
		return true
	}
	return false
}

func getType(rawData []byte) string {
	d := data.LaunchRequest{}

	if err := convert.FromJsonE(rawData, &d); err != nil {
		log.Println(err)
	}

	return d.Request.Type
}

func GetSessionData(rawData []byte) data.Session {
	common := data.RequestPart{}
	if err := convert.FromJsonE(rawData, &common); err != nil {
		log.Println(err)
	}

	return common.Session
}

func NewRequest(rawData []byte) interface{} {
	common := data.RequestPart{}
	if err := convert.FromJsonE(rawData, &common); err != nil {
		log.Println(err)
		return false
	}

	requestType := common.Request.Type
	if requestType == INTENT_REQUEST {
		request := IntentRequest{}
		request.Type = requestType
		request.Common = common
		if err := convert.FromJsonE(rawData, &request.Data); err != nil {
			log.Println(err)
			return false
		}
		request.Dialog = NewDialog(request.Data.Request)

		return request
	} else if requestType == LAUNCH_REQUEST {
		request := LaunchRequest{}
		request.Type = requestType
		request.Common = common
		if err := convert.FromJsonE(rawData, &request.Data); err != nil {
			log.Println(err)
			return false
		}
		return request
	} else if requestType == SESSION_ENDED_REQUEST {
		request := SessionEndedRequest{}
		request.Type = requestType
		request.Common = common
		if err := convert.FromJsonE(rawData, &request.Data); err != nil {
			log.Println(err)
			return false
		}
		return request
	} else {
		if match, _ := regexp.MatchString("^AudioPlayer", requestType); match {
			request := AudioPlayerEventRequest{}
			request.Type = requestType
			request.Common = common
			if err := convert.FromJsonE(rawData, &request.Data); err != nil {
				log.Println(err)
				return false
			}
			return request
		} else if match, _ := regexp.MatchString("^VideoPlayer", requestType); match {
			request := VideoPlayerEventRequest{}
			request.Type = requestType
			request.Common = common
			if err := convert.FromJsonE(rawData, &request.Data); err != nil {
				log.Println(err)
				return false
			}
			return request
		}

		request := EventRequest{}
		request.Type = requestType
		request.Common = common
		if err := convert.FromJsonE(rawData, &request.Data); err != nil {
			log.Println(err)
			return false
		}
		return request
	}
}
