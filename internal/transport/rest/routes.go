package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gofra/internal/config"
	"gofra/internal/storage"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO make resources separated files if this prototype will growth

/*
	would prefer some gin or chi router registering
	or make url pattern builder over that std mux HandleFunc
	but seems okay "as is" for prototype for now
*/

type Routing struct {
	inmemoryQueue storage.InmemoryQueue
	conf          config.AppConfig
}

func NewRouting(inmemoryQueue storage.InmemoryQueue, conf config.AppConfig) Routing {
	return Routing{inmemoryQueue: inmemoryQueue, conf: conf}
}

func RegisterAppRoutes(mux *http.ServeMux, routing Routing) {
	mux.HandleFunc("GET /queue/{name}/", routing.getQueueMessageResource)
	mux.HandleFunc("PUT /queue/{name}/", routing.putQueueMessageResource)

}

func (rr *Routing) getQueueMessageResource(w http.ResponseWriter, r *http.Request) {
	qName := r.PathValue("name")

	// TODO maybe better move it to some middleware
	timeoutParam := r.URL.Query().Get("timeout")
	timeout, err := strconv.Atoi(timeoutParam)
	if err != nil {
		timeout = rr.conf.RouteDefaultTimeoutSec
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	msg, err := rr.inmemoryQueue.Get(ctx, qName)
	if err != nil && errors.Is(err, storage.ErrTimeoutReached) {
		resp, err := makeJsonResp(w, http.StatusNotFound, "queue waiting message timeout occurred")
		if err != nil {
			handleJsonMarshallingErr(w, err)
			return
		}
		w.Write(resp)
		return
	} else if err != nil {
		resp, err := makeJsonResp(w, http.StatusInternalServerError, "ooops, some unhandled error occurred")
		if err != nil {
			handleJsonMarshallingErr(w, err)
			return
		}
		w.Write(resp)
		return
	}

	resp, err := makeJsonResp(w, http.StatusOK, msg)
	if err != nil {
		handleJsonMarshallingErr(w, err)
		return
	}
	w.Write(resp)
	return
}

func (rr *Routing) putQueueMessageResource(w http.ResponseWriter, r *http.Request) {
	var payload QSchemaIn
	qName := r.PathValue("name")

	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			resp, err := makeJsonResp(w, http.StatusBadRequest, "Content-Type header is not application/json")
			if err != nil {
				handleJsonMarshallingErr(w, err)
				return
			}
			w.Write(resp)
			return
		}
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		resp, err := makeJsonResp(w, http.StatusBadRequest, "unsupported message type - it must be string")
		if err != nil {
			handleJsonMarshallingErr(w, err)
			return
		}
		w.Write(resp)
		return
	}

	if payload.Message == nil {
		resp, err := makeJsonResp(w, http.StatusBadRequest, "payload required field missing: message")
		if err != nil {
			handleJsonMarshallingErr(w, err)
			return
		}
		w.Write(resp)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(rr.conf.RouteDefaultTimeoutSec)*time.Second)
	defer cancel()

	err = rr.inmemoryQueue.Put(ctx, qName, *payload.Message)
	if err != nil && errors.Is(err, storage.ErrQueueFull) {
		resp, err := makeJsonResp(w, http.StatusBadRequest, fmt.Sprintf("queue %s is full", qName))
		if err != nil {
			handleJsonMarshallingErr(w, err)
			return
		}
		w.Write(resp)
		return
	} else if err != nil && errors.Is(err, storage.ErrMaxQueuesReached) {
		resp, err := makeJsonResp(w, http.StatusBadRequest, "queue limits reached")
		if err != nil {
			handleJsonMarshallingErr(w, err)
			return
		}
		w.Write(resp)
		return
	} else if err != nil && errors.Is(err, storage.ErrTimeoutReached) {
		resp, err := makeJsonResp(w, http.StatusNotFound, "queue putting message timeout occurred")
		if err != nil {
			handleJsonMarshallingErr(w, err)
			return
		}
		w.Write(resp)
		return
	} else if err != nil {
		resp, err := makeJsonResp(w, http.StatusInternalServerError, "ooops, some unhandled error occurred")
		if err != nil {
			handleJsonMarshallingErr(w, err)
			return
		}
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil) // prefer return empty json body here for consistency resp content-type, but task says return just empty body
	return
}
