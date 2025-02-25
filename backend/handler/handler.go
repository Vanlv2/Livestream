package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"sort"
	"time"
)

var currentSession *LiveInputSession

type LiveInputSession struct {
	SessionID string `json:"sessionId"`
	IngestURL string `json:"ingestUrl"`
	StreamKey string `json:"streamKey"`
}

type LiveInputResponse struct {
	Result struct {
		UID     string `json:"uid"`
		RTMPS   struct {
			URL       string `json:"url"`
			StreamKey string `json:"streamKey"`
		} `json:"rtmps"`
	} `json:"result"`
	Success bool `json:"success"`
}

type Video struct {
	Status struct {
		State string `json:"state"`
	} `json:"status"`
	Created  time.Time `json:"created"`
	Playback struct {
		HLS string `json:"hls"`
	} `json:"playback"`
	Preview string `json:"preview"`
}

type VideosResponse struct {
	Result  []Video `json:"result"`
	Success bool    `json:"success"`
}

func CreateLiveInput(w http.ResponseWriter, r *http.Request) {

	cloudflareAccountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	cloudflareAPIKey := os.Getenv("CLOUDFLARE_API_KEY")

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/stream/live_inputs", cloudflareAccountID)

	reqBody := `{
		"deleteRecordingAfterDays": 45,
		"meta": { "title": "", "host": "" },
		"recording": { "mode": "automatic" }
	}`

	// Chuyển reqBody thành io.Reader để sử dụng cho http.NewRequest
	reader := bytes.NewReader([]byte(reqBody))

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cloudflareAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error contacting Cloudflare", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var cloudflareResponse LiveInputResponse
	if err := json.NewDecoder(resp.Body).Decode(&cloudflareResponse); err != nil {
		http.Error(w, "Error decoding Cloudflare response", http.StatusInternalServerError)
		return
	}

	if !cloudflareResponse.Success {
		http.Error(w, "Cloudflare API error", http.StatusInternalServerError)
		return
	}

	currentSession = &LiveInputSession{
		SessionID: cloudflareResponse.Result.UID,
		IngestURL: cloudflareResponse.Result.RTMPS.URL,
		StreamKey: cloudflareResponse.Result.RTMPS.StreamKey,
	}
	fmt.Printf("Current session details: %+v\n", currentSession)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentSession)
}

func GetVideos(w http.ResponseWriter, r *http.Request) {
	// Lấy roomId từ query parameters
	sessionId := r.URL.Query().Get("sessionId")
	if sessionId == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	cloudflareAccountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	cloudflareAPIKey := os.Getenv("CLOUDFLARE_API_KEY")

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/stream/live_inputs/%s/videos", cloudflareAccountID, sessionId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+cloudflareAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error contacting Cloudflare", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var videoResponse VideosResponse
	if err := json.NewDecoder(resp.Body).Decode(&videoResponse); err != nil {
		http.Error(w, "Error decoding Cloudflare response", http.StatusInternalServerError)
		return
	}

	if !videoResponse.Success {
		http.Error(w, "Cloudflare API error", http.StatusInternalServerError)
		return
	}

	var liveVideos []Video
	for _, video := range videoResponse.Result {
		if video.Status.State == "live-inprogress" {
			liveVideos = append(liveVideos, video)
		}
	}

	if len(liveVideos) == 0 {
		http.Error(w, "No live video found", http.StatusNotFound)
		return
	}

	// Sắp xếp theo thời gian tạo (mới nhất)
	sort.Slice(liveVideos, func(i, j int) bool {
		return liveVideos[i].Created.After(liveVideos[j].Created)
	})

	activeVideo := liveVideos[0]
	if activeVideo.Playback.HLS == "" {
		http.Error(w, "Active video does not have a valid playback URL", http.StatusInternalServerError)
		return
	}

	// Trả về thông tin playback URL và preview
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"playbackUrl": activeVideo.Playback.HLS, // HLS URL
		"preview":     activeVideo.Preview,      // xem video có sẵn
	})
}


