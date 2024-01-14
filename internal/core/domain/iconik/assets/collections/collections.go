package collections

// Collection is the top level data structure that receives the unmarshalled payload
// response from GET collection (/API/assets/v1/collections/{collection-id}).
type Collection struct {
	Title string
}
