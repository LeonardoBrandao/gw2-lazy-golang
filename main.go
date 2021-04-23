package main

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/LeonardoBrandao/gw2-utility/utils"
)

type GithubRelease struct {
	URL    string `json:"url"`
	Assets []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func main() {
	gwpath := path.Join("C:", "Guild Wars 2")

	var wg sync.WaitGroup
	var d912pxyRelease GithubRelease
	var gwradialRelease GithubRelease

	wg.Add(3)

	err := utils.GetJson("https://api.github.com/repos/megai2/d912pxy/releases/latest", &d912pxyRelease)
	if err != nil {
		fmt.Println(err)
	}
	err = utils.GetJson("https://api.github.com/repos/Friendly0Fire/GW2Radial/releases/latest", &gwradialRelease)
	if err != nil {
		fmt.Println(err)
	}

	os.RemoveAll(path.Join(gwpath, "addons_old"))
	err = os.Rename(path.Join(gwpath, "addons"), path.Join(gwpath, "addons_old"))
	if err != nil {
		fmt.Println(err)
	}

	arctmp, err := os.MkdirTemp("", "arc-*")
	defer os.RemoveAll(arctmp)
	d912tmp, err := os.MkdirTemp("", "d912-*")
	defer os.RemoveAll(d912tmp)
	gwradialtmp, err := os.MkdirTemp("", "gwradial-*")
	defer os.RemoveAll(gwradialtmp)

	go func() {
		defer wg.Done()
		utils.DownloadAddon(gwpath, arctmp, "d912pxy", d912pxyRelease.Assets[0].Name, d912pxyRelease.Assets[0].BrowserDownloadURL)
	}()
	go func() {
		defer wg.Done()
		utils.DownloadAddon(gwpath, d912tmp, "gwradial", gwradialRelease.Assets[0].Name, gwradialRelease.Assets[0].BrowserDownloadURL)
	}()
	go func() {
		defer wg.Done()
		utils.DownloadAddon(gwpath, gwradialtmp, "arcdps", "d3d9.dll", "https://www.deltaconnected.com/arcdps/x64/d3d9.dll")
	}()

	wg.Wait()

}
