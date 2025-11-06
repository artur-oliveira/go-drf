package models

type IModel interface {
	TableName() string

	ModuleName() string
}
