package internal

type RedirectServicePort interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
