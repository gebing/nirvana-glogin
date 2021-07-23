package hawkeye

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gmoss/v2"
	"glogin/config"
	"glogin/constant"
	"glogin/internal"
	"glogin/pbs/hawkeye_login"
	"glogin/pbs/hawkeye_register"
	"strings"
)

// AOD特殊处理（AOD台湾审核和内网unity dev等需绕开鹰眼检查）
func channelFilter(channel string) bool {
	filterCfg := config.Field("hawkeye_filter").String()
	filters := strings.Split(filterCfg, "|")
	log.Infow("hawkeye_filter marks ", "channel", channel, "filters", filters)
	for _, value := range filters {
		if value == channel {
			return true
		}
	}
	return false
}

func CheckLogin(req internal.Req) (interface{}, error) {
	if !config.Field("hawkeye_open").Bool() {
		return nil, nil
	}
	if channelFilter(req.Game.Channel) {
		return nil, nil
	}

	loginIn := &hawkeye_login.Login{
		GameCd: req.Game.GameCd,
		DeviceInfo: &hawkeye_login.DeviceInfo{
			SmId:        req.Client.Dhid,
			Ip:          req.IP,
			AppsflyerId: req.Game.AppsflyerId,
		},
		UserInfo: &hawkeye_login.UserInfo{
			BundleId: req.Game.BundleId,
			Account:  req.Account,
		},
	}
	//service := gmoss.MossWithClusterService("yanghaitao_dev", "hawkeye")
	cfgDc := strings.Split(config.Field(constant.HawkEyeDc).String(), "|")
	service := gmoss.MossWithDcClusterService(cfgDc[0], cfgDc[1], constant.HawkEyeService)
	rsp, err := hawkeye_login.HawkeyeLogin(service, loginIn, gmoss.Call, gmoss.DefaultCallOption())
	if err != nil {
		return nil, err
	} else {
		if !rsp.Pass {
			return nil, fmt.Errorf("hawkeye_login HawkeyeLogin : %v", rsp)
		} else {
			return nil, nil
		}
	}
	return nil, nil
}

func CheckRegister(req internal.Req) (interface{}, error) {
	if !config.Field("hawkeye_open").Bool() {
		return nil, nil
	}
	if channelFilter(req.Game.Channel) {
		return nil, nil
	}

	registerIn := &hawkeye_register.Register{
		GameCd:  req.Game.GameCd,
		Subject: hawkeye_register.Register_Account,
		DeviceInfo: &hawkeye_register.DeviceInfo{
			SmId:        req.Client.Dhid,
			Ip:          req.IP,
			AppsflyerId: req.Game.AppsflyerId,
		},
		UserInfo: &hawkeye_register.UserInfo{
			BundleId: req.Game.BundleId,
		},
	}

	cfgDc := strings.Split(config.Field(constant.HawkEyeDc).String(), "|")
	service := gmoss.MossWithDcClusterService(cfgDc[0], cfgDc[1], constant.HawkEyeService)
	rsp, err := hawkeye_register.HawkeyeRegister(service, registerIn, gmoss.Call, gmoss.DefaultCallOption())
	if err != nil {
		return nil, err
	} else {
		if !rsp.Pass {
			return nil, fmt.Errorf("hawkeye_login HawkeyeLogin : %v", rsp)
		} else {
			return nil, nil
		}
	}
	return nil, nil
}
