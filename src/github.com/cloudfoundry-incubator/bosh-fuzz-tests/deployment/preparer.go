package deployment

import (
	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type preparer struct {
	directorInfo   bltaction.DirectorInfo
	cliRunner      bltclirunner.Runner
	fs             boshsys.FileSystem
	assetsProvider bltassets.Provider
}

func NewPreparer(
	directorInfo bltaction.DirectorInfo,
	cliRunner bltclirunner.Runner,
	fs boshsys.FileSystem,
	assetsProvider bltassets.Provider,
) *preparer {
	return &preparer{
		directorInfo:   directorInfo,
		cliRunner:      cliRunner,
		fs:             fs,
		assetsProvider: assetsProvider,
	}
}

func (p *preparer) Prepare() error {
	err := p.cliRunner.TargetAndLogin(p.directorInfo.URL)
	if err != nil {
		return err
	}

	stemcellPath, err := p.assetsProvider.FullPath("stemcell.tgz")
	if err != nil {
		return err
	}

	err = p.cliRunner.RunWithArgs("upload", "stemcell", stemcellPath)
	if err != nil {
		return err
	}

	releaseDir, err := p.fs.TempDir("release-test")
	if err != nil {
		return err
	}
	defer p.fs.RemoveAll(releaseDir)

	releaseSrcPath, err := p.assetsProvider.FullPath("release")
	if err != nil {
		return err
	}

	err = p.fs.CopyDir(releaseSrcPath, releaseDir)
	if err != nil {
		return err
	}

	err = p.cliRunner.RunInDirWithArgs(releaseDir, "create", "release", "--force")
	if err != nil {
		return err
	}

	err = p.cliRunner.RunInDirWithArgs(releaseDir, "upload", "release")
	if err != nil {
		return err
	}
	return nil
}
