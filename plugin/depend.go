package plugin

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/fankane/go-utils/goroutine"
	"gopkg.in/yaml.v3"
)

const (
	connTypeName  = "-"
	MaxPluginSize = 1000 //最多1000个插件
)

type Depender interface {
	// DependsOn 假如一个插件依赖另一个插件，则返回被依赖的插件的列表：数组元素为 type-name 如 [ "database-mysql", "queue-kafka" ]
	DependsOn() []string
}

// pluginInfo 插件信息。
type pluginInfo struct {
	factory Factory
	typ     string
	name    string
	cfg     yaml.Node
}

func InitPlugins(c Config, ignoreErr bool) error {
	pluginInfos, err := loadPlugins(c)
	if err != nil {
		return err
	}
	return setupPlugins(pluginInfos, ignoreErr)
}

func loadPlugins(c Config) ([]pluginInfo, error) {
	result := make([]pluginInfo, 0)
	for typeName, factories := range c {
		for pluginName, conf := range factories {
			if strings.Contains(pluginName, connTypeName) {
				return nil, fmt.Errorf("pluginName:%s contain forbbiden char [%s]", pluginName, connTypeName)
			}
			f := Get(typeName, pluginName)
			if f == nil {
				return nil, fmt.Errorf("%s-%s not register", typeName, pluginName)
			}
			result = append(result, pluginInfo{
				factory: f,
				typ:     typeName,
				name:    pluginName,
				cfg:     conf,
			})
		}
	}
	return result, nil
}

func setupPlugins(pList []pluginInfo, ignoreErr bool) error {
	var (
		pluginChannel = make(chan pluginInfo, MaxPluginSize) // 使用channel初始化插件队列，方便后面按顺序逐个加载插件
		status        = make(map[string]struct{})            // 插件初始化状态，plugin key 存在则表示初始化完成
	)
	fs := make([]func() error, 0)
	lock := &sync.Mutex{}

	for _, tempPlu := range pList {
		info := tempPlu
		hasDeps, err := info.hasDepends()
		if err != nil {
			return err
		}
		if !hasDeps { //没有依赖项，并发注册
			fs = append(fs, func() error {
				setErr := info.factory.Setup(info.name, &info.cfg)
				if setErr != nil {
					if ignoreErr {
						log.Println(fmt.Sprintf("%s setup failed, err:%s", info.name, setErr))
						return nil
					}
					return setErr
				}
				lock.Lock()
				status[info.key()] = struct{}{}
				lock.Unlock()
				log.Println(fmt.Sprintf("%s:%s installed ", info.typ, info.name))
				return nil
			})
		} else { //有依赖项，放入队列
			pluginChannel <- info
		}
	}

	// 先并发执行没有依赖项的插件注册
	if err := goroutine.Exec(fs, goroutine.WithReturnWhenError(true)); err != nil {
		return fmt.Errorf("not depend plugin init failed err:%s", err)
	}

	// 顺序执行有依赖项的插件注册
	pluginNum := len(pluginChannel)
	for pluginNum > 0 {
		for i := 0; i < pluginNum; i++ {
			p := <-pluginChannel
			hasDeps, err := p.checkDepends(status)
			if err != nil {
				return err
			} else if hasDeps { //有依赖，放到队尾，继续初始化其他的
				pluginChannel <- p
				continue
			}
			if err = p.factory.Setup(p.name, &p.cfg); err != nil {
				return err
			}
			status[p.key()] = struct{}{}
		}
		if pluginNum == len(pluginChannel) { // 取出来又原封不动塞回去，说明没有一个插件setup成功，循环依赖了或者依赖项注册失败
			return fmt.Errorf("cycle depends or depends setup failed")
		}
		pluginNum = len(pluginChannel) //一轮遍历后，更新数量
	}
	return nil
}

func (p *pluginInfo) key() string {
	return p.typ + connTypeName + p.name
}

// 校验是否有依赖项
func (p *pluginInfo) hasDepends() (bool, error) {
	ds, ok := p.factory.(Depender)
	if !ok || len(ds.DependsOn()) == 0 {
		return false, nil
	}
	return true, nil
}

// 检查依赖项是否已经初始化
func (p *pluginInfo) checkDepends(status map[string]struct{}) (bool, error) {
	ds, ok := p.factory.(Depender)
	if !ok || len(ds.DependsOn()) == 0 {
		return false, nil
	}

	// 判断依赖项是否已经初始化
	for _, depend := range ds.DependsOn() {
		if p.key() == depend {
			return false, fmt.Errorf("plugin not allowed to depend on itself")
		}
		if _, ok2 := status[depend]; !ok2 {
			return true, nil
		}
	}
	return false, nil
}
