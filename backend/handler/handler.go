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
	SessionID      string `json:"sessionId"`
	WebRTC         string `json:"WebRTC"`
	WebRTCPlayback string `json:"WebRTCPlayback"`
}

type LiveInputResponse struct {
	Result struct {
		UID    string `json:"uid"`
		WebRTC struct {
			URL string `json:"url"`
		} `json:"webRTC"`
		WebRTCPlayback struct {
			URL string `json:"url"`
		} `json:"webRTCPlayback"`
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
		SessionID:      cloudflareResponse.Result.UID,
		WebRTC:         cloudflareResponse.Result.WebRTC.URL,
		WebRTCPlayback: cloudflareResponse.Result.WebRTCPlayback.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentSession)
}

// API để nhận video stream từ client và đẩy lên Cloudflare
func UploadStream(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling UploadStream request")

	if currentSession == nil {
		http.Error(w, "No active streaming session", http.StatusBadRequest)
		fmt.Println("No active streaming session")
		return
	}

	// Struct to decode JSON
	var req struct {
		Offer  string `json:"offer"` // Changed to string to handle SDP
		Tracks []struct {
			TrackName string `json:"trackName"`
			Mid       string `json:"mid"`
			Location  string `json:"location"`
		} `json:"tracks"`
		Force bool `json:"force"`
	}

	// Read and decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Received tracks:", req.Tracks)

	// Prepare track data
	tracks := make([]map[string]interface{}, len(req.Tracks))
	for i, t := range req.Tracks {
		location := t.Location
		if location == "" {
			location = "local"
		}
		tracks[i] = map[string]interface{}{
			"trackName": t.TrackName,
			"location":  location,
			"mid":       t.Mid,
		}
	}
	requestBody := map[string]interface{}{
		"sessionDescription": req.Offer,
		"tracks":             tracks,
	}

	jsonData, _ := json.Marshal(requestBody)

	// Create request body as SDP content
	client := &http.Client{}
	cfReq, err := http.NewRequest("POST", currentSession.WebRTC, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Set correct headers for SDP content
	cfReq.Header.Set("Content-Type", "application/sdp")
	cfReq.Header.Set("Authorization", "Bearer "+os.Getenv("CLOUDFLARE_API_KEY"))

	resp, err := client.Do(cfReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to upload stream", resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
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
