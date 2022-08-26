package assembler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hculpan/kcpu/kasm/operations"
)

func BuildFile(buildConfig BuildConfig) {
	fmt.Printf("Building %s\n", buildConfig.InputFilename)

	lines, err := readLines(buildConfig.InputFilename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	fmt.Print("Loading symbols...")
	symbols := NewSymbolsVisitor()
	for line, lineText := range lines {
		symbols.ProcessLine(lineText, line+1)
	}
	if len(symbols.Errors()) > 0 {
		fmt.Println("Errors:")
		for _, err := range symbols.Errors() {
			fmt.Println("  " + err.ToString())
		}
		return
	}
	fmt.Println("done.")

	fmt.Print("Assembling program...")
	visitor := NewBuildVisitor(&symbols)
	for line, lineText := range lines {
		visitor.ProcessLine(lineText, line+1)
	}
	if len(visitor.Errors()) > 0 {
		fmt.Println("Errors:")
		for _, err := range visitor.Errors() {
			fmt.Println("  " + err.ToString())
		}
	} else {
		fmt.Println("done.")
	}

	if err := writeListFile(buildConfig, &visitor, &symbols); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("- List file written to '%s'\n", buildConfig.InputFilename[:len(buildConfig.InputFilename)-len(filepath.Ext(buildConfig.InputFilename))]+".list")

	if err := writeSymbolsFile(buildConfig, &visitor, &symbols); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("- Symbols file written to '%s'\n", buildConfig.InputFilename[:len(buildConfig.InputFilename)-len(filepath.Ext(buildConfig.InputFilename))]+".sym")

	if len(visitor.Errors()) == 0 && buildConfig.OutputAssembledFile {
		if err := writeAssembledFile(buildConfig, &visitor); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("- Assembled program written to '%s'\n", buildConfig.OutputFilename)
		fmt.Println("Build successful")
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeAssembledFile(buildConfig BuildConfig, b *BuildVisitor) error {
	file, err := os.Create(buildConfig.OutputFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, op := range b.AssembledOps {
		if !operations.IsNoCodeOp(op) {
			bytes := []byte{0, 0, 0, 0}
			bytes[0] = op.Op
			bytes[1] = op.Register
			bytes[2] = op.DataH
			bytes[3] = op.DataL
			_, err := file.Write(bytes)
			if err != nil {
				return errors.New("unable to write to file")
			}
		}
	}

	return nil
}

func writeSymbolsFile(buildConfig BuildConfig, b *BuildVisitor, s *SymbolsVisitor) error {
	filename := buildConfig.InputFilename[:len(buildConfig.InputFilename)-len(filepath.Ext(buildConfig.InputFilename))] + ".sym"

	f, err := os.Create(filename)
	if err != nil {
		return errors.New("unable to create list file")
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	// Write symbol table
	w.WriteString("Symbols:\n")
	symbolLines := s.ToStrings()
	for _, line := range symbolLines {
		_, err = w.WriteString(line + "\n")
		if err != nil {
			return errors.New("unable to write listing file")
		}
	}

	w.Flush()

	return nil
}

func writeListFile(buildConfig BuildConfig, b *BuildVisitor, s *SymbolsVisitor) error {
	listFilename := buildConfig.InputFilename[:len(buildConfig.InputFilename)-len(filepath.Ext(buildConfig.InputFilename))] + ".list"

	f, err := os.Create(listFilename)
	if err != nil {
		return errors.New("unable to create list file")
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	// Write code
	lines := b.ToSrings()
	for _, line := range lines {
		_, err = w.WriteString("  " + line + "\n")
		if err != nil {
			return errors.New("unable to write listing file")
		}
	}

	w.Flush()

	return nil
}
