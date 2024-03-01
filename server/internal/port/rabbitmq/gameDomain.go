package rabbitmq

type gameCreateRequest struct {
	UserName string `json:"user_name"`
}

type gameJoinRequest struct {
	CreatorUserName string `json:"creator_user_name"`
	JoiningUserName string `json:"joining_user_name"`
}

type gameDelRequest struct {
	UserName string `json:"user_name"`
}

type gameDelResponse struct {
	Err string `json:"error,omitempty"`
}

type getAvailableGamesRequest struct{}

type gameResultRequest struct {
	Winner string `json:"winner"`
	Loser  string `json:"loser"`
}

type getStatRequest struct {
	UserName string `json:"user_name"`
}

type getStatResponse struct {
	Rating int    `json:"rating"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Err    string `json:"error,omitempty"`
}

type getAvailableGamesResponse struct {
	Games []string `json:"games"`
	Err   string   `json:"error,omitempty"`
}

type gameCreateResponse struct {
	User2 string `json:"user2,omitempty"`
	Err   string `json:"error,omitempty"`
}

type gameJoinResponse struct {
	Err string `json:"error,omitempty"`
}

type gameResultResponse struct {
	Err string `json:"error,omitempty"`
}

const ( //queue names
	gameCreate        = "game.create"
	gameJoin          = "game.join"
	getAvailableGames = "game.get_available"
	saveGameResult    = "game.save_result"
	getUserStat       = "game.get_user_stat"
	gameDel           = "game.del"
)
