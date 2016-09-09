package linebot

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var requestBody = `{
    "events": [
        {
            "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
            "type": "message",
            "timestamp": 1462629479859,
            "source": {
                "type": "user",
                "userId": "u206d25c2ea6bd87c17655609a1c37cb8"
            },
            "message": {
                "id": "325708",
                "type": "text",
                "text": "Hello, world"
            }
        },
        {
            "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
            "type": "message",
            "timestamp": 1462629479859,
            "source": {
                "type": "group",
                "groupId": "u206d25c2ea6bd87c17655609a1c37cb8",
                "userId": "u206d25c2ea6bd87c17655609a1c37cb8"
            },
            "message": {
                "id": "325708",
                "type": "text",
                "text": "Hello, world"
            }
        },
        {
            "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
            "type": "message",
            "timestamp": 1462629479859,
            "source": {
                "type": "user",
                "userId": "u206d25c2ea6bd87c17655609a1c37cb8"
            },
            "message": {
                "id": "325708",
                "type": "image"
            }
        },
        {
            "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
            "type": "message",
            "timestamp": 1462629479859,
            "source": {
                "type": "user",
                "userId": "u206d25c2ea6bd87c17655609a1c37cb8"
            },
            "message": {
                "id": "325708",
                "type": "location",
                "title": "hello",
                "address": "〒150-0002 東京都渋谷区渋谷２丁目２１−１",
                "latitude": 35.65910807942215,
                "longitude": 139.70372892916203
            }
        },
        {
            "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
            "type": "message",
            "timestamp": 1462629479859,
            "source": {
                "type": "user",
                "userId": "u206d25c2ea6bd87c17655609a1c37cb8"
            },
            "message": {
                "id": "325708",
                "type": "sticker",
                "packageId": "1",
                "stickerId": "1"
            }
        },
        {
            "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
            "type": "follow",
            "timestamp": 1462629479859,
            "source": {
                "type": "user",
                "userId": "u206d25c2ea6bd87c17655609a1c37cb8"
            }
        },
        {
            "type": "unfollow",
            "timestamp": 1462629479859,
            "source": {
                "type": "user",
                "userId": "u206d25c2ea6bd87c17655609a1c37cb8"
            }
        },
        {
            "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
            "type": "join",
            "timestamp": 1462629479859,
            "source": {
                "type": "group",
                "groupId": "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
            }
        },
        {
            "type": "leave",
            "timestamp": 1462629479859,
            "source": {
                "type": "group",
                "groupId": "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
            }
        },
        {
            "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
            "type": "postback",
            "timestamp": 1462629479859,
            "source": {
                "type": "user",
                "userId": "u206d25c2ea6bd87c17655609a1c37cb8"
            },
            "postback": {
                "data": "action=buyItem&itemId=123123&color=red",
                "params": {}
            }
        },
        {
            "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
            "type": "beacon",
            "timestamp": 1462629479859,
            "source": {
                "type": "user",
                "userId": "U012345678901234567890123456789ab"
            },
            "beacon": {
                "hwid":"374591320"
            }
        }
    ]
}
`

var wantEvents = []Event{
	{
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
		Type:       EventTypeMessage,
		Timestamp:  1462629479859,
		Source: &EventSource{
			Type:   EventSourceTypeUser,
			UserID: "u206d25c2ea6bd87c17655609a1c37cb8",
		},
		Message: &TextMessage{
			ID:   "325708",
			Text: "Hello, world",
		},
	},
	{
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
		Type:       EventTypeMessage,
		Timestamp:  1462629479859,
		Source: &EventSource{
			Type:    EventSourceTypeGroup,
			UserID:  "u206d25c2ea6bd87c17655609a1c37cb8",
			GroupID: "u206d25c2ea6bd87c17655609a1c37cb8",
		},
		Message: &TextMessage{
			ID:   "325708",
			Text: "Hello, world",
		},
	},
	{
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
		Type:       EventTypeMessage,
		Timestamp:  1462629479859,
		Source: &EventSource{
			Type:   EventSourceTypeUser,
			UserID: "u206d25c2ea6bd87c17655609a1c37cb8",
		},
		Message: &ImageMessage{
			ID: "325708",
		},
	},
	{
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
		Type:       EventTypeMessage,
		Timestamp:  1462629479859,
		Source: &EventSource{
			Type:   EventSourceTypeUser,
			UserID: "u206d25c2ea6bd87c17655609a1c37cb8",
		},
		Message: &LocationMessage{
			ID:        "325708",
			Title:     "hello",
			Address:   "〒150-0002 東京都渋谷区渋谷２丁目２１−１",
			Latitude:  35.65910807942215,
			Longitude: 139.70372892916203,
		},
	},
	{
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
		Type:       EventTypeMessage,
		Timestamp:  1462629479859,
		Source: &EventSource{
			Type:   EventSourceTypeUser,
			UserID: "u206d25c2ea6bd87c17655609a1c37cb8",
		},
		Message: &StickerMessage{
			ID:        "325708",
			PackageID: "1",
			StickerID: "1",
		},
	},
	{
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
		Type:       EventTypeFollow,
		Timestamp:  1462629479859,
		Source: &EventSource{
			Type:   EventSourceTypeUser,
			UserID: "u206d25c2ea6bd87c17655609a1c37cb8",
		},
	},
	{
		Type:      EventTypeUnfollow,
		Timestamp: 1462629479859,
		Source: &EventSource{
			Type:   EventSourceTypeUser,
			UserID: "u206d25c2ea6bd87c17655609a1c37cb8",
		},
	},
	{
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
		Type:       EventTypeJoin,
		Timestamp:  1462629479859,
		Source: &EventSource{
			Type:    EventSourceTypeGroup,
			GroupID: "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		},
	},
	{
		Type:      EventTypeLeave,
		Timestamp: 1462629479859,
		Source: &EventSource{
			Type:    EventSourceTypeGroup,
			GroupID: "cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		},
	},
	{
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
		Type:       EventTypePostback,
		Timestamp:  1462629479859,
		Source: &EventSource{
			Type:   EventSourceTypeUser,
			UserID: "u206d25c2ea6bd87c17655609a1c37cb8",
		},
		// TODO: postback
	},
	{
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
		Type:       EventTypeBeacon,
		Timestamp:  1462629479859,
		Source: &EventSource{
			Type:   EventSourceTypeUser,
			UserID: "U012345678901234567890123456789ab",
		},
		// TODO: beacon
	},
}

func TestParseRequest(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client, err := New("testsecret", "testtoken")
		if err != nil {
			t.Error(err)
			return
		}
		gotEvents, err := client.ParseRequest(r)
		if err != nil {
			if err == ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
				t.Error(err)
			}
			return
		}
		if len(gotEvents) != len(wantEvents) {
			t.Errorf("Event length %d; want %d", len(gotEvents), len(wantEvents))
		}
		for i, got := range gotEvents {
			want := wantEvents[i]
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Event %d %q; want %q", i, got, want)
			}
		}
	}))
	defer server.Close()
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// invalid signature
	{
		body := []byte(requestBody)
		req, err := http.NewRequest("POST", server.URL, bytes.NewReader(body))
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("X-LINE-Signature", "invalidsignatue")
		res, err := httpClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		if res.StatusCode != 400 {
			t.Errorf("StatusCode %d; want %d", res.StatusCode, 400)
		}
	}

	// valid signature
	{
		body := []byte(requestBody)
		req, err := http.NewRequest("POST", server.URL, bytes.NewReader(body))
		if err != nil {
			t.Error(err)
			return
		}
		// generate signature
		mac := hmac.New(sha256.New, []byte("testsecret"))
		mac.Write(body)

		req.Header.Set("X-LINE-Signature", base64.StdEncoding.EncodeToString(mac.Sum(nil)))
		res, err := httpClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		if res == nil {
			t.Error("response is nil")
			return
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("status: %d", res.StatusCode)
			return
		}
	}
}

func BenchmarkParseRequest(b *testing.B) {
	body := []byte(requestBody)
	client, err := New("testsecret", "testtoken")
	if err != nil {
		b.Error(err)
		return
	}
	mac := hmac.New(sha256.New, []byte("testsecret"))
	mac.Write(body)
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "", bytes.NewReader(body))
		req.Header.Set("X-LINE-Signature", sign)
		_, _ = client.ParseRequest(req)
	}
}
