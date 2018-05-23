//+build windows

package wmic

import (
	"bytes"
	"encoding/csv"
	"errors"
	//	"fmt"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

type WmicCommand struct {
	CmdArgs []string
}

type DiskDrive struct {
	Caption    string
	DeviceID   string
	Model      string
	Partitions uint
	Size       uint64
}

type Port struct {
	Caption       string
	EndingAddress uint
}

type Volume struct {
	Caption    string
	Capacity   uint64
	BootVolume bool
	Compressed bool
}

// build a command line for wmic command and format as csv output
func buildCommand(cmd interface{}) ([]string, error) {
	cmdString := make([]string, 0, 0)
	s := reflect.Indirect(reflect.ValueOf(cmd))
	t := s.Type()
	if s.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("Unknown Interface!")
	}

	cmdString = append(cmdString, t.Name())
	cmdString = append(cmdString, "get")

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	cmdString = append(cmdString, strings.Join(fields, ","))
	cmdString = append(cmdString, "/format:csv")

	return cmdString, nil
}

//Execute the wmic command and return the stdout/stderr
func RunCmd(dst interface{}) (string, error) {
	cmdLineOpt, _ := buildCommand(dst)

	run := exec.Command("wmic", cmdLineOpt...)

	var stdout, stderr bytes.Buffer
	run.Stdout = &stdout
	run.Stderr = &stderr

	err := run.Run()
	if err != nil {
		return string(stderr.Bytes()), err
	} else {
		return string(stdout.Bytes()), err
	}
}

//Parse the csv format output of the RunCmd
func ParseResult(stdout string, dst interface{}) error {
	dv := reflect.ValueOf(dst).Elem()
	t := dv.Type().Elem()

	dv.Set(reflect.MakeSlice(dv.Type(), 0, 0))

	lines := strings.Split(stdout, "\n")
	var header []int = nil

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			v := reflect.New(t)
			r := csv.NewReader(strings.NewReader(line))

			r.FieldsPerRecord = t.NumField() + 1

			records, err := r.ReadAll()
			if err != nil {
				return err
			}
			//Find the field number of the record
			if header == nil {
				header = make([]int, len(records[0]), len(records[0]))
				for i, record := range records[0] {
					for j := 0; j < t.NumField(); j++ {
						if record == t.Field(j).Name {
							header[i] = j
						}
					}
				}
				continue
			} else {
				for i, record := range records[0] {
					f := reflect.Indirect(v).Field(header[i])
					switch t.Field(header[i]).Type.Kind() {
					case reflect.String:
						f.SetString(record)
					case reflect.Uint, reflect.Uint64:
						uintVal, err := strconv.ParseUint(record, 10, 64)
						if err != nil {
							return err
						}
						f.SetUint(uintVal)
					case reflect.Bool:
						bVal, err := strconv.ParseBool(record)
						if err != nil {
							return err
						}
						f.SetBool(bVal)
					default:
						return errors.New("unknown data type!")
					}
				}

			}
			dv.Set(reflect.Append(dv, reflect.Indirect(v)))
		}
	}

	return nil
}

func GetDiskDriveInfo() ([]DiskDrive, error) {
	var disk []DiskDrive

	output, err := RunCmd(disk)

	if err != nil {
		return nil, err
	}

	err = ParseResult(output, &disk)

	return disk, err
}

func GetPortInfo() ([]Port, error) {
	var port []Port
	output, err := RunCmd(port)

	if err != nil {
		return nil, err
	}

	err = ParseResult(output, &port)

	return port, err
}

func GetVolumeInfo() ([]Volume, error) {
	var vol []Volume
	output, err := RunCmd(vol)

	if err != nil {
		return nil, err
	}

	err = ParseResult(output, &vol)

	return vol, err
}
