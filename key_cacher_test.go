package auth0

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gopkg.in/square/go-jose.v2"
)

func TestGet(t *testing.T) {
	entry := keyCacherEntry{time.Now(), jose.JSONWebKey{KeyID: "test1"}}
	m := make(map[string]keyCacherEntry)
	m["key1"] = entry

	tests := []struct {
		name             string
		mkc              *memoryKeyCacher
		key              string
		expectedErrorMsg string
	}{
		{
			name: "pass - persistent cacher",
			mkc: &memoryKeyCacher{
				entries: m,
				maxAge:  time.Duration(-1),
				size:    -1,
			},
			key:              "key1",
			expectedErrorMsg: "",
		},
		{
			name: "fail - invalid key",
			mkc: &memoryKeyCacher{
				entries: m,
				maxAge:  time.Duration(-1),
				size:    -1,
			},
			key:              "invalid key",
			expectedErrorMsg: "no Keys has been found",
		},
		{
			name: "pass - get key for persistent cacher",
			mkc: &memoryKeyCacher{
				entries: m,
				maxAge:  time.Duration(0),
				size:    -1,
			},
			key:              "key1",
			expectedErrorMsg: "",
		},
		{
			name: "fail - no cacher with -1 maxAge",
			mkc: &memoryKeyCacher{
				entries: nil,
				maxAge:  time.Duration(-1),
				size:    0,
			},
			key:              "key1",
			expectedErrorMsg: "no Keys has been found",
		},
		{
			name: "fail - no cacher",
			mkc: &memoryKeyCacher{
				entries: nil,
				maxAge:  time.Duration(0),
				size:    0,
			},
			key:              "key1",
			expectedErrorMsg: "no Keys has been found",
		},
		{
			name: "pass - custom cacher not expired",
			mkc: &memoryKeyCacher{
				entries: m,
				maxAge:  time.Duration(100) * time.Second,
				size:    1,
			},
			key:              "key1",
			expectedErrorMsg: "",
		},
		{
			name: "fail - custom cacher with expired key",
			mkc: &memoryKeyCacher{
				entries: m,
				maxAge:  time.Duration(-100) * time.Second,
				size:    1,
			},
			key:              "key1",
			expectedErrorMsg: "key exists but is expired",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// test.mkc.Add("test1", downloadedKeys)
			_, err := test.mkc.Get(test.key)

			if test.expectedErrorMsg != "" {
				if err == nil {
					t.Errorf("Validation should have failed with error with substring: " + test.expectedErrorMsg)
				} else if !strings.Contains(err.Error(), test.expectedErrorMsg) {
					t.Errorf("Validation should have failed with error with substring: " + test.expectedErrorMsg + ", but got: " + err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Validation should not have failed with error, but got: " + err.Error())
				}
			}
		})
	}
}

func TestAdd(t *testing.T) {
	downloadedKeys := []jose.JSONWebKey{{KeyID: "test1"}, {KeyID: "test2"}, {KeyID: "test3"}}

	tests := []struct {
		name             string
		mkc              *memoryKeyCacher
		addingKey        string
		gettingKey       string
		expectedFoundKey bool
		expectedErrorMsg string
	}{
		{
			name: "pass - persistent cacher",
			mkc: &memoryKeyCacher{
				entries: make(map[string]keyCacherEntry),
				maxAge:  time.Duration(-1),
				size:    -1,
			},
			addingKey:        "test1",
			gettingKey:       "test1",
			expectedFoundKey: true,
			expectedErrorMsg: "",
		},
		{
			name: "fail - invalid key",
			mkc: &memoryKeyCacher{
				entries: make(map[string]keyCacherEntry),
				maxAge:  time.Duration(-1),
				size:    -1,
			},
			addingKey:        "invalid key",
			gettingKey:       "invalid key",
			expectedFoundKey: false,
			expectedErrorMsg: "no Keys has been found",
		},
		{
			name: "pass - add key for persistent cacher",
			mkc: &memoryKeyCacher{
				entries: make(map[string]keyCacherEntry),
				maxAge:  time.Duration(0),
				size:    -1,
			},
			addingKey:        "test1",
			gettingKey:       "test1",
			expectedFoundKey: true,
			expectedErrorMsg: "",
		},
		{
			name: "fail - no cacher",
			mkc: &memoryKeyCacher{
				entries: make(map[string]keyCacherEntry),
				maxAge:  time.Duration(0),
				size:    0,
			},
			addingKey:        "test1",
			gettingKey:       "test1",
			expectedFoundKey: false,
			expectedErrorMsg: "",
		},
		{
			name: "pass - custom cacher get latest added key",
			mkc: &memoryKeyCacher{
				entries: make(map[string]keyCacherEntry),
				maxAge:  time.Duration(100) * time.Second,
				size:    1,
			},
			gettingKey:       "test3",
			expectedFoundKey: true,
			expectedErrorMsg: "",
		},
		{
			name: "fail - custom cacher add invalid key",
			mkc: &memoryKeyCacher{
				entries: make(map[string]keyCacherEntry),
				maxAge:  time.Duration(100) * time.Second,
				size:    1,
			},
			addingKey:        "invalid key",
			gettingKey:       "test1",
			expectedFoundKey: false,
			expectedErrorMsg: "no Keys has been found",
		},
		{
			name: "fail - custom cacher get key not in cache",
			mkc: &memoryKeyCacher{
				entries: make(map[string]keyCacherEntry),
				maxAge:  time.Duration(100) * time.Second,
				size:    1,
			},
			gettingKey:       "test1",
			expectedFoundKey: false,
			expectedErrorMsg: "",
		},
		{
			name: "pass - custom cacher with capacity 3",
			mkc: &memoryKeyCacher{
				entries: make(map[string]keyCacherEntry),
				maxAge:  time.Duration(100) * time.Second,
				size:    3,
			},
			gettingKey:       "test2",
			expectedFoundKey: true,
			expectedErrorMsg: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			if test.addingKey == "" {
				for i := 0; i < 3; i++ {
					_, err = test.mkc.Add(downloadedKeys[i].KeyID, downloadedKeys)
				}
			} else {
				_, err = test.mkc.Add(test.addingKey, downloadedKeys)
			}
			_, ok := test.mkc.entries[test.gettingKey]
			assert.Equal(t, test.expectedFoundKey, ok)

			if test.expectedErrorMsg != "" {
				if err == nil {
					t.Errorf("Validation should have failed with error with substring: " + test.expectedErrorMsg)
				} else if !strings.Contains(err.Error(), test.expectedErrorMsg) {
					t.Errorf("Validation should have failed with error with substring: " + test.expectedErrorMsg + ", but got: " + err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Validation should not have failed with error, but got: " + err.Error())
				}
			}
		})
	}
}

func TestIsExpired(t *testing.T) {
	tests := []struct {
		name         string
		mkc          *memoryKeyCacher
		sleepingTime int
		expectedBool bool
	}{
		{
			name: "true - key is expired",
			mkc: &memoryKeyCacher{
				entries: map[string]keyCacherEntry{},
				maxAge:  time.Duration(1) * time.Second,
				size:    1,
			},
			sleepingTime: 2,
			expectedBool: true,
		},
		{
			name: "false - key not expired",
			mkc: &memoryKeyCacher{
				entries: map[string]keyCacherEntry{},
				maxAge:  time.Duration(2) * time.Second,
				size:    1,
			},
			sleepingTime: 1,
			expectedBool: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mkc.entries["test1"] = keyCacherEntry{time.Now(), jose.JSONWebKey{KeyID: "test1"}}
			time.Sleep(time.Duration(test.sleepingTime) * time.Second)
			if isExpired(test.mkc, "test1") != test.expectedBool {
				t.Errorf("Should have been " + strconv.FormatBool(test.expectedBool) + " but got different")
			}
		})
	}
}

func TestHandleOverflow(t *testing.T) {
	downloadedKeys := []jose.JSONWebKey{{KeyID: "test1"}, {KeyID: "test2"}, {KeyID: "test3"}}

	tests := []struct {
		name           string
		mkc            *memoryKeyCacher
		expectedLength int
	}{
		{
			name: "true - overflowed and delete 1 key",
			mkc: &memoryKeyCacher{
				entries: map[string]keyCacherEntry{},
				maxAge:  time.Duration(2) * time.Second,
				size:    1,
			},
			expectedLength: 1,
		},
		{
			name: "false - no overflow",
			mkc: &memoryKeyCacher{
				entries: map[string]keyCacherEntry{},
				maxAge:  time.Duration(2) * time.Second,
				size:    2,
			},
			expectedLength: 2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mkc.entries["first"] = keyCacherEntry{JSONWebKey: downloadedKeys[0]}
			test.mkc.entries["second"] = keyCacherEntry{JSONWebKey: downloadedKeys[1]}
			handleOverflow(test.mkc)
			if len(test.mkc.entries) != test.expectedLength {
				t.Errorf("Should have been " + strconv.Itoa(test.expectedLength) + "but got different")
			}
		})
	}
}
