package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"strings"

)

func PanicCatcher(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Printf("%s\n", debug.Stack())
		msg := fmt.Sprintf("PANIC! %s", r)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func DecodeJsonBody(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	var err error

	switch s.(type) {
	case *proto.Request:
		c := s.(*proto.Request)
		err = decoder.Decode(c)
	default:
		rt := reflect.TypeOf(s)
		//return fmt.Errorf("DecodeJsonBody: Unhandled request type: %s", rt)
		// XXX Dev Setting
		errMsg := fmt.Sprintf("DecodeJsonBody: Unhandled request type: %s", rt)
		log.Fatal(errMsg)
	}
	if err != nil {
		return err
	}
	return nil
}

func ResultLength(r *somaResult, t ErrorMarker) int {
	switch t.(type) {
	case *proto.Result:
		switch {
		case r.Datacenters != nil:
			return len(r.Datacenters)
		case r.Levels != nil:
			return len(r.Levels)
		case r.Predicates != nil:
			return len(r.Predicates)
		case r.Status != nil:
			return len(r.Status)
		case r.Oncall != nil:
			return len(r.Oncall)
		case r.Teams != nil:
			return len(r.Teams)
		case r.Nodes != nil:
			return len(r.Nodes)
		case r.Views != nil:
			return len(r.Views)
		case r.Servers != nil:
			return len(r.Servers)
		case r.Units != nil:
			return len(r.Units)
		case r.Providers != nil:
			return len(r.Providers)
		case r.Metrics != nil:
			return len(r.Metrics)
		case r.Modes != nil:
			return len(r.Modes)
		case r.Users != nil:
			return len(r.Users)
		case r.Systems != nil:
			return len(r.Systems)
		case r.Capabilities != nil:
			return len(r.Capabilities)
		case r.Properties != nil:
			return len(r.Properties)
		case r.Attributes != nil:
			return len(r.Attributes)
		case r.Repositories != nil:
			return len(r.Repositories)
		case r.Buckets != nil:
			return len(r.Buckets)
		case r.Groups != nil:
			return len(r.Groups)
		case r.Clusters != nil:
			return len(r.Clusters)
		case r.CheckConfigs != nil:
			return len(r.CheckConfigs)
		case r.Validity != nil:
			return len(r.Validity)
		case r.HostDeployments != nil:
			if len(r.Deployments) > len(r.HostDeployments) {
				return len(r.Deployments)
			}
			return len(r.HostDeployments)
		case r.Deployments != nil:
			return len(r.Deployments)
		}
	default:
		return 0
	}
	return 0
}

func DispatchBadRequest(w *http.ResponseWriter, err error) {
	http.Error(*w, err.Error(), http.StatusBadRequest)
}

func DispatchInternalError(w *http.ResponseWriter, err error) {
	http.Error(*w, err.Error(), http.StatusInternalServerError)
}

func DispatchJsonReply(w *http.ResponseWriter, b *[]byte) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(*b)
}

func GetPropertyTypeFromUrl(u *url.URL) (string, error) {
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
