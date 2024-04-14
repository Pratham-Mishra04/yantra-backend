package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/google/uuid"
)

func EmotionExtractionFromOnboarding(content string) ([]string, []float64) {
	// Define the request body
	requestBody, err := json.Marshal(map[string]string{
		"content": content,
	})
	if err != nil {
		LogServerError("Failed to marshal request body", err, "ml_api")
		return nil, nil
	}

	// Make a POST request to the ML_URL
	resp, err := http.Post(initializers.CONFIG.ML_URL+"/emotion_extraction", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		LogServerError("Failed to make POST request", err, "ml_api")
		return nil, nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogServerError("Failed to read response body", err, "ml_api")
		return nil, nil
	}

	// Unmarshal the response body into a map
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		LogServerError("Failed to unmarshal response body", err, "ml_api")
		return nil, nil
	}

	// Extract emotions and scores from the response
	emotionsArr, ok := response["emotions"].([]interface{})
	if !ok {
		LogServerError("Emotions not found in response", nil, "ml_api")
		return nil, nil
	}
	emotions := make([]string, len(emotionsArr))
	for i, e := range emotionsArr {
		emotions[i] = e.(string)
	}

	scoresArr, ok := response["scores"].([]interface{})
	if !ok {
		LogServerError("Scores not found in response", nil, "ml_api")
		return nil, nil
	}
	scores := make([]float64, len(scoresArr))
	for i, s := range scoresArr {
		scores[i] = s.(float64)
	}

	return emotions, scores
}

func NERExtractionFromOnboarding(content string) []string {
	// Define the request body
	requestBody, err := json.Marshal(map[string]string{
		"content": content,
	})
	if err != nil {
		LogServerError("Failed to marshal request body", err, "ml_api")
		return nil
	}

	// Make a POST request to the ML_URL
	resp, err := http.Post(initializers.CONFIG.ML_URL+"/ner_extraction", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		LogServerError("Failed to make POST request", err, "ml_api")
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogServerError("Failed to read response body", err, "ml_api")
		return nil
	}

	// Unmarshal the response body into a map
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		LogServerError("Failed to unmarshal response body", err, "ml_api")
		return nil
	}

	// Extract words from the response
	wordsArr, ok := response["words"].([]interface{})
	if !ok {
		LogServerError("Words not found in response", nil, "ml_api")
		return nil
	}
	words := make([]string, len(wordsArr))
	for i, w := range wordsArr {
		words[i] = w.(string)
	}

	return words
}

func EmotionExtractionFromPage(page *models.Page) {
	// Define the request body
	requestBody, err := json.Marshal(map[string]string{
		"content": page.Content,
	})
	if err != nil {
		LogServerError("Failed to marshal request body", err, "ml_api")
		return
	}

	// Make a POST request to the ML_URL
	resp, err := http.Post(initializers.CONFIG.ML_URL+"/emotion_extraction", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		LogServerError("Failed to make POST request", err, "ml_api")
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogServerError("Failed to read response body", err, "ml_api")
		return
	}

	// Unmarshal the response body into a map
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		LogServerError("Failed to unmarshal response body", err, "ml_api")
		return
	}

	// Extract emotions and scores from the response
	emotionsArr, ok := response["emotions"].([]interface{})
	if !ok {
		LogServerError("Emotions not found in response", nil, "ml_api")
		return
	}
	emotions := make([]string, len(emotionsArr))
	for i, e := range emotionsArr {
		emotions[i] = e.(string)
	}

	scoresArr, ok := response["scores"].([]interface{})
	if !ok {
		LogServerError("Scores not found in response", nil, "ml_api")
		return
	}
	scores := make([]float64, len(scoresArr))
	for i, s := range scoresArr {
		scores[i] = s.(float64)
	}

	page.Emotions = emotions
	page.Scores = scores

	if err := initializers.DB.Save(&page).Error; err != nil {
		LogDatabaseError(err.Error(), err, "db_error")
	}
}

func NERExtractionFromPage(page *models.Page) {
	// Define the request body
	requestBody, err := json.Marshal(map[string]string{
		"content": page.Content,
	})
	if err != nil {
		LogServerError("Failed to marshal request body", err, "ml_api")
		return
	}

	// Make a POST request to the ML_URL
	resp, err := http.Post(initializers.CONFIG.ML_URL+"/ner_extraction", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		LogServerError("Failed to make POST request", err, "ml_api")
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogServerError("Failed to read response body", err, "ml_api")
		return
	}

	// Unmarshal the response body into a map
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		LogServerError("Failed to unmarshal response body", err, "ml_api")
		return
	}

	// Extract words from the response
	wordsArr, ok := response["words"].([]interface{})
	if !ok {
		LogServerError("Words not found in response", nil, "ml_api")
		return
	}
	words := make([]string, len(wordsArr))
	for i, w := range wordsArr {
		words[i] = w.(string)
	}

	page.NER = words

	if err := initializers.DB.Save(&page).Error; err != nil {
		LogDatabaseError(err.Error(), err, "db_error")
	}
}

type UserDominatingEmotions struct {
	UserID                 uuid.UUID `json:"userID"`
	UserDominatingEmotions []string  `json:"emotions"`
}

func GroupDominatingEmotion(group *models.Group) {
	// Get the list of users in the group
	users := getUsersInGroup(group)

	// Create a slice to store the dominating emotion for each user
	var userDominatingEmotionsList []UserDominatingEmotions

	// Iterate over each user
	for _, user := range users {
		// Get the pages of the last 7 days for the user
		pages := getPagesForLast7Days(user)

		// Initialize a map to store the highest scored emotion for each day
		var userDominatingEmotions []string

		for _, page := range pages {
			// Initialize variables to store the highest score and corresponding emotion for the page
			var highestScore float64
			var highestEmotion string

			// Iterate over the emotions and scores for the page
			for i, score := range page.Scores {
				if score > highestScore {
					highestScore = score
					highestEmotion = page.Emotions[i]
				}
			}

			// Store the highest scored emotion for the page
			userDominatingEmotions = append(userDominatingEmotions, highestEmotion)
		}

		// Store the result for the user
		userDominatingEmotionsList = append(userDominatingEmotionsList, UserDominatingEmotions{UserID: user.ID, UserDominatingEmotions: userDominatingEmotions})
	}

	type UserEmotionsRequest struct {
		Data []UserDominatingEmotions `json:"data"`
	}

	requestBody, err := json.Marshal(UserEmotionsRequest{Data: userDominatingEmotionsList})
	if err != nil {
		LogServerError("Failed to marshal request body", err, "ml_api")
		return
	}

	// Make a POST request to the ML_URL
	resp, err := http.Post(initializers.CONFIG.ML_URL+"/ner_extraction", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		LogServerError("Failed to make POST request", err, "ml_api")
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogServerError("Failed to read response body", err, "ml_api")
		return
	}

	// Unmarshal the response body into a map
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		LogServerError("Failed to unmarshal response body", err, "ml_api")
		return
	}

	// Extract words from the response
	emotion, ok := response["emotion"].(string)
	if !ok {
		LogServerError("emotion not found in response", nil, "ml_api")
		return
	}

	group.Emotion = emotion

	if err := initializers.DB.Save(&group).Error; err != nil {
		LogDatabaseError(err.Error(), err, "db_error")
	}
}

func getUsersInGroup(group *models.Group) []models.User {
	var users []models.User
	for _, membership := range group.Memberships {
		users = append(users, membership.User)
	}

	return users
}

func getPagesForLast7Days(user models.User) []models.Page {
	var pages []models.Page
	if err := initializers.DB.Where("journal_id = ? AND created_at >= ?", user.Journal.ID, time.Now().AddDate(0, 0, -7)).Find(&pages).Error; err != nil {
		LogDatabaseError(err.Error(), err, "db_error")
	}
	return pages
}

type PersonBody struct {
	ID        string      `json:"id"`
	Emotions  [][]string  `json:"emotions"`
	Scores    [][]float64 `json:"scores"`
	Locations []string    `json:"locations"`
}

type GroupBody struct {
	ID       string `json:"id"`
	Emotion  string `json:"emotion"`
	Location string `json:"location"`
}

type RecommendationReqBody struct {
	Person PersonBody  `json:"person"`
	Groups []GroupBody `json:"groups"`
}

func GetGroupRecommendations(user *models.User) []models.Group {
	pages := getPagesForLast7Days(*user)

	var emotions [][]string
	var scores [][]float64
	var locations []string

	for _, page := range pages {
		emotions = append(emotions, page.Emotions)
		scores = append(scores, page.Scores)
		locations = append(locations, page.NER...)
	}

	locations = append(locations, user.Location)

	person := PersonBody{
		ID:        user.ID.String(),
		Emotions:  emotions,
		Scores:    scores,
		Locations: locations,
	}

	var groups []models.Group
	initializers.DB.Find(&groups)

	var groupBodies []GroupBody

	for _, group := range groups {
		groupBodies = append(groupBodies, GroupBody{
			ID:       group.ID.String(),
			Emotion:  group.Emotion,
			Location: group.Location,
		})
	}

	requestBody, err := json.Marshal(RecommendationReqBody{Person: person, Groups: groupBodies})
	if err != nil {
		LogServerError("Failed to marshal request body", err, "ml_api")
		return groups
	}

	// Make a POST request to the ML_URL
	resp, err := http.Post(initializers.CONFIG.ML_URL+"/ner_extraction", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		LogServerError("Failed to make POST request", err, "ml_api")
		return groups
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogServerError("Failed to read response body", err, "ml_api")
		return groups
	}

	// Unmarshal the response body into a map
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		LogServerError("Failed to unmarshal response body", err, "ml_api")
		return groups
	}

	// Extract words from the response
	groupIDs, ok := response["groupIDs"].([]string)
	if !ok {
		LogServerError("emotion not found in response", nil, "ml_api")
		return groups
	}

	groupIDMap := make(map[string]bool)
	for _, id := range groupIDs {
		groupIDMap[id] = true
	}

	// Create a new slice to store the filtered groups
	var filteredGroups []models.Group

	// Iterate over the groups and add only those that are not in groupIDs to the filtered slice
	for _, group := range groups {
		if !groupIDMap[group.ID.String()] {
			filteredGroups = append(filteredGroups, group)
		}
	}

	return filteredGroups
}

func GetGroupRecommendationsFromOnboarding(emotions []string, scores []float64, NERs []string, user *models.User) []models.Group {
	var emotions2d [][]string
	var scores2d [][]float64

	emotions2d = append(emotions2d, emotions)
	scores2d = append(scores2d, scores)

	person := PersonBody{
		ID:        user.ID.String(),
		Emotions:  emotions2d,
		Scores:    scores2d,
		Locations: NERs,
	}

	var groups []models.Group
	initializers.DB.Find(&groups)

	var groupBodies []GroupBody

	for _, group := range groups {
		groupBodies = append(groupBodies, GroupBody{
			ID:       group.ID.String(),
			Emotion:  group.Emotion,
			Location: group.Location,
		})
	}

	requestBody, err := json.Marshal(RecommendationReqBody{Person: person, Groups: groupBodies})
	if err != nil {
		LogServerError("Failed to marshal request body", err, "ml_api")
		return groups
	}

	// Make a POST request to the ML_URL
	resp, err := http.Post(initializers.CONFIG.ML_URL+"/ner_extraction", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		LogServerError("Failed to make POST request", err, "ml_api")
		return groups
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogServerError("Failed to read response body", err, "ml_api")
		return groups
	}

	// Unmarshal the response body into a map
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		LogServerError("Failed to unmarshal response body", err, "ml_api")
		return groups
	}

	// Extract words from the response
	groupIDs, ok := response["groupIDs"].([]string)
	if !ok {
		LogServerError("emotion not found in response", nil, "ml_api")
		return groups
	}

	groupIDMap := make(map[string]bool)
	for _, id := range groupIDs {
		groupIDMap[id] = true
	}

	// Create a new slice to store the filtered groups
	var filteredGroups []models.Group

	// Iterate over the groups and add only those that are not in groupIDs to the filtered slice
	for _, group := range groups {
		if !groupIDMap[group.ID.String()] {
			filteredGroups = append(filteredGroups, group)
		}
	}

	return filteredGroups
}
