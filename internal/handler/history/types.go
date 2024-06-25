package history

type BaseParams struct {
	ID     string `param:"id" validate:"required"`
	Region string `param:"region" validate:"required,oneof=americas asia europe"`
	Limit  int32  `query:"limit" validate:"omitempty,min=1,max=50"`
	Offset int32  `query:"offset" validate:"omitempty,min=0"`
}
