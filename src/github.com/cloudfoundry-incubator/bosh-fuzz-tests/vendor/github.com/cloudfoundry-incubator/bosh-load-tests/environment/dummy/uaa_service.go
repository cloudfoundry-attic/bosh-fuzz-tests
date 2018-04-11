package dummy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"os"

	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshhttp "github.com/cloudfoundry/bosh-utils/httpclient"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
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

	if err := scriptTemplate.Execute(buffer, u.options); err != nil {
		return "", err
	}

	renderedScript, err := u.fs.TempFile("start_uaa")
	if err != nil {
		return "", err
	}

	if err = u.fs.WriteFile(renderedScript.Name(), buffer.Bytes()); err != nil {
		return "", err
	}

	// Why? Writing rendered script to file system seems to leave the file open
	// Ref: https://github.com/moby/moby/issues/9547
	// Temporarily proceeding with using copied file hack
	tempCopy := renderedScript.Name() + "-copy"
	if err := u.fs.CopyFile(renderedScript.Name(), tempCopy); err != nil {
		return "", err
	}

	if err := u.fs.Chmod(tempCopy, os.ModePerm); err != nil {
		return "", err
	}

	return tempCopy, nil
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
