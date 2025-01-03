package entities

type EntityRepo interface {
	userRepo
	accountRepo
	entryRepo
	transferRepo
	sessionRepo
}
