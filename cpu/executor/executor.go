package executor

import (
	"bufio"
	"fmt"
	"os"
)

func Execute(f string) error {
	cpu := NewCpu()
	if err := loadProgram(f, &cpu); err != nil {
		return err
	}

	if err := cpu.StartExecution(); err != nil {
		fmt.Println(cpu.ToString())
		return err
	}

	return nil
}

func loadProgram(f string, cpu *Cpu) error {
	file, err := os.Open(f)

	if err != nil {
		return err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	// calculate the bytes size
	var size int64 = info.Size()
	bytes := make([]byte, size)

	// read into buffer
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		return err
	}

	copy(cpu.Memory, bytes)

	return nil
}
