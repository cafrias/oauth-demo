package common

type SaltedHash struct {
	Salt string
	Hash string
}

type Timestamped struct {
	Created string
	Updated string
}
