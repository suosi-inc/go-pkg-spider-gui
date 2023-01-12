unit Unit1;

{$mode objfpc}{$H+}

interface

uses
  Classes, SysUtils, Forms, Controls, Graphics, Dialogs, ExtCtrls, ComCtrls,
  StdCtrls, Menus, Buttons;

type

  { TFormMain }

  TFormMain = class(TForm)
    BtnRequest: TButton;
    BtnRequestDefault: TButton;
    BtnRequestExample: TButton;
    CheckRequestCharset: TCheckBox;
    CheckRequestType: TCheckBox;
    CheckBoxRequestRedirect: TCheckBox;
    EditRequestProxy: TEdit;
    EditRequestRedirect: TEdit;
    EditRequestType: TEdit;
    EditRequestLength: TEdit;
    EditRequestUa: TEdit;
    EditRequestUrl: TEdit;
    GroupBoxRequest: TGroupBox;
    ImageListToolBar: TImageList;
    LabelRequestProxy: TLabel;
    LabelRequestType: TLabel;
    LabelRequestRedirect: TLabel;
    LabelRequestHeader: TLabel;
    LabelRequestUa: TLabel;
    LabelRequestLength: TLabel;
    LabelRequestUrl: TLabel;
    MemoRequestHeader: TMemo;
    MemoRequest: TMemo;
    MemoDebug: TMemo;
    MenuDebugCopy: TMenuItem;
    MenuDebugClear: TMenuItem;
    PageControl: TPageControl;
    PanelRequestView: TPanel;
    PanelRequestController: TPanel;
    PanelDebug: TPanel;
    PopupMenuDebug: TPopupMenu;
    SplitterRequest: TSplitter;
    SplitterDebug: TSplitter;
    TabSheetContent: TTabSheet;
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
    procedure GroupBox1Click(Sender: TObject);
    procedure LabelRequestProxyClick(Sender: TObject);
    procedure LabelRequestLengthClick(Sender: TObject);
    procedure MenuDebugClearClick(Sender: TObject);
    procedure MenuDebugCopyClick(Sender: TObject);
    procedure MenuToolLanguageClick(Sender: TObject);
    procedure ToolBtnDebugClick(Sender: TObject);
  private

  public

  end;

var
  FormMain: TFormMain;

implementation

{$R *.lfm}

{ TFormMain }

procedure TFormMain.MenuToolLanguageClick(Sender: TObject);
begin

end;

procedure TFormMain.ToolBtnDebugClick(Sender: TObject);
begin
  if ToolBtnDebug.Down then
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

procedure TFormMain.MenuDebugCopyClick(Sender: TObject);
begin
  MemoDebug.CopyToClipboard;
end;

procedure TFormMain.MenuDebugClearClick(Sender: TObject);
begin
  MemoDebug.Text:=''
end;

procedure TFormMain.GroupBox1Click(Sender: TObject);
begin

end;

procedure TFormMain.LabelRequestProxyClick(Sender: TObject);
begin

end;

procedure TFormMain.LabelRequestLengthClick(Sender: TObject);
begin

end;

end.

