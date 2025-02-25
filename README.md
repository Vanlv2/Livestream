Xây dựng **phòng livestream bằng của Cloudflare stream**, sử dụng **samrtcontract**.

---

## 🚀 **Tổng Quan Luồng Thực Thi**
- **Streamer (Broadcaster)**: Gửi yêu cầu create live lên Cloudflare.
- **Cloudflare**: Trả **Ingest URL** và **Stream Key**.
- **OBS** ( hoặc webRTC): Truyền **Ingest URL** và **Stream Key** vào OBS để gửi video/audio lên cloudfalre.
- **Viewer (Người xem)**: Lấy video/audio từ Cloudflare bằng **playback**.
- **Backend**:
  - Cung cấp API cho streamer và viewer để kết nối với Cloudflare.  
- **SmartContract**:  
  - Quản lý phòng livestream (tạo, cập nhật, kết thúc).  
  - Lưu trữ thông tin phòng livestream  

---

# 🔥 **Chi Tiết Luồng Xử Lý**
## 1️⃣ **Người livestream (Broadcaster) tạo phòng**
📌 **Backend tạo session trên Cloudflare**.  

### **API Backend - Tạo phòng livestream**
#### **Request**
```http
POST /api/create
Content-Type: application/json
```
#### **Luồng xử lý Backend**
##### **Create live**
1. Gọi API Cloudflare để tạo **Ingest URL** và **Stream Key**.  
2. Trả thông tin **session_id, Ingest URL** và **Stream Key** cho client.    

#### **Gọi API Cloudflare - Tạo Session**
```bash
curl https://api.cloudflare.com/client/v4/accounts/$ACCOUNT_ID/stream/live_inputs \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $CLOUDFLARE_API_KEY" \
    -d '{
      "deleteRecordingAfterDays": 45,
      "meta": {
        "name": "test stream 1"
      },
      "recording": {}
    }'
```
#### **Response**
```json
{
  "SessionID":"21b5db83fd217aac82549c4d4799aeb0"
  "IngestURL":"rtmps://live.cloudflare.com:443/live/"
  "StreamKey":"25f4d723ad441933a4c15f259f45febck21b5db83fd217aac82549c4d4799aeb0"
}
```
#### **Lưu trữ Smartcontrac**
```

{
	"0": "string: sessionId 0b8f3d17b433bc32577ea9b92485ecc4",
	"1": "string: ingestUrl rtmps://live.cloudflare.com:443/live/",
	"2": "string: streamKey 664fa71b4e5b009baab737017bb7fe8fk0b8f3d17b433bc32577ea9b92485ecc4",
	"3": "string: host 42134",
	"4": "string: title uytu",
	"5": "string: status ended",
	"6": "uint256: createdAt 1740459840"
}
```

✅ **Trả về `IngestURL`  và `StreamKey` từ đó có thể đẩy video/audio của streamer lên**.  

---

## 2️⃣ **Người livestream gửi dữ liệu lên Cloudflare**
📌 **Streamer sau khi tạo phòng đã có `IngestURL`  và `StreamKey` để bắt đầu livestream**.

   **Sử dụng OBS để livestream**
- Cloudflare xử lý và phân phối stream đến người xem.  


✅ **Không cần backend xử lý dữ liệu video, Cloudflare lo phần này.**  

---

## 3️⃣ **Người xem (Viewer) tham gia phòng**
📌 **Cung cấp RoomID để client kết nối smartcontract**.  

### Client lấy SessionID trên Smartcontrac thông qua RoomId
### **API Backend - Lấy URL playback**
#### **Request**
```http
GET  http://api/livestream?sessionId=21b5db83fd217aac82549c4d4799aeb0
```
#### **Luồng xử lý Backend**  
1. Gọi API Cloudflare để lấy **playback**.  
3. Trả về **playback** cho client.  

#### **Gọi API Cloudflare - Lấy playback URL**
```bash
curl -X GET "https://api.cloudflare.com/client/v4/accounts/{AccountID/stream/live_inputs/{sessionId}/videos" \
     -H "Authorization: Bearer YOUR_API_KEY"
```
#### **Response**
```json
 "playback": {
                "hls": "https://customer-o3h7qhxkx89jhu3i.cloudflarestream.com/19ef3a1efb5c5d6719af43bc51515df7/manifest/video.m3u8",
                "dash": "https://customer-o3h7qhxkx89jhu3i.cloudflarestream.com/19ef3a1efb5c5d6719af43bc51515df7/manifest/video.mpd"
            },
```
✅ **Trả về `playback` để client kết nối & xem livestream.**  

---

## 4️⃣ **Streamer kết thúc livestream**
📌 **Streamer dừng stream → clinet cập nhật status**.  

#### **Cập nhật Satmartcontract status ended**

✅ **Livestream kết thúc, viewer không thể xem nữa.**  

---

# 🔥 **Tóm tắt nhiệm vụ của Backend**
| **Chức năng** | **API** | **Nhiệm vụ Backend** |
|--------------|--------|-----------------|
| **Tạo phòng livestream** | `POST /api/create` | Gọi API Cloudflare để tạo session, lưu thông tin SMC |
| **Người livestream gửi dữ liệu** | Frontend sử dụng OBS/WebRTC| Không cần backend xử lý |
| **Người xem tham gia** | `GET /api/livestream` | Lấy session từ SMC, gọi API lấy playback URL |
| **Viewer kết nối OBS (WebRTC)** | **(Frontend dùng playback URL)** | Không cần backend xử lý |
| **Kết thúc livestream** | frontend xoá trên SMC| cập nhật trạng thái SMC |

---
