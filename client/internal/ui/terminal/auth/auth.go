package authUI

import (
	"fmt"
	"os"
)

type authService interface {
	Login(login, password string) error
	Register(login, password string) error
	GerUserLogin() string
	Logout()
}

type AuthUI struct {
	auth authService
}

func New(authService authService) *AuthUI {
	return &AuthUI{
		auth: authService,
	}
}

func (a *AuthUI) Authorization() {

	for { // while unauthorized
		var login string
		var password string
		for { // while invalid input
			fmt.Println("Login: ")
			cntScan, err := fmt.Scan(&login)
			if err != nil || cntScan != 1 {
				fmt.Println("Invalid input")
			} else {
				break
			}
		}
		for { // while invalid input
			fmt.Println("Password: ")
			cntScan, err := fmt.Scan(&password)
			if err != nil || cntScan != 1 {
				fmt.Println("Invalid input")
			} else {
				break
			}
		}

		var command int
		for { // while invalid input
			fmt.Println("Select command: \n1. Login\n2. Register\n3. Exit\nEnter number of command: ")
			cntScan, err := fmt.Scan(&command)
			if cntScan == 1 && err == nil {
				switch command {
				case 1:
					err = a.auth.Login(login, password)
				case 2:
					err = a.auth.Register(login, password)
				case 3:
					os.Exit(0)
				default:
					fmt.Println("Invalid number of command")
					continue
				}
				if err != nil {
					fmt.Println(err.Error())
					break
				} else {
					fmt.Println("Success authorization")
					return
				}
			} else {
				fmt.Println("Invalid input")
			}
		}
	}
}

func (a *AuthUI) GetUserName() string {
	return a.auth.GerUserLogin()
}
