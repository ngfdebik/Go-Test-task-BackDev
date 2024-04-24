package user

type User struct {
	//ID             string `json: "id" bson:"_id"`
	GUID string `json: "guid" bson:"_id"`
}

func New(guid string) *User {
	return &User{GUID: guid}
}

type CreateUserDTO struct {
	GUID           string `json: "guid" bson:"guid"`
	RefreshToken   string `json:"RefreshToken"`
	AccessIssuedAt string `json:"AccessIssuedAt"`
}

func NewDTO(guid string, refresh string, issuedAt string) *CreateUserDTO {
	return &CreateUserDTO{GUID: guid, RefreshToken: refresh, AccessIssuedAt: issuedAt}
}
