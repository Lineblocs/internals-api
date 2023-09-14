package utils

import (
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	helpers "github.com/Lineblocs/go-helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	guuid "github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/model"
)

var settings *model.GlobalSettings
var log *logrus.Logger
var microserviceName string

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
	return float64(length) * .000005
}

func CalculateSTTCosts(recordingLength float64) float64 {
	// Google cloud bills .006 per 15 seconds
	billable := recordingLength / 15
	return 0.006 * billable
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
	Log(logrus.InfoLevel, fmt.Sprintf("file uploaded to, %s\n", aws.StringValue(&result.Location)))
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
	var res int
	switch workspace.Plan {
	case "pay-as-you-go", "starter":
		res = 100
	case "pro", "unknown":
		// Default case: leave res as 0 (nil)
	}

	return &res, nil
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

func GetDIDRoutedServer2(rtcOptimized bool) (*helpers.MediaServer, error) {
	servers, err := helpers.CreateMediaServers()

	if err != nil {
		return nil, err
	}

	var result *helpers.MediaServer
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
		expireTime := started.Add(time.Hour * 24 * time.Duration(expireDays))
		if now.After(expireTime) {
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
	return c.JSON(http.StatusInternalServerError, err.Error())
}

func SetSetting(gs model.GlobalSettings) {
	settings = &gs
}

func GetSetting() *model.GlobalSettings {
	return settings
}

/*
Input: level, message
Todo: Log message with level(Info, Warning, Error, Panic)
Output:
*/
func Log(level logrus.Level, message string) {
	helpers.Log(level, "("+microserviceName+") "+message)
}

/*
Store microservice name locally
*/
func SetMicroservice(username string) {
	microserviceName = username
}

/*
Config func to get env value from key ---
*/
func Config(key string) string {
	// load .env file
	loadDotEnv := os.Getenv("USE_DOTENV")
	if loadDotEnv != "off" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Print("Error loading .env file")
		}
	}
	return os.Getenv(key)
}

func CanPlaceAdditionalCalls() (bool, error) {
	return true, nil
}
