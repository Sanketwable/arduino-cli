// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.
package debug

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/arduino/arduino-cli/arduino/cores/packagemanager"
	dbg "github.com/arduino/arduino-cli/rpc/debug"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var customHardware = paths.New("testdata", "custom_hardware")
var dataDir = paths.New("testdata", "data_dir", "packages")
var sketch = "hello"
var sketchPath = paths.New("testdata", sketch)

func TestGetCommandLine(t *testing.T) {
	pm := packagemanager.NewPackageManager(nil, nil, nil, nil)
	pm.LoadHardwareFromDirectory(customHardware)
	pm.LoadHardwareFromDirectory(dataDir)

	// Windows tools have .exe extension
	var toolExtension = ""
	if runtime.GOOS == "windows" {
		toolExtension = ".exe"
	}

	// Arduino Zero has an integrated debugger port, anc it could be debugged directly using USB
	req := &dbg.DebugConfigReq{
		Instance:   &dbg.Instance{Id: 1},
		Fqbn:       "arduino-test:samd:arduino_zero_edbg",
		SketchPath: sketchPath.String(),
		Port:       "none",
	}

	goldCommand := fmt.Sprintf("%s/arduino-test/tools/arm-none-eabi-gcc/7-2017q4/bin//arm-none-eabi-gdb%s", dataDir, toolExtension) +
		" -ex target extended-remote |" +
		fmt.Sprintf(" %s/arduino-test/tools/openocd/0.10.0-arduino7/bin/openocd%s", dataDir, toolExtension) +
		fmt.Sprintf(" -s \"%s/arduino-test/tools/openocd/0.10.0-arduino7/share/openocd/scripts/\"", dataDir) +
		fmt.Sprintf(" --file \"%s/arduino-test/samd/variants/arduino_zero/openocd_scripts/arduino_zero.cfg\"", customHardware) +
		fmt.Sprintf(" -c \"gdb_port pipe\" -c \"telnet_port 0\" -c init -c halt %s/hello.arduino-test.samd.arduino_zero_edbg.elf", sketchPath)

	command, err := getCommandLine(req, pm)
	assert.Nil(t, err)
	commandToTest := strings.Join(command[:], " ")
	assert.Equal(t, filepath.FromSlash(goldCommand), filepath.FromSlash(commandToTest))

	// Other samd boards such as mkr1000 can be debugged using an external tool such as Atmel ICE connected to
	// the board debug port
	req2 := &dbg.DebugConfigReq{
		Instance:   &dbg.Instance{Id: 1},
		Fqbn:       "arduino-test:samd:mkr1000",
		SketchPath: sketchPath.String(),
		Port:       "none",
	}

	goldCommand2 := fmt.Sprintf("%s/arduino-test/tools/arm-none-eabi-gcc/7-2017q4/bin//arm-none-eabi-gdb%s", dataDir, toolExtension) +
		" -ex target extended-remote |" +
		fmt.Sprintf(" %s/arduino-test/tools/openocd/0.10.0-arduino7/bin/openocd%s", dataDir, toolExtension) +
		fmt.Sprintf(" -s \"%s/arduino-test/tools/openocd/0.10.0-arduino7/share/openocd/scripts/\"", dataDir) +
		fmt.Sprintf(" --file \"%s/arduino-test/samd/variants/mkr1000/openocd_scripts/arduino_zero.cfg\"", customHardware) +
		fmt.Sprintf(" -c \"gdb_port pipe\" -c \"telnet_port 0\" -c init -c halt %s/hello.arduino-test.samd.mkr1000.elf", sketchPath)

	command2, err := getCommandLine(req2, pm)
	assert.Nil(t, err)
	commandToTest2 := strings.Join(command2[:], " ")
	assert.Equal(t, filepath.FromSlash(goldCommand2), filepath.FromSlash(commandToTest2))

}
