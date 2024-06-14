package knowledge

import (
	"context"
	"fmt"
	"os"
	"testing"

	lop "github.com/samber/lo/parallel"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/chroma"
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

func TestReadCode(t *testing.T) {
	path := "/Users/meetai/chatcode"
	docs := BuildDocumentIn(path)
	println(len(docs.MustGet()))
}
func TestWriteChroma(t *testing.T) {
	path := "/Users/meetai/chatcode"
	docs := BuildDocumentIn(path)
	writeChroma(docs.MustGet())
}

func TestQueryChroma(t *testing.T) {
	t.Run("should split file into documents", func(t *testing.T) {
		doc := QueryEmbedding("BuildDocument 做了什么工作", 5, 0.6)
		println(len(doc), doc[0].PageContent)
	})

}
func TestReadFileAndSplit(t *testing.T) {
	t.Run("should split file into documents", func(t *testing.T) {
		// given
		filePath := "testdata/test.go"
		splitter := NewCodeSplitter()
		var docs []schema.Document

		// when
		err := readFileAndSplit(filePath, splitter, &docs)

		// then
		assert.NoError(t, err)
		assert.Len(t, docs, 1)
	})

	t.Run("should return error if file does not exist", func(t *testing.T) {
		// given
		filePath := "nonexistent.go"
		splitter := NewCodeSplitter()
		var docs []schema.Document

		// when
		err := readFileAndSplit(filePath, splitter, &docs)

		// then
		assert.Error(t, err)
		assert.Empty(t, docs)
	})
}

func TestRemoveCollection(t *testing.T) {
	t.Run("should remove collection in chroma", func(t *testing.T) {
		llm, _ := openai.New(openai.WithAPIType(openai.APITypeOpenAI),
			openai.WithEmbeddingModel("text-embedding-ada-002"),
		)
		e, qerr := embeddings.NewEmbedder(llm)
		if qerr != nil {
			panic(qerr)
		}
		store, errNs := chroma.New(
			chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
			chroma.WithOpenAIAPIKey(os.Getenv("OPENAI_API_KEY")),
			chroma.WithEmbedder(e),
			chroma.WithNameSpace("b4ebbfae-ecde-4569-8931-ff3ffccd0943"),
		)
		if errNs != nil {
			println(errNs)
		}
		err := store.RemoveCollection()
		if err != nil {
			println(err.Error())
		}
	})
}
