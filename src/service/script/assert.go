package scriptService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	langUtils "github.com/easysoft/zentaoatf/src/utils/lang"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

/**
Get all test script and suite files in current work dir
*/
func LoadAssetFiles() ([]string, []string) {
	caseFiles := make([]string, 0)
	suitesFiles := make([]string, 0)

	GetAllScriptsInDir(constant.ScriptDir, &caseFiles)
	GetAllScriptsInDir(constant.ScriptDir, &suitesFiles)

	return caseFiles, suitesFiles
}

/**
Get all test result histories for specific test script/suite
*/
func LoadTestResults(assert string) []string {
	ret := make([]string, 0)

	dir := constant.LogDir

	mode, name := GetRunModeAndName(assert)
	reg := fmt.Sprintf("%s-%s-(.+)", mode, name)
	myExp := regexp.MustCompile(reg)

	files, _ := ioutil.ReadDir(dir)
	for _, fi := range files {
		if fi.IsDir() {
			arr := myExp.FindStringSubmatch(fi.Name())
			if len(arr) > 1 {
				ret = append(ret, arr[1])
			}
		}
	}

	return ret
}

/**
Run mode: refer to utils/const/enum
*/
func GetRunModeAndName(assert string) (string, string) {
	ext := path.Ext(assert)
	name := strings.Replace(commonUtils.Base(assert), ext, "", -1)

	var mode string
	if ext == ".suite" {
		mode = constant.RunModeSuite.String()
	} else {
		mode = constant.RunModeScript.String()
	}

	return mode, name
}

func GetLogFolder(mode string, name string, date string) string {
	return fmt.Sprintf("%s-%s-%s", mode, name, date)
}

func GetAllScriptsInDir(filePth string, files *[]string) error {
	if !fileUtils.IsDir(filePth) { // not dir
		pass := CheckFileIsScript(filePth)

		if pass {
			id, _, _ := zentaoUtils.GetCaseInfo(filePth)

			if id > 0 {
				*files = append(*files, filePth)
			}
		}

		return nil
	}

	filePth = fileUtils.AbosutePath(filePth)
	sep := string(os.PathSeparator)

	dir, err := ioutil.ReadDir(filePth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if fi.IsDir() { // 目录, 递归遍历
			GetAllScriptsInDir(filePth+name+sep, files)
		} else {
			path := filePth + name
			if CheckFileIsScript(path) {
				*files = append(*files, path)
			}
		}
	}

	return nil
}

func GetScriptByIdsInDir(dirPth string, idMap map[int]string, files *[]string) error {
	dirPth = fileUtils.AbosutePath(dirPth)

	sep := string(os.PathSeparator)

	name := path.Base(dirPth)
	if strings.Index(name, ".") == 0 || name == "bin" || name == "release" || name == "logs" || name == "xdoc" {
		return nil
	}

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if fi.IsDir() { // 目录, 递归遍历
			GetScriptByIdsInDir(dirPth+name+sep, idMap, files)
		} else {
			regx := langUtils.GetSupportLangageRegx()
			pass, _ := regexp.MatchString("^*.\\."+regx+"$", name)

			if !pass {
				continue
			}

			path := dirPth + name
			if CheckFileIsScript(path) {
				id, _, _ := zentaoUtils.GetCaseInfo(path)

				if id > 0 {
					_, ok := idMap[id]

					if ok {
						*files = append(*files, path)
					}
				}
			}
		}
	}

	return nil
}

func GetCaseIdsInSuiteFile(name string, fileIdMap *map[int]string) {
	content := fileUtils.ReadFile(name)

	for _, line := range strings.Split(content, "\n") {
		idStr := strings.TrimSpace(line)
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err == nil {
			(*fileIdMap)[id] = ""
		}
	}
}

func GetFailedCasesDirectlyFromTestResult(resultFile string, cases *[]string) {
	extName := path.Ext(resultFile)

	if extName == "."+constant.ExtNameResult {
		resultFile = strings.Replace(resultFile, extName, "."+constant.ExtNameJson, -1)
	}

	content := fileUtils.ReadFile(resultFile)

	var report model.TestReport
	json.Unmarshal([]byte(content), &report)

	for _, cs := range report.Cases {
		*cases = append(*cases, cs.Path)
	}
}

func CheckFileIsScript(path string) bool {
	content := fileUtils.ReadFile(path)

	pass, _ := regexp.MatchString("<<<TC", content)

	return pass
}
