package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	u "github.com/LeonardoBrandao/gw2-utility/utils"
	"github.com/otiai10/copy"
)

const d912pxyReleasesURL = "https://api.github.com/repos/megai2/d912pxy/releases"
const gwRadialReleasesURL = "https://api.github.com/repos/Friendly0Fire/GW2Radial/releases"
const arcdpsDownloadURL = "https://www.deltaconnected.com/arcdps/x64/d3d9.dll"

func main() {
	gwpath := path.Clean(os.Args[1])
	addons_arg := os.Args[2:]
	var addons_list []u.Addon
	var addons_list_names []string
	var errors_list []error

	for _, addon := range addons_arg {
		var ext string
		var GitRelease u.GithubRelease
		var download_url string

		switch addon {

		case "arcdps":
			ext = ".dll"
			download_url = arcdpsDownloadURL

		case "d912pxy":
			ext = ".zip"
			err := u.GetJson(d912pxyReleasesURL, &GitRelease)
			if err != nil {
				errors_list = append(errors_list, err)
				continue
			}
			download_url = GitRelease[0].Assets[0].BrowserDownloadURL

		case "gwradial":
			ext = ".zip"
			err := u.GetJson(gwRadialReleasesURL, &GitRelease)
			if err != nil {
				errors_list = append(errors_list, err)
				continue
			}
			download_url = GitRelease[0].Assets[0].BrowserDownloadURL

		default:
			errors_list = append(errors_list, errors.New("Addon name not recognized"))
			continue
		}

		tmpdir, err := os.MkdirTemp("", "gwlazy-"+addon+"-*")
		if err != nil {
			errors_list = append(errors_list, err)
			continue
		}

		var ad = u.Addon{
			Name:         addon,
			Tmpdir:       tmpdir,
			Extension:    ext,
			Download_url: download_url,
		}

		addons_list = append(addons_list, ad)
		addons_list_names = append(addons_list_names, addon)
	}

	var wg sync.WaitGroup
	wg.Add(len(addons_list))

	os.RemoveAll(path.Join(gwpath, "addons_old"))
	err := copy.Copy(path.Join(gwpath, "bin64"), path.Join(gwpath, "addons_old"))
	if err != nil {
		fmt.Println(err)
	}

	for _, addon := range addons_list {
		go func(addon u.Addon) {
			defer wg.Done()
			u.DownloadAddon(addon)
		}(addon)
	}

	wg.Wait()

	u.CopyFiles(addons_list, addons_list_names, gwpath)

	fmt.Println("Press the Enter Key to exit")
	fmt.Scanln()
}
