package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HRobotVSwitch struct {
	ID           int64                       `json:"id"`
	Name         string                      `json:"name"`
	VLAN         int64                       `json:"vlan"`
	Cancelled    bool                        `json:"cancelled"`
	Server       []HRobotVSwitchServer       `json:"server"`
	Subnet       []HRobotVSwitchSubnet       `json:"subnet"`
	CloudNetwork []HRobotVSwitchCloudNetwork `json:"cloud_network"`
}

type HRobotVSwitchServer struct {
	ServerIP      string `json:"server_ip"`
	ServerIPv6Net string `json:"server_ipv6_net"`
	ServerNumber  int64  `json:"server_number"`
	Status        string `json:"status"`
}

type HRobotVSwitchSubnet struct {
	IP      string `json:"ip"`
	Mask    int64  `json:"mask"`
	Gateway string `json:"gateway"`
}

type HRobotVSwitchCloudNetwork struct {
	ID      int64  `json:"id"`
	IP      string `json:"ip"`
	Mask    int64  `json:"mask"`
	Gateway string `json:"gateway"`
}

func (c *HetznerRobotClient) GetVSwitch(ctx context.Context, vSwitchID int64) (*HRobotVSwitch, error) {
	if vSwitchID == 0 {
		return nil, fmt.Errorf("vSwitchID cannot be empty.")
	}

	var vSwitch HRobotVSwitch
	responseBytes, err := c.doRequest(ctx, "GET", fmt.Sprintf("/vswitch/%d", vSwitchID), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBytes, &vSwitch)
	if err != nil {
		return nil, err
	}

	return &vSwitch, nil

}
