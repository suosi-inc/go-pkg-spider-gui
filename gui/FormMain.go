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
	BtnRequestOpen         *vcl.TSpeedButton
	PanelRequestBox        *vcl.TPanel
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
	PanelDomainContent     *vcl.TPanel
	PanelDomainSubdomain   *vcl.TPanel
	GridDomainSubdomain    *vcl.TStringGrid
	SplitterDomain         *vcl.TSplitter
	PanelDomain            *vcl.TPanel
	PanelDomainRequest     *vcl.TPanel
	LabelDomain            *vcl.TLabel
	EditDomain             *vcl.TEdit
	BtnDomainRequest       *vcl.TButton
	LabelDomainRetry       *vcl.TLabel
	EditDomainRetry        *vcl.TSpinEdit
	LabelDomainTimeout     *vcl.TLabel
	EditDomainTimeout      *vcl.TEdit
	CheckDomainSubdomain   *vcl.TCheckBox
	GridDomainData         *vcl.TStringGrid
	TabSheetLink           *vcl.TTabSheet
	PanelLink              *vcl.TPanel
	PageControlLink        *vcl.TPageControl
	TabSheetLinkContent    *vcl.TTabSheet
	GridLinkContent        *vcl.TStringGrid
	TabSheetLinkList       *vcl.TTabSheet
	GridLinkList           *vcl.TStringGrid
	TabSheetLinkUnknown    *vcl.TTabSheet
	GridLinkUnknown        *vcl.TStringGrid
	TabSheetLinkNone       *vcl.TTabSheet
	GridLinkNone           *vcl.TStringGrid
	TabSheetLinkFilter     *vcl.TTabSheet
	GridLinkFilter         *vcl.TStringGrid
	TabSheetLinkDomain     *vcl.TTabSheet
	GridLinkDomain         *vcl.TStringGrid
	PanelLinkRequest       *vcl.TPanel
	LabelLinkTimeout       *vcl.TLabel
	EditLinkTimeout        *vcl.TEdit
	LabelLinkRetry         *vcl.TLabel
	CheckLinkStrictDomain  *vcl.TCheckBox
	LabelLinkUrl           *vcl.TLabel
	EditLinkUrl            *vcl.TEdit
	BtnLinkRequest         *vcl.TButton
	BtnLinkOpen            *vcl.TSpeedButton
	EditLinkRetry          *vcl.TSpinEdit
	EditLinkSearch         *vcl.TEdit
	PanelLinkRule          *vcl.TPanel
	EditLinkRuleDomain1    *vcl.TEdit
	EditLinkRuleDomain2    *vcl.TEdit
	ComboLinkRuleType1     *vcl.TComboBox
	LabelLinkRule1         *vcl.TLabel
	LabelLinkRule2         *vcl.TLabel
	EditLinkRuleContent1   *vcl.TEdit
	CheckLinkRule1         *vcl.TCheckBox
	ComboLinkRuleType2     *vcl.TComboBox
	EditLinkRuleContent2   *vcl.TEdit
	CheckLinkRule2         *vcl.TCheckBox
	BtnLinkSearch          *vcl.TButton
	TabSheetNews           *vcl.TTabSheet
	PanelNews              *vcl.TPanel
	PanelNewsRequest       *vcl.TPanel
	EditNewsTitle          *vcl.TEdit
	LabelNewsTitle         *vcl.TLabel
	EditNewsRetry          *vcl.TSpinEdit
	LabelNewsRetry         *vcl.TLabel
	EditNewsTimeout        *vcl.TEdit
	LabelNewsTimeout       *vcl.TLabel
	BtnNewsOpen            *vcl.TSpeedButton
	BtnNewsRequest         *vcl.TButton
	EditNewsUrl            *vcl.TEdit
	LabelNewsUrl           *vcl.TLabel
	RadioNewsContentType   *vcl.TRadioGroup
	LabelNewsContentType   *vcl.TLabel
	GridNewsInfo           *vcl.TStringGrid
	PanelNewsContent       *vcl.TPanel
	MemoNewsContent        *vcl.TMemo
	EditNewsResultTitle    *vcl.TEdit
	EditNewsResultTime     *vcl.TEdit
	SplitterNews           *vcl.TSplitter
	TabSheetTool           *vcl.TTabSheet
	PanelTool              *vcl.TPanel
	PanelToolDomain        *vcl.TPanel
	LabelToolDomain        *vcl.TLabel
	EditToolDomain         *vcl.TEdit
	BtnToolDomainRequest   *vcl.TButton
	EditToolDomainResult   *vcl.TEdit
	PanelToolLang          *vcl.TPanel
	MemoToolLang           *vcl.TMemo
	EditToolLang           *vcl.TEdit
	BtnToolLang            *vcl.TButton
	LabelToolLangTip       *vcl.TLabel
	SplitterDebug          *vcl.TSplitter
	PanelDebug             *vcl.TPanel
	MemoDebug              *vcl.TMemo
	ImageListToolBar       *vcl.TImageList
	PopupMenuDebug         *vcl.TPopupMenu
	MenuDebugCopy          *vcl.TMenuItem
	MenuDebugClear         *vcl.TMenuItem
	ImageListIcon          *vcl.TImageList
	ProgressBarDomain      *vcl.TProgressBar

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
