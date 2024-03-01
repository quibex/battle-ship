package terminalUI

type auth interface {
	Authorization()
	GetUserName() string
}

type game interface {
	StartGame(userName string)
	ReadyToBattle()
	StartBattle() (win bool)
	SendResult(user1, user2 string)
	GetOpponentName() string
}

type TerminalUI struct {
	auth auth
	game game
}

func New(a auth, g game) *TerminalUI {
	return &TerminalUI{
		auth: a,
		game: g,
	}
}

func (f *TerminalUI) MustRun() {
	f.auth.Authorization()
	for {
		f.game.StartGame(f.auth.GetUserName())
		f.game.ReadyToBattle()
		win := f.game.StartBattle()
		if win {
			f.game.SendResult(f.auth.GetUserName(), f.game.GetOpponentName())
		}
	}
}
