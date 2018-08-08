package service

type DBService interface {
	HasElementBeforeID(id string) (bool, error)
	HasElementAfterID(id string) (bool, error)

	Count() (int, error)
}
