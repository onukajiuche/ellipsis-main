package user

import (
	"brief/internal/constant"
	"brief/internal/model"
	"encoding/json"
	"net/http"

	"brief/utility"

	"github.com/go-chi/chi/v5"
)

//	Register
//
// @Summary		register a user
// @Description	register a user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			user	body		model.User	true	"User"
// @Success		201		{object}	utility.Response{data=model.User}
// @Failure		400		{object}	utility.Response
// @Router			/users [post]
func (base *Controller) Register(w http.ResponseWriter, r *http.Request) {
	req := new(model.User)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := base.Validate.Struct(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrValidation, utility.ValidationResponse(err, base.Validate), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	token, err := base.UserService.Register(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", map[string]interface{}{
		"token": token,
		"user":  req,
	})
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusCreated)
	w.Write(res)

}

//	Login
//
// @Summary		log in
// @Description	log in
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			userInfo	body		model.UserLogin	true	"Login Info"
// @Success		201		{object}	utility.Response{data=model.User}
// @Failure		400		{object}	utility.Response
// @Router			/users/login [post]
func (base *Controller) Login(w http.ResponseWriter, r *http.Request) {
	req := new(model.UserLogin)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := base.Validate.Struct(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrValidation, utility.ValidationResponse(err, base.Validate), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	usr, err := base.UserService.Login(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "logged in successfully", usr)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

//	GetMe
//
// @Summary		get me
// @Description	get ne
// @Tags			User
// @Accept			json
// @Produce		json
// @Success		200	{object}	utility.Response{data=model.User}
// @Failure		400	{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/users [get]
// @Security		JWTToken
func (base *Controller) GetMe(w http.ResponseWriter, r *http.Request) {
	uInfo := r.Context().Value(struct{}{}) // fetch user's info from context
	if uInfo == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, "user ID not found", nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	uId := uInfo.(*model.ContextInfo).ID
	usr, err := base.UserService.Get(uId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", usr)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

//	UpdateMe
//
// @Summary		update a user
// @Description	update a user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			update	body		model.User	true	"User Update"
// @Success		200		{object}	utility.Response{data=model.User}
// @Failure		400		{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/users [patch]
// @Security		JWTToken
func (base *Controller) UpdateMe(w http.ResponseWriter, r *http.Request) {

	req := new(model.User)
	uInfo := r.Context().Value(struct{}{}) // fetch user's info from context

	if uInfo == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, "user ID not found", nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	uId := uInfo.(*model.ContextInfo).ID
	err := base.UserService.Update(uId, req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "updated successfully", req)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

//	ResetPassword
//
// @Summary		update a user's password
// @Description	update a user's password
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			update	body		model.ResetPassword	true	"Password Update"
// @Success		200		{object}	utility.Response{data=model.User}
// @Failure		400		{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/users/reset-password [patch]
// @Security		JWTToken
func (base *Controller) ResetPassword(w http.ResponseWriter, r *http.Request) {

	req := new(model.ResetPassword)
	uInfo := r.Context().Value(struct{}{}) // fetch user's info from context

	if uInfo == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, "user ID not found", nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	uId := uInfo.(*model.ContextInfo).ID
	usr, err := base.UserService.ResetPassword(uId, req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "reset password successfully", usr)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// ADMIN ENDPOINTS

//	Get All
//
// @Summary		list all users - Admin
// @Description	list all users - Admin
// @Tags			User - Admin
// @Accept			json
// @Produce		json
// @Success		200	{object}	utility.Response{data=[]model.User}
// @Failure		400	{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/users/get-all [get]
// @Security		JWTToken
func (base *Controller) GetAll(w http.ResponseWriter, r *http.Request) {

	usrs, err := base.UserService.GetAll()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", usrs)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
}

//	Get User By ID or Email
//
// @Summary		get user - Admin
// @Description	get user - Admin
// @Tags			User - Admin
// @Accept			json
// @Produce		json
// @Param			idOrEmail		path		string					true	"User ID or Email"
// @Success		200	{object}	utility.Response{data=model.User}
// @Failure		400	{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/users/{idOrEmail} [get]
// @Security		JWTToken
func (base *Controller) GetUserByIdOrEmail(w http.ResponseWriter, r *http.Request) {
	idOrEmail := chi.URLParam(r, "idOrEmail")
	usr, err := base.UserService.Get(idOrEmail)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", usr)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

//	Lock User
//
// @Summary		lock user - Admin
// @Description	lock user - Admin
// @Tags			User - Admin
// @Accept			json
// @Produce		json
// @Param			idOrEmail		path		string					true	"User ID or Email"
// @Success		200	{object}	utility.Response{data=model.User}
// @Failure		400	{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/users/lock/{idOrEmail} [patch]
// @Security		JWTToken
func (base *Controller) LockUser(w http.ResponseWriter, r *http.Request) {
	idOrEmail := chi.URLParam(r, "idOrEmail")
	user, err := base.UserService.LockUser(idOrEmail)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "user locked successfully", user)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

//	Unlock User
//
// @Summary		unlock user - Admin
// @Description	unlock user - Admin
// @Tags			User - Admin
// @Accept			json
// @Produce		json
// @Param			idOrEmail		path		string					true	"User ID or Email"
// @Success		200	{object}	utility.Response{data=model.User}
// @Failure		400	{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/users/unlock/{idOrEmail} [patch]
// @Security		JWTToken
func (base *Controller) UnlockUser(w http.ResponseWriter, r *http.Request) {
	idOrEmail := chi.URLParam(r, "idOrEmail")
	user, err := base.UserService.UnlockUser(idOrEmail)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "user unlocked successfully", user)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
