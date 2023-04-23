package internal

type Redirect struct {
	Name      string `json:"name" bson:"name" msgpack:"name"`
	Code      string `json:"code" bson:"code" msgpack:"code"`
	Url       string `json:"url" bson:"url" msgpack:"url" validate:"empty=false & format=url"`
	CreatedAt int64  `json:"created_at" bson:"created_at" msgpack:"created_at"`
}
