package utils

// 消息类型
const (
	MESSAGE_TYPE_USER  = iota + 1 // 单聊
	MESSAGE_TYPE_GROUP            // 群聊
)

// 消息内容类型
const (
	TEXT       = iota + 1 // 文字
	FILE                  // 普通文件
	IMAGE                 // 图片
	AUDIO                 // 音频
	VIDEO                 // 视频
	AUDIO_CHAT            // 语音聊天
	VIDEO_CHAT            // 视频聊天
)
