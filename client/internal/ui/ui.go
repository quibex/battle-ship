package ui

import "fmt"

type authService interface {
	Login(login, password string) (string, error)
	Register(login, password string) (string, error)
}

type UI struct {
	auth authService
}

func New(auth authService) *UI {
	return &UI{
		auth: auth,
	}
}

func (ui *UI) MustRun() {

	authCommands := []string{"login", "register"}

	ui.printCommands(authCommands)

	var command string
	if _, err := fmt.Scan(&command); err != nil {
		panic(err)
	}

	switch command {
	case "login":
		var login, password string
		fmt.Println("Enter login:")
		if _, err := fmt.Scan(&login); err != nil {
			panic(err)
		}
		fmt.Println("Enter password:")
		if _, err := fmt.Scan(&password); err != nil {
			panic(err)
		}
		msg, err := ui.auth.Login(login, password)
		if err != nil {
			panic(err)
		}
		fmt.Println(msg)
	case "register":
		var login, password string
		fmt.Println("Enter login:")
		if _, err := fmt.Scan(&login); err != nil {
			panic(err)
		}
		fmt.Println("Enter password:")
		if _, err := fmt.Scan(&password); err != nil {
			panic(err)
		}
		msg, err := ui.auth.Register(login, password)
		if err != nil {
			panic(err)
		}
		fmt.Println(msg)
	default:
		fmt.Println("Unknown command")
	}

}

func (ui *UI) printCommands(commands []string) {
	for i, cmd := range commands {
		println(i, cmd)
	}
}