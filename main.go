package main

import (
	"crypto/tls"
	"time"

	//  "encoding/json"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/tidwall/gjson"
)

var bot *linebot.Client

const profileURL string = "https://api.line.me/v2/bot/profile/"

func test() {
	inpit := "喵 help"
	MegRune := []rune(strings.TrimSpace(inpit))
	i := strings.Index(inpit, "喵")
	fmt.Println(strings.Index(string(MegRune[i+1:]), "help"))
}
func main() {

	//test()

	//fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	//key := os.Getenv("GoogleKey")
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
	/*	test to get map data
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

/*
func selectTest() string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	var GroupID string
	var UserID string
	var UserName string
	// err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	err = conn.QueryRow(context.Background(), `SELECT "GroupID", "UserID", "UserName" from public."GroupProfile"`).Scan(&GroupID, &UserID, &UserName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		//os.Exit(1)
		return err.Err()
	}
	fmt.Println(GroupID, UserID, UserName)
	return GroupID + UserID + UserName
}
*/

func QueryTest() (string, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var sum string
	var count int
	// Send the query to the server. The returned rows MUST be closed
	// before conn can be used again.
	rows, err := conn.Query(context.Background(), `SELECT "GroupID", "UserID", "UserName", "GID" ,"Time" from public."GroupProfile"`)
	if err != nil {
		return "err", err
	}
	sum = "conn.Query success.\n"
	// rows.Close is called by rows.Next when all rows are read
	// or an error occurs in Next or Scan. So it may optionally be
	// omitted if nothing in the rows.Next loop can panic. It is
	// safe to close rows multiple times.
	defer rows.Close()

	// Iterate through the result set
	count = 0
	for rows.Next() {
		count = count + 1
		var GroupID string
		var UserID string
		var UserName string
		var GID string
		var InsertTime time.Time
		err = rows.Scan(&GroupID, &UserID, &UserName, &GID, &InsertTime)
		if err != nil {
			return "err", err
		}
		sum += "idx=" + strconv.Itoa(count) + ": GroupID=" + GroupID + ",UserID=" + UserID + ",UserName=" + UserName + ",GID=" + GID + ",Time=" + InsertTime.String() + "\n"
	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return rows.Err().Error(), err
	}

	return sum, nil
	// No errors found - do something with sum
}

func insertUserProfile(GroupID, UserID, UserName string) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	conn.QueryRow(context.Background(), "INSERT INTO GroupProfile VALUES($1,$2,$3)", GroupID, UserID, UserName)

	fmt.Println(GroupID, UserID, UserName)
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

func getReplyMsg(message string, source *linebot.EventSource) string {
	MegRune := []rune(strings.TrimSpace(message))
	i := strings.Index(message, "喵")
	var replyMsg string
	if i > -1 {
		//replyMsg = getActionMsg(string(MegRune[i+1:]), userID)
		replyMsg = "---功能回覆---\n" + getActionMsg(string(MegRune[i+1:]), source)
		replyMsg += "\n---使用者訊息---\n" + message
		//replyMsg += "\n---UserPorilfe---\n" + getUserProfile(source.UserID)
	} else {
		replyMsg = ""
	}
	return replyMsg
}

/*
--INSERT INTO public."GroupProfile"  ("GroupID", "UserID", "UserName", "Num", "Time") VALUES ('223213212' , '2', 'v小黑',1, '2020/06/16');
--SELECT * FROM public."GroupProfile";
--DELETE FROM  public."GroupProfile" WHERE "UserID" = '2';
*/
func insertTest(source *linebot.EventSource) string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		// os.Exit(1)
		return err.Error()
	}
	defer conn.Close(context.Background())
	UserName := getUserName(source.UserID)
	row := conn.QueryRow(context.Background(), `INSERT INTO public.GroupProfile (GroupID,UserID,UserName,Num,Time) VALUES($1,$2,$3,$4,$5)`, source.GroupID, source.UserID, UserName, 1, time.Now())
	var n string
	if row != nil {
		row.Scan(&n)
	}
	return "InsertGroupID=" + source.GroupID + `return code=` + n
}

//scanType:string  > default
//scanType:int  >
func testSQLCmd(SQLCmd string, scanType string) string {
	if scanType == "" {
		scanType = "string"
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		//os.Exit(1)
		return err.Error()
	}
	defer conn.Close(context.Background())

	var sum string
	// Send the query to the server. The returned rows MUST be closed
	// before conn can be used again.
	rows, err := conn.Query(context.Background(), SQLCmd)
	if err != nil {
		return err.Error()
	}
	// No errors found - do something with sum
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		if scanType == "int" {
			var n int
			err = rows.Scan(&n)
			if err != nil {
				return err.Error()
			}
			sum += strconv.Itoa(n) + "\n"
		} else {
			var n string
			err = rows.Scan(&n)
			if err != nil {
				return err.Error()
			}
			sum += n + "\n"
		}

	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return err.Error()
	}

	return sum
}

func getGroupCount(source *linebot.EventSource) string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
		//return err.Error()
	}
	defer conn.Close(context.Background())
	var count int
	log.Println("enter get GroupCount,source.GroupID=", source.GroupID)
	sqlcmd := `SELECT COUNT("GID") FROM "GroupProfile" WHERE "GroupID"=:'$1'`
	err = conn.QueryRow(context.Background(), sqlcmd, source.GroupID).Scan(&count)
	log.Println("QueryRow end, count=", count, "err=", err.Error())
	return strconv.Itoa(count)
}

func testInsert(source *linebot.EventSource) string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err.Error()
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(context.Background())

	//log.Println("GroupID=", source.GroupID, "UserID=", source.UserID, "UserName=", getUserName(source.UserID), "GID=", source.GroupID+strconv.Itoa(getGroupCount(source)), "time=", time.Now().Format("2006-01-02 15:04:05"))

	nowGroupIP := ""
	if source.GroupID == "" {
		nowGroupIP = "NULL"
	} else {
		nowGroupIP = source.GroupID
	}
	GID := ""
	if getGroupCount(source) == "0" {
		GID = nowGroupIP
	} else {
		GID = nowGroupIP + getGroupCount(source)
	}

	_, err = tx.Exec(context.Background(), `INSERT INTO public."GroupProfile" ("GroupID", "UserID", "UserName", "GID", "Time") VALUES($1,$2,$3,$4,$5)`, nowGroupIP, source.UserID, getUserName(source.UserID), GID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err.Error() + ", InsertGroupID=" + nowGroupIP + ",GID=" + GID
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err.Error() + ", InsertGroupID=" + nowGroupIP + ",GID=" + GID
	}
	return "InsertGroupID=" + nowGroupIP + ",GID=" + GID
}

func getActionMsg(msgTxt string, source *linebot.EventSource) string {
	if strings.Index(msgTxt, "help") == 1 || msgTxt == "" {
		return getHelp()
	} else if strings.Index(msgTxt, "所有人") == 1 {
		return tagUser(source.UserID)
	} else if strings.Index(msgTxt, "測試插入") == 1 {
		return "測試插入:" + testInsert(source)
	} else if strings.Index(msgTxt, "測試查詢") == 1 {
		result, _ := QueryTest()
		return "測試查詢:" + result
	} else if strings.Index(msgTxt, "測試數量") == 1 {
		return "測試數量:" + getGroupCount(source) + "查詢Group:" + source.GroupID
	} else if strings.Index(msgTxt, "DBCMD") == 1 {
		MegRune := []rune(strings.TrimSpace(msgTxt))
		i := strings.Index(msgTxt, "DBCMD")
		// return string(MegRune[i+len("DBCMD"):])
		return testSQLCmd(string(MegRune[i+len("DBCMD"):]), "")
	} else if strings.Index(msgTxt, "INTDBCMD") == 1 {
		MegRune := []rune(strings.TrimSpace(msgTxt))
		i := strings.Index(msgTxt, "INTDBCMD")
		return testSQLCmd(string(MegRune[i+len("INTDBCMD"):]), "int")
	}
	return "no action after getActionMsg"
}

func tagUser(userID string) string {
	//JSONuserProfile := getUserProfile(userID)
	//return `<@[^>]+>` + gjson.Get(JSONuserProfile, "displayName").String()
	//return `<@` + gjson.Get(JSONuserProfile, "displayName").String() +  `>`
	return `<@` + userID + `>`
}

func getHelp() string {
	helpMsg := `請輸入'喵 指令'
	目前指令：
		所有人	標記所有人(Ex: 喵 所有人)`
	return helpMsg
}

func getUserProfile(userID string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", profileURL+userID, nil)
	req.Header.Set("Authorization", "Bearer {"+os.Getenv("ChannelAccessToken")+"}")
	res, _ := client.Do(req)
	s, _ := ioutil.ReadAll(res.Body)
	return string(s)
}

func getUserName(UserID string) string {
	JSONuserProfile := getUserProfile(UserID)
	return gjson.Get(JSONuserProfile, "displayName").String()
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
				//quota, err := bot.GetMessageQuota().Do()

				if err != nil {
					log.Println("Quota err:", err)
				}

				replyMsg := getReplyMsg(message.Text, event.Source)
				if replyMsg == "" {
					log.Println("NO Action")
				} else {
					if _, err = bot.ReplyMessage(
						event.ReplyToken,
						//linebot.NewTextMessage(message.ID+":"+message.Text+" OK! remain message:"+strconv.FormatInt(quota.Value, 10)),
						linebot.NewTextMessage(replyMsg),
					).Do(); err != nil {
						log.Print(err)
					}
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
			// log.Println("name=", name)
			// log.Println("Latitude =", Latitude, ", Longitude =", Longitude)
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
	radius := "200"
	googleURL := "https://maps.googleapis.com/maps/api/place/nearbysearch/json?radius="
	googleURL += radius + "&type=restaurant"
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
