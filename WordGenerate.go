package main
import (
	"math/rand"
	"time"
	
)

func firstRandomWordGenerate() string {
    words := []string{"ヘッドフォン","サイヤ人","なあぜなあぜ","イカサマ","テニス","魚","格差", "身長", "ドラム", "嘘つき", "東京", "逆上がり", "意識", "運命", "永久", "覚悟", "記憶", "孤独", "姿", "性格", "想像", "魂", "強さ", "天才", "時めき", "仲間", "人間", "望み", "儚い", "光", "勇敢", "成果", "高校", "名残"}

    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}

func secondRandomWordGenerate() string {
    words := []string{"ノートパソコン","俺が基準","有頂天","カメラ撮影","当たり前","燃えるごみ","経験値","万歩計", "インド人", "トップランク", "思い出", "苦しみ", "経験", "最強", "真実", "地平線", "欠場"}

    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}

func thirdRandomWordGenerate() string {
	words := []string{"ワイヤレスイヤホン","LINEスタンプ","ワイヤレスキーボード","最強の要塞","モノマネ芸人","ターンテーブル","きんかんのど飴","ショーケース", "ローキック", "おかわり自由", "墾田永年私財法", "素因数分解", "ミッションインポッシブル"}

    // シード値を設定
    rand.Seed(time.Now().UnixNano())
    // 配列からランダムに単語を選択
    randomIndex := rand.Intn(len(words))
    return words[randomIndex]
}