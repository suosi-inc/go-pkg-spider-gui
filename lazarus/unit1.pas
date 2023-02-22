unit Unit1;

{$mode objfpc}{$H+}

interface

uses
  Classes, SysUtils, mysql80conn, Forms, Controls, Graphics, Dialogs, ExtCtrls,
  ComCtrls, StdCtrls, Menus, Buttons, Spin, EditBtn, ValEdit, Grids,
  ListFilterEdit;

type

  { TFormMain }

  TFormMain = class(TForm)
    BtnLinkOpen: TSpeedButton;
    BtnLinkRequest: TButton;
    BtnRequest: TButton;
    BtnRequestDefault: TButton;
    BtnRequestTipHeader: TSpeedButton;
    BtnRequestTipProxy: TSpeedButton;
    BtnToolDomainRequest: TButton;
    BtnLinkSearch: TButton;
    CheckLinkStrictDomain: TCheckBox;
    CheckRequestCharset: TCheckBox;
    CheckRequestClean: TCheckBox;
    CheckRequestRedirect: TCheckBox;
    CheckRequestType: TCheckBox;
    EditLinkSearch: TEdit;
    EditLinkRetry: TSpinEdit;
    EditLinkTimeout: TEdit;
    EditLinkUrl: TEdit;
    EditRequestLength: TEdit;
    EditRequestProxy: TEdit;
    EditRequestRedirect: TSpinEdit;
    EditRequestTimeout: TEdit;
    EditRequestType: TEdit;
    EditRequestUa: TEdit;
    EditToolDomain: TEdit;
    EditToolDomainResult: TEdit;
    EditRequestUrl: TEdit;
    PanelLinkRule: TPanel;
    PanelRequestBox: TPanel;
    LabelRequestHeader: TLabel;
    LabelRequestLength: TLabel;
    LabelRequestProxy: TLabel;
    LabelRequestRedirect: TLabel;
    LabelRequestTimeout: TLabel;
    LabelRequestType: TLabel;
    LabelRequestUa: TLabel;
    MemoRequestHeader: TMemo;
    PanelLinkRequest: TPanel;
    ImageListIcon: TImageList;
    ImageListToolBar: TImageList;
    LabelLinkRetry: TLabel;
    LabelLinkTimeout: TLabel;
    LabelLinkUrl: TLabel;
    LabelToolDomain: TLabel;
    LabelRequestUrl: TLabel;
    MemoRequest: TMemo;
    MemoDebug: TMemo;
    MenuDebugCopy: TMenuItem;
    MenuDebugClear: TMenuItem;
    PageControl: TPageControl;
    PageControlLink: TPageControl;
    PanelToolDomain: TPanel;
    PanelDomain: TPanel;
    PanelLink: TPanel;
    PanelNews: TPanel;
    PanelTool: TPanel;
    PanelRequestView: TPanel;
    PanelRequestController: TPanel;
    PanelDebug: TPanel;
    PopupMenuDebug: TPopupMenu;
    BtnRequestOpen: TSpeedButton;
    SplitterRequest: TSplitter;
    SplitterDebug: TSplitter;
    GridLinkContent: TStringGrid;
    GridLinkList: TStringGrid;
    GridLinkUnknow: TStringGrid;
    GridLinkNone: TStringGrid;
    GridLinkFilter: TStringGrid;
    GridLinkDomain: TStringGrid;
    TabSheetLinkFilter: TTabSheet;
    TabSheetLinkDomain: TTabSheet;
    TabSheetLinkContent: TTabSheet;
    TabSheetLinkList: TTabSheet;
    TabSheetLinkUnknow: TTabSheet;
    TabSheetLinkNone: TTabSheet;
    TabSheetTool: TTabSheet;
    TabSheetNews: TTabSheet;
    TabSheetLink: TTabSheet;
    TabSheetDomain: TTabSheet;
    TabSheetRequest: TTabSheet;
    ToolBar: TToolBar;
    ToolBtnSplit1: TToolButton;
    ToolBtnRequest: TToolButton;
    ToolBtnSplit2: TToolButton;
    ToolBtnDebug: TToolButton;
    ToolBtnDomain: TToolButton;
    ToolBtnSplit3: TToolButton;
    ToolBtnLink: TToolButton;
    ToolBtnSplit4: TToolButton;
    ToolBtnContent: TToolButton;
    ToolBtnSplit5: TToolButton;
    ToolBtnTool: TToolButton;
    ToolBtnSplit6: TToolButton;
    procedure BtnLinkOpenClick(Sender: TObject);
    procedure BtnLinkRequestClick(Sender: TObject);
    procedure BtnLinkSearchClick(Sender: TObject);
    procedure BtnRequestClick(Sender: TObject);
    procedure BtnRequestDefaultClick(Sender: TObject);
    procedure BtnRequestExampleClick(Sender: TObject);
    procedure BtnRequestOpenClick(Sender: TObject);
    procedure BtnRequestTipHeaderClick(Sender: TObject);
    procedure BtnRequestTipProxyClick(Sender: TObject);
    procedure BtnToolDomainRequestClick(Sender: TObject);
    procedure CheckRequestRedirectChange(Sender: TObject);
    procedure CheckRequestTypeChange(Sender: TObject);
    procedure FormCreate(Sender: TObject);
    procedure MenuDebugClearClick(Sender: TObject);
    procedure MenuDebugCopyClick(Sender: TObject);
    procedure ToolBtnContentClick(Sender: TObject);
    procedure ToolBtnDebugClick(Sender: TObject);
    procedure ToolBtnDomainClick(Sender: TObject);
    procedure ToolBtnLinkClick(Sender: TObject);
    procedure ToolBtnRequestClick(Sender: TObject);
    procedure ToolBtnToolClick(Sender: TObject);
    procedure RemoveToolBtnDown(Sender: TObject);
  private

  public

  end;

var
  FormMain: TFormMain;

implementation

{$R *.lfm}

{ TFormMain }

procedure TFormMain.ToolBtnDebugClick(Sender: TObject);
begin
  if PanelDebug.Visible = False then
     begin
        SplitterDebug.Visible:=True;
        PanelDebug.Visible:=True;
        PageControl.AnchorSideBottom.Control:=SplitterDebug
     end
     else
     begin
        SplitterDebug.Visible:=False;
        PanelDebug.Visible:=False;
        PageControl.AnchorSideBottom.Control:=FormMain;
     end;
end;

procedure TFormMain.ToolBtnRequestClick(Sender: TObject);
begin
     RemoveToolBtnDown(Sender);
     ToolBtnRequest.Down:=True;
     PageControl.ActivePageIndex:=0;
end;

procedure TFormMain.ToolBtnDomainClick(Sender: TObject);
begin
    RemoveToolBtnDown(Sender);
    ToolBtnDomain.Down:=True;
    PageControl.ActivePageIndex:=1;
end;

procedure TFormMain.ToolBtnLinkClick(Sender: TObject);
begin
    RemoveToolBtnDown(Sender);
    ToolBtnLink.Down:=True;
    PageControl.ActivePageIndex:=2;
end;

procedure TFormMain.ToolBtnContentClick(Sender: TObject);
begin
    RemoveToolBtnDown(Sender);
    ToolBtnContent.Down:=True;
    PageControl.ActivePageIndex:=3;
end;

procedure TFormMain.ToolBtnToolClick(Sender: TObject);
begin
   RemoveToolBtnDown(Sender);
   ToolBtnTool.Down:=True;
   PageControl.ActivePageIndex:=4;
end;

procedure TFormMain.RemoveToolBtnDown(Sender: TObject);
begin
   ToolBtnRequest.Down:=False;
   ToolBtnDomain.Down:=False;
   ToolBtnLink.Down:=False;
   ToolBtnContent.Down:=False;
   ToolBtnTool.Down:=False;
end;

procedure TFormMain.MenuDebugCopyClick(Sender: TObject);
begin
  MemoDebug.CopyToClipboard;
end;


procedure TFormMain.MenuDebugClearClick(Sender: TObject);
begin
  MemoDebug.Text:=''
end;

procedure TFormMain.CheckRequestTypeChange(Sender: TObject);
begin
    if CheckRequestType.Checked then
       EditRequestType.Enabled:=False
    else
       EditRequestType.Enabled:=True;
end;

procedure TFormMain.FormCreate(Sender: TObject);
begin

end;

procedure TFormMain.CheckRequestRedirectChange(Sender: TObject);
begin
    if CheckRequestRedirect.Checked then
    begin
         EditRequestRedirect.Enabled:=False;
         EditRequestRedirect.EditorEnabled:=False
    end
    else
    begin
         EditRequestRedirect.Enabled:=True;
         EditRequestRedirect.EditorEnabled:=True;
    end;
end;

procedure TFormMain.BtnRequestDefaultClick(Sender: TObject);
begin

end;

procedure TFormMain.BtnRequestClick(Sender: TObject);
begin

end;

procedure TFormMain.BtnLinkRequestClick(Sender: TObject);
begin

end;

procedure TFormMain.BtnLinkSearchClick(Sender: TObject);
begin

end;

procedure TFormMain.BtnLinkOpenClick(Sender: TObject);
begin

end;

procedure TFormMain.BtnRequestExampleClick(Sender: TObject);
begin

end;

procedure TFormMain.BtnRequestOpenClick(Sender: TObject);
begin

end;

procedure TFormMain.BtnRequestTipHeaderClick(Sender: TObject);
begin

end;

procedure TFormMain.BtnRequestTipProxyClick(Sender: TObject);
begin

end;

procedure TFormMain.BtnToolDomainRequestClick(Sender: TObject);
begin

end;

end.

