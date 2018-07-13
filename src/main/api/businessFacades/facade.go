package businessFacades

import (
	"main/proofs/builder"
	"main/proofs/interpreter"
)

const (
	// BOOT_ADDRESS = 0
	BOOT_SECTOR = 0
	SECTOR_SIZE = 0
)

type ComputerFacade struct {
	Insert *builder.InsertData
	ram    *interpreter.Memory
	hd     *interpreter.HardDrive
}

func NewComputerFacade() *ComputerFacade {
	return &ComputerFacade{new(builder.InsertData), new(interpreter.Memory), new(interpreter.HardDrive)}
}

// func (c *ComputerFacade) Start(boot_address int) {
// 	c.processor.Freeze()
// 	c.ram.Load(boot_address, c.hd.Read(BOOT_SECTOR, SECTOR_SIZE))
// 	c.processor.Jump(boot_address)
// 	c.processor.Execute()
// }

func (c *ComputerFacade) End() {
	// c.processor.Freeze()
	// // c.ram.Load(BOOT_ADDRESS, c.hd.Read(BOOT_SECTOR, SECTOR_SIZE))
	// // c.processor.Jump(BOOT_ADDRESS)
	// c.processor.Execute()

	c.Insert.TDPInsert("hash", "secret", "profileId", "rootHash")
	c.hd.Read(10, 30)
}
