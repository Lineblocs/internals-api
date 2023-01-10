package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	lineblocs "github.com/Lineblocs/go-helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	guuid "github.com/google/uuid"
	logrus_cloudwatchlogs "github.com/kdar/logrus-cloudwatchlogs"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"lineblocs.com/api/model"
)

var settings *model.GlobalSettings
var log = logrus.New()

func CreateAPIID(prefix string) string {
	id := guuid.New()
	return prefix + "-" + id.String()
}

func LookupBestCallRate(number string, typeRate string) *model.CallRate {
	return &model.CallRate{CallRate: 9.99}
}

func ToCents(dollars float64) int {
	result := dollars * 100
	return int(result)
}

func CalculateTTSCosts(length int) float64 {
	var result float64 = float64(length) * .000005
	return result
}

func CalculateSTTCosts(recordingLength float64) float64 {
	// Google cloud bills .006 per 15 seconds
	billable := recordingLength / 15
	var result float64 = 0.006 * billable
	return result
}

func CreateS3URL(folder string, id string) string {
	return "https://lineblocs.s3.ca-central-1.amazonaws.com/" + folder + "/" + id
}

func UploadS3(folder string, name string, file multipart.File) error {
	bucket := "lineblocs"
	key := folder + "/" + name
	// The session the S3 Uploader will use
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1")})
	if err != nil {
		return fmt.Errorf("S3 session err: %s", err)
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(session)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return nil
}

func GetPlanRecordingLimit(workspace *model.Workspace) (int, error) {
	if workspace.Plan == "pay-as-you-go" {
		return 1024, nil
	} else if workspace.Plan == "starter" {
		return 1024 * 2, nil
	} else if workspace.Plan == "pro" {
		return 1024 * 32, nil
	} else if workspace.Plan == "starter" {
		return 1024 * 128, nil
	}
	return 0, nil
}

func GetPlanFaxLimit(workspace *model.Workspace) (*int, error) {
	var res *int
	if workspace.Plan == "pay-as-you-go" {
		*res = 100
	} else if workspace.Plan == "starter" {
		*res = 100
	} else if workspace.Plan == "pro" {
		res = nil
	} else if workspace.Plan == "starter" {
		res = nil
	}
	return res, nil
}

func CheckRouteMatches(from string, to string, prefix string, prepend string, match string) (bool, error) {
	full := prefix + match
	valid, err := regexp.MatchString(full, to)
	if err != nil {
		return false, err
	}
	return valid, err
}

func CheckCIDRMatch(sourceIp string, ipWithCidr string) (bool, error) {
	// remove port if needed

	ipSlice1 := strings.Split(sourceIp, ":")
	ip1 := ipSlice1[0]

	// check if there is port
	ipSlice2 := strings.Split(ipWithCidr, ":")
	var ip2 string
	if len(ipSlice2) > 1 {
		cidr := strings.Split(ipSlice2[1], "/")
		ip2 = ipSlice2[0] + "/" + cidr[1]
	} else {
		ip2 = ipSlice2[0]
	}

	_, net1, err := net.ParseCIDR(ip1 + "/32")
	if err != nil {
		return false, err
	}
	_, net2, err := net.ParseCIDR(ip2)
	if err != nil {
		return false, err
	}

	return net2.Contains(net1.IP), nil
}

func GetDIDRoutedServer2(rtcOptimized bool) (*lineblocs.MediaServer, error) {
	servers, err := lineblocs.CreateMediaServers()

	if err != nil {
		return nil, err
	}

	var result *lineblocs.MediaServer
	for _, server := range servers {
		//if result == nil || result != nil && server.LiveCallCount < result.LiveCallCount && rtcOptimized == server.RtcOptimized {
		if result == nil || result != nil && server.LiveCPUPCTUsed < result.LiveCPUPCTUsed && rtcOptimized == server.RtcOptimized {
			result = server
		}
	}
	return result, nil
}

func CheckFreeTrialStatus(plan string, started time.Time) string {
	if plan == "trial" {
		now := time.Now()
		//make configurable
		expireDays := 10
		expireHours := expireDays * 24
		started.Add(time.Hour * time.Duration(expireHours))
		if started.After(now) {
			return "expired"
		}
		return "pending-expiry"
	}
	return "not-applicable"
}

func LookupSIPAddresses(host string) (*[]net.IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, errors.New("No IP match found..")
	}
	return &ips, nil
}

// get first match
func LookupSIPAddress(host string) (*string, error) {
	ips, err := LookupSIPAddresses(host)
	if err != nil {
		return nil, err
	}
	ip := (*ips)[0].String()
	return &ip, nil
}

func CheckSIPServerHealth(routingSIPURI string) (bool, error) {
	return true, nil
}

func ReadEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func CheckIfCarrier(token string) bool {
	return true
}

func HandleInternalErr(msg string, err error, c echo.Context) error {
	Log(logrus.FatalLevel, msg)
	Log(logrus.FatalLevel, err.Error())
	return c.JSON(http.StatusInternalServerError, NewError(err))
}

func SetSetting(gs model.GlobalSettings) {
	settings = &gs
}

func GetSetting() *model.GlobalSettings {
	return settings
}

// Configure Logrus
func InitLogrus() {
	log = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "%lvl%: %time% - %msg%\n",
		},
	}

	if os.Getenv("USE_LOG_TYPE") == "file" { // If logging type is file
		log.Out = os.Stdout
		file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	} else if os.Getenv("USE_LOG_TYPE") == "cloud watch" { // If you want to use AWS Cloud Watch
		group := os.Getenv("AWS_CLOUDWATCHLOGS_GROUP_NAME")
		stream := os.Getenv("AWS_CLOUDWATCHLOGS_STREAM_NAME")

		// logs.us-east-1.amazonaws.com
		// Define the session - using SharedConfigState which forces file or env creds
		sess, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            aws.Config{Region: aws.String("us-east-1")},
		})
		if err != nil {
			panic("Not going to be able to write to cloud watch if you cant create a session")
		}

		// Determine if we are authorized to access AWS with the credentials provided. This does not mean you have access to the
		// services required however.
		_, err = sts.New(sess).GetCallerIdentity(&sts.GetCallerIdentityInput{})
		if err != nil {
			panic("Couldn't Validate our aws credentials")
		}

		hook, err := logrus_cloudwatchlogs.NewHook(group, stream, sess)
		if err != nil {
			log.Fatal(err)
		}

		log.Hooks.Add(hook)
		log.Out = io.Discard
		log.Formatter = logrus_cloudwatchlogs.NewProdFormatter()
	} else if os.Getenv("USE_LOG_TYPE") == "logstash" { // If you want to use Logstash
		conn, err := net.Dial("tcp", "logstash.mycompany.net:8911")
		if err != nil {
			log.Fatal(err)
		}
		hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "myappName"}))

		log.Hooks.Add(hook)
		ctx := log.WithFields(logrus.Fields{
			"method": "main",
		})
		ctx.Info("Hello World!")
	}
}

func Log(level logrus.Level, message string) {
	log.Log(level, message)
}
