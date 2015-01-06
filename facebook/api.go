package facebook

type API struct {
	graphURL string
}

func NewAPI(conf Config) API {
	return API{
		graphURL: "https://graph.facebook.com/" + conf.ApiVersion,
	}
}
