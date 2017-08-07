package Commtools

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

type configFile struct {
	fileName_    string                 `The configuration file name`
	moduleCount_ uint32                 `The total number of modules`
	itemCount_   uint32                 `The total number of configuration items`
	configItems_ map[string]interface{} `Configuration items storage structure`
}

/*create a new config instance*/
func NewConfigInstance(filename string) (cfile *configFile, err error) {
	cfile = &configFile{fileName_: filename}
	cfile.configItems_ = make(map[string]interface{})
	fd, err := os.Open(filename)
	if nil != err {
		return nil, err
	}
	defer fd.Close()
	err = cfile.dealFileContentByReader(bufio.NewReader(fd))
	return
}

/*return the config file name*/
func (cfile *configFile) GetFileName() string {
	return cfile.fileName_
}

/*return module count*/
func (cfile *configFile) GetModuleCount() uint32 {
	return cfile.moduleCount_
}

/*return items count*/
func (cfile *configFile) GetItemCount() uint32 {
	return cfile.itemCount_
}

/*if does not contain the configuration items, returns an empty*/
func (cfile *configFile) GetStrConfItem(module, name string) (val string) {
	optModule, ok := cfile.configItems_[module]
	if false == ok {
		return ""
	}
	mapModule := optModule.(map[string]interface{})
	value, ok := mapModule[name]
	if false == ok {
		return ""
	}
	return value.(string) //, nil
}

/*if does not contain the configuration items, returns -1*/
func (cfile *configFile) GetIntegerConfItem(module, name string) (val int64) {
	_val := cfile.GetStrConfItem(module, name)
	if "" == _val {
		return -1
	}
	val, err := strconv.ParseInt(_val, 10, 64)
	if nil != err {
		return -1
	}
	return
}

/*if does not contain the configuration items, returns -1*/
func (cfile *configFile) GetFloatConfItem(module, name string) (val float64) {
	_val := cfile.GetStrConfItem(module, name)
	if "" == _val {
		return -1
	}
	val, err := strconv.ParseFloat(_val, 64)
	if nil != err {
		return -1
	}
	return
}

func (cfile *configFile) dealFileContentByReader(fileReader *bufio.Reader) error {
	if nil == fileReader {
		return errors.New("the reader is nil")
	}
	var newModuleName string
	newModule := make(map[string]interface{})
	for {
		line, _, err := fileReader.ReadLine()
		if nil != err {
			if err != io.EOF {
				return err
			}
			break
		}
		strLine := strings.Trim(string(line), " ")
		if len(strLine) < 3 {
			continue
		}
		switch strLine[0] {
		case '#':
			continue
		case '[':
			endIndex := strings.IndexByte(strLine, ']')
			if endIndex < 0 {
				continue
			}
			if newModuleName != "" {
				cfile.configItems_[newModuleName] = newModule
				newModule = make(map[string]interface{})
			}
			newModuleName = strings.Trim(strLine[1:endIndex], " ")
			cfile.moduleCount_++
		default:
			sepIndex := strings.IndexByte(strLine, '=')
			if sepIndex < 0 {
				continue
			}
			newModule[strings.Trim(strLine[:sepIndex], " ")] = strings.Trim(strLine[sepIndex+1:], " ")
			cfile.itemCount_++
		}
	}
	if newModuleName != "" {
		cfile.configItems_[newModuleName] = newModule
	}
	return nil
}
