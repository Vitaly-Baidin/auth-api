package redis

type Option func(*Redis)

func Addresses(addr []string) Option {
	return func(r *Redis) {
		addrMap := make(map[string]string)
		for i := 0; i < len(addr); i += 2 {
			addrMap[addr[i]] = addr[i+1]
		}

		r.addresses = addrMap
	}
}
