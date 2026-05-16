package v1

import (
	"gin_mall_tmp/serializer"
	aiservice "gin_mall_tmp/service/ai"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SyncKnowledge syncs policy documents and products into the vector database.
func SyncKnowledge(c *gin.Context) {
	var req aiservice.KnowledgeSyncRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	policyCount, productCount, totalChunks, err := aiservice.SyncKnowledgeBase(c.Request.Context(), req.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "同步知识库失败",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, serializer.Response{
		Status: http.StatusOK,
		Msg:    "知识库同步完成",
		Data: aiservice.KnowledgeSyncPayload{
			PolicyCount:  policyCount,
			ProductCount: productCount,
			TotalChunks:  totalChunks,
		},
	})
}

// SearchKnowledge performs vector search on the knowledge base.
func SearchKnowledge(c *gin.Context) {
	var req aiservice.KnowledgeSearchRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 5
	}

	results, err := aiservice.RetrieveKnowledge(c.Request.Context(), req.Query, limit, req.Source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "知识库搜索失败",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, serializer.Response{
		Status: http.StatusOK,
		Msg:    "ok",
		Data:   results,
	})
}

// GetKnowledgeStats returns the total number of chunks in the vector database.
func GetKnowledgeStats(c *gin.Context) {
	count, err := aiservice.GetChunkCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "获取知识库统计失败",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, serializer.Response{
		Status: http.StatusOK,
		Msg:    "ok",
		Data: gin.H{
			"total_chunks": count,
		},
	})
}
