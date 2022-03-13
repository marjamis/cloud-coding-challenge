package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/*
Creates the router used for determing which handler is used on a requests
*/
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

/*
Wrapper for all request to ensure consistent logging and context creations
*/
func Logger(mh http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.NewV4()
		ctx := context.WithValue(context.Background(), "requestId", requestId)

		w.Header().Set("metis-request-id", requestId.String())
		log.WithFields(log.Fields{
			"Function": "Logger#1",
		}).Info(r)
		mh.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*
Will curl the metadata service to get the instances id or if not found report this.
*/
func GetInstanceId() (instanceId string) {
	instanceId = "NOT_FOUND"

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "GetInstanceId#1",
		}).Error(err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "GetInstanceId#2",
		}).Error(err)
	}
	instanceId = string(data)

	return
}

/*
Reads the most recent entry from the database, which in this case means the last access to the website.
*/
func GetDBData() (raw DBData, err error) {
	data, err := DB_CONNECTION.Query("select id,instanceId,requesterIp,addedDate from metis order by id desc limit 1")
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

	err = data.Scan(&raw.Id, &raw.InstanceId, &raw.RequesterIp, &raw.AddedDate)
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "GetDBData#2",
		}).Error(err)
		return raw, err
	}

	return raw, nil
}

/*
Writes to the DB some information, such as instanceId, requester IP address and the timestamp to be used on subsequent gets.
*/
func WriteData(r *http.Request) (err error) {
	data, err := DB_CONNECTION.Query("insert into metis values(\"\",\"" + INSTANCE_ID + "\",\" " + r.Header.Get("X-Forwarded-For") + " \",\"" + time.Now().String() + "\")")
	defer data.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "WriteData#1",
		}).Error(err)
		return err
	}
	return
}
