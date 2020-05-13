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
	//  "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/tidwall/gjson"
)

var bot *linebot.Client

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	//key := os.Getenv("GoogleKey")
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
	/*
			/*
			   mapData := getMapDate()
			   fmt.Println("--------------------------------")
			   fmt.Println(mapData)

		/*
			results := gjson.Get(getMapDate(), "results")
			if results.IsArray() {
				for i := 0; i < len(results.Array()); i++ {
					nowJson := results.Array()[i].String()
					business_status := gjson.Get(nowJson, "business_status")
					if business_status.String() == "OPERATIONAL" {
						name := gjson.Get(nowJson, "name")
						geometry := gjson.Get(nowJson, "geometry")
						fmt.Println("name=", name)
						fmt.Println("geometry=", geometry)
						fmt.Println("====================")

					}

				}
			}
	*/

	//oneRestaurant := getRestaurantTest()
	//log.Println(*oneRestaurant)
}

func getRestaurantTest() *restaurant {
	mapData := getMapDate()
	oneRestaurant := getOneRestaurant(string(mapData))
	log.Println("intest", oneRestaurant)
	return oneRestaurant
}

func getMapDataTest() string {
	input := `{"html_attributions":[],"next_page_token":"CqQCHwEAAH5IE5FDWxe87UnGZwQClsJQRntmtMK_rsgW5w8AWaT90q_5ASx4PoRsdZNRU7xx3lp-8Zy7NqwwRUCkL3HQd7HDKqXMu6IXYMX3a9Fhjx9Be0Q-sOP6emq3twhlazTU3pYo3Mg3_tpWQO7Kx0BqtJMd0jb5PMr-hmCxevLxagtscw4h6yz4068j9AXPEcYK2ek4h-wEJDXJQlck5OMyA71El_ispAQUZKu83FbZJXl7trioqujZyBswQFi8DSmFWzzzz8JR0nVqH2LTcEUk-hb9wZiwxxmTZp6y16OjtrbC1md7Vwd2twKbegUyFQrPkyPR12AYsh4k3pIncfPtwKSNB01AJ8EYI5wEHtGVtFoCDp9JQRb7TPYkqtlIXJ9GdBIQekkOU7YT_6Pizn53SE4YaBoUzJuM8zTUqP5C0qZB__jEb033a1g","results":[{"business_status":"OPERATIONAL","geometry":{"location":{"lat":24.961639,"lng":121.226577},"viewport":{"northeast":{"lat":24.9629762802915,"lng":121.2279073302915},"southwest":{"lat":24.9602783197085,"lng":121.2252093697085}}},"icon":"https://maps.gstatic.com/mapfiles/place_api/icons/restaurant-71.png","id":"059fe723f7c1c7984643eb4a23dd912ce6bf0b5c","name":"Yuanshao","opening_hours":{"open_now":false},"photos":[{"height":3024,"html_attributions":["<a href=\"https://maps.google.com/maps/contrib/110488068015271633600\">温淑君</a>"],"photo_reference":"CmRaAAAAmD-N78Ud66EFQrD_s_8iKlvTezhj-7YIXzsrfCTQ6ZxX8TBNFyyOJLDElpAHvWCamyMjOh-PnRfzfKSmykQGX66bgnGi6uT21ZRBMFW_45mxHR-giaQ5i0CUVl1frkHbEhCbtUfphSEvNHmoilskYS0EGhSQSzelrbBmOfrTo_8QgE86MnfgYA","width":4032}],"place_id":"ChIJDaPANjciaDQRQc3LojUgGe8","plus_code":{"compound_code":"X66G+MJ Taiwan, Zhongli District, 新街里","global_code":"7QP3X66G+MJ"},"price_level":2,"rating":4.2,"reference":"ChIJDaPANjciaDQRQc3LojUgGe8","scope":"GOOGLE","types":["restaurant","food","point_of_interest","establishment"],"user_ratings_total":938,"vicinity":"2F, No. 245號, Yuanhua Road, Zhongli District"},{"business_status":"OPERATIONAL","geometry":{"location":{"lat":24.9608113,"lng":121.2259249},"viewport":{"northeast":{"lat":24.9621300802915,"lng":121.2273143302915},"southwest":{"lat":24.9594321197085,"lng":121.2246163697085}}},"icon":"https://maps.gstatic.com/mapfiles/place_api/icons/restaurant-71.png","id":"65faec8cd4ccfdea93a83b266d7ee37359b74f04","name":"湄南小鎮泰國菜","opening_hours":{"open_now":false},"photos":[{"height":1365,"html_attributions":["<a href=\"https://maps.google.com/maps/contrib/109724900483908027869\">劉志宏</a>"],"photo_reference":"CmRaAAAAnpsfw1ou4kGXAa_o3Gb2_1Feo8DDc1yqjxwqJZDNMaUamxx30WaCZVklGbgg-CRPqoxpE4Buel6urz1cZRBqVj9Xs3NQ5ZvKtOm9KHDHbscIaVxhMxwRJnwgWMOtKb4BEhBtrUQ--92fAp-fAINzaiVBGhSoJGSfK429ihY8arSaotEs4L58xQ","width":2048}],"place_id":"ChIJLXSLGDciaDQRTGz66CqkGi4","plus_code":{"compound_code":"X66G+89 Taiwan, Zhongli District, 新街里","global_code":"7QP3X66G+89"},"price_level":2,"rating":4.2,"reference":"ChIJLXSLGDciaDQRTGz66CqkGi4","scope":"GOOGLE","types":["restaurant","food","point_of_interest","establishment"],"user_ratings_total":674,"vicinity":"No. 306號, Yanping Road, Zhongli District"}],"status":"OK"}`
	return input
}

func getMapDate() []byte {
	file, err := os.Open("mapData.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)

	log.Println(b)
	log.Println("--------------------------------")
	return b
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
				if _, err = bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(message.ID+":"+message.Text+" OK! remain message:"+strconv.FormatInt(quota.Value, 10)),
				).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.LocationMessage:
				resResult := *getRestaurant(message.Latitude, message.Longitude)
				log.Println("Restaurant result > ")
				log.Println(resResult)
				if _, err := bot.ReplyMessage(
					event.ReplyToken,
					//linebot.NewTextMessage("Name = "+resResult.name+"Latitude = "+resResult.Latitude+"Longitude = "+resResult.Longitude),
					linebot.NewLocationMessage(resResult.name, resResult.address, resResult.Latitude, resResult.Longitude),

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

type restaurant struct {
	name      string
	Latitude  float64
	Longitude float64
	address   string
}

func getRestaurant(Latitude, Longitude float64) *restaurant {
	//var jsonObj map[string]interface{}
	//json.Unmarshal(getJSONFromLocation(Latitude, Longitude), &jsonObj)
	//Todo:https://ithelp.ithome.com.tw/articles/10205062?sc=iThelpR
	mapData := getJSONFromLocation(Latitude, Longitude)
	oneRestaurant := getOneRestaurant(mapData)
	return oneRestaurant
}

func getOneRestaurant(mapData string) *restaurant {
	oneRestaurant := restaurant{}
	results := gjson.Get(mapData, "results")
	if results.IsArray() {
		nowJSON := results.Array()[rand.Intn(len(results.Array()))].String()
		fmt.Println(nowJSON)
		businessStatus := gjson.Get(nowJSON, "business_status")
		if businessStatus.String() == "OPERATIONAL" {
			name := gjson.Get(nowJSON, "name")
			Latitude := gjson.Get(nowJSON, "geometry.location.lat")
			Longitude := gjson.Get(nowJSON, "geometry.location.lng")
			address := gjson.Get(nowJSON, "vicinity")
			//geometry := gjson.Get(nowJson ,"geometry")
			log.Println("name=", name)
			log.Println("Latitude =", Latitude, ", Longitude =", Longitude)
			Lat, err := strconv.ParseFloat(Latitude.String(), 8)
			Lon, err := strconv.ParseFloat(Longitude.String(), 8)
			if err != nil {
				return &oneRestaurant
			}
			oneRestaurant.name = name.String()
			oneRestaurant.Latitude = Lat
			oneRestaurant.Longitude = Lon
			oneRestaurant.address = address.String()
		}
	}
	/*
	   for i := 0 ; i < len(results.Array()) ; i++{
	       nowJson := results.Array()[i].String()
	       business_status:= gjson.Get(nowJson ,"business_status")
	       if business_status.String() == "OPERATIONAL" {
	           name := gjson.Get(nowJson ,"name")
	           geometry := gjson.Get(nowJson ,"geometry")
	           fmt.Println("name=",name)
	           fmt.Println("geometry=",geometry)
	           fmt.Println("====================")

	       }

	   }
	*/
	return &oneRestaurant
}

func getJSONFromLocation(Latitude, Longitude float64) string {
	googleURL := "https://maps.googleapis.com/maps/api/place/nearbysearch/json?radius=500&type=restaurant"
	googleURL += "&location=" + fmt.Sprintf("%f", Latitude) + "," + fmt.Sprintf("%f", Longitude)
	googleURL += "&key=" + os.Getenv("GoogleKey")

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
	//gitfmt.Printf("%s", sitemap)
	status := gjson.Get(string(sitemap), "status")
	var mapResult string
	if status.String() == "OK" {
		mapResult = string(sitemap)
	} else {
		mapResult = ""
	}
	return mapResult
}
