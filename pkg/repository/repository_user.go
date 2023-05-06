package repository

import (
	db "backend/internal/config"
	"backend/internal/model"
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserEmailExists = errors.New("user email exists")

// GetUserByID retrieves a user from the database by their ID
func GetUserByID(userID int) (model.User, error) {
	var user model.User

	row := db.Db.QueryRow("SELECT * FROM user WHERE id = ?", userID)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.TimeZoneID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

// UpdateUser updates a user's information in the database
func UpdateUser(user model.User) error {
	var exists bool
	db.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM user WHERE userEmail = ?)", user.Email).Scan(&exists)
	if exists {
		return ErrUserEmailExists
	}

	_, err := db.Db.Exec("UPDATE user SET userFirstName = ?, userLastName = ?, userTimeZoneId = ?, userEmail = ? WHERE id = ?", user.FirstName, user.LastName, user.TimeZoneID, user.Email, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

// UpdateUser updates a user's information in the database from GitHub
func UpdateUserFromGit(user model.User) error {
	var exists bool
	db.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM user WHERE userEmail = ? && id != ?)", user.Email, user.ID).Scan(&exists)
	if exists {
		return ErrUserEmailExists
	}

	// Get existing user followers from database
	existingFollowers, err_fol := GetUserFollowers(user.ID)
	if err_fol != nil {
		return err_fol
	}

	// Compare new followers to existing followers and add to the table
	if user.Followers != nil {
		if existingFollowers != nil {
			var followers []string
			for _, f := range existingFollowers {
				followers = append(followers, *f)
			}
			for _, follower := range user.Followers {
				if !containsString(followers, follower) {
					db.Db.Exec("INSERT INTO follower (userId, followerUserName) VALUES (?, ?)", user.ID, follower)
				}
			}
		} else {
			for _, follower := range user.Followers {
				db.Db.Exec("INSERT INTO follower (userId, followerUserName) VALUES (?, ?)", user.ID, follower)
			}
		}

	} else {
		db.Db.Exec("DELETE FROM follower WHERE userId = ?", user.ID)
	}

	// Get existing user repos from database
	existingRepos, err_rep := GetUserRepos(user.ID)
	if err_rep != nil {
		return err_rep
	}

	if user.Repos != nil {
		if existingRepos != nil {
			var repos []string

			for _, r := range existingRepos {
				repos = append(repos, *r)
			}

			for _, repo := range user.Repos {
				if !containsString(repos, repo) {
					db.Db.Exec("INSERT INTO user_repo (userId, userRepos) VALUES (?, ?)", user.ID, repo)
				}
			}
		} else {
			for _, repo := range user.Repos {
				db.Db.Exec("INSERT INTO user_repo (userId, userRepos) VALUES (?, ?)", user.ID, repo)
			}
		}
	}

	_, err := db.Db.Exec("UPDATE user SET userName = ?, userEmail = ?, userFirstName = ?, userLastName = ?, userTimeZoneId = ? WHERE id = ?", user.Username, user.Email, user.FirstName, user.LastName, user.TimeZoneID, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetUserFollowers retrieves a list of followers for a given user ID
func GetUserFollowers(userID int) ([]*string, error) {
	var followers []*string

	rows, err := db.Db.Query("SELECT followerUserName FROM follower WHERE userId = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			var followerUserName *string
			err := rows.Scan(&followerUserName)
			if err != nil {
				return nil, err
			}
			followers = append(followers, followerUserName)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return followers, nil
}

// GetUserRepos retrieves a list of repos for a given user ID
func GetUserRepos(userID int) ([]*string, error) {
	var repos []*string

	rows, err := db.Db.Query("SELECT userRepos FROM user_repo WHERE userId = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			var userRepos *string
			err := rows.Scan(&userRepos)
			if err != nil {
				return nil, err
			}
			repos = append(repos, userRepos)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func containsString(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
