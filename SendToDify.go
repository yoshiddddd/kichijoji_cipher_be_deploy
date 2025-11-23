package main
import(
	"encoding/json"
	"fmt"
	"io"
	// "log"
	"net/http"
	"os"
	// "github.com/joho/godotenv"
	"bytes"
)
func sendToDify(answers map[*Client]AnswerMessage) (DifyResponse, error) {
	
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
        return DifyResponse{}, fmt.Errorf("AnswerMessage が2つ必要ですが、%d つしかありません", len(data))
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
        return DifyResponse{}, fmt.Errorf("error encoding data to JSON: %v", err)
    }

    req, err := http.NewRequest("POST", "https://api.dify.ai/v1/chat-messages", bytes.NewBuffer(requestBody))
    if err != nil {
        return DifyResponse{}, fmt.Errorf("error creating HTTP request: %v", err)
    }
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
    req.Header.Set("Content-Type", "application/json")
    // リクエストの送信
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return DifyResponse{}, fmt.Errorf("error sending request to Dify: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return DifyResponse{}, fmt.Errorf("failed to send data to Dify, status code: %d", resp.StatusCode)
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return DifyResponse{}, fmt.Errorf("error reading response body: %v", err)
    }

    // レスポンスボディ全体をログ出力
    fmt.Printf("Dify Response Body: %s\n", string(body))

	// 1段階目: Dify APIの外側のレスポンスをパース
	var apiResponse DifyAPIResponse
    if err := json.Unmarshal(body, &apiResponse); err != nil {
        return DifyResponse{}, fmt.Errorf("error unmarshalling API response: %v", err)
    }

	// 2段階目: answerフィールドの中のJSON文字列を動的にパース
	var answerData map[string]interface{}
	if err := json.Unmarshal([]byte(apiResponse.Answer), &answerData); err != nil {
		return DifyResponse{}, fmt.Errorf("error unmarshalling answer JSON: %v", err)
	}

	// 動的フィールド名からデータを抽出
	difyResponse := DifyResponse{
		Winner:    getStringValue(answerData, "winner"),
		User1Name: getStringValue(answerData, "user1Name"),
		User2Name: getStringValue(answerData, "user2Name"),
		Feedback:  getStringValue(answerData, "feedback"),
	}

	// ユーザー名をキーとした回答とポイントを取得
	user1Name := difyResponse.User1Name
	user2Name := difyResponse.User2Name

	// User1の回答とポイント
	if user1Name != "" {
		// まず "user1Name_answer" の形式を試す
		difyResponse.User1Answer = getStringValue(answerData, user1Name+"_answer")
		difyResponse.User1Point = getIntValue(answerData, user1Name+"_point")
	}

	// User2の回答とポイント
	if user2Name != "" {
		// まず "user2Name_answer" の形式を試す
		difyResponse.User2Answer = getStringValue(answerData, user2Name+"_answer")
		difyResponse.User2Point = getIntValue(answerData, user2Name+"_point")

		// 見つからない場合は "user2Name2_answer" の形式を試す
		if difyResponse.User2Answer == "" {
			difyResponse.User2Answer = getStringValue(answerData, user2Name+"2_answer")
			difyResponse.User2Point = getIntValue(answerData, user2Name+"2_point")
		}
	}

    // パース結果をログ出力
    fmt.Printf("Parsed Dify Response - Winner: %s, User1: %s, User2: %s\n",
		difyResponse.Winner, difyResponse.User1Name, difyResponse.User2Name)

    return difyResponse, nil
}

// ヘルパー関数: mapから文字列値を取得
func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// ヘルパー関数: mapから整数値を取得
func getIntValue(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0
}