package conf

//go:generate mockgen -destination=mocks/runtime.go -package=mocks . IConfigRuntime
type IConfigRuntime interface {
	NeedCvtURLToBase64() bool
}
