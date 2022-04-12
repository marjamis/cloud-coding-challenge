package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/marjamis/cloud-coding-challenge/metis/internal/pkg/database"
	"github.com/marjamis/cloud-coding-challenge/metis/internal/pkg/instance"

	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
)

// NewRouter creates the router used for determing which handler is used on a requests
func NewRouter() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = Logger(http.HandlerFunc(NotFound), "NotFound")

	for _, route := range routes {
		var handler http.Handler
		handler = Logger(route.HandlerFunc, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return
}

// Logger is a wrapper for all request to ensure consistent logging and context creations
func Logger(mh http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, _ := uuid.NewV4()
		ctx := context.WithValue(context.Background(), "requestId", requestID)

		w.Header().Set("metis-request-id", requestID.String())
		log.WithFields(log.Fields{
			"Function": "Logger#1",
		}).Info(r)
		mh.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetDBData reads the most recent entry from the database, which in this case means the last access to the website.
func GetDBData() (raw DBData, err error) {
	data, err := database.DBConnection.Query("select id,instanceId,requesterIp,addedDate from metis order by id desc limit 1")
	defer data.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "GetDBData#1",
		}).Error(err)
		return raw, err
	}

	if !data.Next() {
		return raw, nil
	}

	err = data.Scan(&raw.ID, &raw.InstanceID, &raw.RequesterIP, &raw.AddedDate)
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "GetDBData#2",
		}).Error(err)
		return raw, err
	}

	return raw, nil
}

// WriteData writes to the DB some information, such as instanceId, requester IP address and the timestamp to be used on subsequent gets.
func WriteData(r *http.Request) (err error) {
	data, err := database.DBConnection.Query("insert into metis values(\"\",\"" + instance.InstanceID + "\",\" " + r.Header.Get("X-Forwarded-For") + " \",\"" + time.Now().String() + "\")")
	defer data.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "WriteData#1",
		}).Error(err)
		return err
	}
	return
}
