package agent

import (
	"time"

	"github.com/go-resty/resty/v2"
	rmm "github.com/wh1te909/rmmagent/shared"
)

func (a *WindowsAgent) InstallChoco() {

	var result rmm.ChocoInstalled
	result.AgentID = a.AgentID
	result.Installed = false

	rClient := resty.New()
	rClient.SetTimeout(30 * time.Second)

	url := "/api/v3/choco/"
	r, err := rClient.R().Get("https://chocolatey.org/install.ps1")
	if err != nil {
		a.Logger.Debugln(err)
		a.rClient.R().SetBody(result).Post(url)
		return
	}
	if r.IsError() {
		a.rClient.R().SetBody(result).Post(url)
		return
	}

	_, _, exitcode, err := a.RunScript(string(r.Body()), "powershell", []string{}, 900)
	if err != nil {
		a.Logger.Debugln(err)
		a.rClient.R().SetBody(result).Post(url)
		return
	}

	if exitcode != 0 {
		a.rClient.R().SetBody(result).Post(url)
		return
	}

	result.Installed = true
	a.rClient.R().SetBody(result).Post(url)
}

func (a *WindowsAgent) InstallWithChoco(name, version string) (string, error) {
	out, err := CMD("choco.exe", []string{"install", name, "--version", version, "--yes"}, 900, false)
	if err != nil {
		a.Logger.Errorln(err)
		return err.Error(), err
	}
	if out[1] != "" {
		return out[1], nil
	}
	return out[0], nil
}
