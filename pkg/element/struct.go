package element

type BiuStruct struct {
	Named
	Fields                []*BiuField
	WithBefore, WithAfter bool
}
