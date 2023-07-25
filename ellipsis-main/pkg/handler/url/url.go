package url

import (
	"brief/internal/constant"
	"brief/internal/model"
	"brief/utility"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Redirect - /api/v1/{hash} - GET
func (base *Controller) Redirect(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	url, err := base.UrlService.Redirect(hash)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	http.Redirect(w, r, url.LongURL, http.StatusTemporaryRedirect)
}

//	Shorten
//
// @Summary		shorten a url
// @Description	shorten a url
// @Tags			URL
// @Accept			json
// @Produce		json
// @Param			url	body		model.URL	true	"URL"
// @Success		201		{object}	utility.Response{data=model.URL}
// @Failure		400		{object}	utility.Response
// @Router			/url/shorten [post]
// @Security		JWTToken
func (base *Controller) Shorten(w http.ResponseWriter, r *http.Request) {
	req := new(model.URL)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	// Fetch user information from context
	uInfo := r.Context().Value(struct{}{})
	ctxInfo, ok := uInfo.(*model.ContextInfo)
	if !ok {
		ctxInfo = nil
	}

	if err := base.Validate.Struct(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrValidation, utility.ValidationResponse(err, base.Validate), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := base.UrlService.Shorten(req, ctxInfo, r); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "successfully created url", req)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

//	Get Url's
//
// @Summary		get all my urls
// @Description	get all my urls
// @Tags			URL
// @Accept			json
// @Produce		json
// @Success		200	{object}	utility.Response{data=[]model.URL}
// @Failure		400	{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/url [get]
// @Security		JWTToken
func (base *Controller) GetUrls(w http.ResponseWriter, r *http.Request) {
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
	urls, err := base.UrlService.GetURLs(uId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", urls)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

//	Delete
//
// @Summary		delete my url
// @Description	delete my url
// @Tags			URL
// @Accept			json
// @Produce		json
// @Param			id		path		string					true	"url ID"
// @Success		200	{object}	utility.Response{data=[]model.URL}
// @Failure		400	{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/url/{id} [delete]
// @Security		JWTToken
func (base *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	uInfo := r.Context().Value(struct{}{}) // fetch user's info from context

	if uInfo == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, "user ID not found", nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	uContextInfo := uInfo.(*model.ContextInfo)
	url, err := base.UrlService.Delete(uContextInfo, urlId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successfully deleted url", url)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// ADMIN ENDPOINTS

//	Get All
//
// @Summary		list all urls - Admin
// @Description	list all urls - Admin
// @Tags			URL - Admin
// @Accept			json
// @Produce		json
// @Success		200	{object}	utility.Response{data=[]model.User}
// @Failure		400	{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/url/get-all [get]
// @Security		JWTToken
func (base *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	urls, err := base.UrlService.GetAll()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", urls)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
}

//	Get Urls by User ID
//
// @Summary		get urls by a user - Admin
// @Description	get urls by a user - Admin
// @Tags			URL - Admin
// @Accept			json
// @Produce		json
// @Param			user-id		path		string					true	"user ID"
// @Success		200	{object}	utility.Response{data=[]model.User}
// @Failure		400	{object}	utility.Response
// @Failure		401		{object}	utility.Response
// @Router			/url/get-all/{user-id} [get]
// @Security		JWTToken
func (base *Controller) GetUrlsByUserID(w http.ResponseWriter, r *http.Request) {
	uID := chi.URLParam(r, "user-id")
	urls, err := base.UrlService.GetURLs(uID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", urls)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
