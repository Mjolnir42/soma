package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
	"github.com/mjolnir42/soma/lib/proto"
)

func PanicCatcher(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Printf("%s\n", debug.Stack())
		msg := fmt.Sprintf("PANIC! %s", r)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func DecodeJSONBody(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	var err error

	switch s.(type) {
	case *proto.Request:
		c := s.(*proto.Request)
		err = decoder.Decode(c)
	case *auth.Kex:
		c := s.(*auth.Kex)
		err = decoder.Decode(c)
	default:
		rt := reflect.TypeOf(s)
		err = fmt.Errorf("DecodeJSONBody: Unhandled request type: %s", rt)
	}
	return err
}

func DispatchBadRequest(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Error(*w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func DispatchUnauthorized(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusUnauthorized)
		return
	}
	http.Error(*w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func DispatchForbidden(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusForbidden)
		return
	}
	http.Error(*w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

func DispatchNotFound(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusNotFound)
		return
	}
	http.Error(*w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func DispatchConflict(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusConflict)
		return
	}
	http.Error(*w, http.StatusText(http.StatusConflict), http.StatusConflict)
}

func DispatchInternalError(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(*w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func DispatchNotImplemented(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusNotImplemented)
		return
	}
	http.Error(*w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func DispatchJSONReply(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*b)
}

func DispatchOctetReply(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set("Content-Type", `application/octet-stream`)
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*b)
}

func GetPropertyTypeFromURL(u *url.URL) (string, error) {
	// strip surrounding / and skip first path element `property|filter`
	el := strings.Split(strings.Trim(u.Path, "/"), "/")[1:]
	if el[0] == "property" {
		// looks like the path was /filter/property/...
		el = el[1:]
	}
	switch el[0] {
	case "service":
		switch el[1] {
		case "team":
			return "service", nil
		case "global":
			return "template", nil
		default:
			return "", errors.New("Unknown service property type")
		}
	default:
		return el[0], nil
	}
}

// extractAddress extracts the IP address part of the IP:port string
// set as net/http.Request.RemoteAddr. It handles IPv4 cases like
// 192.0.2.1:48467 and IPv6 cases like [2001:db8::1%lo0]:48467
func extractAddress(str string) string {
	var addr string

	switch {
	case strings.Contains(str, `]`):
		// IPv6 address [2001:db8::1%lo0]:48467
		addr = strings.Split(str, `]`)[0]
		addr = strings.Split(addr, `%`)[0]
		addr = strings.TrimLeft(addr, `[`)
	default:
		// IPv4 address 192.0.2.1:48467
		addr = strings.Split(str, `:`)[0]
	}
	return addr
}

func msgRequest(l *log.Logger, q *msg.Request) {
	l.Printf(LogStrSRq,
		q.Section,
		q.Action,
		q.AuthUser,
		q.RemoteAddr,
	)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
