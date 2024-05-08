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
	Interface              string `json:"iface,omitempty"`
	Type                   string `json:"type,omitempty"`
	Address                string `json:"address,omitempty"`
	Address6               string `json:"address6,omitempty"`
	AutoStart              int64  `json:"autostart,omitempty"`
	BondPrimary            string `json:"bond-primary,omitempty"`
	BondMode               string `json:"bond_mode,omitempty"`
	BondTransmitHashPolicy string `json:"bond_xmit_hash_policy,omitempty"`
	BridgePorts            string `json:"bridge_ports,omitempty"`
	BridgeVlanAware        int64  `json:"bridge_vlan_aware,omitempty"`
	CIDR                   string `json:"cidr,omitempty"`
	CIDR6                  string `json:"cidr6,omitempty"`
	Comments               string `json:"comments,omitempty"`
	Comments6              string `json:"comments6,omitempty"`
	Gateway                string `json:"gateway,omitempty"`
	Gateway6               string `json:"gateway6,omitempty"`
	MTU                    int64  `json:"mtu,omitempty"`
	Netmask                string `json:"netmask,omitempty"`
	Netmask6               string `json:"netmask6,omitempty"`
	OVSBonds               string `json:"ovs_bonds,omitempty"`
	OVSBridge              string `json:"ovs_bridge,omitempty"`
	OVSOptions             string `json:"ovs_options,omitempty"`
	OVSPorts               string `json:"ovs_ports,omitempty"`
	OVSTag                 string `json:"ovs_tag,omitempty"`
	Slaves                 string `json:"slaves,omitempty"`
	VlanID                 string `json:"vlan-id,omitempty"`
	VlanRawDevice          string `json:"vlan-raw-device,omitempty"`
}

// Network The structure that represents a Proxmox node network
type Network struct {
	Interface   string   `json:"iface,omitempty"`
	Type        string   `json:"type,omitempty"`
	Address     string   `json:"address,omitempty"`
	Autostart   int      `json:"autostart,omitempty"`
	Gateway     string   `json:"gateway,omitempty"`
	Families    []string `json:"families,omitempty"`
	Method6     string   `json:"method6,omitempty"`
	BridgeFd    string   `json:"bridge_fd,omitempty"`
	Netmask     string   `json:"netmask,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Active      int      `json:"active,omitempty"`
	Method      string   `json:"method,omitempty"`
	BridgeStp   string   `json:"bridge_stp,omitempty"`
	Cidr        string   `json:"cidr,omitempty"`
	BridgePorts string   `json:"bridge_ports,omitempty"`
}
