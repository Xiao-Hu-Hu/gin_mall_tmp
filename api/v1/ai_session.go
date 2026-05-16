package v1

import (
	"gin_mall_tmp/cache"
	"gin_mall_tmp/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetAISessions 获取用户的 AI 会话列表
func GetAISessions(c *gin.Context) {
	token := util.ExtractTokenFromRequest(c.Request)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "未提供令牌",
		})
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌校验失败",
			"error":  err.Error(),
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	sessions, err := cache.GetUserSessions(claims.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "获取会话列表失败",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"msg":      "ok",
		"sessions": sessions,
	})
}

// GetAISessionHistory 获取指定会话的聊天历史
func GetAISessionHistory(c *gin.Context) {
	token := util.ExtractTokenFromRequest(c.Request)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "未提供令牌"})
		return
	}
	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "令牌校验失败"})
		return
	}
	_ = claims // verify token ownership

	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "会话 ID 不能为空"})
		return
	}

	session, err := cache.GetAISession(sessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "msg": "会话不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"msg":      "ok",
		"session":  session,
		"history":  session.History,
	})
}

// DeleteAISession 删除指定会话
func DeleteAISession(c *gin.Context) {
	token := util.ExtractTokenFromRequest(c.Request)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "未提供令牌",
		})
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "令牌校验失败",
			"error":  err.Error(),
		})
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "会话 ID 不能为空",
		})
		return
	}

	err = cache.DeleteAISession(sessionID, claims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "删除会话失败",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "会话已删除",
	})
}
