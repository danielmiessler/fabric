package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/danielmiessler/fabric/plugins/ai/exolab"
	"github.com/danielmiessler/fabric/plugins/strategy"

	"github.com/samber/lo"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins"
	"github.com/danielmiessler/fabric/plugins/ai"
	"github.com/danielmiessler/fabric/plugins/ai/anthropic"
	"github.com/danielmiessler/fabric/plugins/ai/azure"
	"github.com/danielmiessler/fabric/plugins/ai/dryrun"
	"github.com/danielmiessler/fabric/plugins/ai/gemini"
	"github.com/danielmiessler/fabric/plugins/ai/lmstudio"
	"github.com/danielmiessler/fabric/plugins/ai/ollama"
	"github.com/danielmiessler/fabric/plugins/ai/openai"
	"github.com/danielmiessler/fabric/plugins/ai/openai_compatible"
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"github.com/danielmiessler/fabric/plugins/template"
	"github.com/danielmiessler/fabric/plugins/tools"
	"github.com/danielmiessler/fabric/plugins/tools/jina"
	"github.com/danielmiessler/fabric/plugins/tools/lang"
	"github.com/danielmiessler/fabric/plugins/tools/youtube"
)

func NewPluginRegistry(db *fsdb.Db) (ret *PluginRegistry, err error) {
	ret = &PluginRegistry{
		Db:             db,
		VendorManager:  ai.NewVendorsManager(),
		VendorsAll:     ai.NewVendorsManager(),
		PatternsLoader: tools.NewPatternsLoader(db.Patterns),
		YouTube:        youtube.NewYouTube(),
		Language:       lang.NewLanguage(),
		Jina:           jina.NewClient(),
		Strategies:     strategy.NewStrategiesManager(),
	}

	var homedir string
	if homedir, err = os.UserHomeDir(); err != nil {
		return
	}
	ret.TemplateExtensions = template.NewExtensionManager(filepath.Join(homedir, ".config/fabric"))

	ret.Defaults = tools.NeeDefaults(ret.GetModels)

	// Create a vendors slice to hold all vendors (order doesn't matter initially)
	vendors := []ai.Vendor{}

	// Add non-OpenAI compatible clients
	vendors = append(vendors,
		openai.NewClient(),
		ollama.NewClient(),
		azure.NewClient(),
		gemini.NewClient(),
		anthropic.NewClient(),
		lmstudio.NewClient(),
		exolab.NewClient(),
	)

	// Add all OpenAI-compatible providers
	for providerName := range openai_compatible.ProviderMap {
		provider, _ := openai_compatible.GetProviderByName(providerName)
		vendors = append(vendors, openai_compatible.NewClient(provider))
	}

	// Sort vendors by name for consistent ordering (case-insensitive)
	sort.Slice(vendors, func(i, j int) bool {
		return strings.ToLower(vendors[i].GetName()) < strings.ToLower(vendors[j].GetName())
	})

	// Add all sorted vendors to VendorsAll
	ret.VendorsAll.AddVendors(vendors...)
	_ = ret.Configure()

	return
}

func (o *PluginRegistry) ListVendors(out io.Writer) error {
	vendors := lo.Map(o.VendorsAll.Vendors, func(vendor ai.Vendor, _ int) string {
		return vendor.GetName()
	})
	fmt.Fprint(out, "Available Vendors:\n\n")
	for _, vendor := range vendors {
		fmt.Fprintf(out, "%s\n", vendor)
	}
	return nil
}

type PluginRegistry struct {
	Db *fsdb.Db

	VendorManager      *ai.VendorsManager
	VendorsAll         *ai.VendorsManager
	Defaults           *tools.Defaults
	PatternsLoader     *tools.PatternsLoader
	YouTube            *youtube.YouTube
	Language           *lang.Language
	Jina               *jina.Client
	TemplateExtensions *template.ExtensionManager
	Strategies         *strategy.StrategiesManager
}

func (o *PluginRegistry) SaveEnvFile() (err error) {
	// Now create the .env with all configured VendorsController info
	var envFileContent bytes.Buffer

	o.Defaults.Settings.FillEnvFileContent(&envFileContent)
	o.PatternsLoader.SetupFillEnvFileContent(&envFileContent)
	o.Strategies.SetupFillEnvFileContent(&envFileContent)

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
	groupsPlugins := common.NewGroupsItemsSelector("Available plugins (please configure all required plugins):",
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

	groupsPlugins.AddGroupItems("Tools", o.Defaults, o.Jina, o.Language, o.PatternsLoader, o.Strategies, o.YouTube)

	for {
		groupsPlugins.Print(false)

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
				var vendor ai.Vendor
				if vendor, ok = plugin.(ai.Vendor); ok {
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

func (o *PluginRegistry) ConfigureVendors() {
	o.VendorManager.Clear()
	for _, vendor := range o.VendorsAll.Vendors {
		if vendorErr := vendor.Configure(); vendorErr == nil {
			o.VendorManager.AddVendors(vendor)
		}
	}
}

func (o *PluginRegistry) GetModels() (ret *ai.VendorsModels, err error) {
	o.ConfigureVendors()
	ret, err = o.VendorManager.GetModels()
	return
}

// Configure buildClient VendorsController based on the environment variables
func (o *PluginRegistry) Configure() (err error) {
	o.ConfigureVendors()
	_ = o.Defaults.Configure()
	_ = o.PatternsLoader.Configure()

	//YouTube and Jina are not mandatory, so ignore not configured error
	_ = o.YouTube.Configure()
	_ = o.Jina.Configure()
	_ = o.Language.Configure()
	return
}

func (o *PluginRegistry) GetChatter(model string, modelContextLength int, strategy string, stream bool, dryRun bool) (ret *Chatter, err error) {
	ret = &Chatter{
		db:     o.Db,
		Stream: stream,
		DryRun: dryRun,
	}

	defaultModel := o.Defaults.Model.Value
	defaultModelContextLength, err := strconv.Atoi(o.Defaults.ModelContextLength.Value)
	defaultVendor := o.Defaults.Vendor.Value
	vendorManager := o.VendorManager

	if err != nil {
		defaultModelContextLength = 0
		err = nil
	}

	ret.modelContextLength = modelContextLength
	if ret.modelContextLength == 0 {
		ret.modelContextLength = defaultModelContextLength
	}

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
		var errMsg string
		if defaultModel == "" || defaultVendor == "" {
			errMsg = "Please run, fabric --setup, and select default model and vendor."
		} else {
			errMsg = "could not find vendor."
		}
		err = fmt.Errorf(
			" Requested Model = %s\n Default Model = %s\n Default Vendor = %s.\n\n%s",
			model, defaultModel, defaultVendor, errMsg)
		return
	}
	ret.strategy = strategy
	return
}
