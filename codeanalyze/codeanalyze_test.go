package codeanalyze

import (
	"bytes"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCodeAnalyze(t *testing.T) {
	goParser := GoParser{}
	pyParser := PythonParser{}

	Convey("Given some codebase ", t, func() {
		path := "/Users/meetai/chatcode"

		Convey("When the path is go", func() {
			module := goParser.ParseDirectory(path)
			So(len(module.MustGet().Submodules), ShouldEqual, 1)
			So(len(module.MustGet().SrcFiles), ShouldEqual, 1)
		})
		//		Convey("When the path is python", func() {
		//			module := pyParser.ParseDirectory(path)
		//			So(len(module.MustGet().Submodules), ShouldEqual, 2)
		//			So(len(module.MustGet().SrcFiles), ShouldEqual, 4)
		//
		//		})
		Convey("When the path is py", func() {
			path := "/Users/meetai/codefuse-chatbot"
			module := pyParser.ParseDirectory(path)
			So(len(module.MustGet().Submodules), ShouldEqual, 1)
			So(len(module.MustGet().SrcFiles), ShouldEqual, 1)
		})

	})

}

func TestReadDirEntries(t *testing.T) {
	Convey("ReadDir depth:", t, func() {
		entries, err := os.ReadDir("/Users")
		if err != nil {
			panic(err.Error())
		}

		for _, en := range entries {
			println(string(en.Name()))
		}
	})
}

func TestCodeAnalyzePy(t *testing.T) {
	pyParser := PythonParser{}

	Convey("Given some codebase ", t, func() {
		path := "/Users/meetai/codefuse-chatbot"
		Convey("When the path is py", func() {
			module := pyParser.ParseDirectory(path)
			So(len(module.MustGet().Submodules), ShouldEqual, 1)
			So(len(module.MustGet().SrcFiles), ShouldEqual, 1)
		})

	})

}

func TestDirTree(t *testing.T) {
	Convey("Given some codebase ", t, func() {
		path := "/Users/meetai/codefuse-chatbot"
		Convey("When the path is py", func() {
			var treeBuff = bytes.NewBuffer([]byte(""))
			dirTree(treeBuff, path, true, 5)
			println(treeBuff.String())
			So(len(treeBuff.String()), ShouldBeGreaterThan, 10)
		})

	})
}
