package main

import (
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	out         = flag.String("out", "vendor", "Directory to generate files in")
	import_path = flag.String("import_path", "", "Specify import path/prefix for nested types")
)

func writeCode(fullname string, code string) error {
	nameComponents := strings.Split(fullname, "/")
	pkgDir := filepath.Join(*out, nameComponents[0])
	if _, err := os.Stat(pkgDir); os.IsNotExist(err) {
		err = os.MkdirAll(pkgDir, os.ModeDir|os.FileMode(0775))
		if err != nil {
			return err
		}
	}
	filename := filepath.Join(pkgDir, nameComponents[1]+".go")

	res, err := format.Source([]byte(code))
	if err != nil {
		return fmt.Errorf("Error formatting generated code: %+v", err)
	}

	return ioutil.WriteFile(filename, res, os.FileMode(0664))
}

func main() {
	flag.Parse()
	if _, err := os.Stat(*out); os.IsNotExist(err) {
		err = os.MkdirAll(*out, os.ModeDir|os.FileMode(0775))
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	if flag.NArg() < 2 {
		fmt.Println("USAGE: gengo [-out=] [-import_path=] msg|srv <NAME> [<FILE>]")
		os.Exit(-1)
	}

	rosPkgPath := os.Getenv("ROS_PACKAGE_PATH")

	context, err := NewMsgContext(strings.Split(rosPkgPath, ":"))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	mode := flag.Arg(0)
	fullname := flag.Arg(1)

	fmt.Printf("Generating %v...", fullname)

	if mode == "msg" {
		var spec *MsgSpec
		var err error
		if flag.NArg() == 2 {
			spec, err = context.LoadMsg(fullname)
		} else {
			spec, err = context.LoadMsgFromFile(flag.Arg(2), fullname)
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
		if flag.NArg() == 2 {
			spec, err = context.LoadSrv(fullname)
		} else {
			spec, err = context.LoadSrvFromFile(flag.Arg(2), fullname)
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
		fmt.Println("USAGE: gengo <MSG>")
		os.Exit(-1)
	}
	fmt.Println("Done")
}
