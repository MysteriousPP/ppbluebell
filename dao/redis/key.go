package redis

//redis key
//redis key尽量使用命名空间的方式，方便查询和拆分

const (
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "post:time"   // zset;帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  // zset;帖子及投票数的分数
	KeyPostVotedZsetPF = "post:voted:" //zset;记录用户及投票的类型;参数是post id
	KeyCommunitySetPF  = "community:"  //set保存分区下帖子的id
)

func getRedisKey(key string) string {
	return Prefix + key
}
