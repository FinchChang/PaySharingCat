package userunit

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"paySharingCat/vendor/github.com/line/line-bot-sdk-go/linebot"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/tidwall/gjson"
)

type UserRecord struct {
	Time        time.Time
	UserID      string
	UserName    string
	Mesage      string
	MessageType string
}

const profileURL string = "https://api.line.me/v2/bot/profile/"

func getUserProfile(source *linebot.EventSource) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", profileURL+source.UserID, nil)
	req.Header.Set("Authorization", "Bearer {"+os.Getenv("ChannelAccessToken")+"}")
	res, _ := client.Do(req)
	s, _ := ioutil.ReadAll(res.Body)
	return string(s)
}

func RecordInsert(event *linebot.Event, source *linebot.EventSource, Message string) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(context.Background())

	//log.Println("GroupID=", source.GroupID, "UserID=", source.UserID, "UserName=", getUserName(source.UserID), "GID=", source.GroupID+strconv.Itoa(getGroupCount(source)), "time=", time.Now().Format("2006-01-02 15:04:05"))
	/*
		nowGroupIP := ""
		if source.GroupID == "" {
			nowGroupIP = "NULL"
		} else {
			nowGroupIP = source.GroupID
		}
		GID := ""
		GroupCount := ""
	*/
	UserProfileJSON := getUserProfile(source)
	_, err = tx.Exec(context.Background(), `INSERT INTO public."MsgRecord" ("UserID", "UserPhoto",  "UserName", "Message","MessageType","IPAddress", "Time") VALUES($1,$2,$3,$4,$5,$6,$7)`, source.UserID, gjson.Get(UserProfileJSON, "phpictureUrloto").String(), gjson.Get(UserProfileJSON, "displayName").String(), Message, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	//*output = "InsertGroupID=" + nowGroupIP + "\nGID=" + GID + "\nGroupCount = " + GroupCount
	return nil
}
