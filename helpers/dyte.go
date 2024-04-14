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

type CreateMeetingResponse struct {
	Success bool        `json:"success"`
	Data    MeetingData `json:"data"`
	Message string      `json:"message"`
}

type MeetingData struct {
	ID                string `json:"id"`
	Title             string `json:"title"`
	PreferredRegion   string `json:"preferred_region"`
	CreatedAt         string `json:"created_at"`
	RecordOnStart     bool   `json:"record_on_start"`
	UpdatedAt         string `json:"updated_at"`
	LiveStreamOnStart bool   `json:"live_stream_on_start"`
	PersistChat       bool   `json:"persist_chat"`
	Status            string `json:"status"`
	// RecordingConfig   Recording `json:"recording_config"`
}

func CreateDyteMeeting(event *models.Event) error {
	URL := "https://api.dyte.io/v2/meetings"

	CONFIG := map[string]interface{}{
		"title":                event.ID.String(),
		"preferred_region":     "ap-south-1",
		"record_on_start":      false,
		"live_stream_on_start": false,
		"recording_config": map[string]interface{}{
			"max_seconds":      60,
			"file_name_prefix": "string",
			"video_config": map[string]interface{}{
				"codec":  "H264",
				"width":  1280,
				"height": 720,
				"watermark": map[string]interface{}{
					"url":      "http://example.com",
					"size":     map[string]interface{}{"width": 1, "height": 1},
					"position": "left top",
				},
				"export_file": true,
			},
			"audio_config": map[string]interface{}{
				"codec":       "AAC",
				"channel":     "stereo",
				"export_file": true,
			},
			"storage_config": map[string]interface{}{
				"type":        "aws",
				"access_key":  "string",
				"secret":      "string",
				"bucket":      "string",
				"region":      "us-east-1",
				"path":        "string",
				"auth_method": "KEY",
				"username":    "string",
				"password":    "string",
				"host":        "string",
				"port":        0,
				"private_key": "string",
			},
			"dyte_bucket_config": map[string]interface{}{"enabled": true},
			"live_streaming_config": map[string]interface{}{
				"rtmp_url": "rtmp://a.rtmp.youtube.com/live2",
			},
		},
		"persist_chat":     false,
		"summarize_on_end": false,
	}

	jsonValue, _ := json.Marshal(CONFIG)

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return AppError{Code: 500, Message: "Error creating request for Dyte SDK", LogMessage: err.Error(), Err: err}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+initializers.CONFIG.DYTE_TOKEN)
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
		return err
	}
	return nil
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

	URL := "https://api.dyte.io/v2/meetings/" + event.DyteID + "/participants"

	CONFIG := map[string]interface{}{
		"name":                  user.Name,
		"picture":               "https://i.imgur.com/test.jpg",
		"preset_name":           preset,
		"custom_participant_id": user.ID.String(),
	}

	jsonValue, _ := json.Marshal(CONFIG)

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", AppError{Code: 500, Message: "Error creating request for Dyte SDK", LogMessage: err.Error(), Err: err}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+initializers.CONFIG.DYTE_TOKEN)
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
		return response.Data.Token, nil
	} else {
		return "", err
	}
}

// TODO test this
func GetDyteMeetingParticipants(meetingID string) ([]models.User, error) {
	URL := fmt.Sprintf("https://api.dyte.io/v2/meetings/%s/participants", meetingID)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, AppError{Code: fiber.StatusInternalServerError, Message: "Error creating request for Dyte SDK", LogMessage: err.Error(), Err: err}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+initializers.CONFIG.DYTE_TOKEN)

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
