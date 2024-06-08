package proxmox

// The models in this file represent the Network endpoint in Proxmox
// The response from the Proxmox server is wrapped in an JSON object called "data".
// The Response structs deal with unwrapping these objects.

// NetworksResponse The response from Proxmox when a list of networks is returned
type NetworksResponse struct {
	Data []Network `json:"data"`
}

// NetworkResponse The response from Proxmox when a single network is returned
type NetworkResponse struct {
	Data Network `json:"data"`
}

// The NetworkRequest struct handles the create and update requests we need to send to the Proxmox server.
// This is different from the Network structure because the API endpoints don't accept the same fields they return.
// I have written more about it here: https://github.com/clincha-org/proxmox-api/issues/5

// NetworkRequest The request that Proxmox expects when creating and modifying Networks
type NetworkRequest struct {
	Interface       string `json:"iface,omitempty"`
	Type            string `json:"type,omitempty"`
	Address         string `json:"address,omitempty"`
	AutoStart       *bool  `json:"autostart,omitempty"`
	BridgePorts     string `json:"bridge_ports,omitempty"`
	BridgeVlanAware *bool  `json:"bridge_vlan_aware,omitempty"`
	CIDR            string `json:"cidr,omitempty"`
	Comments        string `json:"comments,omitempty"`
	Gateway         string `json:"gateway,omitempty"`
	MTU             int64  `json:"mtu,omitempty"`
	Netmask         string `json:"netmask,omitempty"` // Subnet Mask Notation
	VlanID          int64  `json:"vlan-id,string,omitempty"`
}

// Network The structure that represents a Proxmox node network
type Network struct {
	Interface       string `json:"iface,omitempty"`
	Type            string `json:"type,omitempty"`
	Address         string `json:"address,omitempty"`
	Autostart       int64  `json:"autostart,omitempty"`
	BridgePorts     string `json:"bridge_ports,omitempty"`
	BridgeVlanAware int64  `json:"bridge_vlan_aware,omitempty"`
	CIDR            string `json:"cidr,omitempty"`
	Comments        string `json:"comments,omitempty"`
	Gateway         string `json:"gateway,omitempty"`
	MTU             int64  `json:"mtu,string,omitempty"`
	Netmask         string `json:"netmask,omitempty"` // CIDR Notation
	VlanID          int64  `json:"vlan-id,string,omitempty"`
	// Computed values
	Families []string `json:"families,omitempty"`
	Method   string   `json:"method,omitempty"`
	Active   int64    `json:"active,omitempty"`
}
