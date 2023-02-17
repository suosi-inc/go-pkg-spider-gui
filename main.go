// 由res2go IDE插件自动生成。
package main

import (
	"github.com/suosi-inc/go-pkg-spider/spider-gui/gui"
	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
)

func main() {
	vcl.Application.SetScaled(true)
	vcl.Application.SetTitle("spider-gui")
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&gui.FormMain)
	vcl.Application.Run()
}
