package userunit

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/jackc/pgx/v4"
)

type UserRecord struct {
	Time        time.Time
	UserID      string
	UserName    string
	Message     string
	MessageType string
	PictureURL  string
	Latitude    string
	Longitude   string
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

func RecordInsert(c *gin.Context, MsgData UserRecord) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	log.Println("MsgData=", MsgData)
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `INSERT INTO public."MsgRecord" ("UserID", "UserPhoto",  "UserName", "Message","MessageType","IPAddress","Latitude, "Longitude","Time") VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`, MsgData.UserID, MsgData.PictureURL, MsgData.UserName, MsgData.Message, MsgData.MessageType, c.ClientIP(), MsgData.Latitude, MsgData.Longitude, time.Now().Format("2006-01-02 15:04:05"))
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
