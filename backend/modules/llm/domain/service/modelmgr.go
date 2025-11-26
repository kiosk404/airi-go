package service

//go:generate mockgen -destination=mocks/manage.go -package=mocks . ModelManager
type ModelManager interface {
}
