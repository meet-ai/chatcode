package knowledge

import (
	"chatcode/consts"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/samber/mo"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

type EmbeddingVector [][]float32

func BatchCreateEmbedding(content []string) mo.Result[EmbeddingVector] {
	//By default, the length of the embedding vector will be 1536 for text-embedding-3-small
	//or 3072 for text-embedding-3-large . You can reduce the dimensions of the ...
	//	模型	每美元粗略页数	BEIR 搜索评估的示例性能
	//text-embedding-ada-002	3000	53.9
	//*-davinci-*-001	6	52.8
	//*-curie-*-001	60	50.9
	//*-babbage-*-001	240	50.4
	//*-ada-*-001	300	49.0
	opts := []openai.Option{
		openai.WithModel("gpt-3.5-turbo-0125"),
		openai.WithEmbeddingModel("text-embedding-ada-002"),
	}
	llm, err := openai.New(opts...)
	if err != nil {
		return mo.Err[EmbeddingVector](err)
	}
	ctx := context.Background()
	embedings, err := llm.CreateEmbedding(ctx, content)
	if err != nil {
		return mo.Err[EmbeddingVector](err)
	}
	return mo.Ok[EmbeddingVector](embedings)
}
func CreateKnowledge(path string) mo.Result[string] {
	return writeChroma(BuildDocumentIn(path).OrEmpty())
}
func writeChroma(docs []schema.Document) mo.Result[string] {
	embeddingModelName := "text-embedding-ada-002"
	llm, _ := openai.New(
		openai.WithAPIType(openai.APITypeOpenAI),
		openai.WithEmbeddingModel(embeddingModelName),
	)
	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return mo.Err[string](err)
	}
	store, errNs := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithOpenAIAPIKey(os.Getenv("OPENAI_API_KEY")),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace(consts.CollectionName),
		chroma.WithEmbedder(e),
	)
	if errNs != nil {
		return mo.Err[string](err)
	}
	_, err = store.AddDocuments(context.TODO(), docs)
	if err != nil {
		return mo.Err[string](err)
	}
	return mo.Ok[string](consts.CollectionName)
}

func BuildDocumentIn(folderPath string) mo.Result[[]schema.Document] {
	splitter := NewCodeSplitter()
	var docs []schema.Document
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".py")) {
			return readFileAndSplit(path, splitter, &docs)
		}
		return nil
	})

	if err != nil {
		return mo.Err[[]schema.Document](err)
	}

	return mo.Ok[[]schema.Document](docs)
}

func readFileAndSplit(path string, splitter *CodeSplitter, docs *[]schema.Document) error {
	log.Printf("Reading file %s", path)
	fd, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %v", path, err)
		return err
	}

	ctx := context.TODO()
	docsFromFile, err := documentloaders.NewText(fd).LoadAndSplit(ctx, splitter)
	if err != nil {
		log.Printf("Error splitting file %s: %v", path, err)
		return err
	}

	for i := range docsFromFile {
		docsFromFile[i].Metadata["source"] = path
		*docs = append(*docs, docsFromFile[i])
	}

	log.Printf("Read %d documents from %s", len(docsFromFile), path)
	return nil
}

func QueryEmbedding(target string, topK int, threshold float32) []schema.Document {
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
		chroma.WithNameSpace(consts.CollectionName),
	)
	if errNs != nil {
		log.Fatalf("new: %v\n", errNs)
	}

	docs, err := store.SimilaritySearch(context.Background(),
		target, topK,
		vectorstores.WithScoreThreshold(threshold))
	if err != nil {
		log.Fatalf("search: %v\n", err)
		println(err.Error())
	}
	return docs
}

func BuildContext(docs []schema.Document) string {
	var buf strings.Builder
	for i, doc := range docs {
		_, err := fmt.Fprintf(&buf, "\n %d: %s\n", i, doc.PageContent)
		if err != nil {
			log.Fatalf("fprint: %v\n", err)
		}
	}
	return buf.String()
}
