package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Title             string    `json:"title"`
	PreferredRegion   string    `json:"preferred_region"`
	RecordOnStart     bool      `json:"record_on_start"`
	LiveStreamOnStart bool      `json:"live_stream_on_start"`
	RecordingConfig   Recording `json:"recording_config"`
	PersistChat       bool      `json:"persist_chat"`
}

type Recording struct {
	MaxSeconds          int           `json:"max_seconds"`
	FileNamePrefix      string        `json:"file_name_prefix"`
	DyteBucketConfig    DyteBucket    `json:"dyte_bucket_config"`
	LiveStreamingConfig LiveStreaming `json:"live_streaming_config"`
}

type DyteBucket struct {
	Enabled bool `json:"enabled"`
}

type LiveStreaming struct {
	RTMPURL string `json:"rtmp_url"`
}

type CreateMeetingResponse struct {
	Success bool        `json:"success"`
	Data    MeetingData `json:"data"`
}

type MeetingData struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	PreferredRegion   string    `json:"preferred_region"`
	CreatedAt         string    `json:"created_at"`
	RecordOnStart     bool      `json:"record_on_start"`
	UpdatedAt         string    `json:"updated_at"`
	LiveStreamOnStart bool      `json:"live_stream_on_start"`
	PersistChat       bool      `json:"persist_chat"`
	Status            string    `json:"status"`
	RecordingConfig   Recording `json:"recording_config"`
}

func CreateDyteMeeting(event *models.Event) error {
	URL := "https://api.dyte.io/v2/meetings"
	CONFIG := Config{
		Title:             event.ID.String(),
		PreferredRegion:   "ap-south-1",
		RecordOnStart:     false,
		LiveStreamOnStart: false,
		RecordingConfig: Recording{
			MaxSeconds:     60,
			FileNamePrefix: "string",
			DyteBucketConfig: DyteBucket{
				Enabled: true,
			},
			LiveStreamingConfig: LiveStreaming{
				RTMPURL: "rtmp://a.rtmp.youtube.com/live2",
			},
		},
		PersistChat: false,
	}

	jsonValue, _ := json.Marshal(CONFIG)

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return AppError{Code: 500, Message: "Error creating request for Dyte SDK", LogMessage: err.Error(), Err: err}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+initializers.CONFIG.DYTE_API_KEY)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AppError{Code: 500, Message: "Error making request to Dyte SDK", LogMessage: err.Error(), Err: err}
	}
	defer resp.Body.Close()

	var response CreateMeetingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if response.Success {
		event.DyteID = response.Data.ID
		if err := initializers.DB.Save(&event).Error; err != nil {
			return AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	} else {
		//TODO
		return err
	}
	return nil
}

type ParticipantRequest struct {
	Name                string `json:"name"`
	Picture             string `json:"picture"`
	PresetName          string `json:"preset_name"`
	CustomParticipantID string `json:"custom_participant_id"`
}

type ParticipantResponse struct {
	Success bool            `json:"success"`
	Data    ParticipantData `json:"data"`
}

type ParticipantsResponse struct {
	Success bool              `json:"success"`
	Data    []ParticipantData `json:"data"`
}

type ParticipantData struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Picture             string `json:"picture"`
	CustomParticipantID string `json:"custom_participant_id"`
	PresetName          string `json:"preset_name"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	Token               string `json:"token"`
}

func GetDyteMeetingAuthToken(event *models.Event, user *models.User) (string, error) {
	preset := "group_call_participant"
	if event.Group.Moderator.UserID == user.ID {
		preset = "group_call_host"
	}

	URL := "https://api.dyte.io/v2/meetings/meeting_id/participants"
	participant := ParticipantRequest{
		Name:                user.Name,
		Picture:             "https://i.imgur.com/test.jpg",
		PresetName:          preset,
		CustomParticipantID: user.ID.String(),
	}

	jsonValue, _ := json.Marshal(participant)

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", AppError{Code: 500, Message: "Error creating request for Dyte SDK", LogMessage: err.Error(), Err: err}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic 123")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", AppError{Code: 500, Message: "Error making request to Dyte SDK", LogMessage: err.Error(), Err: err}
	}
	defer resp.Body.Close()

	var response ParticipantResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if response.Success {
		event.DyteID = response.Data.ID
		return response.Data.Token, nil
	} else {
		//TODO
		return "", err
	}
}

func GetDyteMeetingParticipants(meetingID string) ([]models.User, error) {
	URL := fmt.Sprintf("https://api.dyte.io/v2/meetings/%s/participants", meetingID)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, AppError{Code: fiber.StatusInternalServerError, Message: "Error creating request for Dyte SDK", LogMessage: err.Error(), Err: err}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic 123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, AppError{Code: fiber.StatusInternalServerError, Message: "Error making request to Dyte SDK", LogMessage: err.Error(), Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, AppError{Code: resp.StatusCode, Message: "Failed to get Dyte meeting participants"}
	}

	var participantResponse ParticipantsResponse
	if err := json.NewDecoder(resp.Body).Decode(&participantResponse); err != nil {
		return nil, AppError{Code: fiber.StatusInternalServerError, Message: "Error decoding Dyte SDK response", LogMessage: err.Error(), Err: err}
	}

	var users []models.User

	for _, participant := range participantResponse.Data {
		userID := participant.CustomParticipantID

		var user models.User
		if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
			LogDatabaseError("Invalid User In Callback", err, "meeting-"+meetingID)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}
