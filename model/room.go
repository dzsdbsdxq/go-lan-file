package model

type RoomOnline struct {
	RoomId     string //房间ID
	RoomTitle  string // 房间标题
	RoomDesc   string // 房间信息
	RoomOwner  string // 房主
	CreateTime uint64
	ExpireTime uint64
}
type RoomContentOnline struct {
	RoomId     string
	RoomText   string
	CreateTime uint64
	UpdateTime uint64
	RoomOwner  string
	ExpireTime uint64
}
