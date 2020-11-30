package bolt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func createContext(w http.ResponseWriter, r *http.Request) *Ctx {
	path, method := r.URL.Path, r.Method

	return &Ctx{
		Path:       path,
		Method:     method,
		Status:     200,
		PathParams: make(map[string]string),
		Headers:    make(map[string]string),
		R:          r,
		W:          w,
	}
}

func ToStruct(ctx *Ctx, i interface{}) error {
	if len(*ctx.JSON) != 0 {
		if err := json.Unmarshal(*ctx.JSON, &i); err != nil {
			e := fmt.Sprintf("error parsing body: %v", err)
			return errors.New(e)
		}
	}
	return nil
}

func requestBodyMethods(method string) bool {
	if method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodOptions {
		return false
	}
	return true
}
