package interfaces

type BaseSerializer interface {
	Serialize(value *struct{}) *struct{}
}