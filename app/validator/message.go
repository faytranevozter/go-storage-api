package validator

import "fmt"

// RequiredErr Error message for required
func RequiredErr(attribute string) string {
	return fmt.Sprintf("The %s field is required", attribute)
}

// InArrayErr Error message for in array
func InArrayErr(attribute string) string {
	return fmt.Sprintf("The selected %s is invalid.", attribute)
}

// GreaterThanErr Error message for in array
func GreaterThanErr(attribute, val string) string {
	return fmt.Sprintf("The selected %s must greater than "+val+".", attribute)
}
