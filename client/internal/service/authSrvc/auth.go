package authSrvc

type authMQ interface {
	Login(login, password string) error
	Register(login, password string) error
}

type Auth struct {
	mq    authMQ
	login string
}

func New(mq authMQ) *Auth {
	return &Auth{
		mq: mq,
	}
}

func (a *Auth) Login(login, password string) error {
	err := a.mq.Login(login, password)
	if err != nil {
		return err
	}
	a.login = login
	return nil
}

func (a *Auth) Register(login, password string) error {
	err := a.mq.Register(login, password)
	if err != nil {
		return err
	}
	a.login = login
	return nil
}

func (a *Auth) GerUserLogin() string {
	return a.login
}

func (a *Auth) Logout() {
	a.login = ""
}
