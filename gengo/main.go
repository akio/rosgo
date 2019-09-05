package main

import (
	"fmt"
	"github.com/edwinhayes/rosgo/libgengo"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func writeCode(fullname string, code string) error {
	nameComponents := strings.Split(fullname, "/")
	pkgDir := filepath.Join("vendor", nameComponents[0])
	if _, err := os.Stat(pkgDir); os.IsNotExist(err) {
		err = os.MkdirAll(pkgDir, os.ModeDir|os.FileMode(0775))
		if err != nil {
			return err
		}
	}
	filename := filepath.Join(pkgDir, nameComponents[1]+".go")

	return ioutil.WriteFile(filename, []byte(code), os.FileMode(0664))
}

func main() {
	if _, err := os.Stat("vendor"); os.IsNotExist(err) {
		err = os.Mkdir("vendor", os.ModeDir|os.FileMode(0775))
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	if len(os.Args) < 3 {
		fmt.Println("USAGE: gengo msg|srv <NAME> [<FILE>]")
		os.Exit(-1)
	}

	rosPkgPath := os.Getenv("ROS_PACKAGE_PATH")

	context, err := NewMsgContext(strings.Split(rosPkgPath, ":"))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	mode := os.Args[1]
	fullname := os.Args[2]

	fmt.Printf("Generating %v...", fullname)

	if mode == "msg" {
		var spec *MsgSpec
		var err error
		if len(os.Args) == 3 {
			spec, err = context.LoadMsg(fullname)
		} else {
			spec, err = context.LoadMsgFromFile(os.Args[3], fullname)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		var code string
		code, err = GenerateMessage(context, spec)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		err = writeCode(fullname, code)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	} else if mode == "srv" {
		var spec *SrvSpec
		var err error
		if len(os.Args) == 3 {
			spec, err = context.LoadSrv(fullname)
		} else {
			spec, err = context.LoadSrvFromFile(os.Args[3], fullname)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		srvCode, reqCode, resCode, err := GenerateService(context, spec)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		err = writeCode(fullname, srvCode)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		err = writeCode(spec.Request.FullName, reqCode)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		err = writeCode(spec.Response.FullName, resCode)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	} else {
		fmt.Println("USAGE: genmsg <MSG>")
		os.Exit(-1)
	}
	fmt.Println("Done")
}
