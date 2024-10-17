package core

import (
	"bytes"
	"fmt"
	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins/ai/azure"
	"github.com/danielmiessler/fabric/plugins/tools"
	"github.com/samber/lo"
	"strconv"

	"github.com/danielmiessler/fabric/plugins"
	"github.com/danielmiessler/fabric/plugins/ai"
	"github.com/danielmiessler/fabric/plugins/ai/anthropic"
	"github.com/danielmiessler/fabric/plugins/ai/dryrun"
	"github.com/danielmiessler/fabric/plugins/ai/gemini"
	"github.com/danielmiessler/fabric/plugins/ai/groq"
	"github.com/danielmiessler/fabric/plugins/ai/mistral"
	"github.com/danielmiessler/fabric/plugins/ai/ollama"
	"github.com/danielmiessler/fabric/plugins/ai/openai"
	"github.com/danielmiessler/fabric/plugins/ai/openrouter"
	"github.com/danielmiessler/fabric/plugins/ai/siliconcloud"
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"github.com/danielmiessler/fabric/plugins/tools/jina"
	"github.com/danielmiessler/fabric/plugins/tools/lang"
	"github.com/danielmiessler/fabric/plugins/tools/youtube"
)

func NewPluginRegistry(db *fsdb.Db) (ret *PluginRegistry) {
	ret = &PluginRegistry{
		Db:             db,
		VendorManager:  ai.NewVendorsManager(),
		VendorsAll:     ai.NewVendorsManager(),
		PatternsLoader: tools.NewPatternsLoader(db.Patterns),
		YouTube:        youtube.NewYouTube(),
		Language:       lang.NewLanguage(),
		Jina:           jina.NewClient(),
	}

	ret.Defaults = tools.NeeDefaults(ret.VendorManager.GetModels)

	ret.VendorsAll.AddVendors(openai.NewClient(), ollama.NewClient(), azure.NewClient(), groq.NewClient(),
		gemini.NewClient(), anthropic.NewClient(), siliconcloud.NewClient(), openrouter.NewClient(), mistral.NewClient())
	_ = ret.Configure()

	return
}

type PluginRegistry struct {
	Db *fsdb.Db

	VendorManager  *ai.VendorsManager
	VendorsAll     *ai.VendorsManager
	Defaults       *tools.Defaults
	PatternsLoader *tools.PatternsLoader
	YouTube        *youtube.YouTube
	Language       *lang.Language
	Jina           *jina.Client
}

func (o *PluginRegistry) SaveEnvFile() (err error) {
	// Now create the .env with all configured VendorsController info
	var envFileContent bytes.Buffer

	o.Defaults.Settings.FillEnvFileContent(&envFileContent)
	o.PatternsLoader.SetupFillEnvFileContent(&envFileContent)

	for _, vendor := range o.VendorManager.Vendors {
		vendor.SetupFillEnvFileContent(&envFileContent)
	}

	o.YouTube.SetupFillEnvFileContent(&envFileContent)
	o.Jina.SetupFillEnvFileContent(&envFileContent)
	o.Language.SetupFillEnvFileContent(&envFileContent)

	err = o.Db.SaveEnv(envFileContent.String())
	return
}

func (o *PluginRegistry) Setup() (err error) {
	setupQuestion := plugins.NewSetupQuestion("Enter the number of the plugin to setup")
	groupsPlugins := common.NewGroupsItemsSelector[plugins.Plugin]("Available plugins",
		func(plugin plugins.Plugin) string {
			var configuredLabel string
			if plugin.IsConfigured() {
				configuredLabel = " (configured)"
			} else {
				configuredLabel = ""
			}
			return fmt.Sprintf("%v%v", plugin.GetSetupDescription(), configuredLabel)
		})

	groupsPlugins.AddGroupItems("AI Vendors [at least one, required]", lo.Map(o.VendorsAll.Vendors,
		func(vendor ai.Vendor, _ int) plugins.Plugin {
			return vendor
		})...)

	groupsPlugins.AddGroupItems("Tools", o.Defaults, o.PatternsLoader, o.YouTube, o.Language, o.Jina)

	for {
		groupsPlugins.Print()

		if answerErr := setupQuestion.Ask("Plugin Number"); answerErr != nil {
			break
		}

		if setupQuestion.Value == "" {
			break
		}
		number, parseErr := strconv.Atoi(setupQuestion.Value)
		setupQuestion.Value = ""

		if parseErr == nil {
			var plugin plugins.Plugin
			if _, plugin, err = groupsPlugins.GetGroupAndItemByItemNumber(number); err != nil {
				return
			}

			if pluginSetupErr := plugin.Setup(); pluginSetupErr != nil {
				println(pluginSetupErr.Error())
			} else {
				if err = o.SaveEnvFile(); err != nil {
					break
				}
			}

			if _, ok := o.VendorManager.VendorsByName[plugin.GetName()]; !ok {
				if vendor, ok := plugin.(ai.Vendor); ok {
					o.VendorManager.AddVendors(vendor)
				}
			}
		} else {
			break
		}
	}

	err = o.SaveEnvFile()

	return
}

func (o *PluginRegistry) SetupVendor(vendorName string) (err error) {
	if err = o.VendorsAll.SetupVendor(vendorName, o.VendorManager.VendorsByName); err != nil {
		return
	}
	err = o.SaveEnvFile()
	return
}

// Configure buildClient VendorsController based on the environment variables
func (o *PluginRegistry) Configure() (err error) {
	for _, vendor := range o.VendorsAll.Vendors {
		if vendorErr := vendor.Configure(); vendorErr == nil {
			o.VendorManager.AddVendors(vendor)
		}
	}
	_ = o.Defaults.Configure()
	_ = o.PatternsLoader.Configure()

	//YouTube and Jina are not mandatory, so ignore not configured error
	_ = o.YouTube.Configure()
	_ = o.Jina.Configure()
	_ = o.Language.Configure()
	return
}

func (o *PluginRegistry) GetChatter(model string, stream bool, dryRun bool) (ret *Chatter, err error) {
	ret = &Chatter{
		db:     o.Db,
		Stream: stream,
		DryRun: dryRun,
	}

	defaultModel := o.Defaults.Model.Value
	defaultVendor := o.Defaults.Vendor.Value
	vendorManager := o.VendorManager

	if dryRun {
		ret.vendor = dryrun.NewClient()
		ret.model = model
		if ret.model == "" {
			ret.model = defaultModel
		}
	} else if model == "" {
		ret.vendor = vendorManager.FindByName(defaultVendor)
		ret.model = defaultModel
	} else {
		var models *ai.VendorsModels
		if models, err = vendorManager.GetModels(); err != nil {
			return
		}
		ret.vendor = vendorManager.FindByName(models.FindGroupsByItemFirst(model))
		ret.model = model
	}

	if ret.vendor == nil {
		err = fmt.Errorf(
			"could not find vendor.\n Model = %s\n Model = %s\n Vendor = %s",
			model, defaultModel, defaultVendor)
		return
	}
	return
}
