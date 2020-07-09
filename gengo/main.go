package main

import (
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/edwinhayes/rosgo/libgengo"
)

var (
	out        = flag.String("out", "vendor", "Directory to generate files in")
	importPath = flag.String("import_path", "", "Specify import path/prefix for nested types")
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
		fmt.Println("USAGE: gengo [-out=] [-import_path=] msg|srv|action <NAME> [<FILE>]")
		os.Exit(-1)
	}

	rosPkgPath := os.Getenv("ROS_PACKAGE_PATH")

	context, err := libgengo.NewMsgContext(strings.Split(rosPkgPath, ":"))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	mode := flag.Arg(0)
	fullname := flag.Arg(1)

	fmt.Printf("Generating %v...", fullname)

	if mode == "msg" {
		var spec *libgengo.MsgSpec
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
		code, err = libgengo.GenerateMessage(context, spec, false)
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
		var spec *libgengo.SrvSpec
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
		srvCode, reqCode, resCode, err := libgengo.GenerateService(context, spec)
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
	} else if mode == "action" {
		var spec *libgengo.ActionSpec
		var err error

		if len(os.Args) == 3 {
			spec, err = context.LoadAction(fullname)
		} else {
			spec, err = context.LoadActionFromFile(os.Args[3], fullname)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		actionCode, codeMap, err := libgengo.GenerateAction(context, spec)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		err = writeCode(fullname, actionCode)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		for name, code := range codeMap {
			err = writeCode(name, code)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		}

	} else {
		fmt.Println("USAGE: gengo <MSG>")
		os.Exit(-1)
	}
	fmt.Println("Done")
}
