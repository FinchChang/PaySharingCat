// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	//	"encoding/json"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var bot *linebot.Client

func main() {
	/********************************
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	//key := os.Getenv("GoogleKey")
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
	*/
	mapData := getMapDate()
	fmt.Println("--------------------------------")
	fmt.Println(mapData)
}

func getMapDate() []byte {
	file, err := os.Open("mapData.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)

	fmt.Println(b)
	fmt.Println("--------------------------------")
	return []byte(b)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.LocationMessage:
				getRestaurant(message.Latitude, message.Longitude)
				if _, err := bot.ReplyMessage(
					event.ReplyToken,
					//linebot.NewLocationMessage(message.Title, message.Address, message.Latitude, message.Longitude),
					//linebot.NewTextMessage(message.Title, message.Address, message.Latitude, message.Longitude),
				).Do(); err != nil {
					//return err
					log.Print(err)
				}
				//return nil
			}

		}
	}
}

func getRestaurant(Latitude, Longitude float64) {
	//var jsonObj map[string]interface{}
	//json.Unmarshal(getJSONFromLocation(Latitude, Longitude), &jsonObj)
	//Todo:https://ithelp.ithome.com.tw/articles/10205062?sc=iThelpR

}

/*
type Restaurant struct {
	business_status String `json:"status"`
	geometry        string `json:name`
}
*/
func getJSONFromLocation(Latitude, Longitude float64) {
	googleURL := "https://maps.googleapis.com/maps/api/place/nearbysearch/json?radius=1500&type=restaurant"
	googleURL += "&location=" + fmt.Sprintf("%f", Latitude) + "," + fmt.Sprintf("%f", Longitude)
	googleURL += "&key=" + os.Getenv("GoogleKey")
	fmt.Println("googleURL=", googleURL)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Get(googleURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", sitemap)

	//var jsonObj map[string]interface{}
	//results := jsonObj["results"].(map[string]interface{})
}
