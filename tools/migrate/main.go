package main

import _ "github.com/ServiceComb/service-center/server/core/registry/etcd"
import (
	"context"
	"fmt"
	"github.com/ServiceComb/service-center/pkg/lager"
	"github.com/ServiceComb/service-center/server/core"
	"github.com/ServiceComb/service-center/server/core/mux"
	pb "github.com/ServiceComb/service-center/server/core/proto"
	"github.com/ServiceComb/service-center/server/core/registry"
	"github.com/ServiceComb/service-center/server/service/microservice"
	"github.com/ServiceComb/service-center/util"
	"github.com/ServiceComb/service-center/version"
	"github.com/astaxie/beego"
	"os"
	"strings"
)

const MIN_VERSION = "0.1.1"

func init() {
	logFormatText, err := beego.AppConfig.Bool("LogFormatText")
	loggerFile := os.ExpandEnv(beego.AppConfig.String("logfile"))
	enableRsyslog, err := beego.AppConfig.Bool("EnableRsyslog")
	if err != nil {
		enableRsyslog = false
	}
	enableStdOut := beego.AppConfig.DefaultString("runmode", "prod") == "dev"
	util.InitLogger("migrate", &lager.Config{
		LoggerLevel:   beego.AppConfig.String("loglevel"),
		LoggerFile:    loggerFile,
		EnableRsyslog: enableRsyslog,
		LogFormatText: logFormatText,
		EnableStdOut:  enableStdOut,
	})
}

func main() {
	lock, err := mux.Lock(mux.GLOBAL_LOCK)
	if err != nil {
		util.LOGGER.Errorf(err, "create global lock failed")
		os.Exit(1)
	}

	if !needUpgrade() {
		lock.Unlock()
		util.LOGGER.Infof("load the version is %s, PASS.", core.GetSystemConfig().Version)
		os.Exit(0)
	}

	err = ChangeIncompatibleKeysStore()
	if err != nil {
		lock.Unlock()

		util.LOGGER.Errorf(err, "migrate keys failed")
		os.Exit(1)
	}

	err = upgradeSystemConfig()
	if err != nil {
		lock.Unlock()

		util.LOGGER.Errorf(err, "upgrade system config failed")
		os.Exit(1)
	}
	lock.Unlock()

	os.Exit(0)
}

func needUpgrade() bool {
	if core.GetSystemConfig() == nil {
		err := core.LoadSystemConfig()
		if err != nil {
			util.LOGGER.Errorf(err, "check version failed, can not load the system config")
			return false
		}
	}
	return !microservice.VersionMatchRule(core.GetSystemConfig().Version,
		fmt.Sprintf("%s+", MIN_VERSION))
}

func upgradeSystemConfig() error {
	cfg := core.GetSystemConfig()
	cfg.Version = version.Ver().Version
	return core.UpgradeSystemConfig()
}

func ChangeIncompatibleKeysStore() error {
	// get all domain/project
	domainProject := map[string]struct{}{}

	projResp, err := registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
		Action:     registry.GET,
		Key:        util.StringToBytesWithNoCopy(core.GetServiceRootKey("")),
		WithPrefix: true,
		KeyOnly:    true,
	})
	if err != nil {
		util.LOGGER.Errorf(err, "load all domain/projects failed")
		return err
	}

	for _, projKv := range projResp.Kvs {
		key := util.BytesToStringWithNoCopy(projKv.Key)
		arr := strings.Split(key, "/")
		str := arr[4] + "/" + arr[5]
		if _, ok := domainProject[str]; !ok {
			domainProject[str] = struct{}{}
		}
	}

	util.LOGGER.Infof("load all domain/projects(%d) [OK]", len(domainProject))
	for domain := range domainProject {
		// tag
		resp, err := registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
			Action:     registry.GET,
			Key:        util.StringToBytesWithNoCopy(GetOldServiceTagRootKey(domain)),
			WithPrefix: true,
		})
		if err != nil {
			util.LOGGER.Errorf(err, "%s: load all old tags failed", domain)
			return err
		}

		util.LOGGER.Infof("%s: load all old tags(%d) [OK]", domain, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			key := util.BytesToStringWithNoCopy(kv.Key)
			serviceId := key[strings.LastIndex(key, "/")+1:]
			newKey := core.GenerateServiceTagKey(domain, serviceId)
			_, err := registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
				Action: registry.PUT,
				Key:    util.StringToBytesWithNoCopy(newKey),
				Value:  kv.Value,
			})
			if err != nil {
				util.LOGGER.Errorf(err, "%s: put new tags failed", domain)
				return err
			}
			util.LOGGER.Infof("%s: migrate tag %s to %s", domain, key, newKey)
		}
		// rule
		resp, err = registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
			Action:     registry.GET,
			Key:        util.StringToBytesWithNoCopy(GetOldServiceRuleRootKey(domain)),
			WithPrefix: true,
		})
		if err != nil {
			util.LOGGER.Errorf(err, "%s: get all old rules failed", domain)
			return err
		}

		util.LOGGER.Infof("%s: load all old rules(%d) [OK]", domain, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			key := util.BytesToStringWithNoCopy(kv.Key)
			arr := strings.Split(key, "/")
			l := len(arr)
			newKey := core.GenerateServiceRuleKey(domain, arr[l-2], arr[l-1])
			_, err := registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
				Action: registry.PUT,
				Key:    util.StringToBytesWithNoCopy(newKey),
				Value:  kv.Value,
			})
			if err != nil {
				util.LOGGER.Errorf(err, "%s: put new rules failed", domain)
				return err
			}
			util.LOGGER.Infof("%s: migrate rule %s to %s", domain, key, newKey)
		}
		// rule index
		resp, err = registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
			Action:     registry.GET,
			Key:        util.StringToBytesWithNoCopy(GetOldServiceRuleIndexRootKey(domain)),
			WithPrefix: true,
		})
		if err != nil {
			util.LOGGER.Errorf(err, "%s: get all old rule indexes failed", domain)
			return err
		}

		util.LOGGER.Infof("%s: load all old rule indexes(%d) [OK]", domain, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			key := util.BytesToStringWithNoCopy(kv.Key)
			arr := strings.Split(key, "/")
			l := len(arr)
			newKey := core.GenerateRuleIndexKey(domain, arr[l-3], arr[l-2], arr[l-1])
			_, err := registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
				Action: registry.PUT,
				Key:    util.StringToBytesWithNoCopy(newKey),
				Value:  kv.Value,
			})
			if err != nil {
				util.LOGGER.Errorf(err, "%s: put new rule indexes failed", domain)
				return err
			}
			util.LOGGER.Infof("%s: migrate rule index %s to %s", domain, key, newKey)
		}
		// dependency
		for _, t := range []string{"/p/", "/c/"} {
			resp, err = registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
				Action:     registry.GET,
				Key:        util.StringToBytesWithNoCopy(GetOldServiceDependencyRootKey(domain) + t),
				WithPrefix: true,
			})
			if err != nil {
				util.LOGGER.Errorf(err, "%s: get all old dependencies failed", domain)
				return err
			}

			util.LOGGER.Infof("%s: load all old dependencies(%d) [OK]", domain, len(resp.Kvs))
			for _, kv := range resp.Kvs {
				key := util.BytesToStringWithNoCopy(kv.Key)
				arr := strings.Split(key, "/")
				l := len(arr)
				newKey := core.GenerateServiceDependencyKey(arr[l-3], domain, arr[l-2], arr[l-1])
				_, err := registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
					Action: registry.PUT,
					Key:    util.StringToBytesWithNoCopy(newKey),
					Value:  kv.Value,
				})
				if err != nil {
					util.LOGGER.Errorf(err, "%s: put new dependencies failed", domain)
					return err
				}
				util.LOGGER.Infof("%s: migrate dependencies %s to %s", domain, key, newKey)
			}
		}
		// dependency rule
		resp, err = registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
			Action:     registry.GET,
			Key:        util.StringToBytesWithNoCopy(GetOldServiceDependencyRuleRootKey(domain)),
			WithPrefix: true,
		})
		if err != nil {
			util.LOGGER.Errorf(err, "%s: get all old dependency rules failed", domain)
			return err
		}

		util.LOGGER.Infof("%s: load all old dependency rules(%d) [OK]", domain, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			key := util.BytesToStringWithNoCopy(kv.Key)
			arr := strings.Split(key, "/")
			l := len(arr)
			newKey := core.GenerateServiceDependencyRuleKey(arr[l-5], domain, &pb.MicroServiceKey{
				AppId:       arr[l-4],
				Stage:       arr[l-3],
				ServiceName: arr[l-2],
				Version:     arr[l-1],
			})
			_, err := registry.GetRegisterCenter().Do(context.Background(), &registry.PluginOp{
				Action: registry.PUT,
				Key:    util.StringToBytesWithNoCopy(newKey),
				Value:  kv.Value,
			})
			if err != nil {
				util.LOGGER.Errorf(err, "%s: put new dependency rules failed", domain)
				return err
			}
			util.LOGGER.Infof("%s: migrate dependency rule %s to %s", domain, key, newKey)
		}
	}

	util.LOGGER.Infof("changed all incompatible keys store")
	return nil
}

func GetOldServiceRuleRootKey(tenant string) string {
	return util.StringJoin([]string{
		core.GetDomainProjectRootKey(tenant),
		core.REGISTRY_SERVICE_KEY,
		core.REGISTRY_RULE_KEY,
	}, "/")
}

func GetOldServiceRuleIndexRootKey(tenant string) string {
	return util.StringJoin([]string{
		core.GetDomainProjectRootKey(tenant),
		core.REGISTRY_RULE_KEY,
		core.REGISTRY_INDEX,
	}, "/")
}

func GetOldServiceTagRootKey(tenant string) string {
	return util.StringJoin([]string{
		core.GetDomainProjectRootKey(tenant),
		core.REGISTRY_SERVICE_KEY,
		core.REGISTRY_TAG_KEY,
	}, "/")
}

func GetOldServiceDependencyRuleRootKey(tenant string) string {
	return util.StringJoin([]string{
		core.GetDomainProjectRootKey(tenant),
		core.REGISTRY_SERVICE_KEY,
		core.REGISTRY_DEPENDENCY_KEY,
		"rule",
	}, "/")
}

func GetOldServiceDependencyRootKey(tenant string) string {
	return util.StringJoin([]string{
		core.GetDomainProjectRootKey(tenant),
		core.REGISTRY_SERVICE_KEY,
		core.REGISTRY_DEPENDENCY_KEY,
	}, "/")
}
