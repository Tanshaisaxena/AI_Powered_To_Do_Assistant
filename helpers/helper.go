package helpers

import (
	models "To_DO_Assistant/Models"
	"To_DO_Assistant/constants"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// For building the answer in a string format
func BuildAnswer(question string, tasks []models.ScoreTasks) string {
	if len(tasks) == 0 {
		return "I couldn't find any tasks related to your question."
	}

	var sb strings.Builder
	_, err := sb.WriteString("The response to the question on tasks:")
	if err != nil {
		fmt.Println("Error while writing to string builder")
	}
	for _, task := range tasks {
		sb.WriteString("\nId:")
		sb.WriteString(strconv.Itoa(task.Task.Id))
		sb.WriteString("\nTitle:")
		sb.WriteString(task.Task.Title)
		sb.WriteString("\nDescription:")
		sb.WriteString(task.Task.Description)
		sb.WriteString("\nStatus:")
		sb.WriteString(task.Task.Status)
	}
	return sb.String()
}

func StopWordsCheck(word string) bool {
	_, exists := models.Stopwords[word]
	return exists
}

func GetEmbeddings(Taskstr string) ([]float64, error) {
	body := map[string]string{
		"model":  "nomic-embed-text",
		"prompt": Taskstr,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		fmt.Println(constants.GetEmbeddings + " Error Unmarshalling: " + err.Error())
	}
	resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(constants.GetEmbeddings + " Error during Http post: " + err.Error())
	}

	var Results struct {
		Embedding []float64 `json:"embedding"`
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&Results); err != nil {
		fmt.Println(constants.GetEmbeddings + " Error during Decode: " + err.Error())
	}

	return Results.Embedding, err
}

func BuildTaskString(task models.Task) string {
	return fmt.Sprintf("Title: %s Description:%s Tags:%s Status:%s Priority:%s Notes:%s ", task.Title, task.Description, task.Tags, task.Status, task.Priority, task.Notes)
}

func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float64
	var normA float64
	var normB float64

	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func BuildRAGPrompt(question string, tasks []models.TaskSimilarityScore) string {
	var sb strings.Builder

	sb.WriteString("You are a helpful task assistant.\n\n")
	sb.WriteString("Answer only using the tasks provided below.\n")
	sb.WriteString("If the answer is not present, say you don't know.\n\n")

	sb.WriteString("Relevant tasks:\n")

	for i, t := range tasks {
		sb.WriteString(fmt.Sprintf(
			"%d. %s\nNotes: %s\n\n",
			i+1,
			t.Task.Title,
			t.Task.Notes,
		))
	}

	sb.WriteString("User question:\n")
	sb.WriteString(question)

	return sb.String()
}

func GetRAGResponse(prompt string) (string, error) {
	var Responsefromrag struct {
		Response string `json:"response"`
	}
	payload := map[string]any{
		"model":  "mistral",
		"prompt": prompt,
		"stream": false,
	}
	body, err := json.Marshal(&payload)
	if err != nil {
		fmt.Println(constants.GenerateRAGResponse+"Error while marshalling payload", err.Error())
	}
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(constants.GenerateRAGResponse+"Error in POST request", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&Responsefromrag)
	if err != nil {
		fmt.Println(constants.GenerateRAGResponse+"Error in Decoding ", err.Error())
		return "", err
	}
	return Responsefromrag.Response, err
}

func TaskEvaluation(tasks []models.Task, words []string) (Results []models.ScoreTasks) {
	fmt.Println(constants.TaskEvaluation + " Starting Task Evaluation")
	for _, task := range tasks {
		ID := task.Id
		fmt.Println(constants.TaskEvaluation+" ID ", ID)
		Title := strings.ToLower(task.Title)
		fmt.Println(constants.TaskEvaluation+" Title ", Title)
		Description := strings.ToLower(task.Description)
		fmt.Println(constants.TaskEvaluation+" Description ", Description)
		Tags := strings.ToLower(task.Tags)
		fmt.Println(constants.TaskEvaluation+" Tags ", Tags)
		Notes := strings.ToLower(task.Notes)
		fmt.Println(constants.TaskEvaluation+" Notes ", Notes)
		Priority := strings.ToLower(task.Priority)
		fmt.Println(constants.TaskEvaluation+" Priority ", Priority)

		for _, word := range words {
			if StopWordsCheck(word) {
				continue
			}
			var counter = 0
			fmt.Println(constants.TaskEvaluation+" Meaningfull Word", word)
			if len(word) >= 3 {
				if strings.Contains(Title, word) {
					counter++
				}
				if strings.Contains(Description, word) {
					counter++
				}
				if strings.Contains(Tags, word) {
					counter++
				}
				if strings.Contains(Notes, word) {
					counter++
				}
				if strings.Contains(Priority, word) {
					counter++
				}
				if counter > 0 {
					Results = append(Results, models.ScoreTasks{Score: counter, Task: task})
					fmt.Println(constants.TaskEvaluation+" Found in task:", task.Id)
					break
				}
			}

		}

	}
	return
}
