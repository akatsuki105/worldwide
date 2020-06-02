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
	isCGB := cpu.Cartridge.IsCGB

	if err := encoder.Encode(cpu.Reg); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.RAM); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.ROMBank.ptr); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.RAMBank.ptr); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.RAMBank); err != nil {
		return nil, err
	}
	if err := encoder.Encode(cpu.bankMode); err != nil {
		return nil, err
	}

	if isCGB {
		// ゲームボーイカラー
		if err := encoder.Encode(cpu.WRAMBankPtr); err != nil {
			return nil, err
		}
		if err := encoder.Encode(cpu.WRAMBank); err != nil {
			return nil, err
		}
	}

	if err := encoder.Encode(cpu.GPU); err != nil {
		return nil, err
	}

	if isCGB {
		if err := encoder.Encode(cpu.RTC); err != nil {
			return nil, err
		}
	}

	return w.Bytes(), nil
}

func (cpu *CPU) gobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	isCGB := cpu.Cartridge.IsCGB

	if err := decoder.Decode(&cpu.Reg); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.RAM); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.ROMBank.ptr); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.RAMBank.ptr); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.RAMBank); err != nil {
		return err
	}
	if err := decoder.Decode(&cpu.bankMode); err != nil {
		return err
	}

	if isCGB {
		if err := decoder.Decode(&cpu.WRAMBankPtr); err != nil {
			return err
		}
		if err := decoder.Decode(&cpu.WRAMBank); err != nil {
			return err
		}
	}

	if err := decoder.Decode(&cpu.GPU); err != nil {
		return err
	}

	if isCGB {
		if err := decoder.Decode(&cpu.RTC); err != nil {
			return err
		}
	}

	return nil
}

func (cpu *CPU) dumpData() {
	time.Sleep(time.Millisecond * 200)
	dumpname := fmt.Sprintf("%s/%s.dmp", cpu.romdir, cpu.Cartridge.Title)
	dumpfile, err := os.Create(dumpname)
	if err != nil {
		fmt.Println("Failed to dump: ", err.Error())
		return
	}
	defer dumpfile.Close()

	data, err := cpu.gobEncode()
	if err != nil {
		fmt.Println("Failed to dump: ", err.Error())
		return
	}

	_, err = dumpfile.Write(data)
	if err != nil {
		fmt.Println("Failed to dump: ", err.Error())
		return
	}
}

func (cpu *CPU) loadData() {
	time.Sleep(time.Millisecond * 200)
	dumpname := fmt.Sprintf("%s/%s.sav", cpu.romdir, cpu.Cartridge.Title)

	data, err := ioutil.ReadFile(dumpname)
	if err != nil {
		fmt.Println("Failed to load: ", err.Error())
		return
	}

	if err = cpu.gobDecode(data); err != nil {
		fmt.Println("Failed to load: ", err.Error())
		return
	}
}
