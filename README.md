Xây dựng phòng livestream bằng WHIP/WHEP của Cloudflare, sử dụng MongoDB làm database.

🚀 Tổng Quan Luồng Thực Thi
Streamer (Broadcaster): Gửi video/audio lên Cloudflare bằng WHIP.
Viewer (Người xem): Nhận video/audio từ Cloudflare bằng WHEP.
Backend:
Quản lý phòng livestream (tạo, cập nhật, kết thúc).
Cung cấp API cho streamer và viewer để kết nối với Cloudflare.
Lưu trữ thông tin phòng livestream trong MongoDB.
🔥 Chi Tiết Luồng Xử Lý
1️⃣ Người livestream (Broadcaster) tạo phòng
📌 Backend tạo session trên Cloudflare & lưu vào MongoDB.

API Backend - Tạo phòng livestream
Request
POST /api/livestream/create
Content-Type: application/json

{
  "title": "Lớp học online",
  "host": "user_123"
}
Luồng xử lý Backend
Gọi API Cloudflare để tạo WHIP session.
Lưu thông tin room_id, session_id, whip_url vào MongoDB.
Trả về whip_url để streamer gửi dữ liệu.
Gọi API Cloudflare - Tạo WHIP Session
curl https://api.cloudflare.com/client/v4/accounts/$ACCOUNT_ID/stream/live_inputs \
    -H 'Content-Type: application/json' \
    -H "X-Auth-Email: $CLOUDFLARE_EMAIL" \
    -H "X-Auth-Key: $CLOUDFLARE_API_KEY" \
    -d '{
      "deleteRecordingAfterDays": 45,
      "meta": {
        "name": "test stream 1"
      },
      "recording": {}
    }'
Response
{
  "id": "session_abcxyz",
  "metadata": {
    "title": "Lớp học online"
  },
  "whip_url": "https://customer-m033z5x00ks6nunl.cloudflarestream.com/b236bde30eb07b9d3ad/webRTC/publish"
}
Lưu vào MongoDB
db.livestreams.insert_one({
    "room_id": "live_123456",
    "session_id": "session_abcxyz",
    "whip_url": "https://cloudflare.com/whip/session_abcxyz",
    "host": "user_123",
    "status": "active",
    "created_at": datetime.utcnow()
})
✅ Trả về whip_url từ đó có thể đẩy video/audio của streamer lên.

2️⃣ Người livestream gửi dữ liệu lên Cloudflare
📌 Streamer sau khi tạo phòng đã có whip_url để bắt đầu livestream.

Cloudflare xử lý và phân phối stream đến người xem.
✅ Không cần backend xử lý dữ liệu video, Cloudflare lo phần này.

3️⃣ Người xem (Viewer) tham gia phòng
📌 Backend cung cấp WHEP URL để viewer kết nối.

API Backend - Lấy URL WHEP
Request
GET /api/livestream/join?room_id=live_123456
Luồng xử lý Backend
Lấy session_id từ MongoDB.
Gọi API Cloudflare để lấy WHEP URL.
Trả về whep_url cho viewer.
Tìm session trong MongoDB
session = db.livestreams.find_one({"room_id": "live_123456"})
session_id = session["session_id"]
Gọi API Cloudflare - Lấy WHEP URL
curl -X GET "https://api.cloudflare.com/client/v4/accounts/{account_id}/calls/sessions/session_abcxyz/whep" \
     -H "Authorization: Bearer YOUR_API_KEY"
Response
{
  "whep_url": "https://cloudflare.com/whep/session_abcxyz"
}
✅ Trả về whep_url để viewer kết nối & xem livestream.

4️⃣ Viewer sử dụng WebRTC để xem livestream
📌 Frontend kết nối đến whep_url bằng WebRTC API.

Code Client-side WebRTC
const videoElement = document.getElementById("video");
const peerConnection = new RTCPeerConnection();
const videoTrack = new MediaStream();

// Lấy WHEP URL từ API backend
fetch("/api/livestream/join?room_id=live_123456")
  .then(response => response.json())
  .then(data => {
      const offer = data.sdp;
      peerConnection.setRemoteDescription(new RTCSessionDescription({type: "offer", sdp: offer}));
      peerConnection.createAnswer().then(answer => peerConnection.setLocalDescription(answer));
  });

peerConnection.ontrack = (event) => {
    videoTrack.addTrack(event.track);
    videoElement.srcObject = videoTrack;
};
✅ Viewer có thể xem livestream trực tiếp trên trình duyệt.

5️⃣ Streamer kết thúc livestream
📌 Streamer dừng stream → Backend xóa session trên Cloudflare & MongoDB.

API Backend - Kết thúc livestream
Request
POST /api/livestream/end
Content-Type: application/json

{
  "room_id": "live_123456"
}
Luồng xử lý Backend
Lấy session_id từ MongoDB.
Gọi API Cloudflare để xóa session.
Cập nhật MongoDB trạng thái "ended".
Gọi API Cloudflare - Xóa Session
curl -X DELETE "https://api.cloudflare.com/client/v4/accounts/{account_id}/calls/sessions/session_abcxyz" \
     -H "Authorization: Bearer YOUR_API_KEY"
Cập nhật MongoDB
db.livestreams.update_one(
    {"room_id": "live_123456"},
    {"$set": {"status": "ended"}}
)
✅ Livestream kết thúc, viewer không thể xem nữa.

🔥 Tóm tắt nhiệm vụ của Backend
Chức năng	API	Nhiệm vụ Backend
Tạo phòng livestream	POST /api/livestream/create	Tạo session WHIP, lưu MongoDB
Người livestream gửi dữ liệu	Thực hiện ở phía front end	Không cần backend xử lý
Người xem tham gia	GET /api/livestream/join	Lấy session từ MongoDB, gọi API lấy WHEP URL
Viewer kết nối WebRTC	(Frontend dùng WHEP URL)	Không cần backend xử lý
Kết thúc livestream	POST /api/livestream/end	Xóa session trên Cloudflare, cập nhật MongoDB
