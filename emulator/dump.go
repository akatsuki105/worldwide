package emulator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func (cpu *CPU) gobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(cpu.Reg); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.RAM); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.Header); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.ROMBankPtr); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.ROMBank); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.RAMBankPtr); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.RAMBank); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.bankMode); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.tileCache); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.mapCache); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.VRAMCache); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.VRAMModified); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (cpu *CPU) gobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(&cpu.Reg); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.RAM); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.Header); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.ROMBankPtr); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.ROMBank); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.RAMBankPtr); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.RAMBank); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.bankMode); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.tileCache); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.mapCache); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.VRAMCache); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.VRAMModified); err != nil {
		return err
	}
	return nil
}

func (cpu *CPU) dumpData() {
	fmt.Println("dump...")
	time.Sleep(time.Millisecond * 200)

	dumpname := fmt.Sprintf("./dump/%s.dmp", cpu.Header.Title)
	dumpfile, err := os.Create(dumpname)
	if err != nil {
		fmt.Println("dump failed.")
	}
	defer dumpfile.Close()

	data, err := cpu.gobEncode()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = dumpfile.Write(data)
	if err != nil {
		fmt.Println("dump failed.")
	} else {
		fmt.Println("dump success!")
	}
}

func (cpu *CPU) loadData() {
	fmt.Println("loading dumpfile...")
	time.Sleep(time.Millisecond * 200)

	dumpname := fmt.Sprintf("./dump/%s.dmp", cpu.Header.Title)
	data, err := ioutil.ReadFile(dumpname)
	if err != nil {
		fmt.Println("loading dumpfile failed.")
		return
	}

	if err = cpu.gobDecode(data); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("loading dumpfile success!")
}
