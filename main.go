package main

import (
	"fmt"
	"os"
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

	os.RemoveAll("tmp")
	os.RemoveAll("addons_old")
	err = os.Rename("addons", "addons_old")
	if err != nil {
		fmt.Println(err)
	}
	err = os.Mkdir("addons", 0755)
	err = os.Mkdir("tmp", 0755)
	err = os.Mkdir("tmp/arcdps", 0755)
	err = os.Mkdir("tmp/d912pxy", 0755)
	err = os.Mkdir("tmp/gwradial", 0755)
	defer os.RemoveAll("tmp")

	go func() {
		defer wg.Done()
		utils.DownloadFile(d912pxyRelease.Assets[0].Name, "tmp/d912pxy/"+d912pxyRelease.Assets[0].Name, d912pxyRelease.Assets[0].BrowserDownloadURL, "d912pxy")
	}()
	go func() {
		defer wg.Done()
		utils.DownloadFile(gwradialRelease.Assets[0].Name, "tmp/gwradial/"+gwradialRelease.Assets[0].Name, gwradialRelease.Assets[0].BrowserDownloadURL, "gwradial")
	}()
	go func() {
		defer wg.Done()
		utils.DownloadFile("d3d9.dll", "tmp/arcdps/d3d9.dll", "https://www.deltaconnected.com/arcdps/x64/d3d9.dll", "arcdps")
	}()

	wg.Wait()

}
