package entities

type Repo interface {
	UserRepo
}

type RepoStub struct {
	UserRepoStub
}
