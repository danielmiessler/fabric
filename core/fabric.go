package core

import (
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/danielmiessler/fabric/plugins/ai"
	"github.com/danielmiessler/fabric/plugins/ai/anthropic"
	"github.com/danielmiessler/fabric/plugins/ai/azure"
	"github.com/danielmiessler/fabric/plugins/ai/dryrun"
	"github.com/danielmiessler/fabric/plugins/ai/gemini"
	"github.com/danielmiessler/fabric/plugins/ai/groq"
	"github.com/danielmiessler/fabric/plugins/ai/mistral"
	"github.com/danielmiessler/fabric/plugins/ai/ollama"
	"github.com/danielmiessler/fabric/plugins/ai/openai"
	"github.com/danielmiessler/fabric/plugins/ai/openrouter"
	"github.com/danielmiessler/fabric/plugins/ai/siliconcloud"
	core2 "github.com/danielmiessler/fabric/plugins/core"
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"github.com/danielmiessler/fabric/plugins/tools/jina"
	"github.com/danielmiessler/fabric/plugins/tools/lang"
	"github.com/danielmiessler/fabric/plugins/tools/youtube"
	"github.com/pkg/errors"
	"os"
)

const NoSessionPatternUserMessages = "no session, pattern or user messages provided"

func NewFabric(db *fsdb.Db) (ret *Fabric, err error) {
	ret = NewFabricBase(db)
	err = ret.Configure()
	return
}

func NewFabricForSetup(db *fsdb.Db) (ret *Fabric) {
	ret = NewFabricBase(db)
	_ = ret.Configure()
	return
}

// NewFabricBase Create a new Fabric from a list of already configured VendorsController
func NewFabricBase(db *fsdb.Db) (ret *Fabric) {

	ret = &Fabric{
		VendorManager:  ai.NewVendorsManager(),
		Db:             db,
		Defaults:       core2.NeeDefaults(),
		VendorsAll:     ai.NewVendorsManager(),
		PatternsLoader: core2.NewPatternsLoader(db.Patterns),
		YouTube:        youtube.NewYouTube(),
		Language:       lang.NewLanguage(),
		Jina:           jina.NewClient(),
	}

	ret.VendorsAll.AddVendors(openai.NewClient(), azure.NewClient(), ollama.NewClient(), groq.NewClient(),
		gemini.NewClient(), anthropic.NewClient(), siliconcloud.NewClient(), openrouter.NewClient(), mistral.NewClient())

	return
}

type Fabric struct {
	VendorManager  *ai.VendorsManager
	VendorsAll     *ai.VendorsManager
	PatternsLoader *core2.PatternsLoader
	YouTube        *youtube.YouTube
	Language       *lang.Language
	Jina           *jina.Client

	Db       *fsdb.Db
	Defaults *core2.Defaults
}

type ChannelName struct {
	channel chan []string
	name    string
}

func (o *Fabric) SaveEnvFile() (err error) {
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

func (o *Fabric) Setup() (err error) {
	if err = o.SetupVendors(); err != nil {
		return
	}

	if err = o.Defaults.Setup(o.VendorManager.GetModels()); err != nil {
		return
	}
	if err = o.SaveEnvFile(); err != nil {
		return
	}

	_ = o.YouTube.SetupOrSkip()

	if err = o.Jina.SetupOrSkip(); err != nil {
		return
	}

	if err = o.PatternsLoader.Setup(); err != nil {
		return
	}

	if err = o.Language.SetupOrSkip(); err != nil {
		return
	}

	err = o.SaveEnvFile()

	return
}

func (o *Fabric) SetupVendors() (err error) {
	o.VendorManager.Models = nil
	if o.VendorManager.Vendors, err = o.VendorsAll.Setup(); err != nil {
		return
	}

	if !o.VendorManager.HasVendors() {
		err = errors.New("No vendors configured")
		return
	}

	err = o.SaveEnvFile()

	return
}

func (o *Fabric) SetupVendor(vendorName string) (err error) {
	if err = o.VendorsAll.SetupVendor(vendorName, o.VendorManager.Vendors); err != nil {
		return
	}
	err = o.SaveEnvFile()
	return
}

// Configure buildClient VendorsController based on the environment variables
func (o *Fabric) Configure() (err error) {
	if err = o.Defaults.Configure(); err != nil {
		return
	}

	for _, vendor := range o.VendorsAll.Vendors {
		if vendorErr := vendor.Configure(); vendorErr == nil {
			o.VendorManager.AddVendors(vendor)
		}
	}
	if err = o.PatternsLoader.Configure(); err != nil {
		return
	}

	//YouTube and Jina are not mandatory, so ignore not configured error
	_ = o.YouTube.Configure()
	_ = o.Jina.Configure()
	_ = o.Language.Configure()

	return
}

func (o *Fabric) GetChatter(model string, stream bool, dryRun bool) (ret *Chatter, err error) {
	ret = &Chatter{
		db:     o.Db,
		Stream: stream,
		DryRun: dryRun,
	}

	if dryRun {
		ret.vendor = dryrun.NewClient()
		ret.model = model
		if ret.model == "" {
			ret.model = o.Defaults.Model.Value
		}
	} else if model == "" {
		ret.vendor = o.VendorManager.FindByName(o.Defaults.Vendor.Value)
		ret.model = o.Defaults.Model.Value
	} else {
		ret.vendor = o.VendorManager.FindByName(o.VendorManager.GetModels().FindVendorsByModelFirst(model))
		ret.model = model
	}

	if ret.vendor == nil {
		err = fmt.Errorf(
			"could not find vendor.\n Model = %s\n Model = %s\n Vendor = %s",
			model, o.Defaults.Model.Value, o.Defaults.Vendor.Value)
		return
	}
	return
}

func (o *Fabric) CopyToClipboard(message string) (err error) {
	if err = clipboard.WriteAll(message); err != nil {
		err = fmt.Errorf("could not copy to clipboard: %v", err)
	}
	return
}

func (o *Fabric) CreateOutputFile(message string, fileName string) (err error) {
	var file *os.File
	if file, err = os.Create(fileName); err != nil {
		err = fmt.Errorf("error creating file: %v", err)
		return
	}
	defer file.Close()
	if _, err = file.WriteString(message); err != nil {
		err = fmt.Errorf("error writing to file: %v", err)
	}
	return
}
