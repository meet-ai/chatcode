package knowledge

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	chroma_go "github.com/amikos-tech/chroma-go/types"
	"github.com/google/uuid"
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

func writeChroma(docs []schema.Document) {
	llm, err := openai.New()
	if err != nil {
		panic(err.Error())
	}
	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		panic(err.Error())
	}
	store, errNs := chroma.New(
		chroma.WithChromaURL("http://127.0.0.1:7892"),
		chroma.WithOpenAIAPIKey(os.Getenv("OPENAI_API_KEY")),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace(uuid.New().String()),
		chroma.WithEmbedder(e),
	)
	if errNs != nil {
		panic(errNs.Error())
	}
	_, err = store.AddDocuments(context.TODO(), docs)
	if err != nil {
		println(err.Error())
	}

	// type meta = map[string]any
	//
	// // Add documents to the vector store.
	//
	//	_, errAd := store.AddDocuments(context.Background(), []schema.Document{
	//		{PageContent: "Tokyo", Metadata: meta{"population": 9.7, "area": 622}},
	//		{PageContent: "Kyoto", Metadata: meta{"population": 1.46, "area": 828}},
	//		{PageContent: "Hiroshima", Metadata: meta{"population": 1.2, "area": 905}},
	//		{PageContent: "Kazuno", Metadata: meta{"population": 0.04, "area": 707}},
	//		{PageContent: "Nagoya", Metadata: meta{"population": 2.3, "area": 326}}, {PageContent: "Toyota", Metadata: meta{"population": 0.42, "area": 918}},
	//		{PageContent: "Fukuoka", Metadata: meta{"population": 1.59, "area": 341}},
	//		{PageContent: "Paris", Metadata: meta{"population": 11, "area": 105}},
	//		{PageContent: "London", Metadata: meta{"population": 9.5, "area": 1572}},
	//		{PageContent: "Santiago", Metadata: meta{"population": 6.9, "area": 641}},
	//		{PageContent: "Buenos Aires", Metadata: meta{"population": 15.5, "area": 203}},
	//		{PageContent: "Rio de Janeiro", Metadata: meta{"population": 13.7, "area": 1200}},
	//		{PageContent: "Sao Paulo", Metadata: meta{"population": 22.6, "area": 1523}},
	//	})
	//
	//	if errAd != nil {
	//		log.Fatalf("AddDocument: %v\n", errAd)
	//	}
	//
	// ctx := context.TODO()
	//
	//	type exampleCase struct {
	//		name         string
	//		query        string
	//		numDocuments int
	//		options      []vectorstores.Option
	//	}
	//
	// type filter = map[string]any
	//
	//	exampleCases := []exampleCase{
	//		{
	//			name:         "Up to 5 Cities in Japan",
	//			query:        "Which of these are cities are located in Japan?",
	//			numDocuments: 5,
	//			options: []vectorstores.Option{
	//				vectorstores.WithScoreThreshold(0.8),
	//			},
	//		},
	//		{
	//			name:         "A City in South America",
	//			query:        "Which of these are cities are located in South America?",
	//			numDocuments: 1,
	//			options: []vectorstores.Option{
	//				vectorstores.WithScoreThreshold(0.8),
	//			},
	//		},
	//		{
	//			name:         "Large Cities in South America",
	//			query:        "Which of these are cities are located in South America?",
	//			numDocuments: 100,
	//			options: []vectorstores.Option{
	//				vectorstores.WithFilters(filter{
	//					"$and": []filter{
	//						{"area": filter{"$gte": 1000}},
	//						{"population": filter{"$gte": 13}},
	//					},
	//				}),
	//			},
	//		},
	//	}
	//
	// // run the example cases
	// results := make([][]schema.Document, len(exampleCases))
	//
	//	for ecI, ec := range exampleCases {
	//		docs, errSs := store.SimilaritySearch(ctx, ec.query, ec.numDocuments, ec.options...)
	//		if errSs != nil {
	//			log.Fatalf("query1: %v\n", errSs)
	//		}
	//		results[ecI] = docs
	//	}
	//
	// // print out the results of the run
	// fmt.Printf("Results:\n")
	//
	//	for ecI, ec := range exampleCases {
	//		texts := make([]string, len(results[ecI]))
	//		for docI, doc := range results[ecI] {
	//			texts[docI] = doc.PageContent
	//		}
	//		fmt.Printf("%d. case: %s\n", ecI+1, ec.name)
	//		fmt.Printf("    result: %s\n", strings.Join(texts, ", "))
	//	}
}

func readCodeFromFolder(folderPath string) ([]schema.Document, error) {
	var alldocs []schema.Document
	splitter := NewCodeSplitter()

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".py")) {
			fd, err := os.Open(path)
			if err != nil {
				return err
			}

			ctx := context.TODO()
			docs, err := documentloaders.NewText(fd).LoadAndSplit(ctx, splitter)
			if err != nil {
				println(err)
			}
			for i := range docs {
				docs[i].Metadata["source"] = path
			}
			alldocs = append(alldocs, docs...)
		} else if info.IsDir() {
			doc, _ := readCodeFromFolder(filepath.Join(folderPath, info.Name()))
			alldocs = append(alldocs, doc...)
		}
		//TODO add more dir
		return nil
	})

	return alldocs, err
}

// Function to split text into chunks

// Function to generate embeddings using OpenAI API
//func generateEmbeddings(contents []string, step int) mo.Result[EmbeddingVector] {
//	var allEmbeddings [][]float32
//
//	for i := 0; i < len(contents); i += step {
//		end := i + step
//		if end > len(contents) {
//			end = len(contents)
//		}
//		batch := contents[i:end]
//		embeddings := BatchCreateEmbedding(batch)
//		allEmbeddings = append(allEmbeddings, embeddings.OrEmpty()...)
//	}
//	return mo.Ok[EmbeddingVector](allEmbeddings)
//}

// Function to write embeddings to Chroma
//func writeToChroma(embeddings [][]float32, filePaths []string, chunkIds []string) error {
//	client, err := chroma.NewClient(chroma.Config{
//		Host: "http://localhost:8000", // Replace with your Chroma server address
//	})
//	if err != nil {
//		return err
//	}
//
//	for i, embedding := range embeddings {
//		err := client.IndexDocument(chunkIds[i], embedding)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

func queryEmbedding() {
	llm, _ := openai.New(
		openai.WithAPIType(openai.APITypeOpenAI),
	)
	e, qerr := embeddings.NewEmbedder(llm)
	if qerr != nil {
		panic(qerr)
	}
	store, errNs := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithOpenAIAPIKey(os.Getenv("OPENAI_API_KEY")),
		chroma.WithDistanceFunction(chroma_go.COSINE),
		chroma.WithNameSpace(uuid.New().String()),
		chroma.WithEmbedder(e),
	)
	if errNs != nil {
		log.Fatalf("new: %v\n", errNs)
	}
	docs, err := store.SimilaritySearch(context.Background(),
		"这篇代码讲什么功能", 10,
		vectorstores.WithScoreThreshold(0.8))
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("meta:%v,content:%s ,score:%f\n", docs[0].Metadata, docs[0].PageContent, docs[0].Score)
}
