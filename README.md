# go-pkg-spider-gui

## 构建

```
# Windows 下隐藏 cmd 窗口，并且将 liblcl 打包到可执行文件内，
go build -ldflags "-H windowsgui" -tags tempdll -o bin/win64

# 缩小构建大小：-s 去掉符号信息, -w 去掉 DWARF 调试信息
go build -ldflags "-s -w -H windowsgui" -tags tempdll -o bin/win64

# MacOS
go build -o bin/macos64
```