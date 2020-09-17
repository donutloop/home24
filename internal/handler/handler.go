package handler

import (
	"encoding/json"
	"github.com/donutloop/home24/pkg/webscrapper"
	"github.com/go-chi/chi"
	"io"
	"io/ioutil"
	"net/http"
)

type GetWebsiteDataRequest struct {
	WebsiteURL string `json:"website_url"`
}

type WebsiteDataService struct {
	Webscrapper *webscrapper.WebScrapper
}

func (websiteDataService *WebsiteDataService) Extract(ctx Context, req interface{}) (interface{}, error) {

	getWebsiteDataRequest := req.(*GetWebsiteDataRequest)

	ctx.LogInfo("url ", getWebsiteDataRequest.WebsiteURL)

	getWebsiteDataResponse, err := websiteDataService.Webscrapper.Extract(getWebsiteDataRequest.WebsiteURL)
	if err != nil {
		return nil, err
	}
	return getWebsiteDataResponse, nil
}

// HandlerFunc is an http.HandlerFunc with an Context.
type HandlerDomainFunc func(Context, interface{}) (interface{}, error)

type Context struct {
	LogError    func(v ...interface{})
	LogInfo     func(v ...interface{})
	RouteParams chi.RouteParams
}

func WrapGetWebsiteDataHandler(HandlerDomainFunc HandlerDomainFunc, ctx Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		getWebsiteDataRequest := new(GetWebsiteDataRequest)
		if err := UnmarshalRequest(r.Body, getWebsiteDataRequest); err != nil {
			ctx.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		getWebsiteDataResponse, err := HandlerDomainFunc(ctx, getWebsiteDataRequest)
		if err != nil {
			ctx.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(getWebsiteDataResponse)
		if err != nil {
			ctx.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(respBody)
		if err != nil {
			ctx.LogError(err)
		}
	})
}

func UnmarshalRequest(requestBody io.ReadCloser, v interface{}) error {
	b, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}
