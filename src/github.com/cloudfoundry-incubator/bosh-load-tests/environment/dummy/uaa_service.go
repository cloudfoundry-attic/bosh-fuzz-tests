package dummy

import (
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshhttp "github.com/cloudfoundry/bosh-utils/http"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

type UAAServiceOptions struct {
	AssetsPath            string
	TomcatPath            string
	UaaHttpPort           int
	UaaServerPort         int
	UaaAccessLogDirectory string
}

type UAAService struct {
	options        UAAServiceOptions
	cmdRunner      boshsys.CmdRunner
	process        boshsys.Process
	assetsProvider bltassets.Provider
	fs             boshsys.FileSystem
	logger         boshlog.Logger
}

func NewUAAService(
	options UAAServiceOptions,
	cmdRunner boshsys.CmdRunner,
	assetsProvider bltassets.Provider,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) *UAAService {
	return &UAAService{
		options:        options,
		cmdRunner:      cmdRunner,
		assetsProvider: assetsProvider,
		fs:             fs,
		logger:         logger,
	}
}

func (u *UAAService) Start() error {
	scriptTemplatePath, err := u.assetsProvider.FullPath("start_uaa.template")
	if err != nil {
		return err
	}

	uaaScriptPath, err := u.renderStartScript(scriptTemplatePath)
	if err != nil {
		return err
	}

	u.process, err = u.cmdRunner.RunComplexCommandAsync(bltcom.CreateCommand(uaaScriptPath))
	if err != nil {
		return bosherr.WrapError(err, "starting uaa")
	}

	return u.waitForServiceToStart()
}

func (u *UAAService) Stop() {
	u.process.TerminateNicely(5 * time.Second)
}

func (u *UAAService) renderStartScript(uaaScriptTemplatePath string) (string, error) {
	scriptTemplate := template.Must(template.ParseFiles(uaaScriptTemplatePath))
	buffer := bytes.NewBuffer([]byte{})

	err := scriptTemplate.Execute(buffer, u.options)
	if err != nil {
		return "", err
	}

	renderedScript, err := u.fs.TempFile("start_uaa")
	if err != nil {
		return "", err
	}

	err = u.fs.WriteFile(renderedScript.Name(), buffer.Bytes())
	if err != nil {
		return "", err
	}

	u.fs.Chmod(renderedScript.Name(), os.ModePerm)

	return renderedScript.Name(), nil
}

func (u *UAAService) waitForServiceToStart() error {
	uaaURL := fmt.Sprintf("http://localhost:%d/uaa/oauth/token?client_id=test&grant_type=client_credentials", u.options.UaaHttpPort)

	retryHTTPClient := boshhttp.NewRetryClient(
		http.DefaultClient,
		uint(30),
		1*time.Second,
		u.logger,
	)

	httpRequest, err := http.NewRequest(http.MethodGet, uaaURL, nil)
	if err != nil {
		return err
	}
	httpRequest.SetBasicAuth("test", "secret")

	response, err := retryHTTPClient.Do(httpRequest)
	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	u.logger.Debug("UAAService: Test Response Body:", string(responseData))

	if !strings.Contains(string(responseData), "access_token") {
		return bosherr.Error("Response from UAA does not contain the string 'access_token")
	}

	return nil
}
