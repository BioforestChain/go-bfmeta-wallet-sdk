package broadcast

type Req struct {
	Key interface{} `json:"key"`
}
type Params map[string]interface{}
