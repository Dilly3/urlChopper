package internal

type RedirectRepositoryPort interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
