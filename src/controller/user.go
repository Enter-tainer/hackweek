package controller

import (
	"log"
	"net/http"
	"tree-hole/config"
	"tree-hole/model"
	"tree-hole/util"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
)

type paramUserGetToken struct {
	Email    string `query:"email" validate:"required"`
	Password string `query:"password" validate:"required"`
}

type responseUserGetToken struct {
	Token  string `json:"token"`
	ID     string `json:"_id"`
	Expire int64  `json:"expire_time"`
}

func UserGetToken(context echo.Context) error {
	var param paramUserGetToken
	if err := context.Bind(&param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}

	user, found, err := model.GetUserWithEmail(param.Email)
	if !found {
		return util.ErrorResponse(context, http.StatusBadRequest, "user not found")
	}
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	if user.Password != util.PasswordHash(param.Password) {
		return util.ErrorResponse(context, http.StatusForbidden, "email and password don't match")
	}
	if !user.Verified {
		return util.ErrorResponse(context, http.StatusForbidden, "please verify your email")
	}

	token, expireTime, err := util.NewJWTToken(user.ID.Hex())
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	return util.SuccessResponse(context, http.StatusOK, responseUserGetToken{
		Token:  token,
		ID:     user.ID.Hex(),
		Expire: expireTime.Unix(),
	})
}

type paramUserRegister struct {
	// Username string `json:"username" validate:"required"`
	// Password string `json:"password" validate:"required"`
	// Phone    string `json:"phone" validate:"required,numeric"`
	Email string `json:"email" validate:"required,email"`
}

type responseUserRegister struct {
	ID string `json:"_id"`
}

func UserRegister(context echo.Context) error {
	var param paramUserRegister
	if err := context.Bind(&param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}

	_, found, err := model.GetUserWithEmail(param.Email)
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	if found {
		return util.ErrorResponse(context, http.StatusBadRequest, "email already exists")
	}

	user := model.User{
		// Username: param.Username,
		Password: "",
		// Phone:    param.Phone,
		Email: param.Email,
		// IsAdmin:  false,
		Verified: false,
	}
	idHex, err := model.AddUser(user)

	verifyCode := util.RandomString(config.Config.App.VerifyCodeLength)
	log.Printf("code for %s: %s", param.Email, verifyCode)
	err = util.SendEmail(param.Email, "注册邮箱验证码", "您的邮箱验证码为：<code>"+verifyCode+"</code>")
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	err = model.AddVerifyCode(verifyCode, idHex)
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessResponse(context, http.StatusCreated, responseUserRegister{
		ID: idHex,
	})
}

type paramUserVerify struct {
	ID       string `json:"_id" validate:"required"`
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

func UserVerify(context echo.Context) error {
	var param paramUserVerify
	if err := context.Bind(&param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}

	id, found, err := model.GetVerifyCode(param.Code)
	if !found {
		return util.ErrorResponse(context, http.StatusBadRequest, "verify code not found or expired")
	}
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	if id != param.ID {
		return util.ErrorResponse(context, http.StatusBadRequest, "verify code doesn't match")
	}
	err = model.DeleteVerifyCode(param.Code)
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}

	err = model.UpdateUser(param.ID, bson.M{"verified": true, "password": util.PasswordHash(param.Password)})
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	return util.SuccessResponse(context, http.StatusOK, nil)
}

// type paramUserUpdateInfo struct {
// 	ID       string `json:"_id" validate:"required"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// 	Phone    string `json:"phone" validate:"omitempty,numeric"`
// 	Email    string `json:"email" validate:"omitempty,email"`
// }

// func UserUpdateInfo(context echo.Context) error {
// 	var param paramUserUpdateInfo
// 	if err := context.Bind(&param); err != nil {
// 		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
// 	}
// 	if err := context.Validate(param); err != nil {
// 		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
// 	}

// 	id := util.MustGetIDFromContext(context)
// 	user, _, err := model.GetUserWithID(id)
// 	if err != nil {
// 		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
// 	}

// 	info := make(bson.M)
// 	if param.Username != "" {
// 		_, found, err := model.GetUserWithUsername(param.Username)
// 		if err != nil {
// 			return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
// 		}
// 		if param.Username != user.Username && found {
// 			return util.ErrorResponse(context, http.StatusBadRequest, "username already exists")
// 		}
// 		info["username"] = param.Username
// 	}
// 	if param.Password != "" {
// 		info["password"] = base64.StdEncoding.EncodeToString([]byte(param.Password))
// 	}
// 	if param.Email != "" {
// 		verifyCode := util.RandomString(config.Config.App.VerifyCodeLength)
// 		err := util.SendEmail(param.Email, "注册邮箱验证码", "您的邮箱验证码为：<code>"+verifyCode+"</code>")
// 		if err != nil {
// 			return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
// 		}
// 		err = model.AddVerifyCode(verifyCode, param.ID)
// 		if err != nil {
// 			return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
// 		}
// 		info["email"] = param.Email
// 		info["verified"] = false
// 	}
// 	if param.Phone != "" {
// 		info["phone"] = param.Phone
// 	}

// 	err = model.UpdateUser(param.ID, info)
// 	if err != nil {
// 		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
// 	}
// 	return util.SuccessResponse(context, http.StatusOK, nil)
// }

type paramUserDelete struct {
	ID string `query:"_id" validate:"required"`
}

func UserDelete(context echo.Context) error {
	var param paramUserDelete
	if err := context.Bind(&param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}
	if param.ID != util.MustGetIDFromContext(context) {
		return util.ErrorResponse(context, http.StatusForbidden, "you can not delete others' account")
	}

	// id := util.MustGetIDFromContext(context)
	// isAdmin, err := model.IsUserAdmin(id)
	// if err != nil {
	// 	return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	// }
	// if !isAdmin {
	// 	return util.ErrorResponse(context, http.StatusForbidden, "you are not admin")
	// }

	err := model.DeleteUser(param.ID)
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}
	return util.SuccessResponse(context, http.StatusOK, nil)
}

type paramUserGetInfo struct {
	ID string `query:"_id" validate:"required"`
}

type responseUserGetInfo struct {
	ID    string `json:"_id"`
	Email string `json:"email"`
}

func UserGetInfo(context echo.Context) error {
	var param paramUserGetInfo
	if err := context.Bind(&param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(param); err != nil {
		return util.ErrorResponse(context, http.StatusBadRequest, err.Error())
	}

	id := util.MustGetIDFromContext(context)
	// isAdmin, err := model.IsUserAdmin(id)
	// if err != nil {
	// 	return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	// }

	if param.ID != id {
		return util.ErrorResponse(context, http.StatusForbidden, "you are not admin")
	}

	user, found, err := model.GetUserWithID(param.ID)
	if !found {
		return util.ErrorResponse(context, http.StatusBadRequest, "user not found")
	}
	if err != nil {
		return util.ErrorResponse(context, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessResponse(context, http.StatusOK, responseUserGetInfo{ID: user.ID.Hex(), Email: user.Email})
}
