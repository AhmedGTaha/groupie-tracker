package main

import (
	"fmt"
	"learning-go/helper"
	"sync"
	"time"
)

// const is used for values that should not change while the program runs.
const conferenceTickets int = 50

// These are package-level variables, so every function in this file can use them.
var conferenceName = "Go Conference"
var remainingTickets uint = 50

// bookings is a slice of UserData structs.
// make creates an empty slice that is ready for append.
var bookings = make([]UserData, 0)

// A struct lets us group related values into one custom type.
// One booking has a first name, last name, email, and ticket count.
type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
}

// WaitGroup lets main wait for goroutines to finish before the program exits.
var wg = sync.WaitGroup{}

func main() {

	greetUsers()

	// This loop keeps accepting bookings until tickets are sold out.
	for {
		firstName, lastName, email, userTickets := getUserInput()

		// The validation function lives in the helper package.
		// Because it starts with a capital letter, main can call it.
		isValidName, isValidEmail, isValidTicketNumber := helper.ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

		if isValidName && isValidEmail && isValidTicketNumber {

			bookTicket(userTickets, firstName, lastName, email)

			// Add 1 before starting the goroutine so WaitGroup knows work is coming.
			wg.Add(1)

			// go starts sendTicket in a goroutine, which runs at the same time as main.
			go sendTicket(userTickets, firstName, lastName, email)

			firstNames := getFirstNames()
			fmt.Printf("The first names of bookings are: %v\n", firstNames)

			if remainingTickets == 0 {
				fmt.Println("Our conference is booked out. Come back next year.")
				break
			}
		} else {
			if !isValidName {
				fmt.Println("first name or last name you entered is too short")
			}
			if !isValidEmail {
				fmt.Println("email address you entered doesn't contain @ sign")
			}
			if !isValidTicketNumber {
				fmt.Println("number of tickets you entered is invalid")
			}
		}
	}

	// Wait blocks here until every wg.Done has been called.
	wg.Wait()
}

func greetUsers() {
	fmt.Printf("Welcome to %v booking application\n", conferenceName)
	fmt.Printf("We have total of %v tickets and %v are still available.\n", conferenceTickets, remainingTickets)
	fmt.Println("Get your tickets here to attend")
}

func getFirstNames() []string {
	firstNames := []string{}

	// range loops over all bookings.
	// We use _ because we do not need the index number.
	for _, booking := range bookings {
		firstNames = append(firstNames, booking.firstName)
	}
	return firstNames
}

func getUserInput() (string, string, string, uint) {
	var firstName string
	var lastName string
	var email string
	var userTickets uint

	// Scan stores the typed value into the variable.
	// The & gives Scan the variable's memory address.
	fmt.Println("Enter your first name: ")
	fmt.Scan(&firstName)

	fmt.Println("Enter your last name: ")
	fmt.Scan(&lastName)

	fmt.Println("Enter your email address: ")
	fmt.Scan(&email)

	fmt.Println("Enter number of tickets: ")
	fmt.Scan(&userTickets)

	return firstName, lastName, email, userTickets
}

func bookTicket(userTickets uint, firstName string, lastName string, email string) {
	remainingTickets = remainingTickets - userTickets

	// This creates one UserData value and fills its fields.
	var userData = UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}

	bookings = append(bookings, userData)
	fmt.Printf("List of bookings is %v\n", bookings)

	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", remainingTickets, conferenceName)
}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {
	// defer runs wg.Done at the end of this function, even if the function changes later.
	defer wg.Done()

	// Sleep pretends that sending an email takes time.
	time.Sleep(5 * time.Second)

	var ticket = fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)
	fmt.Println("#################")
	fmt.Printf("Sending ticket:\n %v \nto email address %v\n", ticket, email)
	fmt.Println("#################")
}
