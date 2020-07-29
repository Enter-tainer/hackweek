package controller

import (
	"net/http"
	"tree-hole/model"
	"tree-hole/util"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type responsePostGetAllSingle struct {
	ID        string `json:"_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Title     string `json:"title"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
}

type responsePostGetAll struct {
	PostsInfo []responsePostGetAllSingle `json:"posts_info"`
}

func PostGetAll(context echo.Context) error {
	_ = util.MustGetIDFromContext(context)
	result := responsePostGetAll{PostsInfo: make([]responsePostGetAllSingle, 0)}
	posts, err := model.GetAllPost()
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	for _, v := range posts {
		result.PostsInfo = append(result.PostsInfo, responsePostGetAllSingle{ID: v.ID.Hex(), CreatedAt: v.CreatedAt.Unix(), UpdatedAt: v.UpdatedAt.Unix(), Title: v.Title, Content: v.Content, UserID: util.UserIDHash(v.Salt, v.User.Hex())})
	}
	return util.SuccessResponse(context, http.StatusOK, result)
}

type responsePostGetWithIDReply struct {
	ID        string `json:"_id"`
	UserID    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Content   string `json:"content"`
}

type responsePostGetWithID struct {
	ID        string                       `json:"_id"`
	CreatedAt int64                        `json:"created_at"`
	UpdatedAt int64                        `json:"updated_at"`
	Title     string                       `json:"title"`
	UserID    string                       `json:"user_id"`
	Content   string                       `json:"content"`
	Reply     []responsePostGetWithIDReply `json:"reply"`
}

func PostGetWithId(context echo.Context) error {
	_ = util.MustGetIDFromContext(context)
	id := context.Param("id")
	post, found, err := model.GetPostWithID(id)
	if !found {
		return util.ErrorResponse(context, http.StatusBadRequest, "post not found")
	}
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	var reply []responsePostGetWithIDReply
	for _, v := range post.Reply {
		reply = append(reply, responsePostGetWithIDReply{
			Content:   v.Content,
			UpdatedAt: v.UpdatedAt.Unix(),
			CreatedAt: v.CreatedAt.Unix(),
			UserID:    util.UserIDHash(post.Salt, v.User.Hex()),
			ID:        v.ID.Hex(),
		})
	}
	return util.SuccessResponse(context, http.StatusOK, responsePostGetWithID{
		Title:     post.Title,
		ID:        post.ID.Hex(),
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Unix(),
		UpdatedAt: post.UpdatedAt.Unix(),
		UserID:    util.UserIDHash(post.Salt, post.User.Hex()),
		Reply:     reply,
	})
}

type paramPostNew struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type responsePostNew struct {
	PostID string `json:"_id"`
}

func PostNew(context echo.Context) error {
	var param paramPostNew
	if err := context.Bind(&param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}
	id := util.MustGetIDFromContext(context)
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	postID, err := model.AddPost(model.Post{
		Title:   param.Title,
		Content: param.Content,
		User:    userID,
	})
	if err != nil {
		util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	return util.SuccessResponse(context, http.StatusOK, responsePostNew{PostID: postID})
}

type paramPostNewComment struct {
	Content string `json:"content" validate:"required"`
}

func PostNewComment(context echo.Context) error {
	var param paramPostNewComment
	postID := context.Param("id")
	if err := context.Bind(&param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}
	userTmpID := util.MustGetIDFromContext(context)
	userID, err := primitive.ObjectIDFromHex(userTmpID)
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	_, err = model.AddReplyWithPostID(postID, model.Comment{
		User:    userID,
		Content: param.Content,
	})
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	return util.SuccessResponse(context, http.StatusOK, nil)
}
