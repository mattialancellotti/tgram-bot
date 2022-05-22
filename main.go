package main

import (
     "bytes"
     "encoding/json"
     "errors"
     "fmt"
     "net/http"
     "strings"

     "gopkg.in/ini.v1"
)

type TelegramWebHook struct {
     Message struct {
          Text string `json:"text"`
          Chat struct {
               ID int64 `json:"id"`
          } `json:"chat"`
     } `json:"message"`
}

func Handler(res http.ResponseWriter, req *http.Request) {
     /* Decoding the JSON response */
     body := &TelegramWebHook{}
     if err := json.NewDecoder(req.Body).Decode(body); err != nil {
          fmt.Println("Couldn't decode request body", err)
          return
     }

     /* This checks whether the messages contains or not the word athenae */
     if !strings.Contains(strings.ToLower(body.Message.Text), "athenae") {
          return
     }

     /* Building message and sending it to the client */
     if err := bot.makeRequest("sendMessage", sayAthenae(body.Message.Chat.ID)); err != nil {
          fmt.Println("Error in sending reply:", err)
          return
     }

     fmt.Println("reply sent")
}

type BotMessage struct {
     ChatID int64 `json:"chat_id"`
     Text string `json:"text"`
}

func sayAthenae(ID int64) *BotMessage {
     /* Building bot's message */
     return &BotMessage{
          ChatID: ID,
          Text: "Athenae!",
     }
}

func (bot *BotClient) makeRequest(API string, v interface{}) error {
     /* Packing everything in JSON */
     reqBytes, err := json.Marshal(v)
     if err != nil {
          return err
     }

     /* Making the response */
     res, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%d:%s/%s",
                              bot.BotID, bot.Token, API), "application/json",
                              bytes.NewBuffer(reqBytes))
     if err != nil {
          return err
     }

     /* Checking if the request return OK */
     if res.StatusCode != http.StatusOK {
          return errors.New("Unexpected status: " + res.Status)
     }

     return nil
}

type BotClient struct {
     BotID uint64 `ini:"Bot"`
     Token string `ini:"Token"`
}

var bot *BotClient

func main() {
     /* Loading and reading the INI file */
     bot = new(BotClient)
     err := ini.MapTo(bot, "private.ini")
     fmt.Println(bot.BotID)
     if err != nil {
          fmt.Printf("Failed to load private.ini: %v", err)
          return
     }

     http.ListenAndServe(":3000", http.HandlerFunc(Handler))
}
