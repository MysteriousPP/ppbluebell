package redis

import "bluebell/models"

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//1. 从redis获取id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2. 确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	//3. ZREVRANGE 查询
	return client.ZRevRange(key, start, end).Result()
}
