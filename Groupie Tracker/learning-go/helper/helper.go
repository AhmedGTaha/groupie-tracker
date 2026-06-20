package helper

import "strings"

// ValidateUserInput starts with a capital letter, so it is exported.
// Exported names can be used from another package, like main.
func ValidateUserInput(firstName string, lastName string, email string, userTickets uint, remainingTickets uint) (bool, bool, bool) {
	// Each variable stores one true/false validation result.
	isValidName := len(firstName) >= 2 && len(lastName) >= 2
	isValidEmail := strings.Contains(email, "@")
	isValidTicketNumber := userTickets > 0 && userTickets <= remainingTickets

	// Go functions can return multiple values.
	return isValidName, isValidEmail, isValidTicketNumber
}
