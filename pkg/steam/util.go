package steam

const id64Base int64 = 76561197960265728

// ID32To64 converts a steam id 32 to steam id 64
func ID32To64(id int32) int64 {
	return id64Base + int64(id)
}

// ID64To32 converts a steam id 64 to steam id 32
func ID64To32(id int64) int32 {
	return int32(id)
}
