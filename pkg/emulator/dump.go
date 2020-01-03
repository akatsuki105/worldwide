package emulator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/sqweek/dialog"
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
	if err := encoder.Encode(cpu.ROMBankPtr); err != nil {
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
	if err := decoder.Decode(&cpu.ROMBankPtr); err != nil {
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
	var dumpname string
	if runtime.GOOS == "windows" {
		tmp, err := dialog.File().Filter("save file(.dmp)", "dmp").Title("Save data into file").Save()
		if err != nil {
			dialog.Message("%s", "dump data failed.").Title("Error").Error()
			return
		}
		dumpname = tmp
	} else {
		dumpname = fmt.Sprintf("./dump/%s.dmp", cpu.Cartridge.Title)
	}

	dumpfile, err := os.Create(dumpname)
	if err != nil {
		dialog.Message("%s", err.Error()).Title("Error").Error()
		return
	}
	defer dumpfile.Close()

	data, err := cpu.gobEncode()
	if err != nil {
		dialog.Message("%s", err.Error()).Title("Error").Error()
		return
	}

	_, err = dumpfile.Write(data)
	if err != nil {
		dialog.Message("%s", err.Error()).Title("Error").Error()
	}
}

func (cpu *CPU) loadData() {
	time.Sleep(time.Millisecond * 200)
	dumpname := cpu.selectData()
	if dumpname == "" {
		return
	}

	data, err := ioutil.ReadFile(dumpname)
	if err != nil {
		dialog.Message("%s", "load dumpfile failed.").Title("Error").Error()
		return
	}

	if err = cpu.gobDecode(data); err != nil {
		dialog.Message("%s", err.Error()).Title("Error").Error()
		return
	}
}

func (cpu *CPU) selectData() string {
	var filepath string
	switch runtime.GOOS {
	case "windows":
		tmp, err := dialog.File().Filter("GameBoy Dump File", "dmp").Load()
		if err != nil {
			return ""
		}
		filepath = tmp
	default:
		filepath = fmt.Sprintf("./dump/%s.dmp", cpu.Cartridge.Title)
	}
	return filepath
}
