package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/service/script"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"strconv"
	"strings"
)

var CurrRun string
var CurrResult string

var runViews []string

func showRun(g *gocui.Gui, v *gocui.View) error {
	DestoryContentPanel()
	HighlightTab(v.Name(), tabs)

	h := vari.MainViewHeight / 2
	maxX, _ := g.Size()

	panelResultList := NewPanelWidget("panelResultList", constant.LeftWidth, 2, 50, h, "")
	ViewMap["testing"] = append(ViewMap["testing"], panelResultList.Name())
	runViews = append(runViews, panelResultList.Name())

	panelCaseList := NewPanelWidget("panelCaseList", constant.LeftWidth, h+2, 50, vari.MainViewHeight-h, "")
	ViewMap["testing"] = append(ViewMap["testing"], panelCaseList.Name())
	runViews = append(runViews, panelCaseList.Name())

	panelCaseResult := NewPanelWidget("panelCaseResult", constant.LeftWidth+50, 2,
		maxX-constant.LeftWidth-51, vari.MainViewHeight, "")
	ViewMap["testing"] = append(ViewMap["testing"], panelCaseResult.Name())
	runViews = append(runViews, panelCaseResult.Name())

	for idx, v := range runViews {
		if idx < 3 {
			setViewScroll(v)
		}

		if idx < 2 {
			setViewLineHighlight(v)
		}
	}

	setViewLineSelected("panelResultList", selectResultEvent)
	setViewLineSelected("panelCaseList", selectCaseEvent)

	results := scriptService.LoadTestResults(CurrAsset)
	fmt.Fprintln(panelResultList, strings.Join(results, "\n"))

	return nil
}

func init() {

}

func selectResultEvent(g *gocui.Gui, v *gocui.View) error {
	clearPanelCaseResult()

	v.Highlight = true

	line, _ := GetSelectedLine(v, ".*")
	CurrResult = line
	//content := scriptService.GetTestResultForDisplay(CurrAsset, line)

	content := make([]string, 0)
	report := testingService.GetTestTestReportForSubmit(CurrAsset, line)
	for _, cs := range report.Cases {
		id := cs.Id
		title := cs.Title
		result := cs.Status

		str := fmt.Sprintf("%d-%s: %s", id, title, result)
		content = append(content, str)
	}

	panelCaseList, _ := g.View("panelCaseList")
	panelCaseList.Clear()
	fmt.Fprintln(panelCaseList, strings.Join(content, "\n"))

	maxX, _ := g.Size()
	uploadButton := NewButtonWidgetAutoWidth("uploadButton", maxX-35, 0, "[Upload Result]", toUploadResult)
	uploadButton.Frame = false
	runViews = append(runViews, uploadButton.Name())

	return nil
}

func selectCaseEvent(g *gocui.Gui, v *gocui.View) error {
	v.Highlight = true

	caseLine, _ := GetSelectedLine(v, ".*")
	caseIdStr := strings.Split(caseLine, "-")[0]
	caseId, _ := strconv.Atoi(caseIdStr)

	content := make([]string, 0)
	report := testingService.GetTestTestReportForSubmit(CurrAsset, CurrResult)
	for _, cs := range report.Cases {
		if cs.Id == caseId {
			for _, step := range cs.Steps {
				content = append(content, testingService.GetStepText(step))
				content = append(content, "")
			}
		}
	}

	panelCaseResult, _ := g.View("panelCaseResult")
	panelCaseResult.Clear()
	fmt.Fprintln(panelCaseResult, strings.Join(content, "\n"))

	// show submit bug button
	maxX, _ := g.Size()
	bugButton := NewButtonWidgetAutoWidth("bugButton", maxX-18, 0, "[Report Bug]", toReportBug)
	bugButton.Frame = false
	runViews = append(runViews, bugButton.Name())

	return nil
}

func clearPanelCaseResult() {
	panelCaseResult, _ := vari.Cui.View("panelCaseResult")
	if panelCaseResult != nil {
		panelCaseResult.Clear()
	}
	vari.Cui.DeleteView("bugButton")
}

func toUploadResult(g *gocui.Gui, v *gocui.View) error {
	zentaoService.SubmitResult(CurrAsset, CurrResult)

	return nil
}

func toReportBug(g *gocui.Gui, v *gocui.View) error {
	InitReportBugPage()

	return nil
}

func DestoryRunPanel() {
	for _, v := range runViews {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}
}
