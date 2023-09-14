package utils

import (
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"lineblocs.com/api/model"
)

func Test_CreateAPIID(t *testing.T) {
	t.Run("Should create a valid API ID with the given prefix", func(t *testing.T) {
		prefix := "api"
		result := CreateAPIID(prefix)

		assert.True(t, strings.HasPrefix(result, prefix+"-"))

		_, err := uuid.Parse(result[len(prefix)+1:])
		assert.NoError(t, err)
	})
}

func Test_LookupBestCallRate(t *testing.T) {
	t.Run("Should return the correct CallRate struct", func(t *testing.T) {
		expectedRate := 9.99
		result := LookupBestCallRate("someNumber", "someType")

		assert.Equal(t, expectedRate, result.CallRate)
	})
}

func Test_ToCents(t *testing.T) {
	t.Run("Should convert dollars to cents correctly", func(t *testing.T) {
		dollars := 12.34
		expectedCents := 1234
		result := ToCents(dollars)

		assert.Equal(t, expectedCents, result)
	})
}

func Test_CalculateTTSCosts(t *testing.T) {
	t.Run("Should calculate TTS costs correctly", func(t *testing.T) {
		length := 1000
		expectedCost := 0.005
		result := CalculateTTSCosts(length)

		assert.InDelta(t, expectedCost, result, 0.000001)
	})
}

func Test_CalculateSTTCosts(t *testing.T) {
	t.Run("Should calculate STT costs correctly", func(t *testing.T) {
		recordingLength := 75.0
		expectedCost := 0.03 // 75 seconds / 15 seconds * 0.006 = 0.03
		result := CalculateSTTCosts(recordingLength)

		assert.InDelta(t, expectedCost, result, 0.000001)
	})
}

func Test_CreateS3URL(t *testing.T) {
	t.Run("Should create a valid S3 URL", func(t *testing.T) {
		folder := "myfolder"
		id := "myid"
		expectedURL := "https://lineblocs.s3.ca-central-1.amazonaws.com/myfolder/myid"
		result := CreateS3URL(folder, id)

		assert.Equal(t, expectedURL, result)
	})
}

func Test_GetPlanRecordingLimit(t *testing.T) {
	t.Run("Should return correct recording limit for different plans", func(t *testing.T) {
		plans := []string{"pay-as-you-go", "starter", "pro", "unknown"}
		expectedLimits := []int{1024, 1024 * 2, 1024 * 32, 0}

		for i, plan := range plans {
			workspace := &model.Workspace{Plan: plan}
			result, err := GetPlanRecordingLimit(workspace)

			assert.NoError(t, err)
			assert.Equal(t, expectedLimits[i], result)
		}
	})
}

func Test_GetPlanFaxLimit(t *testing.T) {
	t.Run("Should return correct fax limit for different plans", func(t *testing.T) {
		// Test with different plans
		plans := []string{"pay-as-you-go", "starter", "pro", "unknown"}
		expectedLimits := []int{100, 100, 0, 0}

		for i, plan := range plans {
			workspace := &model.Workspace{Plan: plan}
			result, err := GetPlanFaxLimit(workspace)

			assert.NoError(t, err)
			assert.Equal(t, expectedLimits[i], *result)
		}
	})
}

func Test_CheckRouteMatches(t *testing.T) {

	fromExample := "source.example.com"
	toExample := "destination.example.com"

	t.Run("Should match route correctly", func(t *testing.T) {
		from := fromExample
		to := toExample
		prefix := "destination."
		prepend := "prepend"
		match := "example.com"

		result, err := CheckRouteMatches(from, to, prefix, prepend, match)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("Should return false for mismatched route", func(t *testing.T) {
		from := fromExample
		to := "different.example.com"
		prefix := "destination."
		prepend := "prepend"
		match := "example.com"

		result, err := CheckRouteMatches(from, to, prefix, prepend, match)

		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("Should return error for invalid regex", func(t *testing.T) {
		from := fromExample
		to := toExample
		prefix := "destination["
		prepend := "prepend"
		match := "example.com"

		_, err := CheckRouteMatches(from, to, prefix, prepend, match)

		assert.Error(t, err)
	})
}

func Test_CheckCIDRMatch(t *testing.T) {

	sourceIpTest := "192.168.1.1"
	ipWithCidr := "192.168.1.0/24"

	t.Run("Should match CIDR correctly", func(t *testing.T) {
		ipWithCidr := ipWithCidr

		result, err := CheckCIDRMatch(sourceIpTest, ipWithCidr)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("Should return false for non-matching CIDR", func(t *testing.T) {
		ipWithCidr := "192.168.2.0/24"

		result, err := CheckCIDRMatch(sourceIpTest, ipWithCidr)

		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("Should return error for invalid IP", func(t *testing.T) {
		sourceIp := "invalid.ip.address"
		ipWithCidr := ipWithCidr

		_, err := CheckCIDRMatch(sourceIp, ipWithCidr)

		assert.Error(t, err)
	})

	t.Run("Should return error for invalid CIDR", func(t *testing.T) {
		sourceIp := sourceIpTest
		ipWithCidr := "invalid.cidr"

		_, err := CheckCIDRMatch(sourceIp, ipWithCidr)

		assert.Error(t, err)
	})
}

func Test_CheckFreeTrialStatus(t *testing.T) {
	t.Run("Should return 'expired' for trial started 10 days ago", func(t *testing.T) {
		started := time.Now().Add(-10 * 24 * time.Hour)
		result := CheckFreeTrialStatus("trial", started)
		assert.Equal(t, "expired", result)
	})

	t.Run("Should return 'pending-expiry' for trial started 5 days ago", func(t *testing.T) {
		started := time.Now().Add(-5 * 24 * time.Hour)
		result := CheckFreeTrialStatus("trial", started)
		assert.Equal(t, "pending-expiry", result)
	})

	t.Run("Should return 'not-applicable' for non-trial plan", func(t *testing.T) {
		result := CheckFreeTrialStatus("pro", time.Now())
		assert.Equal(t, "not-applicable", result)
	})
}

func Test_LookupSIPAddresses(t *testing.T) {
	t.Run("Should return IP addresses for a valid host", func(t *testing.T) {
		host := "example.com"
		ips, err := LookupSIPAddresses(host)
		assert.NoError(t, err)
		assert.NotNil(t, ips)
		assert.NotEmpty(t, *ips)
		assert.IsType(t, &[]net.IP{}, ips)
	})

	t.Run("Should return an error for an invalid host", func(t *testing.T) {
		host := "invalidhost"
		ips, err := LookupSIPAddresses(host)
		assert.Error(t, err)
		assert.Nil(t, ips)
	})

	t.Run("Should return an error for an empty host", func(t *testing.T) {
		host := ""
		ips, err := LookupSIPAddresses(host)
		assert.Error(t, err)
		assert.Nil(t, ips)
	})
}

func Test_LookupSIPAddress(t *testing.T) {
	t.Run("Should return the first IP address as a string for a valid host", func(t *testing.T) {
		host := "example.com"
		ip, err := LookupSIPAddress(host)
		assert.NoError(t, err)
		assert.NotNil(t, ip)
		assert.NotEmpty(t, *ip)
	})

	t.Run("Should return an error for an invalid host", func(t *testing.T) {
		host := "invalidhost"
		ip, err := LookupSIPAddress(host)
		assert.Error(t, err)
		assert.Nil(t, ip)
	})
}

func Test_CheckSIPServerHealth(t *testing.T) {
	t.Run("Should always return true for health check", func(t *testing.T) {
		routingSIPURI := "sip:example.com"
		result, err := CheckSIPServerHealth(routingSIPURI)
		assert.NoError(t, err)
		assert.True(t, result)
	})
}

func Test_ReadEnv(t *testing.T) {
	t.Run("Should return the environment variable if it exists", func(t *testing.T) {
		key := "EXISTING_ENV_VARIABLE"
		expectedValue := "somevalue"
		os.Setenv(key, expectedValue)

		value := ReadEnv(key, "fallback")
		assert.Equal(t, expectedValue, value)
	})

	t.Run("Should return the fallback value if the environment variable does not exist", func(t *testing.T) {
		key := "NON_EXISTING_ENV_VARIABLE"
		fallbackValue := "fallback"

		value := ReadEnv(key, fallbackValue)
		assert.Equal(t, fallbackValue, value)
	})
}

func Test_Config(t *testing.T) {
	t.Run("Should return the value of a valid key from the environment", func(t *testing.T) {
		key := "VALID_KEY"
		expectedValue := "somevalue"
		os.Setenv(key, expectedValue)

		value := Config(key)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("Should return an empty string for an invalid key from the environment", func(t *testing.T) {
		key := "INVALID_KEY"
		os.Unsetenv(key)

		value := Config(key)
		assert.Empty(t, value)
	})

	t.Run("Should return an empty string if .env is not loaded", func(t *testing.T) {
		os.Setenv("USE_DOTENV", "off")

		key := "SOME_KEY"
		value := Config(key)
		assert.Empty(t, value)
	})
}

func Test_CanPlaceAdditionalCalls(t *testing.T) {
	t.Run("Should always return true without errors", func(t *testing.T) {
		result, err := CanPlaceAdditionalCalls()
		assert.NoError(t, err)
		assert.True(t, result)
	})
}

func Test_CheckIfCarrier(t *testing.T) {
	t.Run("Should return true for a valid carrier token", func(t *testing.T) {
		token := "validToken"
		result := CheckIfCarrier(token)
		assert.True(t, result)
	})
}

func Test_SetSetting(t *testing.T) {
	t.Run("Should set the global setting", func(t *testing.T) {
		expectedSetting := model.GlobalSettings{ValidateCallerId: true}
		SetSetting(expectedSetting)
		actualSetting := GetSetting()
		assert.Equal(t, expectedSetting, *actualSetting)
	})
}
