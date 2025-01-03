package tokens

type TokenMaker interface {
	CreateToken(username string, duration int64) (string, error)
	VerifyToken(token string) (*Payload, error)
}
