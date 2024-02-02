package auth


type authMQ interface {
	Login(login, password string) error
	Register(login, password string) error
}

type Auth struct {
	mq  authMQ
}

func New(mq authMQ) *Auth {
	return &Auth{
		mq:  mq,
	}
}

func (a *Auth) Login(login, password string) (string, error) {
	err := a.mq.Login(login, password)
	if err != nil {
		return err.Error(), nil
	}
	return "ok", nil
}

func (a *Auth) Register(login, password string) (string, error) {
	err := a.mq.Register(login, password)
	if err != nil {
		return err.Error(), nil
	}
	return "ok", nil
}
