package service

import (
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

type ParentCommentInfoList struct {
	Lock  *sync.Mutex
	IdMap map[string]*model.ParentCommentInfo
}

type SubCommentInfoList struct {
	Lock  *sync.Mutex
	IdMap map[string]*model.CommentInfo
}

// Get comment list.
func CommentList(evaluationId uint32, limit, offset int32, userId uint32, visitor bool) (*[]model.ParentCommentInfo, uint32, error) {
	// Get parent comments from database
	parentComments, count, err := model.GetParentComments(evaluationId, limit, offset)
	if err != nil {
		return nil, count, err
	}

	var parentIds []string
	for _, parentComment := range *parentComments {
		parentIds = append(parentIds, parentComment.Id)
	}

	parentCommentInfoList := ParentCommentInfoList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[string]*model.ParentCommentInfo, len(*parentComments)),
	}

	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	for _, parentComment := range *parentComments {
		wg.Add(1)
		go func(parentComment *model.ParentCommentModel) {
			defer wg.Done()

			// 获取父评论详情
			parentCommentInfo, err := GetParentCommentInfo(parentComment.Id, userId, visitor)
			if err != nil {
				errChan <- err
				return
			}

			parentCommentInfoList.Lock.Lock()
			defer parentCommentInfoList.Lock.Unlock()

			parentCommentInfoList.IdMap[parentCommentInfo.Id] = parentCommentInfo

		}(&parentComment)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, count, err
	}

	var infos []model.ParentCommentInfo
	for _, id := range parentIds {
		infos = append(infos, *parentCommentInfoList.IdMap[id])
	}

	return &infos, count, nil
}

// Get the response data information of a parent comment.
func GetParentCommentInfo(id string, userId uint32, visitor bool) (*model.ParentCommentInfo, error) {
	// Get comment from Database
	comment := &model.ParentCommentModel{Id: id}
	if err := comment.GetById(); err != nil {
		return nil, err
	}

	// Get the user of the parent comment
	userInfo, err := GetUserInfoById(comment.UserId)
	if err != nil {
		return nil, err
	}

	// Get like state
	var isLike = false
	if !visitor {
		isLike = model.CommentHasLiked(userId, comment.Id)
	}

	// Get subComments' infos
	subCommentInfos, err := GetSubCommentInfosByParentId(comment.Id, userId, visitor)
	if err != nil {
		return nil, err
	}

	data := &model.ParentCommentInfo{
		Id:              comment.Id,
		Content:         comment.Content,
		LikeNum:         comment.LikeNum,
		IsLike:          isLike,
		Time:            comment.Time,
		IsAnonymous:     comment.IsAnonymous,
		UserInfo:        userInfo,
		SubCommentsNum:  comment.SubCommentNum,
		SubCommentsList: subCommentInfos,
	}

	return data, nil
}

// Get subComments' infos by parent id.
func GetSubCommentInfosByParentId(id string, userId uint32, visitor bool) (*[]model.CommentInfo, error) {
	// Get subComments from Database
	comments, err := model.GetSubCommentsByParentId(id)
	if err != nil {
		return nil, err
	}

	var commentIds []string
	for _, comment := range *comments {
		commentIds = append(commentIds, comment.Id)
	}

	subCommentInfoList := SubCommentInfoList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[string]*model.CommentInfo, len(*comments)),
	}

	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// 并发获取子评论详情列表
	for _, comment := range *comments {
		wg.Add(1)

		go func(comment *model.SubCommentModel) {
			defer wg.Done()

			// Get a subComment's info by its id
			info, err := GetSubCommentInfoById(comment.Id, userId, visitor)
			if err != nil {
				errChan <- err
				return
			}

			subCommentInfoList.Lock.Lock()
			defer subCommentInfoList.Lock.Unlock()

			subCommentInfoList.IdMap[info.Id] = info

		}(&comment)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, err
	}

	var commentInfos []model.CommentInfo
	for _, id := range commentIds {
		commentInfos = append(commentInfos, *subCommentInfoList.IdMap[id])
	}

	return &commentInfos, nil
}

// Get the response information of a subComment by id.
func GetSubCommentInfoById(id string, userId uint32, visitor bool) (*model.CommentInfo, error) {
	// Get comment from Database
	comment := &model.SubCommentModel{Id: id}
	if err := comment.GetById(); err != nil {
		return nil, err
	}

	// Get the user of the subComment
	commentUser, err := GetUserInfoById(comment.UserId)
	if err != nil {
		return nil, err
	}

	// Get the target user of the subComment
	targetUser, err := GetUserInfoById(comment.TargetUserId)
	if err != nil {
		return nil, err
	}

	// Get like state
	var isLike = false
	if !visitor {
		isLike = model.CommentHasLiked(userId, comment.Id)
	}

	data := &model.CommentInfo{
		Id:             comment.Id,
		Content:        comment.Content,
		LikeNum:        comment.LikeNum,
		IsLike:         isLike,
		Time:           comment.Time,
		UserInfo:       commentUser,
		TargetUserInfo: targetUser,
	}

	return data, nil
}

// Update liked number of a comment after liking or canceling it.
func UpdateCommentLikeNum(commentId string, num int) (uint32, error) {
	subComment, ok := model.IsSubComment(commentId)
	if ok {
		err := subComment.UpdateLikeNum(num)
		return subComment.LikeNum, err
	}

	parentComment := &model.ParentCommentModel{Id: commentId}
	if err := parentComment.GetById(); err != nil {
		return 0, err
	}

	err := parentComment.UpdateLikeNum(num)
	return parentComment.LikeNum, err
}