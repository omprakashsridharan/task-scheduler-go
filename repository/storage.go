package repository

type Storage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	delete(key string) error
	exists(key string) (bool, error)
}
