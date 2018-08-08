package graphql

import (
	"strconv"
	"time"
)

//GraphQLResponseArray represents an array in a gql response
type ResponseArray struct {
	Response []interface{}
}

//GraphQLResponse represents a general object in a graphql response
type Response struct {
	Response interface{}
}

func (r *ResponseArray) Get(index int) *Response {
	return &Response{r.Response[index]}
}

func (r *ResponseArray) Len() int {
	return len(r.Response)
}

func (r *Response) GetObject(key string) *Response {
	return &Response{(r.Response.(map[string]interface{}))[key]}
}

func (r *Response) GetArray(key string) *ResponseArray {
	return &ResponseArray{r.Response.(map[string]interface{})[key].([]interface{})}
}

func (r *Response) GetString(key string) string {
	return (r.Response.(map[string]interface{}))[key].(string)
}

func (r *Response) GetInt64(key string) (int64, error) {
	return strconv.ParseInt(r.GetString(key), 10, 64)
}

func JSTimestampToTime(timestamp int64) time.Time {
	seconds := timestamp / 1000
	nSeconds := (timestamp % 1000) * 1e6

	return time.Unix(seconds, nSeconds)
}
