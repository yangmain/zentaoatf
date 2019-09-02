package zentaoUtils

import (
	"fmt"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	dateUtils "github.com/easysoft/zentaoatf/src/utils/date"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GenSuperApiUri(model string, methd string, params [][]string) string {
	var sep string
	if vari.RequestType == constant.RequestTypePathInfo {
		sep = ","
	} else {
		sep = "&"
	}

	paramStr := ""
	i := 0
	for _, p := range params {
		if i > 0 {
			paramStr += sep
		}
		paramStr += p[0] + "=" + p[1]
		i++
	}

	var uri string
	if vari.RequestType == constant.RequestTypePathInfo {
		uri = fmt.Sprintf("api-getmodel-%s-%s-%s.json", model, methd, paramStr)
	} else {
		uri = fmt.Sprintf("?m=api&f=getmodel&model=%s&methodName=%s&params=%s", model, methd, paramStr)
	}

	logUtils.PrintToCmd(uri, -1)
	return uri
}

func GenApiUri(module string, methd string, param string) string {
	if vari.RequestType == constant.RequestTypePathInfo {
		return fmt.Sprintf("%s-%s-%s.json", module, methd, param)
	}

	return ""
}

func ScriptToLogDir() string {
	logDir := vari.WorkDir + constant.LogDir
	fileUtils.MkDirIfNeeded(logDir)

	return logDir + vari.RunDir
}

func ScriptToLogName(file string) string {
	logDir := vari.WorkDir + constant.LogDir
	fileUtils.MkDirIfNeeded(logDir)

	nameWithSuffix := commonUtils.Base(file)

	logFile := logDir + vari.RunDir + nameWithSuffix + ".log"

	return logFile
}

func ScriptToExpectName(file string) string {
	fileSuffix := path.Ext(file)
	expectName := strings.TrimSuffix(file, fileSuffix) + ".ex"

	return expectName
}

func RunDateFolder() string {
	runName := dateUtils.DateTimeStrFmt(time.Now(), "2006-01-02T150405") + string(os.PathSeparator)

	return runName
}

func GetCaseInfo(file string) (int, int, string) {
	content := fileUtils.ReadFile(file)

	var caseId int
	var productId int
	var title string

	myExp := regexp.MustCompile(`[\S\s]*caseId:\s*([^\n]*?)\s*\n`)
	arr := myExp.FindStringSubmatch(content)
	if len(arr) > 1 {
		caseId, _ = strconv.Atoi(arr[1])
	}

	myExp = regexp.MustCompile(`[\S\s]*productId:\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(content)
	if len(arr) > 1 {
		productId, _ = strconv.Atoi(arr[1])
	}

	myExp = regexp.MustCompile(`[\S\s]*title:\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(content)
	if len(arr) > 1 {
		title = arr[1]
	}

	return caseId, productId, title
}

func ReadExpect(file string) [][]string {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(`<<<TC[\S\s]*expects:[^\n]*\n+([\S\s]*?)(readme:|TC)`)
	arr := myExp.FindStringSubmatch(content)

	str := ""
	if len(arr) > 1 {
		expects := arr[1]

		if strings.Index(expects, "@file") > -1 {
			str = fileUtils.ReadFile(ScriptToExpectName(file))
		} else {
			str = commonUtils.RemoveBlankLine(expects)
		}
	}

	ret := GenExpectArr(str)

	return ret
}

func ReadCheckpointSteps(file string) []string {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(`<<<TC[\S\s]*steps:[^\n]*\n*([\S\s]*)\n+expects:`)
	arr := myExp.FindStringSubmatch(content)

	str := ""
	if len(arr) > 1 {
		checkpoints := arr[1]
		str = commonUtils.RemoveBlankLine(checkpoints)
	}

	ret := GenCheckpointStepArr(str)

	return ret
}

func GenCheckpointStepArr(str string) []string {
	ret := make([]string, 0)
	for _, line := range strings.Split(str, "\n") {
		line := strings.TrimSpace(line)

		if strings.Index(line, "@") == 0 {
			ret = append(ret, line)
		}
	}

	return ret
}

func GenExpectArr(str string) [][]string {
	_, arr := GenArr(str, false)
	return arr
}

func GenLogArr(str string) (bool, [][]string) {
	skip, arr := GenArr(str, true)
	return skip, arr
}
func GenArr(str string, checkSkip bool) (bool, [][]string) {
	ret := make([][]string, 0)
	indx := -1
	for _, line := range strings.Split(str, "\n") {
		line := strings.TrimSpace(line)

		if checkSkip && strings.ToLower(line) == "skip" {
			return true, nil
		}

		if strings.Index(line, "#") == 0 {
			ret = append(ret, make([]string, 0))
			indx++
		} else if indx > -1 {
			if len(line) > 0 && indx < len(ret) {
				ret[indx] = append(ret[indx], line)
			}
		}
	}

	return false, ret
}

func ReadLog(logFile string) (bool, [][]string) {
	str := fileUtils.ReadFile(logFile)

	skip, ret := GenLogArr(str)
	return skip, ret
}
