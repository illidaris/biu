package element

type BiuFile struct {
	Path       string
	Name       string
	Package    *BiuPackage
	BiuStructs []*BiuStruct
}
