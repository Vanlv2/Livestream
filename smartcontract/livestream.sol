// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract LivestreamRoom {
    struct Livestream {
        string roomId;
        string sessionId;
        string ingestUrl;
        string streamKey;
        string host;
        string title;
        string status; // "active", "ended", etc.
        uint256 createdAt;
    }

    mapping(string => Livestream) public livestreams;

    event LivestreamCreated(
        string roomId,
        string sessionId,
        string ingestUrl,
        string streamKey,
        string title,
        string host
    );
    event TitleUpdated(string roomId, string newTitle);
    event LivestreamEnded(string roomId, string sessionId);

    // Tạo phòng livestream mới (sẽ được gọi bởi streamer)
    function createLivestream(
        string memory _sessionId,
        string memory _ingestUrl,
        string memory _streamKey,
        string memory _host,
        string memory _title
    ) public {
        // Tạo roomId bằng cách băm địa chỉ của người gọi, sessionId và block.timestamp
        bytes32 roomIdHash = keccak256(
            abi.encodePacked(msg.sender, _sessionId, block.timestamp)
        );
        string memory _roomId = toString(roomIdHash);

        // Kiểm tra xem roomId đã tồn tại hay chưa
        require(
            bytes(livestreams[_roomId].roomId).length == 0,
            "Room already exists"
        );

        Livestream memory newLivestream = Livestream({
            roomId: _roomId,
            sessionId: _sessionId,
            ingestUrl: _ingestUrl,
            streamKey: _streamKey,
            host: _host,
            title: _title,
            status: "active",
            createdAt: block.timestamp
        });

        livestreams[_roomId] = newLivestream;
        emit LivestreamCreated(
            _roomId,
            _sessionId,
            _ingestUrl,
            _streamKey,
            _title,
            _host
        );
    }

    // Lấy thông tin phòng livestream (Viewer gọi để lấy session và URL)
    function getLivestreamInfo(string memory _roomId)
        public
        view
        returns (
            string memory sessionId,
            string memory ingestUrl,
            string memory streamKey,
            string memory host,
            string memory title,
            string memory status,
            uint256 createdAt
        )
    {
        Livestream storage livestream = livestreams[_roomId];
        require(bytes(livestream.roomId).length != 0, "Room does not exist");
        return (
            livestream.sessionId,
            livestream.ingestUrl,
            livestream.streamKey,
            livestream.host,
            livestream.title,
            livestream.status,
            livestream.createdAt
        );
    }

    function updateTitle(string memory _roomId, string memory _newTitle)
        public
    {
        require(
            bytes(livestreams[_roomId].roomId).length != 0,
            "Livestream does not exist"
        );

        livestreams[_roomId].title = _newTitle;

        emit TitleUpdated(_roomId, _newTitle);
    }

    // Kết thúc livestream (Streamer gọi khi kết thúc phiên livestream)
    function endLivestream(string memory _roomId) public {
        Livestream storage livestream = livestreams[_roomId];
        require(bytes(livestream.roomId).length != 0, "Room does not exist");
        require(
            keccak256(bytes(livestream.status)) == keccak256(bytes("active")),
            "Livestream already ended"
        );

        livestream.status = "ended";
        emit LivestreamEnded(_roomId, livestream.sessionId);
    }

    function toString(bytes32 data) internal pure returns (string memory) {
        // Chuyển đổi bytes32 sang chuỗi string
        bytes memory alphabet = "0123456789abcdef";
        bytes memory str = new bytes(64);
        for (uint256 i = 0; i < 32; i++) {
            str[i * 2] = alphabet[uint256(uint8(data[i] >> 4))];
            str[1 + i * 2] = alphabet[uint256(uint8(data[i] & 0x0f))];
        }
        return string(str);
    }
}
