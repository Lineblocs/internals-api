package store

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	helpers "github.com/Lineblocs/go-helpers"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

/*
Implementation of Logger Store
*/

type LoggerStore struct {
	db *sql.DB
}

func NewLoggerStore(db *sql.DB) *LoggerStore {
	return &LoggerStore{
		db: db,
	}
}

/*
Input: Log model
Todo : Create log model and store to db, send log email
Output: First Value: LastInsertId, Second Value: error
If success return (logId, nil) else return (nil, err)
*/
func (ls *LoggerStore) StartLogRoutine(workspace *model.Workspace, log *model.LogRoutine) (*string, error) {
	var user *helpers.User

	user, err := helpers.GetUserFromDB(log.UserId)
	if err != nil {
		utils.Log(logrus.ErrorLevel, "Could not get user..")
		return nil, err
	}

	now := time.Now()
	apiId := utils.CreateAPIID("log")
	stmt, err := ls.db.Prepare("INSERT INTO debugger_logs (`from`, `to`, `title`, `report`, `workspace_id`, `level`, `api_id`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )")

	if err != nil {
		utils.Log(logrus.ErrorLevel, "Could not prepare query..")
		return nil, err
	}

	defer stmt.Close()
	res, err := stmt.Exec(log.From, log.To, log.Title, log.Report, workspace.Id, log.Level, apiId, now, now)
	if err != nil {
		utils.Log(logrus.ErrorLevel, "Could not execute query..")
		return nil, err
	}

	logId, err := res.LastInsertId()
	if err != nil {
		utils.Log(logrus.ErrorLevel, "Could not get insert id..")
		return nil, err
	}
	logIdStr := strconv.FormatInt(logId, 10)

	go sendLogRoutineEmail(log, user, workspace)

	return &logIdStr, err
}

/*
Input: LogRoutine model, User model, Workspace model
Todo : Send Log email.
Output: If success return nil else return err
*/
func sendLogRoutineEmail(log *model.LogRoutine, user *helpers.User, workspace *model.Workspace) error {
	mg := mailgun.NewMailgun(utils.Config("MAILGUN_DOMAIN"), utils.Config("MAILGUN_API_KEY"))
	m := mg.NewMessage(
		"Lineblocs <monitor@lineblocs.com>",
		"Debug Monitor",
		"Debug Monitor",
		user.Email)
	m.AddCC("contact@lineblocs.com")
	//m.AddBCC("bar@example.com")

	body := `<html>
<head></head>
<body>
	<h1>Lineblocs Monitor Report</h1>
	<h5>` + log.Title + `</h5>
	<p>` + log.Report + `</p>
</body>
</html>`

	m.SetHtml(body)
	//m.AddAttachment("files/test.jpg")
	//m.AddAttachment("files/test.txt")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := mg.Send(ctx, m)
	if err != nil {
		return err
	}
	return nil
}
