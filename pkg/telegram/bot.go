package telegram

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func RunBot() {
	var bot *tgbotapi.BotAPI
	var err error
	if os.Getenv("ENV") == "PROD" {
		bot, err = tgbotapi.NewBotAPI(os.Getenv("TGTOKEN"))
		if err != nil {
			log.Panic(err)
		}
	} else if os.Getenv("ENV") == "DOCKER" {
		bot, err = tgbotapi.NewBotAPI(os.Getenv("TGTOKENENVFILE"))
		if err != nil {
			log.Panic(err)
		}
	} else {
		bot, err = tgbotapi.NewBotAPI(os.Getenv("TGTOKENENVFILE"))
		if err != nil {
			log.Panic(err)
		}
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			text := update.Message.Text
			s := strings.Split(text, " ")
			if len(s) != 2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong amount of parameters, please use \"COIN timeframe\" format. For"+
					" example \"BTCUSDT 30m\" or \"sdd 20220403\"")
				bot.Send(msg)
			} else {
				switch {
				case strings.Contains(text, "sdd"):
					resp, _ := http.Get(fmt.Sprintf("https://heroku-tool.herokuapp.com/tg/sdd/%s", s[1]))
					body, _ := io.ReadAll(resp.Body)
					resp.Body.Close()
					symbolMap := make(map[string]float64)
					json.Unmarshal(body, &symbolMap)

					s := make(SDDUPairList, len(symbolMap))

					i := 0
					for k, v := range symbolMap {
						s[i] = SDDUPair{k, v}
						i++
					}

					sort.Sort(sort.Reverse(s))

					var resString, resString1 string
					for i, k := range s {
						if i <= len(s)/2 {
							resString += fmt.Sprintf("%22s", fmt.Sprintf("%s %s%% ", k.Key, strconv.FormatFloat(k.Value, 'f', 1, 64)))
						} else {
							resString1 += fmt.Sprintf("%22s", fmt.Sprintf("%s %s%% ", k.Key, strconv.FormatFloat(k.Value, 'f', 1, 64)))
							//fmt.Sprintf(k.Key + " " + strconv.FormatFloat(k.Value, 'f', 1, 64) + "%\n")
						}
					}

					/*keys := make([]string, 0, len(symbolMap))
					for k := range symbolMap {
						keys = append(keys, k)
					}
					sort.Strings(keys)
					var resString string
					var resString1 string
					for i, k := range keys {
						if i <= len(keys)/2 {
							resString += k + strconv.FormatFloat(symbolMap[k], 'f', 1, 64) + "% || "
						} else {
							resString1 += k + strconv.FormatFloat(symbolMap[k], 'f', 1, 64) + "% || "
						}
					}*/

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, resString)
					msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, resString1)

					bot.Send(msg)
					bot.Send(msg1)
				case strings.Contains(text, "sdu"):
					resp, _ := http.Get(fmt.Sprintf("https://heroku-tool.herokuapp.com/tg/sdu/%s", s[1]))
					body, _ := io.ReadAll(resp.Body)
					resp.Body.Close()
					symbolMap := make(map[string]float64)
					json.Unmarshal(body, &symbolMap)

					s := make(SDDUPairList, len(symbolMap))

					i := 0
					for k, v := range symbolMap {
						s[i] = SDDUPair{k, v}
						i++
					}

					sort.Sort(sort.Reverse(s))

					var resString, resString1 string
					for i, k := range s {
						if i <= len(s)/2 {
							resString += fmt.Sprintf("%22s", fmt.Sprintf("%s %s%% ", k.Key, strconv.FormatFloat(k.Value, 'f', 1, 64)))
						} else {
							resString1 += fmt.Sprintf("%22s", fmt.Sprintf("%s %s%% ", k.Key, strconv.FormatFloat(k.Value, 'f', 1, 64)))
							//fmt.Sprintf(k.Key + " " + strconv.FormatFloat(k.Value, 'f', 1, 64) + "%\n")
						}
					}

					/*keys := make([]string, 0, len(symbolMap))
					for k := range symbolMap {
						keys = append(keys, k)
					}
					sort.Strings(keys)
					var resString string
					var resString1 string
					for i, k := range keys {
						if i <= len(keys)/2 {
							resString += k + strconv.FormatFloat(symbolMap[k], 'f', 1, 64) + "% || "
						} else {
							resString1 += k + strconv.FormatFloat(symbolMap[k], 'f', 1, 64) + "% || "
						}
					}*/

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, resString)
					msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, resString1)

					bot.Send(msg)
					bot.Send(msg1)
				default:
					s0 := strings.ToUpper(s[0])
					s1 := strings.ToLower(s[1])
					resp, _ := http.Get(fmt.Sprintf("https://heroku-tool.herokuapp.com/tg/%s/%s", s0, s1))
					body, _ := io.ReadAll(resp.Body)
					resp.Body.Close()
					var test *CandleResponse
					json.Unmarshal(body, &test)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, test.Candle.String())

					bot.Send(msg)
				}
			}
		}
	}
}

type SDDUPair struct {
	Key   string
	Value float64
}

type SDDUPairList []SDDUPair

func (s SDDUPairList) Len() int           { return len(s) }
func (s SDDUPairList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SDDUPairList) Less(i, j int) bool { return s[i].Value < s[j].Value }

type CandleResponse struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Candle  *Candle `json:"candle"`
}

func (c Candle) String() string {
	return fmt.Sprintf("Last %v candle of %v pair closed at %v with:\nopen: %v\nhigh: %v\nlow: %v \nclose: %v", c.Timeframe, c.Coin, c.UTCCloseTime,
		c.Open, c.High, c.Low, c.Close)
}

type Candle struct {
	//ID                       primitive.ObjectID `bson:"_id"`
	ID                       int64     `bson:"id" json:"id" pg:"id"`
	MyID                     string    `bson:"my_id" json:"my_id" pg:"my_id,use_zero"`
	CoinTF                   string    `bson:"coin_tf" json:"coin_tf" pg:"coin_tf,use_zero"`
	Coin                     string    `bson:"coin" json:"coin" pg:"coin,use_zero"`
	Timeframe                string    `bson:"timeframe" json:"timeframe" pg:"timeframe,use_zero"`
	UTCOpenTime              time.Time `bson:"utc_open_time" json:"utc_open_time" pg:"utc_open_time,use_zero"`
	OpenTime                 int64     `bson:"open_time" json:"open_time" pg:"open_time,use_zero"`
	Open                     float64   `bson:"open" json:"open" pg:"open,use_zero"`
	High                     float64   `bson:"high" json:"high" pg:"high,use_zero"`
	Low                      float64   `bson:"low" json:"low" pg:"low,use_zero"`
	Close                    float64   `bson:"close" json:"close" pg:"close,use_zero"`
	Volume                   float64   `bson:"volume" json:"volume" pg:"volume,use_zero"`
	UTCCloseTime             time.Time `bson:"utc_close_time" json:"utc_close_time" pg:"utc_close_time,use_zero"`
	CloseTime                int64     `bson:"close_time" json:"close_time" pg:"close_time,use_zero"`
	QuoteAssetVolume         float64   `bson:"quote_asset_volume" json:"quote_asset_volume" pg:"quote_asset_volume,use_zero"`
	NumberOfTrades           int64     `bson:"number_of_trades" json:"number_of_trades" pg:"number_of_trades,use_zero"`
	TakerBuyBaseAssetVolume  float64   `bson:"taker_buy_base_asset_volume" json:"taker_buy_base_asset_volume" pg:"taker_buy_base_asset_volume,use_zero"`
	TakerBuyQuoteAssetVolume float64   `bson:"taker_buy_quote_asset_volume" json:"taker_buy_quote_asset_volume" pg:"taker_buy_quote_asset_volume,use_zero"`
	MA50                     float64   `bson:"ma50" json:"ma50" pg:"ma50,use_zero"`
	MA50Trend                bool      `bson:"ma50trend" json:"ma50trend" pg:"ma50trend,use_zero"`
	MA100                    float64   `bson:"ma100" json:"ma100" pg:"ma100,use_zero"`
	MA100Trend               bool      `bson:"ma100trend" json:"ma100trend" pg:"ma100trend,use_zero"`
	MA200                    float64   `bson:"ma200" json:"ma200" pg:"ma200,use_zero"`
	MA200Trend               bool      `bson:"ma200trend" json:"ma200trend" pg:"ma200trend,use_zero"`
}
