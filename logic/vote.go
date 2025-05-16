package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

//投票功能：
//1.用户投票的数据
//

/*
投票的几种情况：
direction=1时，有两种情况：

	1.之前没有投过票，现在投赞成票
	2.之前投反对票，现在改投赞成票

direction=0时，有两种情况：

	1.之前反对，现在取消
	2.之前赞成，现在取消

投票的限制：
每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许投票了。

	1.到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2.到期之后删除那个KeyPostVotedZSetPF
*/
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
