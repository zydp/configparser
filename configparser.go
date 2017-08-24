package configparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

/*the handle*/
type configFile struct {
	fileName_     string                 `The configuration file name`
	moduleCount_  uint32                 `The total number of modules`
	itemCount_    uint32                 `The total number of configuration items`
	itemsRwMutex_ sync.RWMutex           `Lock the items`
	configItems_  map[string]interface{} `Configuration items storage structure`
}

/*create a new config instance*/
func NewConfigInstance(filename string) (cfile *configFile, err error) {
	cfile = &configFile{fileName_: filename, configItems_: make(map[string]interface{})}
	cfile.itemsRwMutex_.Lock()
	defer cfile.itemsRwMutex_.Unlock()
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return cfile, nil
	} else {
		fd, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if nil != err {
			return cfile, err
		}
		err = cfile.dealFileContentByReader(bufio.NewReader(fd))
		defer fd.Close()
	}
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
	cfile.itemsRwMutex_.RLock()
	defer cfile.itemsRwMutex_.RUnlock()
	optModule, ok := cfile.configItems_[module]
	if !ok {
		return ""
	}
	mapModule := optModule.(map[string]interface{})
	value, ok := mapModule[name]
	if !ok {
		return ""
	}
	val, ok = value.(string) //, nil
	if ok {
		return val
	}
	return ""
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

/*parse the file*/
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
		case ';':
			fallthrough
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

/*change a item value or create a new module and item */
func (cfile *configFile) SetItemValue(module, name string, value interface{}) (err error) {
	if nil == value || "" == module || "" == name {
		return errors.New("name or value is invalid ,please check")
	}
	cfile.itemsRwMutex_.Lock()
	defer cfile.itemsRwMutex_.Unlock()
	Module := cfile.getModule(module)
	mapModule := Module.(map[string]interface{})
	mapModule[name] = fmt.Sprint(value)
	return
}

/*check the module , if not exist create a new module*/
func (cfile *configFile) getModule(module string) interface{} {
	Module, ok := cfile.configItems_[module]
	if !ok {
		newModule := make(map[string]interface{})
		cfile.configItems_[module] = newModule
		return newModule
	}
	return Module
}

/*delete a item*/
func (cfile *configFile) DelItem(module, name string) (err error) {
	if "" == module || "" == name {
		return errors.New("name is invalid ,please check")
	}
	cfile.itemsRwMutex_.Lock()
	defer cfile.itemsRwMutex_.Unlock()
	Module, ok := cfile.configItems_[module]
	if !ok {
		return errors.New("module is not exist")
	}
	mapModule := Module.(map[string]interface{})
	delete(mapModule, name)
	return
}

/*delete a module*/
func (cfile *configFile) DelModule(module string) (err error) {
	if "" == module {
		return errors.New("name is invalid ,please check")
	}
	cfile.itemsRwMutex_.Lock()
	defer cfile.itemsRwMutex_.Unlock()
	delete(cfile.configItems_, module)
	return
}

/*save a new configuration file*/
func (cfile *configFile) SaveToFile(filename string) (err error) {
	if "" == filename {
		return errors.New("name can't set empty")
	}
	var tmpFileName = filename
	if !strings.HasSuffix(tmpFileName, ".conf") {
		tmpFileName += ".conf"
	}
	bakName := tmpFileName
	var newfile *os.File = nil
	var ok bool = false
	for i := 0; i < 100; i++ {
		if _, err = os.Stat(tmpFileName); os.IsNotExist(err) {
			ok = true
			break
		}
		if i > 0 {
			tmpFileName = fmt.Sprint(bakName, ".new.", i)
		} else {
			tmpFileName = fmt.Sprint(bakName, ".new")
		}
	}
	if !ok {
		return errors.New("cannot use this name to save")
	}
	newfile, err = os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE, 0666)
	if nil != err {
		return
	}
	return cfile.saveToFile(newfile)
}
func (cfile *configFile) saveToFile(filehandle *os.File) (err error) {
	cfile.itemsRwMutex_.Lock()
	defer cfile.itemsRwMutex_.Unlock()
	for moduleName, module := range cfile.configItems_ {
		filehandle.WriteString("[" + moduleName + "]\n")
		for k, v := range module.(map[string]interface{}) {
			filehandle.WriteString(k + " = " + v.(string) + "\n")
		}
	}
	return
}
