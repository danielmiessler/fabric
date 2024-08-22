package common

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

const AnswerReset = "reset"

type Configurable struct {
	Settings
	SetupQuestions

	Label         string
	EnvNamePrefix string

	ConfigureCustom func() error
}

func (o *Configurable) GetName() string {
	return o.Label
}

func (o *Configurable) AddSetting(name string, required bool) (ret *Setting) {
	ret = NewSetting(fmt.Sprintf("%v%v", o.EnvNamePrefix, BuildEnvVariable(name)), required)
	o.Settings = append(o.Settings, ret)
	return
}

func (o *Configurable) AddSetupQuestion(name string, required bool) (ret *SetupQuestion) {
	return o.AddSetupQuestionCustom(name, required, "")
}

func (o *Configurable) AddSetupQuestionCustom(name string, required bool, question string) (ret *SetupQuestion) {
	setting := o.AddSetting(name, required)
	ret = &SetupQuestion{Setting: setting, Question: question}
	if ret.Question == "" {
		ret.Question = fmt.Sprintf("Enter your %v %v", o.Label, strings.ToUpper(name))
	}
	o.SetupQuestions = append(o.SetupQuestions, ret)
	return
}

func (o *Configurable) Configure() (err error) {
	if err = o.Settings.Configure(); err != nil {
		return
	}

	if o.ConfigureCustom != nil {
		err = o.ConfigureCustom()
	}
	return
}

func (o *Configurable) Setup() (err error) {
	if err = o.Ask(o.Label); err != nil {
		return
	}

	err = o.Configure()
	return
}

func (o *Configurable) SetupOrSkip() (err error) {
	if err = o.Setup(); err != nil {
		fmt.Printf("[%v] skipped\n", o.GetName())
	}
	return
}

func (o *Configurable) SetupFillEnvFileContent(fileEnvFileContent *bytes.Buffer) {
	o.Settings.FillEnvFileContent(fileEnvFileContent)
}

func NewSetting(envVariable string, required bool) *Setting {
	return &Setting{
		EnvVariable: envVariable,
		Required:    required,
	}
}

type Setting struct {
	EnvVariable string
	Value       string
	Required    bool
}

func (o *Setting) IsValid() bool {
	return o.IsDefined() || !o.Required
}

func (o *Setting) IsValidErr() (err error) {
	if !o.IsValid() {
		err = fmt.Errorf("%v=%v, is not valid", o.EnvVariable, o.Value)
	}
	return
}

func (o *Setting) IsDefined() bool {
	return o.Value != ""
}

func (o *Setting) Configure() error {
	if o.Value == "" {
		o.Value = os.Getenv(o.EnvVariable)
	}
	return o.IsValidErr()
}

func (o *Setting) FillEnvFileContent(buffer *bytes.Buffer) {
	if o.IsDefined() {
		buffer.WriteString(o.EnvVariable)
		buffer.WriteString("=")
		//buffer.WriteString("\"")
		buffer.WriteString(o.Value)
		//buffer.WriteString("\"")
		buffer.WriteString("\n")
	}
	return
}

func (o *Setting) Print() {
	fmt.Printf("%v: %v\n", o.EnvVariable, o.Value)
}

type SetupQuestion struct {
	*Setting
	Question string
}

func (o *SetupQuestion) Ask(label string) (err error) {
	var prefix string

	if label != "" {
		prefix = fmt.Sprintf("[%v] ", label)
	} else {
		prefix = ""
	}

	fmt.Println()
	if o.Value != "" {
		fmt.Printf("%v%v (leave empty for '%s' or type '%v' to remove the value):\n",
			prefix, o.Question, o.Value, AnswerReset)
	} else {
		fmt.Printf("%v%v (leave empty to skip):\n", prefix, o.Question)
	}

	var answer string
	fmt.Scanln(&answer)
	answer = strings.TrimRight(answer, "\n")
	if answer == "" {
		answer = o.Value
	} else if strings.ToLower(answer) == AnswerReset {
		answer = ""
	}
	err = o.OnAnswer(answer)
	return
}

func (o *SetupQuestion) OnAnswer(answer string) (err error) {
	o.Value = answer
	err = o.IsValidErr()
	return
}

type Settings []*Setting

func (o Settings) IsConfigured() (ret bool) {
	ret = true
	for _, setting := range o {
		if ret = setting.IsValid(); !ret {
			break
		}
	}
	return
}

func (o Settings) Configure() (err error) {
	for _, setting := range o {
		if err = setting.Configure(); err != nil {
			break
		}
	}
	return
}

func (o Settings) FillEnvFileContent(buffer *bytes.Buffer) {
	for _, setting := range o {
		setting.FillEnvFileContent(buffer)
	}
	return
}

type SetupQuestions []*SetupQuestion

func (o SetupQuestions) Ask(label string) (err error) {
	fmt.Println()
	fmt.Printf("[%v]\n", label)
	for _, question := range o {
		if err = question.Ask(""); err != nil {
			break
		}
	}
	return
}

func BuildEnvVariablePrefix(name string) (ret string) {
	ret = BuildEnvVariable(name)
	if ret != "" {
		ret += "_"
	}
	return
}

func BuildEnvVariable(name string) string {
	name = strings.TrimSpace(name)
	return strings.ReplaceAll(strings.ToUpper(name), " ", "_")
}
