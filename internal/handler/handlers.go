package handlers

import (
	"backend/internal/model"
	"backend/pkg/repository"
	services "backend/pkg/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// var producer sarama.SyncProducer

// updateUserHandler updates additional user information
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from request URL
	userID, err := strconv.Atoi(r.URL.Path[len("/users/"):])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse request data from request body
	var requestData model.RequestData
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the request data
	if err := services.ValidateRequestData(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve user from database
	user, err := repository.GetUserByID(userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update user information with request data
	if requestData.Data.Email != "" {
		user.Email = requestData.Data.Email
	}
	if requestData.Data.FirstName != "" {
		user.FirstName = requestData.Data.FirstName
	}
	if requestData.Data.LastName != "" {
		user.LastName = requestData.Data.LastName
	}
	if requestData.Data.TimeZoneID != "" {
		user.TimeZoneID = requestData.Data.TimeZoneID
	}

	// Update user in database
	err = repository.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// produceUserEventHandler produces a user event to 'users' topic based on updated user information
func ProduceUserEventHandler(w http.ResponseWriter, r *http.Request) {

	// Parse user ID from request URL
	userID, err := strconv.Atoi(r.URL.Path[len("/produce/"):])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve user from database
	user, err := repository.GetUserByID(userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the user information from Github.
	userURL := fmt.Sprintf("https://api.github.com/users/%s", user.Username)
	resp, err := http.Get(userURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to retrieve user information from Github: %s", resp.Status), resp.StatusCode)
		return
	}

	var githubUser model.GitHubUser

	var empty []string

	// Decode user data response
	err = json.NewDecoder(resp.Body).Decode(&githubUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update user information with Github data
	if githubUser.Email != nil {
		user.Email = *githubUser.Email
	}
	if githubUser.UserName != "" {
		user.Username = githubUser.UserName
	}

	res_fol, err_fol := handleUserFollowersFromGit(user.Username)

	if err_fol != nil {
		http.Error(w, err_fol.Error(), http.StatusBadRequest)
		return
	}

	if res_fol != nil {
		user.Followers = empty
	} else {
		user.Followers = res_fol
	}

	res_rep, err_rep := handleUserReposFromGit(user.Username)

	if err_rep != nil {
		http.Error(w, err_rep.Error(), http.StatusBadRequest)
		return
	}

	if res_rep != nil {
		user.Repos = empty
	} else {
		user.Repos = res_rep
	}

	// Update user in database
	err = repository.UpdateUserFromGit(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfoChanged := model.UserInfoChanged{
		Meta: struct {
			Type      string `json:"type"`
			EventID   string `json:"event_id"`
			CreatedAt int64  `json:"created_at"`
			TraceID   string `json:"trace_id"`
			ServiceID string `json:"service_id"`
		}{
			Type:      "UserInfoChanged",
			EventID:   uuid.New().String(),
			CreatedAt: time.Now().UnixNano(),
			TraceID:   uuid.New().String(),
			ServiceID: "user-service",
		},
		Payload: struct {
			ID         int      `json:"id"`
			Username   string   `json:"userName"`
			Followers  []string `json:"followers,omitempty"`
			Repos      []string `json:"repos,omitempty"`
			Email      string   `json:"userEmail"`
			FirstName  string   `json:"userFirstName"`
			LastName   string   `json:"userLastName"`
			TimeZoneID string   `json:"userTimeZoneId"`
		}{
			ID:         user.ID,
			Username:   user.Username,
			Followers:  user.Followers,
			Repos:      user.Repos,
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			TimeZoneID: user.TimeZoneID,
		},
	}

	// Marshal event to JSON
	eventJSON, err := json.Marshal(userInfoChanged)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the success response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(eventJSON)
}

// Get user followers from GitHub
func handleUserFollowersFromGit(userName string) ([]string, error) {
	userFollowersURL := fmt.Sprintf("https://api.github.com/users/%s/followers", userName)
	resp_fol, err_fol := http.Get(userFollowersURL)
	if err_fol != nil {
		return nil, err_fol

	}

	defer resp_fol.Body.Close()

	var githubUserFollowers model.GitHubFollowers
	var followers []string

	if resp_fol.StatusCode != http.StatusOK {
		return nil, nil
	}

	// Decode user followers response
	err_fol = json.NewDecoder(resp_fol.Body).Decode(&githubUserFollowers)
	if err_fol != nil {
		return nil, err_fol
	}

	for _, fol := range githubUserFollowers {
		followers = append(followers, fol.UserName)
	}

	if len(followers) == 0 {
		return nil, nil
	}

	return followers, nil
}

// Get user repos from GitHub
func handleUserReposFromGit(userName string) ([]string, error) {
	userReposURL := fmt.Sprintf("https://api.github.com/users/%s/repos", userName)
	resp_rep, err_rep := http.Get(userReposURL)
	if err_rep != nil {
		return nil, err_rep
	}
	defer resp_rep.Body.Close()

	if resp_rep.StatusCode != http.StatusOK {
		return nil, nil
	}

	var githubUserRepos model.GitHubRepos
	var repos []string

	if resp_rep.StatusCode != http.StatusOK {
		return nil, nil
	}

	// Decode user followers response
	err_rep = json.NewDecoder(resp_rep.Body).Decode(&githubUserRepos)
	if err_rep != nil {
		return nil, err_rep
	}

	for _, rep := range githubUserRepos {
		repos = append(repos, rep.Name)
	}

	if len(repos) == 0 {
		return nil, nil
	}

	return repos, nil
}
