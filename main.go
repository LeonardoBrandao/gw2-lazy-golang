package main

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

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
	for _, addon := range addons_arg {
		switch addon {

		case "arcdps":
			arctmp, err := os.MkdirTemp("", "gwlazy-arc-*")
			if err != nil {
				fmt.Println(err)
			}
			defer os.RemoveAll(arctmp)

			var ad = u.Addon{
				Name:         "arcdps",
				Tmpdir:       arctmp,
				Extension:    ".dll",
				Download_url: arcdpsDownloadURL,
			}
			addons_list = append(addons_list, ad)
			break

		case "d912pxy":
			var d912pxyRelease u.GithubRelease
			d912tmp, err := os.MkdirTemp("", "gwlazy-d912pxy-*")
			if err != nil {
				fmt.Println(err)
			}
			defer os.RemoveAll(d912tmp)
			err = u.GetJson(d912pxyReleasesURL, &d912pxyRelease)
			if err != nil {
				fmt.Println(err)
			}

			var ad = u.Addon{
				Name:         "d912pxy",
				Tmpdir:       d912tmp,
				Extension:    ".zip",
				Download_url: d912pxyRelease[0].Assets[0].BrowserDownloadURL,
			}
			addons_list = append(addons_list, ad)

		case "gwradial":
			var gwradialRelease u.GithubRelease

			gwradialtmp, err := os.MkdirTemp("", "gwlazy-gwradial-*")
			if err != nil {
				fmt.Println(err)
			}
			defer os.RemoveAll(gwradialtmp)
			err = u.GetJson(gwRadialReleasesURL, &gwradialRelease)
			if err != nil {
				fmt.Println(err)
			}

			var ad = u.Addon{
				Name:         "gwradial",
				Tmpdir:       gwradialtmp,
				Extension:    ".zip",
				Download_url: gwradialRelease[0].Assets[0].BrowserDownloadURL,
			}
			addons_list = append(addons_list, ad)
		}
	}

	var wg sync.WaitGroup

	wg.Add(len(addons_list))

	os.RemoveAll(path.Join(gwpath, "addons_old"))
	err := copy.Copy(path.Join(gwpath, "addons"), path.Join(gwpath, "addons_old"))
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

	go counter()
	fmt.Println("Press the Enter Key to stop anytime")
	fmt.Scanln()
}

func counter() {
	i := 0
	for {
		time.Sleep(time.Second * 1)
		i++
	}
}
