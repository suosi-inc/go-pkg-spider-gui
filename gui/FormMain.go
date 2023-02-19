// 由res2go IDE插件自动生成，不要编辑。
package gui

import (
	_ "embed"
	"github.com/ying32/govcl/vcl"
)

type TFormMain struct {
	*vcl.TForm
	ToolBar                *vcl.TToolBar
	ToolBtnSplit1          *vcl.TToolButton
	ToolBtnRequest         *vcl.TToolButton
	ToolBtnSplit2          *vcl.TToolButton
	ToolBtnDomain          *vcl.TToolButton
	ToolBtnSplit3          *vcl.TToolButton
	ToolBtnLink            *vcl.TToolButton
	ToolBtnSplit4          *vcl.TToolButton
	ToolBtnContent         *vcl.TToolButton
	ToolBtnSplit5          *vcl.TToolButton
	ToolBtnTool            *vcl.TToolButton
	ToolBtnSplit6          *vcl.TToolButton
	ToolBtnDebug           *vcl.TToolButton
	PageControl            *vcl.TPageControl
	TabSheetRequest        *vcl.TTabSheet
	PanelRequestView       *vcl.TPanel
	MemoRequest            *vcl.TMemo
	PanelRequestController *vcl.TPanel
	LabelRequestUrl        *vcl.TLabel
	EditRequestUrl         *vcl.TEdit
	BtnRequest             *vcl.TButton
	BtnRequestLink         *vcl.TSpeedButton
	GroupBoxRequest        *vcl.TGroupBox
	LabelRequestUa         *vcl.TLabel
	EditRequestUa          *vcl.TEdit
	LabelRequestLength     *vcl.TLabel
	LabelRequestType       *vcl.TLabel
	EditRequestType        *vcl.TEdit
	LabelRequestRedirect   *vcl.TLabel
	CheckRequestCharset    *vcl.TCheckBox
	CheckRequestType       *vcl.TCheckBox
	MemoRequestHeader      *vcl.TMemo
	LabelRequestHeader     *vcl.TLabel
	CheckRequestRedirect   *vcl.TCheckBox
	LabelRequestProxy      *vcl.TLabel
	EditRequestProxy       *vcl.TEdit
	BtnRequestDefault      *vcl.TButton
	EditRequestRedirect    *vcl.TSpinEdit
	EditRequestLength      *vcl.TEdit
	CheckRequestClean      *vcl.TCheckBox
	LabelRequestTimeout    *vcl.TLabel
	EditRequestTimeout     *vcl.TEdit
	BtnRequestTipProxy     *vcl.TSpeedButton
	BtnRequestTipHeader    *vcl.TSpeedButton
	SplitterRequest        *vcl.TSplitter
	TabSheetDomain         *vcl.TTabSheet
	TabSheetLink           *vcl.TTabSheet
	TabSheetContent        *vcl.TTabSheet
	TabSheetTool           *vcl.TTabSheet
	PanelTool              *vcl.TPanel
	SplitterDebug          *vcl.TSplitter
	PanelDebug             *vcl.TPanel
	MemoDebug              *vcl.TMemo
	ImageListToolBar       *vcl.TImageList
	PopupMenuDebug         *vcl.TPopupMenu
	MenuDebugCopy          *vcl.TMenuItem
	MenuDebugClear         *vcl.TMenuItem
	ImageListIcon          *vcl.TImageList
	PanelContent           *vcl.TPanel
	PanelLink              *vcl.TPanel
	PanelDomain            *vcl.TPanel

	//::private::
	TFormMainFields
}

var FormMain *TFormMain

// vcl.Application.CreateForm(&FormMain)

func NewFormMain(owner vcl.IComponent) (root *TFormMain) {
	vcl.CreateResForm(owner, &root)
	return
}

//go:embed resources/FormMain.gfm
var formMainBytes []byte

// 注册Form资源
var _ = vcl.RegisterFormResource(FormMain, &formMainBytes)
