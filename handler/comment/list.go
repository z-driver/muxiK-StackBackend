package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

type commentListResponse struct {
	ParentCommentSum  uint32                     `json:"parent_comment_sum"`
	ParentCommentList *[]model.ParentCommentInfo `json:"parent_comment_list"`
}

// 获取评论列表
// @Summary 获取评论列表
// @Tags comment
// @Param token header string false "游客登录则不需要此字段或为空"
// @Param id path string true "评课id"
// @Param limit query integer true "最大的一级评论数量"
// @Param pageNum query integer true "翻页页码，默认为0"
// @Success 200 {object} comment.commentListResponse
// @Router /evaluation/{id}/comments/ [get]
func GetComments(c *gin.Context) {
	log.Info("GetComments function is called.")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
	}

	size := c.DefaultQuery("limit", "20")
	limit, err := strconv.ParseInt(size, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	pageNum := c.DefaultQuery("pageNum", "0")
	num, err := strconv.ParseInt(pageNum, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	// userId获取与游客模式判断
	var userId uint32
	visitor := false

	userIdInterface, ok := c.Get("id")
	if !ok {
		visitor = true
	} else {
		userId = userIdInterface.(uint32)
		log.Info("User auth successful.")
	}

	list, count, err := service.CommentList(uint32(id), int32(limit), int32(num*limit), userId, visitor)
	if err != nil {
		handler.SendError(c, errno.ErrCommentList, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, commentListResponse{
		ParentCommentSum:  count,
		ParentCommentList: list,
	})
}