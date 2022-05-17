package lang

type Instruction struct {
	prgrm   *Program
	exprsns []*Expression
}

type Instructions struct {
	instrctns []*Instruction
}
