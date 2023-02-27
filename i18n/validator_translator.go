package i18n

import (
	"fmt"

	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
)

const (
	validatorTranslatorKeyFormat = "valid_%v"
)

// 封装validator使用的翻译
// 默认注册validator的翻译时tag很短很容易冲突，包装翻译和添加方法给key加前缀

type validatorTranslator struct {
	ut.Translator
}

func NewValidatorTranslator(trans ut.Translator) *validatorTranslator {
	return &validatorTranslator{Translator: trans}
}

func (t *validatorTranslator) Add(key interface{}, text string, override bool) error {
	return t.Translator.Add(t.getKey(key), text, override)
}

func (t *validatorTranslator) AddCardinal(key interface{}, text string, rule locales.PluralRule, override bool) error {
	return t.Translator.AddCardinal(t.getKey(key), text, rule, override)
}

func (t *validatorTranslator) AddOrdinal(key interface{}, text string, rule locales.PluralRule, override bool) error {
	return t.Translator.AddOrdinal(t.getKey(key), text, rule, override)
}

func (t *validatorTranslator) AddRange(key interface{}, text string, rule locales.PluralRule, override bool) error {
	return t.Translator.AddRange(t.getKey(key), text, rule, override)
}

func (t *validatorTranslator) T(key interface{}, params ...string) (string, error) {
	return t.Translator.T(t.getKey(key), params...)
}

func (t *validatorTranslator) C(key interface{}, num float64, digits uint64, param string) (string, error) {
	return t.Translator.C(t.getKey(key), num, digits, param)
}

func (t *validatorTranslator) O(key interface{}, num float64, digits uint64, param string) (string, error) {
	return t.Translator.O(t.getKey(key), num, digits, param)
}

func (t *validatorTranslator) R(key interface{}, num1 float64, digits1 uint64, num2 float64, digits2 uint64, param1, param2 string) (string, error) {
	return t.Translator.R(t.getKey(key), num1, digits1, num2, digits2, param1, param2)
}

func (t *validatorTranslator) getKey(key interface{}) interface{} {
	return fmt.Sprintf(validatorTranslatorKeyFormat, key)
}
