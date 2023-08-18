package prometheus

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fankane/go-utils/plugin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

const (
	pluginType = "monitor"
	pluginName = "prometheus"
)

var (
	defaultPath    = "/metrics"
	DefaultFactory = &Factory{}
	mu             sync.RWMutex

	CollGauge     CollectType = "gauge"
	CollCounter   CollectType = "counter"
	CollHistogram CollectType = "histogram"
	CollSummary   CollectType = "summary"

	collectorList = make(map[CollectType]map[string]prometheus.Collector)
)

var okCollTypeList = []CollectType{
	CollGauge,
	CollCounter,
	CollHistogram,
	CollSummary,
}

type CollectType string

type Config struct {
	Port           int              `yaml:"port" validate:"gt=0,lte=65535"`
	Path           string           `yaml:"path"`
	CustomCollects []*CustomCollect `yaml:"custom_collects"`
}

type CustomCollect struct {
	CollType string                  `yaml:"coll_type" validate:"required,oneof=gauge counter"`
	Info     map[string]*CollectInfo `yaml:"info"` //key: 上报名称
}

type CollectInfo struct {
	Help       string              `yaml:"help"`
	Labels     []string            `yaml:"labels"`
	Buckets    []float64           `yaml:"buckets"`    //Histogram 配置
	Objectives map[float64]float64 `yaml:"objectives"` //summary 配置
}

func init() {
	plugin.Register(pluginName, DefaultFactory)
}

type Factory struct {
}

// Type 日志插件类型
func (f *Factory) Type() string {
	return pluginType
}

// Setup 启动加载log配置 并注册日志
func (f *Factory) Setup(name string, node *yaml.Node) error {
	conf := &Config{}
	if err := node.Decode(&conf); err != nil {
		return fmt.Errorf("decode err:%s", err)
	}
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("prometheus conf invalided err:%s", err)
	}
	path := conf.Path
	if strings.TrimSpace(path) == "" {
		path = defaultPath
	}
	http.Handle(path, promhttp.Handler())
	if err := InitCollects(conf.CustomCollects); err != nil {
		return err
	}

	var err error
	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	}()
	time.Sleep(time.Millisecond * 3) //等待3毫秒，让 ListenAndServe 完成, 如果失败，3ms足够返回error
	if err != nil {
		return fmt.Errorf("prometheus listen err:%s", err)
	}
	return nil
}

// RegisteredCollTypeList 通过配置文件注册的采集器类型列表
func RegisteredCollTypeList() []CollectType {
	list := make([]CollectType, 0)
	for collectType, _ := range collectorList {
		list = append(list, collectType)
	}
	return list
}

// RegisteredCollNameList  通过配置文件注册的采集器里面的指标名称列表
func RegisteredCollNameList(collType CollectType) []string {
	names := make([]string, 0)
	for collectType, nameMap := range collectorList {
		if collectType == collType {
			for name, _ := range nameMap {
				names = append(names, name)
			}
			return names
		}
	}
	return []string{}
}

func GetGaugeVec(name string) *prometheus.GaugeVec {
	res := getCollector(name, CollGauge)
	if res == nil {
		return nil
	}
	return res.(*prometheus.GaugeVec)
}
func GetCounterVec(name string) *prometheus.CounterVec {
	res := getCollector(name, CollCounter)
	if res == nil {
		return nil
	}
	return res.(*prometheus.CounterVec)
}
func GetHistogram(name string) *prometheus.HistogramVec {
	res := getCollector(name, CollHistogram)
	if res == nil {
		return nil
	}
	return res.(*prometheus.HistogramVec)
}
func GetSummary(name string) *prometheus.SummaryVec {
	res := getCollector(name, CollSummary)
	if res == nil {
		return nil
	}
	return res.(*prometheus.SummaryVec)
}

func getCollector(name string, collType CollectType) prometheus.Collector {
	existMap, ok := collectorList[collType]
	if !ok || existMap == nil {
		return nil
	}
	g, ok := existMap[name]
	if !ok {
		return nil
	}
	return g
}

func InitCollects(collects []*CustomCollect) error {
	for _, collect := range collects {
		if err := checkCollType(collect.CollType); err != nil {
			return err
		}
		createCollVec(CollectType(collect.CollType), collect.Info)
	}
	return nil
}

func checkCollType(collType string) error {
	for _, collectType := range okCollTypeList {
		if collectType == CollectType(collType) {
			return nil
		}
	}
	return fmt.Errorf("unsuport collect type:%s", collType)
}

func createCollVec(collType CollectType, collect map[string]*CollectInfo) {
	for name, info := range collect {
		var cs prometheus.Collector
		switch collType {
		case CollGauge:
			cs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Name: name,
				Help: info.Help,
			}, info.Labels)
		case CollCounter:
			cs = prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: name,
				Help: info.Help,
			}, info.Labels)
		case CollHistogram:
			cs = prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name:    name,
				Help:    info.Help,
				Buckets: info.Buckets,
			}, info.Labels)
		case CollSummary:
			cs = prometheus.NewSummaryVec(prometheus.SummaryOpts{
				Name:       name,
				Help:       info.Help,
				Objectives: info.Objectives,
			}, info.Labels)
		}
		setCollector(collType, name, cs)
		prometheus.MustRegister(cs)
	}
}

func setCollector(collType CollectType, name string, coll prometheus.Collector) {
	if _, ok := collectorList[collType]; !ok {
		collectorList[collType] = map[string]prometheus.Collector{
			name: coll,
		}
	} else {
		collectorList[collType][name] = coll
	}
}
