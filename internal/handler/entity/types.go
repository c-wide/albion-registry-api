package entity

type BaseParams struct {
	Region string `param:"region" validate:"required,oneof=americas asia europe"`
	ID     string `param:"id" validate:"required"`
}
