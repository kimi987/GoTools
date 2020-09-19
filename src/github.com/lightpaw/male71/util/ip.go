package util

func ToU32Ip(addr [4]byte) uint32 {
	/* 127.0.0.1 = 0x100007f */
	var ip32 uint32
	ip32 |= uint32(addr[0])
	ip32 |= uint32(addr[1]) << 8
	ip32 |= uint32(addr[2]) << 16
	ip32 |= uint32(addr[3]) << 24

	return ip32
}