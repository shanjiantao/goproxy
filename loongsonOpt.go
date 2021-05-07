package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/goproxyio/goproxy/v2/proxy"
	"golang.org/x/mod/module"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type loongsonOpt struct {}



func (*loongsonOpt)NewContext(r *http.Request) (context.Context, error) {
	return context.Background(),nil
}

func (*loongsonOpt)List(ctx context.Context, mpath string) (proxy.File, error) {
	escMod, err := module.EscapePath(mpath)
	if err != nil {
		return nil, err
	}
	file := filepath.Join(downloadRoot, escMod, "@v", "list")
	if info, err := os.Stat(file); err == nil && time.Since(info.ModTime()) < listExpire {
		return os.Open(file)
	}
	var list struct {
		Path     string
		Versions []string
	}
	if err := goJSON(&list, "go", "list", "-m", "-json", "-versions", mpath+"@latest"); err != nil {
		return nil, err
	}
	if list.Path != mpath {
		return nil, fmt.Errorf("go list -m: asked for %s but got %s", mpath, list.Path)
	}
	data := []byte(strings.Join(list.Versions, "\n") + "\n")
	if len(data) == 1 {
		data = nil
	}
	err = os.MkdirAll(path.Dir(file), os.ModePerm)
	if err != nil {
		log.Printf("make cache dir failed, err: %v.", err)
		return nil, err
	}
	if err := ioutil.WriteFile(file, data, 0666); err != nil {
		return nil, err
	}

	return os.Open(file)
}

func (*loongsonOpt)Latest(ctx context.Context, path string) (proxy.File, error) {
	d, err := download(module.Version{Path: path, Version: "latest"})
	if err != nil {
		return nil, err
	}
	return os.Open(d.Info)
}

func (*loongsonOpt)	Info(ctx context.Context, m module.Version) (proxy.File, error) {
	strs:= strings.Split(m.Version,"-")
	if len(strs)!= 3 {
		log.Printf("error format version\n")
		return nil,errors.New("error format version\n")
	}
	t,_ := time.Parse("20060102150405",strs[1])

	return proxy.NewInfo(m,t),nil
}

func (*loongsonOpt) GoMod(ctx context.Context, m module.Version) (proxy.File, error) {
	return proxy.NewGoMod(m),nil
}

func (*loongsonOpt)	Zip(ctx context.Context, m module.Version) (proxy.File, error) {
	//TODO 判断存在


	return proxy.NewZip(filepath.Join(cacheDir,"loongson",m.Path),m),nil
}