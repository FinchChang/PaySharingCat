package userunit

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/jackc/pgx/v4"
)

// UserRecord used for record all user messages
type UserRecord struct {
	Time        time.Time
	UserID      string
	UserName    string
	Message     string
	MessageType string
	PictureURL  string
	Latitude    float64
	Longitude   float64
}

const UserProfileURL string = "https://api.line.me/v2/bot/profile/{userId}"
const GroupProfileURL string = "https://api.line.me/v2/bot/group/{groupId}/member/{userId}"

// GetUserProfile is used to get user profile by URL
// Can check source type to decide the type is group , room or user only
func GetUserProfile(source *linebot.EventSource) string {
	var MemberURL string

	if source.Type == "group" {
		MemberURL = GroupProfileURL
		log.Println("GroupID=" + source.GroupID)
		strings.ReplaceAll(MemberURL, "{groupId}", source.GroupID)
	} else if source.Type == "room" {

	} else if source.Type == "user" {
		MemberURL = UserProfileURL
	}
	log.Println("GroupID=" + source.GroupID)
	log.Println("UserID=" + source.UserID)
	strings.ReplaceAll(MemberURL, "{userId}", source.UserID)
	log.Println("MemberURL=" + MemberURL)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", MemberURL, nil)
	req.Header.Set("Authorization", "Bearer {"+os.Getenv("ChannelAccessToken")+"}")
	res, _ := client.Do(req)
	s, _ := ioutil.ReadAll(res.Body)
	log.Println(string(s))
	return string(s)
}

func RecordInsert(c *gin.Context, MsgData UserRecord) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	//log.Println("MsgData=", MsgData)
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `INSERT INTO public."MsgRecord" ("UserID", "UserPhoto",  "UserName", "Message","MessageType","Latitude", "Longitude", "IPAddress","Time") VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`, MsgData.UserID, MsgData.PictureURL, MsgData.UserName, MsgData.Message, MsgData.MessageType, MsgData.Latitude, MsgData.Longitude, c.ClientIP(), time.Now().Format("2006-01-02 15:04:05"))
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
