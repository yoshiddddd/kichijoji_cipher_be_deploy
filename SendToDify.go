package main
import(
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"bytes"
)
func sendToDify(answers map[*Client]AnswerMessage) (string, error) {
	
	// err := godotenv.Load()
    // if err != nil {
    //     log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
    // }
    token := os.Getenv("DIFY_APIKEY")
    fmt.Println(token)

    // AnswerMessage を格納するスライスを作成
    var data []AnswerMessage
    for _, answer := range answers {
        data = append(data, answer)
    }

    // AnswerMessage が2つあることを確認
    if len(data) < 2 {
        return "", fmt.Errorf("AnswerMessage が2つ必要ですが、%d つしかありません", len(data))
    }

    // 送信するクエリの内容を作成
    query := fmt.Sprintf("keyword: %s name(%s) Answer: %s CountTime: %d, name(%s) Answer: %s CountTime: %d",
        data[0].Data.Keyword,
        data[0].Data.Name, data[0].Data.Answer, data[0].Data.CountTime,
        data[1].Data.Name, data[1].Data.Answer, data[1].Data.CountTime)

	payload := DifyRequestPayload{
        Inputs:         map[string]interface{}{}, 
        Query:          query,
        ResponseMode:   "blocking",
        ConversationID: "",         
        User:           "abc-123",  
        Files: []File{
            {
                Type:           "image",
                TransferMethod: "remote_url",
                URL:            "https://cloud.dify.ai/logo/logo-site.png",
            },
        },
    }

	requestBody, err := json.Marshal(payload)
	// requestBody, err := json.Marshal(data)
    if err != nil {
        return "",fmt.Errorf("error encoding data to JSON: %v", err)
    }

    req, err := http.NewRequest("POST", "https://api.dify.ai/v1/chat-messages", bytes.NewBuffer(requestBody))
    if err != nil {
        return "",fmt.Errorf("error creating HTTP request: %v", err)
    }
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
    req.Header.Set("Content-Type", "application/json")
    // リクエストの送信
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "",fmt.Errorf("error sending request to Dify: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "",fmt.Errorf("failed to send data to Dify, status code: %d", resp.StatusCode)
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "",fmt.Errorf("error reading response body: %v", err)
    }
	var difyResponse DifyResponse
    if err := json.Unmarshal(body, &difyResponse); err != nil {
        return "",fmt.Errorf("error unmarshalling response: %v", err)
    }

    // `answer`フィールドの確認とログ出力
    return  difyResponse.Answer,nil
}