package i18n

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/richkeyu/gocommons/config"
	"github.com/richkeyu/gocommons/i18n/locales/cn"
	"github.com/richkeyu/gocommons/i18n/locales/tw"
	"github.com/richkeyu/gocommons/plog"
	"github.com/richkeyu/gocommons/server"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/ar"
	"github.com/go-playground/locales/bg"
	"github.com/go-playground/locales/cs"
	"github.com/go-playground/locales/de"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/he"
	"github.com/go-playground/locales/id"
	"github.com/go-playground/locales/it"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/ko"
	"github.com/go-playground/locales/pl"
	"github.com/go-playground/locales/pt"
	"github.com/go-playground/locales/ro"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/locales/th"
	"github.com/go-playground/locales/tr"
	"github.com/go-playground/locales/vi"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	validator_ar "github.com/go-playground/validator/v10/translations/ar"
	validator_en "github.com/go-playground/validator/v10/translations/en"
	validator_es "github.com/go-playground/validator/v10/translations/es"
	validator_fr "github.com/go-playground/validator/v10/translations/fr"
	validator_id "github.com/go-playground/validator/v10/translations/id"
	validator_it "github.com/go-playground/validator/v10/translations/it"
	validator_ja "github.com/go-playground/validator/v10/translations/ja"
	validator_pt "github.com/go-playground/validator/v10/translations/pt"
	validator_ru "github.com/go-playground/validator/v10/translations/ru"
	validator_tr "github.com/go-playground/validator/v10/translations/tr"
	validator_vi "github.com/go-playground/validator/v10/translations/vi"
	validator_zh "github.com/go-playground/validator/v10/translations/zh"
	validator_zh_tw "github.com/go-playground/validator/v10/translations/zh_tw"
)

// 根据开源项目实现
// https://github.com/go-playground/locales
// https://github.com/go-playground/universal-translator/

const (
	configKeyName         = "i18n"
	validatorFieldTagName = "name"
)

const (
	ErrorReport = "error"
	ErrorInfo   = "info"
	ErrorIgnore = "ignore"
)

var (
	trans          *ut.UniversalTranslator         // 全局翻译包含全部语言
	defaultTrans   Translator                      // 默认语言
	validatorTrans map[string]*validatorTranslator // validator 翻译
)

var conf Config

type Config struct {
	BaseUri  string `json:"base_uri" yaml:"base_uri"`
	Category string `json:"category" yaml:"category"`
	Default  string `json:"default" yaml:"default"`
	Support  string `json:"support" yaml:"support"`
	Error    string `json:"error" yaml:"error"`
}

func (c Config) GetSupportList() []string {
	return strings.Split(strings.Trim(c.Support, ","), ",")
}

// Init 初始化全部翻译 远程加载翻译内容 启动更新翻译内容协程
func Init() {
	// 配置
	err := config.Load(configKeyName, &conf)
	if err != nil {
		panic(fmt.Sprintf("load i18n config fail: %s", err))
	}

	// 默认翻译
	defaultLocalTrans := getLocalTranslator(conf.Default)
	if defaultLocalTrans == nil {
		panic(fmt.Sprintf("default lang not found: %s", conf.Default))
	}

	// 初始化翻译器
	trans = ut.New(defaultLocalTrans, defaultLocalTrans)
	for _, lang := range conf.GetSupportList() {
		_, isFound := trans.GetTranslator(lang)
		if !isFound {
			localTrans := getLocalTranslator(lang)
			if localTrans == nil {
				panic(fmt.Sprintf("support lang not found: %s", lang))
			}
			err = trans.AddTranslator(localTrans, true)
			if err != nil {
				panic(fmt.Sprintf("support lang add fail: %s; lang:%s", err, lang))
			}
		}
	}

	// 校验
	err = trans.VerifyTranslations()
	if err != nil {
		panic(fmt.Sprintf("erify translation fail: %s", conf.Default))
	}

	// 默认语言
	defaultTranslator, _ := trans.GetTranslator(conf.Default)
	defaultTrans = newTranslator(defaultTranslator)

	// validator 翻译 仅初始化时执行一次
	validatorTrans = make(map[string]*validatorTranslator, len(conf.GetSupportList()))
	for _, lang := range conf.GetSupportList() {
		t, isFound := trans.GetTranslator(lang)
		if isFound {
			if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
				validatorTrans[t.Locale()] = NewValidatorTranslator(t)
				err = registerValidatorTranslations(lang, v, validatorTrans[t.Locale()])
				if err != nil {
					fmt.Println("registerValidatorTranslations fail: ", err)
				}
			}
		}
	}
	// validator 自定义字段名称
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("name"), ",", 2)[0]
			if len(name) > 0 {
				return name
			}
			return fld.Name
		})
	}

	// 异步加载语言内容 防止阻塞启动进程
	go func() {
		ctx := server.NewContext(context.Background(), &gin.Context{})
		loadFromBase(ctx, conf)
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				loadFromBase(ctx, conf)
			}
		}
	}()
}

func loadFromBase(ctx context.Context, conf Config) {
	allBaseTrans, err := getTransFromBase(ctx, conf)
	if err != nil {
		plog.Infof(ctx, "i18n load from base service fail: %s \n", err)
		return
	}
	for lang, keys := range allBaseTrans {
		// 获取翻译
		etTranslator, isFound := trans.GetTranslator(lang)
		if !isFound {
			plog.Errorf(ctx, "translator not found: %s", lang)
			continue
		}

		// 加载数据库翻译
		for key, value := range keys {
			err = etTranslator.Add(key, value, true)
			if err != nil {
				plog.Infof(nil, "translator add fail: %s;key: %s; value: %s", err, key, value)
			}
		}
	}
}

func T(ctx context.Context, str string, params ...string) string {
	request := server.FromContext(ctx)
	if request != nil {
		if transInter, ok := request.Get(translatorContextKey); ok {
			if t, ok := transInter.(Translator); ok {
				return t.T(ctx, str, params...)
			}
		}
	}
	return defaultTrans.T(ctx, str, params...)
}

func GetValidatorTranslator(lang string) *validatorTranslator {
	return validatorTrans[lang]
}

// 系统增加新语言支持时需要增加这里
func getLocalTranslator(lang string) locales.Translator {
	switch lang {
	case "en":
		return en.New()
	case "zh":
		return zh.New()
	case "cn": // 自定义的名称
		return cn.New()
	case "ar":
		return ar.New()
	case "bg":
		return bg.New()
	case "pl":
		return pl.New()
	case "de":
		return de.New()
	case "ru":
		return ru.New()
	case "fr":
		return fr.New()
	case "ko":
		return ko.New()
	case "cs":
		return cs.New()
	case "ro":
		return ro.New()
	case "pt":
		return pt.New()
	case "ja":
		return ja.New()
	case "th":
		return th.New()
	case "tr":
		return tr.New()
	case "es":
		return es.New()
	case "he":
		return he.New()
	case "id":
		return id.New()
	case "it":
		return it.New()
	case "vi":
		return vi.New()
	case "tw": // 自定义的名称
		return tw.New()
	}
	return nil
}

// 系统增加新语言支持时需要增加这里
func registerValidatorTranslations(lang string, v *validator.Validate, trans ut.Translator) error {
	switch lang {
	case "en":
		return validator_en.RegisterDefaultTranslations(v, trans)
	case "zh":
		return validator_zh.RegisterDefaultTranslations(v, trans)
	case "cn": // 自定义的名称
		return validator_zh.RegisterDefaultTranslations(v, trans)
	case "ar":
		return validator_ar.RegisterDefaultTranslations(v, trans)
	case "bg": // validator 不存在翻译的用英语
		return validator_en.RegisterDefaultTranslations(v, trans)
	case "pl":
		return validator_en.RegisterDefaultTranslations(v, trans)
	case "de":
		return validator_en.RegisterDefaultTranslations(v, trans)
	case "ru":
		return validator_ru.RegisterDefaultTranslations(v, trans)
	case "fr":
		return validator_fr.RegisterDefaultTranslations(v, trans)
	case "ko":
		return validator_en.RegisterDefaultTranslations(v, trans)
	case "cs":
		return validator_en.RegisterDefaultTranslations(v, trans)
	case "ro":
		return validator_en.RegisterDefaultTranslations(v, trans)
	case "pt":
		return validator_pt.RegisterDefaultTranslations(v, trans)
	case "ja":
		return validator_ja.RegisterDefaultTranslations(v, trans)
	case "th":
		return validator_en.RegisterDefaultTranslations(v, trans)
	case "tr":
		return validator_tr.RegisterDefaultTranslations(v, trans)
	case "es":
		return validator_es.RegisterDefaultTranslations(v, trans)
	case "he":
		return validator_en.RegisterDefaultTranslations(v, trans)
	case "id":
		return validator_id.RegisterDefaultTranslations(v, trans)
	case "it":
		return validator_it.RegisterDefaultTranslations(v, trans)
	case "vi":
		return validator_vi.RegisterDefaultTranslations(v, trans)
	case "tw": // 自定义的名称
		return validator_zh_tw.RegisterDefaultTranslations(v, trans)
	default:
		return validator_en.RegisterDefaultTranslations(v, trans)
	}
	return nil
}
