package interpreter

import "fmt"

type Memory struct{}

func (m *Memory) Load(position int, data []byte) {
	fmt.Println("Memory.Load()")
}

type HardDrive struct{}

func (hd *HardDrive) Read(lba int, size int) []byte {
	fmt.Println("HardDrive.Read()")
	return make([]byte, 0)
}
