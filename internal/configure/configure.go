package configure

import (
	"encoding/json"
	"github.com/castbox/nirvana-gcore/consul"
	log "github.com/castbox/nirvana-gcore/glog"
	"sync"
)

var (
	dynamicCfg DynamicPub
	pubConfig  sync.Map
)

type DynamicPub struct {
}

func (d DynamicPub) Reload(appName string, jsonData []byte) {
	bInfo := PubCfg{}
	err2 := json.Unmarshal(jsonData, &bInfo)
	if err2 != nil {
		log.Errorw("WatchPubCfgCallBack ", "path2", err2)
		return
	}
	//fmt.Println(bInfo)
	pubConfig.Store(appName, bInfo)
	log.Infow("realod dir", "appName", appName, "data", bInfo)
}

func GetCfg(appName string, vsn string, ip string) ClusterCfg {
	//log.Infow("query CfgHandler GetCfg", "appName", appName, "vsn", vsn, "ip", ip)
	v, ok := pubConfig.Load(appName)
	if !ok {
		log.Infow("GetPubCfg BundleId error", "BundleId", appName)
		return ClusterCfg{}
	}
	//fmt.Println(v)
	p := v.(PubCfg)
	if len(p.WhiteList) > 0 {
		for _, wIp := range p.WhiteList {
			if ip == wIp {
				return p.ConnCfg.Test
			}
		}
	}
	if p.CheckVsn == vsn {
		return p.ConnCfg.Check
	}
	return p.ConnCfg.Stable
}

// WatchDynamicPubDir
// http://127.0.0.1:2000/v1/kv/app_dynamic_cfg/lwk_dev/glogin/pub_cfg/?recurse=true
// dynamicPubUrl := fmt.Sprintf("kv/app_dynamic_cfg/%s/pub_cfg/?recurse=true", cluster)
//[{"LockIndex":0,"Key":"app_dynamic_cfg/lwk_dev/pub_cfg/com.dh.bpc.gp","Flags":0,"Value":"ewogICAgImNvbm5fY2ZnIjp7CiAgICAgICAgInN0YWJsZSI6ewoJCQkJCQkiY2x1c3Rlcl90eXBlIiA6IDEsCiAgICAgICAgICAgICJ1Z2F0ZV9hZGRyIjoiMTAuMC4yNDAuMjM0OjE4ODg5IiwKICAgICAgICAgICAgInZzbl9hZGRyIjoiMTAuMC4yNDAuMjM0OjE5ODg5IiwKICAgICAgICAgICAgInVwYXlfYWRkciI6IjEwLjAuMjQwLjE5OjgwODgiLAogICAgICAgICAgICAidWNoYXRfYWRkciI6ICJhb2QtZGV2LXVjaGF0LmRoZ2FtZXMuY246MTg4ODciLAogICAgICAgICAgICAiY29tbXVuaXR5X3dlYl9hZGRyIjoiaHR0cDovLzEwLjAuMC4xOTo3NzcwL2FvZC92MS9pbmRleC5odG1sIy9ob21lIiwKICAgICAgICAgICAgImNvbW11bml0eV9zcnZfYWRkciI6Imh0dHA6Ly8xODIuMTUwLjIyLjYxOjI3Nzc4IiwKICAgICAgICAgIAkiYWljc193c19hZGRyIjogIndzczovL2FpY3MtY2xpLmRldi1kaC5jb20vZGhfd3Mvd3NfYXBwIiwKICAgICAgICAgICAgImFpY3NfaHR0cF9hZGRyIjogImh0dHBzOi8vYWljcy1jbGkuZGV2LWRoLmNvbSIKICAgICAgICB9LAogICAgICAgICJjaGVjayI6ewogICAgICAgICAgICAiY2x1c3Rlcl90eXBlIjogMiwKICAgICAgICAgICAgInVnYXRlX2FkZHIiOiIiLAogICAgICAgICAgICAidnNuX2FkZHIiOiAiIiwKICAgICAgICAgICAgInVwYXlfYWRkciI6IiIsCiAgICAgICAgICAgICJ1Y2hhdF9hZGRyIjogIiIKICAgICAgICB9CiAgICB9LAogICAgImNoZWNrX3ZzbiI6IiIKfQ==","CreateIndex":16921314,"ModifyIndex":16921314},{"LockIndex":0,"Key":"app_dynamic_cfg/lwk_dev/pub_cfg/com.droidhang.bpc.ios","Flags":0,"Value":"ewogICAgImNvbm5fY2ZnIjp7CiAgICAgICAgInN0YWJsZSI6ewoJCQkJCQkiY2x1c3Rlcl90eXBlIiA6IDEsCiAgICAgICAgICAgICJ1Z2F0ZV9hZGRyIjoiMTAuMC4yNDAuMjM0OjE4ODg5IiwKICAgICAgICAgICAgInZzbl9hZGRyIjoiMTAuMC4yNDAuMjM0OjE5ODg5IiwKICAgICAgICAgICAgInVwYXlfYWRkciI6Imh0dHA6Ly8xMC4wLjI0MC4xOTo4MDg4IiwKICAgICAgICAgICAgInVjaGF0X2FkZHIiOiAiYW9kLWRldi11Y2hhdC5kaGdhbWVzLmNuOjE4ODg3IiwKICAgICAgICAgICAgImNvbW11bml0eV93ZWJfYWRkciI6Imh0dHA6Ly8xMC4wLjAuMTk6Nzc3MC9hb2QvdjEvaW5kZXguaHRtbCMvaG9tZSIsCiAgICAgICAgICAgICJjb21tdW5pdHlfc3J2X2FkZHIiOiJodHRwOi8vMTgyLjE1MC4yMi42MToyNzc3OCIsCiAgICAgICAgICAJImFpY3Nfd3NfYWRkciI6ICJ3c3M6Ly9haWNzLWNsaS5kZXYtZGguY29tL2RoX3dzL3dzX2FwcCIsCiAgICAgICAgICAgICJhaWNzX2h0dHBfYWRkciI6ICJodHRwczovL2FpY3MtY2xpLmRldi1kaC5jb20iCiAgICAgICAgfSwKICAgICAgICAiY2hlY2siOnsKICAgICAgICAgICAgImNsdXN0ZXJfdHlwZSI6IDIsCiAgICAgICAgICAgICJ1Z2F0ZV9hZGRyIjoiIiwKICAgICAgICAgICAgInZzbl9hZGRyIjogIiIsCiAgICAgICAgICAgICJ1cGF5X2FkZHIiOiIiLAogICAgICAgICAgICAidWNoYXRfYWRkciI6ICIiCiAgICAgICAgfQogICAgfSwKICAgICJjaGVja192c24iOiIiCn0=","CreateIndex":16821432,"ModifyIndex":16821481}]
//	[
//	{
//		"LockIndex": 0,
//		"Key": "app_static_cfg/lwk_dev/glogin/pub_cfg/com.dh.bpc.gp",
//		"Flags": 0,
//		"Value": "ewogICAgICAgICAgICAiZmFjZWJvb2tfb2F1dGhfdXJsIjogImh0dHBzOi8vZ3JhcGguZmFjZWJvb2suY29tL2RlYnVnX3Rva2VuP2FjY2Vzc190b2tlbj02MjI0MDAyNjUzNjMyMDIlN0MyMjY4ODgxNDEwYjNmYTA5MzNjYjJiMGVjZWFhNTE2ZSZpbnB1dF90b2tlbj0iLAogICAgICAgICAgICAiYXBwc2ZseWVyX0FORFJPSUQiOiAiaHR0cHM6Ly9hcGkyLmFwcHNmbHllci5jb20vaW5hcHBldmVudC9jb20uZHJvaWRoYW5nLmFvZC5ncCIsCiAgICAgICAgICAgICJhcHBzZmx5ZXJfb3BlbiI6IHRydWUsCiAgICAgICAgICAgICJhcHBzZmx5ZXJfSU9TIjogImh0dHBzOi8vYXBpMi5hcHBzZmx5ZXIuY29tL2luYXBwZXZlbnQvaWQxMTUzNDYxOTE1IiwKICAgICAgICAgICAgImFwcHNmbHllcl9BdXRoZW50aWNhdGlvbiI6ICIzNkZmTmsyNDR4aTlCQ3hFVVJxYTVuIiwKICAgICAgICAgICAgImFwcHNmbHllcl9yZWdpc3RyYXRpb25JZCI6IDEyCiB9",
//		"CreateIndex": 21120571,
//		"ModifyIndex": 21120586
//	}
//]
func WatchDynamicPubDir() {
	// 本地环境
	//var err error
	//sInfo := &consul.ServiceInfo{
	//	Cluster: "lwk_dev",
	//	Service: "glogin",
	//	Index:   2,
	//}

	//部署环境
	sInfo, err := consul.GetServiceInfoByPath()
	if err != nil {
		log.Infow("consul.GetServiceInfoByPath sInfo error", "err", err)
		return
	}
	dir := "lwk_dev/pub_cfg"
	if sInfo != nil {
		dir = sInfo.Cluster + "/pub_cfg"
	}
	log.Infow("WatchDynamicPubDir dir", "dir", dir)
	if err = consul.WatchDir(dir, dynamicCfg); err != nil {
		log.Fatalw("failed to watch dir", "err", err)
	}
}
