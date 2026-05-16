package cache

import (
	"fmt"
	"strconv"
)

const (
	RankKey         = "rank"          // 浏览量排行
	RankPurchaseKey = "rank:purchase" // 购买量排行
	RankFavoriteKey = "rank:favorite" // 收藏量排行
	RankHotKey      = "rank:hot"      // 综合热度排行
)

func ProductViewKey(id uint) string {
	return fmt.Sprintf("view:product:%s", strconv.Itoa(int(id)))
}
