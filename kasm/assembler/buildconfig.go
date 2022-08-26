package assembler

import "path/filepath"

type BuildConfig struct {
	InputFilename       string
	OutputFilename      string
	OutputAssembledFile bool
}

func (b *BuildConfig) SetDefaultOutputFilename() {
	b.OutputFilename = b.InputFilename[:len(b.InputFilename)-len(filepath.Ext(b.InputFilename))] + ".kcpu"
}
