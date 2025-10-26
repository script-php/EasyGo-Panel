package auth

import (
	"errors"
	"os"
	"os/user"

	"github.com/msteinert/pam"
)

// User represents a system user
type User struct {
	Username string
	UID      int
	GID      int
	Home     string
	Shell    string
}

// GetCurrentUser returns the current system user
func GetCurrentUser() (*User, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	uid := os.Getuid()
	gid := os.Getgid()

	return &User{
		Username: currentUser.Username,
		UID:      uid,
		GID:      gid,
		Home:     currentUser.HomeDir,
		Shell:    "/bin/bash", // Default shell
	}, nil
}

// AuthenticateUser authenticates a user using PAM
func AuthenticateUser(username, password string) error {
	t, err := pam.StartFunc("login", username, func(s pam.Style, msg string) (string, error) {
		switch s {
		case pam.PromptEchoOff:
			return password, nil
		case pam.PromptEchoOn:
			return username, nil
		case pam.ErrorMsg:
			return "", errors.New(msg)
		case pam.TextInfo:
			return "", nil
		}
		return "", errors.New("unrecognized PAM message style")
	})

	if err != nil {
		return err
	}

	if err := t.Authenticate(0); err != nil {
		return err
	}

	return nil
}

// IsRoot checks if the current user is root
func IsRoot() bool {
	return os.Getuid() == 0
}

// RequireRoot ensures the user is running as root
func RequireRoot() error {
	if !IsRoot() {
		return errors.New("this operation requires root privileges")
	}
	return nil
}

// GetSystemUsers returns all system users
func GetSystemUsers() ([]*User, error) {
	var users []*User
	
	// Implementation would read /etc/passwd
	// For now, return current user as example
	currentUser, err := GetCurrentUser()
	if err != nil {
		return nil, err
	}
	
	users = append(users, currentUser)
	return users, nil
}