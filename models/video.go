package models

// 视频模型
type Video struct {
	Aid         int64  `json:"aid"`          // 视频aid
	Bvid        string `json:"bvid"`         //视频id
	Title       string `json:"title"`        // 视频标题
	Subtitle    string `json:"subtitle"`     // 副标题
	Bullet      int64  `json:"video_review"` //弹幕数
	Created     int64  `json:"created"`      // 视频创建时间
	Poster      string `json:"poster"`       // 视频封面
	Description string `json:"description"`  // 视频描述
	Duration    string `json:"length"`       // 视频时长
	Typeid      int    `json:"typeid"`
	Comment     uint64 `json:"comment"`   // 视频评论数
	Play        uint64 `json:"play"`      //视频播放量
	Copyright   string `json:"copyright"` // 版权
}
