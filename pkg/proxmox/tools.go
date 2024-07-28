package proxmox

import (
	"fmt"
	"strconv"
)

func ConvertCIDRToNetmask(cidr *string) (*string, error) {
	cidrInt, err := strconv.Atoi(*cidr)
	if err != nil {
		return cidr, err
	}
	var mask uint32 = 0xFFFFFFFF << (32 - uint32(cidrInt))
	var netmask = fmt.Sprintf("%d.%d.%d.%d", byte(mask>>24), byte(mask>>16), byte(mask>>8), byte(mask))
	return &netmask, nil
}
