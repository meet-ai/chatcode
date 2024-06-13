package knowledge

import (
	"context"
	"fmt"
	"os"
	"testing"

	lop "github.com/samber/lo/parallel"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
)

func TestDocLoadSplit(t *testing.T) {
	fd, _ := os.Open("/Users/meetai/chatcode/go.sum")
	ctx := context.TODO()
	splitter := NewCodeSplitter()
	docs, err := documentloaders.NewText(fd).LoadAndSplit(ctx, splitter)
	if err != nil {
		println(err)
	}
	fmt.Printf("meta:%v,content:%s ,score:%f\n", docs[0].Metadata, docs[0].PageContent, docs[0].Score)

	lop.ForEach(docs, func(doc schema.Document, _ int) {
		doc.Metadata["source"] = "path-------------"
	})
	println(len(docs[0].PageContent))
}

func TestWriteChroma(t *testing.T) {
	path := "/Users/meetai/chatcode"
	docs, err := readCodeFromFolder(path)
	if err != nil {
		panic(err)
	}
	writeChroma(docs)
}

func TestQueryChroma(t *testing.T) {
	queryEmbedding()
}
