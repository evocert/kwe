package contexting

type ContextGroup struct {
	contxts map[string]*Context
}

type ContextGrouping struct {
	ctxgrps map[string]*ContextGroup
}
