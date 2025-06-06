package constants

import "fmt"

var USER_ALREADY_AUTHENTICATED = func(username string) string {
	return fmt.Sprintf("User %s is already authenticated", username)
}

var USER_SUCCESS_AUTHENTICATED = func(username string) string {
	return fmt.Sprintf("Succesfully authenticated user %s", username)
}
