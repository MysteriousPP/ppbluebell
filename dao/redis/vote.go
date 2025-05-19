package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //一票的分数
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
)

func CreatePost(postID int64) error {
	pipeline := client.TxPipeline()
	//贴子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipeline.Exec()
	return err
}

/*
投一票加432分
投票的几种情况：
direction=1时，有两种情况：

	1.之前没有投过票，现在投赞成票   差值的绝对值：1 +432
	2.之前投反对票，现在改投赞成票   差值的绝对值：2 +432*2

direction=0时，有两种情况：

	1.之前反对，现在取消            差值的绝对值：1 -432
	2.之前赞成，现在取消			差值的绝对值：1	+432

direction=-1时，有两种情况：

	1.之前没有投过票，现在投反对票		差值的绝对值：1	-432
	2.之前投赞成票，现在改投反对票		差值的绝对值：2 -432*2

投票的限制：
每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许投票了。

	1.到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2.到期之后删除那个KeyPostVotedZSetPF
*/
func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制
	// 去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2和3需要放到一个pipeline事务中

	//2.更新贴子的分数
	//先查之前的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZsetPF+postID), userID).Val()

	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value)
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	//3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZsetPF+postID), postID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZsetPF+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err

}
