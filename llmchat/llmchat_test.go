package llmchat

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChatCode(t *testing.T) {
	Convey("Given some codebase ", t, func() {
		path := "/Users/meetai/chatcode"

		Convey("When the path is go", func() {
			content := ChatCode(path)
			So(content.IsOk(), ShouldEqual, true)
		})
	})
}

func TestChatCodeDir(t *testing.T) {
	Convey("Given some codebase ", t, func() {
		path := "/Users/meetai/codefuse-chatbot"

		Convey("When the path is go", func() {
			content := ChatCodeDir(path)
			So(content.IsOk(), ShouldEqual, true)
		})
	})
}
