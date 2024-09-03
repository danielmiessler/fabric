package core

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/db"
	"github.com/danielmiessler/fabric/vendors/anthropic"
	"github.com/danielmiessler/fabric/vendors/azure"
	"github.com/danielmiessler/fabric/vendors/dryrun"
	"github.com/danielmiessler/fabric/vendors/gemini"
	"github.com/danielmiessler/fabric/vendors/groc"
	"github.com/danielmiessler/fabric/vendors/ollama"
	"github.com/danielmiessler/fabric/vendors/openai"
	"github.com/danielmiessler/fabric/youtube"
	"github.com/pkg/errors"
)

const DefaultPatternsGitRepoUrl = "https://github.com/danielmiessler/fabric.git"
const DefaultPatternsGitRepoFolder = "patterns"

func NewFabric(db *db.Db) (ret *Fabric, err error) {
	ret = NewFabricBase(db)
	err = ret.Configure()
	return
}

func NewFabricForSetup(db *db.Db) (ret *Fabric) {
	ret = NewFabricBase(db)
	_ = ret.Configure()
	return
}

// NewFabricBase Create a new Fabric from a list of already configured VendorsController
func NewFabricBase(db *db.Db) (ret *Fabric) {

	ret = &Fabric{
		VendorsManager: NewVendorsManager(),
		Db:             db,
		VendorsAll:     NewVendorsManager(),
		PatternsLoader: NewPatternsLoader(db.Patterns),
		YouTube:        youtube.NewYouTube(),
	}

	label := "Default"
	ret.Configurable = &common.Configurable{
		Label:           label,
		EnvNamePrefix:   common.BuildEnvVariablePrefix(label),
		ConfigureCustom: ret.configure,
	}

	ret.DefaultVendor = ret.AddSetting("Vendor", true)
	ret.DefaultModel = ret.AddSetupQuestionCustom("Model", true,
		"Enter the index the name of your default model")

	ret.VendorsAll.AddVendors(openai.NewClient(), azure.NewClient(), ollama.NewClient(), groc.NewClient(),
		gemini.NewClient(), anthropic.NewClient())

	return
}

type Fabric struct {
	*common.Configurable
	*VendorsManager
	VendorsAll *VendorsManager
	*PatternsLoader
	*youtube.YouTube

	Db *db.Db

	DefaultVendor *common.Setting
	DefaultModel  *common.SetupQuestion
}

type ChannelName struct {
	channel chan []string
	name    string
}

func (o *Fabric) SaveEnvFile() (err error) {
	// Now create the .env with all configured VendorsController info
	var envFileContent bytes.Buffer

	o.Settings.FillEnvFileContent(&envFileContent)
	o.PatternsLoader.SetupFillEnvFileContent(&envFileContent)

	for _, vendor := range o.Vendors {
		vendor.SetupFillEnvFileContent(&envFileContent)
	}

	o.YouTube.SetupFillEnvFileContent(&envFileContent)

	err = o.Db.SaveEnv(envFileContent.String())
	return
}

func (o *Fabric) Setup() (err error) {
	if err = o.SetupVendors(); err != nil {
		return
	}

	if err = o.SetupDefaultModel(); err != nil {
		return
	}

	_ = o.YouTube.SetupOrSkip()

	if err = o.PatternsLoader.Setup(); err != nil {
		return
	}

	err = o.SaveEnvFile()

	return
}

func (o *Fabric) SetupDefaultModel() (err error) {
	vendorsModels := o.GetModels()

	vendorsModels.Print()

	if err = o.Ask(o.Label); err != nil {
		return
	}

	index, parseErr := strconv.Atoi(o.DefaultModel.Value)
	if parseErr == nil {
		o.DefaultVendor.Value, o.DefaultModel.Value = vendorsModels.GetVendorAndModelByModelIndex(index)
	} else {
		o.DefaultVendor.Value = vendorsModels.FindVendorsByModelFirst(o.DefaultModel.Value)
	}

	//verify
	vendorNames := vendorsModels.FindVendorsByModel(o.DefaultModel.Value)
	if len(vendorNames) == 0 {
		err = errors.Errorf("You need to chose an available default model.")
		return
	}

	fmt.Println()
	o.DefaultVendor.Print()
	o.DefaultModel.Print()

	err = o.SaveEnvFile()

	return
}

func (o *Fabric) SetupVendors() (err error) {
	o.Models = nil
	if o.Vendors, err = o.VendorsAll.Setup(); err != nil {
		return
	}

	if !o.HasVendors() {
		err = errors.New("No vendors configured")
		return
	}

	err = o.SaveEnvFile()

	return
}

// Configure buildClient VendorsController based on the environment variables
func (o *Fabric) configure() (err error) {
	for _, vendor := range o.VendorsAll.Vendors {
		if vendorErr := vendor.Configure(); vendorErr == nil {
			o.AddVendors(vendor)
		}
	}
	if err = o.PatternsLoader.Configure(); err != nil {
		return
	}

	//YouTube is not mandatory, so ignore not configured error
	_ = o.YouTube.Configure()

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
			ret.model = o.DefaultModel.Value
		}
	} else if model == "" {
		ret.vendor = o.FindByName(o.DefaultVendor.Value)
		ret.model = o.DefaultModel.Value
	} else {
		ret.vendor = o.FindByName(o.GetModels().FindVendorsByModelFirst(model))
		ret.model = model
	}

	if ret.vendor == nil {
		err = fmt.Errorf(
			"could not find vendor.\n Model = %s\n DefaultModel = %s\n DefaultVendor = %s",
			model, o.DefaultModel.Value, o.DefaultVendor.Value)
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

func (o *Chat) BuildChatSession() (ret *db.Session, err error) {
	// new messages will be appended to the session and used to send the message
	if o.Session != nil {
		ret = o.Session
	} else {
		ret = &db.Session{}
	}

	systemMessage := strings.TrimSpace(o.Context) + strings.TrimSpace(o.Pattern)

	if systemMessage != "" {
		ret.Append(&common.Message{Role: "system", Content: systemMessage})
	}

	userMessage := strings.TrimSpace(o.Message)
	if userMessage != "" {
		ret.Append(&common.Message{Role: "user", Content: userMessage})
	}

	if ret.IsEmpty() {
		ret = nil
		err = fmt.Errorf("no session, pattern or user messages provided")
	}
	return
}
