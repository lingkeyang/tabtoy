package main

import (
	"flag"
	"github.com/davyxu/tabtoy/v3/compiler"
	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/gen/gosrc"
	"github.com/davyxu/tabtoy/v3/gen/jsondata"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"os"
)

type V3GenEntry struct {
	f    gen.GenFunc
	name *string
}

// v3新增
var (
	paramIndexFile = flag.String("index", "", "input multi-files configs")

	v3GenList = []V3GenEntry{
		{gosrc.Generate, paramGoOut},
		{jsondata.Generate, paramJsonOut},
	}
)

func selectFileLoader(globals *model.Globals, para bool) {
	globals.IndexGetter = new(helper.SyncFileLoader)
	globals.Para = para
	if globals.Para {
		// 缓冲文件
		asyncLoader := helper.NewAsyncFileLoader()

		for _, pragma := range globals.IndexList {
			asyncLoader.AddFile(pragma.TableFileName)
		}

		asyncLoader.Commit()

		globals.TableGetter = asyncLoader
	} else {
		globals.TableGetter = globals.IndexGetter
	}
}

func GenFile(globals *model.Globals) error {
	for _, entry := range v3GenList {

		if *entry.name == "" {
			continue
		}

		filename := *entry.name

		if data, err := entry.f(globals); err != nil {
			return err
		} else {

			report.Log.Infoln(filename)

			err = helper.WriteFile(filename, data)

			if err != nil {
				return err
			}

		}
	}

	return nil
}

func V3Entry() {
	globals := model.NewGlobals()
	globals.Version = Version_v3

	globals.IndexFile = *paramIndexFile
	globals.PackageName = *paramPackageName
	globals.CombineStructName = *paramCombineStructName

	selectFileLoader(globals, *paramPara)

	var err error

	err = compiler.Compile(globals)

	if err != nil {
		goto Exit
	}

	report.Log.Debugln("Generate files...")
	err = GenFile(globals)
	if err != nil {
		goto Exit
	}

	return
Exit:
	report.Log.Errorln(err)
	os.Exit(1)
}
