package rest

const (
	AuthAPIToken        = AuthType(1)
	AuthOAuth           = AuthType(2)
	AuthClusterAPIToken = AuthType(3)
)

type AuthType uint8

func (a AuthType) APIToken() bool {
	return a&AuthAPIToken != 0
}

func (a AuthType) OAuth() bool {
	return a&AuthOAuth != 0
}

func (a AuthType) ClusterAPIToken() bool {
	return a&AuthClusterAPIToken != 0
}
