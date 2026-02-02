//TODO
//1️⃣ Add citations (task titles used)
// 2️⃣ Stream responses
// 3️⃣ Switch to llama3 and compare quality
// 4️⃣ Add fallback if no tasks pass threshold
// 5️⃣ Move embeddings to pgvector

package handler

import (
	models "To_DO_Assistant/Models"
	"To_DO_Assistant/constants"
	"To_DO_Assistant/helpers"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	tasks   []models.Task
	counter = 0
)

func Post(ctx *gin.Context) {
	var task models.Task

	err := ctx.ShouldBindJSON(&task)
	if err != nil {
		log.Print("Error Occured", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	task.Status = "Todo"
	task.Id = counter
	task.CreatedAt = time.Now()
	task.DueDate = time.Now().AddDate(0, 0, 1)
	taskstr := helpers.BuildTaskString(task)
	responseback, err := helpers.GetEmbeddings(taskstr)
	if err != nil {
		fmt.Println(constants.Post + "Error While getting embeddings" + err.Error())
	}
	// fmt.Println(constants.AskRag+"Got Embedding from local RAG agent", responseback)
	task.Embedding = responseback
	tasks = append(tasks, task)
	counter++
	fmt.Println(constants.Post+"Length of Embedding:", len(task.Embedding))
	ctx.JSON(http.StatusCreated, task)
}

func Get(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, tasks)
}

func Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		err := fmt.Errorf("Error in fetching id")
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	for i, task := range tasks {
		if task.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
		}
	}
	ctx.JSON(http.StatusNoContent, "Task deleted")
}

func SearchQuery(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		err := fmt.Errorf("Error in fetching query")
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	query = strings.ToLower(query)
	var results []models.Task
	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Title), query) ||
			strings.Contains(strings.ToLower(task.Description), query) ||
			strings.Contains(strings.ToLower(task.Tags), query) ||
			strings.Contains(strings.ToLower(task.Notes), query) ||
			strings.Contains(strings.ToLower(task.Priority), query) {
			results = append(results, task)
		}
	}
	ctx.JSON(http.StatusFound, results)
}

func Ask(ctx *gin.Context) {
	var req models.Ask
	var results []models.ScoreTasks
	// var FinalQuewords []string

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.Question) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "question is required"})
		return
	}

	rq := strings.ToLower(req.Question)
	words := strings.Fields(rq)

	// fmt.Println("[ASK] Entering the Words segregation section")
	// for i, word := range words {
	// 	fmt.Println(i, word)
	// 	if _, exists := models.Stopwords[word]; exists {
	// 		FinalQuewords = append(FinalQuewords, word)
	// 	}
	// }
	fmt.Println("[ASK] Initiating Task Evaluation Section")
	results = helpers.TaskEvaluation(tasks, words)
	fmt.Print("[Ask] Results", results)
	responseback := helpers.BuildAnswer(req.Question, results)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	ctx.JSON(http.StatusOK, gin.H{
		"tasks":  results,
		"Result": responseback,
	})
}

func AskRAG(ctx *gin.Context) {
	var req models.Ask
	var results, filteredtasks []models.TaskSimilarityScore
	// var FinalQuewords []string

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.Question) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "question is required"})
		return
	}

	QuesEmbedding, err := helpers.GetEmbeddings(req.Question)
	if err != nil {
		fmt.Println(constants.AskRag + "Error While getting Question embedding" + err.Error())
	}
	fmt.Println(constants.Post+"Length of Question Embedding:", len(QuesEmbedding))
	for _, task := range tasks {
		similarity := helpers.CosineSimilarity(QuesEmbedding, task.Embedding)
		fmt.Printf("\n"+constants.AskRag+"Similarity: \t Task%d Value%f", task.Id, similarity)
		results = append(results, models.TaskSimilarityScore{Task: task, SimilarityScore: similarity})
	}
	//similarity score check
	for _, result := range results {
		if constants.SimilarityThreshold <= result.SimilarityScore {
			filteredtasks = append(filteredtasks, result)
		}
	}

	//sorting on basis of similarity score
	sort.Slice(filteredtasks, func(i, j int) bool {
		return filteredtasks[i].SimilarityScore > filteredtasks[j].SimilarityScore
	})
	//topk check
	if len(filteredtasks) > constants.TopK {
		filteredtasks = filteredtasks[:constants.TopK]
	}

	prompt := helpers.BuildRAGPrompt(req.Question, filteredtasks)
	Response, err := helpers.GetRAGResponse(prompt)
	if err != nil {
		fmt.Println(constants.AskRag + "Error While response from GetRAGResponse" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Answer": Response,
	})
}
