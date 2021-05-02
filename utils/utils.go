package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/otiai10/copy"
)

type GithubRelease []struct {
	URL    string `json:"url"`
	Assets []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type Addon struct {
	Name         string
	Tmpdir       string
	Extension    string
	Download_url string
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func GetJson(url string, target interface{}) error {
	myClient := &http.Client{Timeout: 10 * time.Second}
	response, err := myClient.Get(url)
	if err != nil {
		return err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(responseData, &target)
	return err
}

func DownloadAddon(addon Addon) {
	var fn strings.Builder
	fn.WriteString(addon.Name)
	fn.WriteString(addon.Extension)
	filename := fn.String()
	filepath := path.Join(addon.Tmpdir, filename)

	// Get the data
	resp, err := http.Get(addon.Download_url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Downloaded: " + filename)

	iszip, err := regexp.MatchString(".zip", filename)
	if iszip {
		fmt.Println("unzipping " + addon.Name)
		Unzip(path.Join(addon.Tmpdir, filename), addon.Tmpdir)
	}
}

func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

// Find takes a slice and looks for an element in it.
func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func CopyFiles(addons_list []Addon, addons_list_names []string, gwpath string) bool {
	if _, err := os.Stat(path.Join(gwpath, "bin64")); os.IsNotExist(err) {
		os.MkdirAll(path.Join(gwpath, "bin64"), 777)
	}

	os.RemoveAll(path.Join(gwpath, "bin64", "d3d9.dll"))
	os.RemoveAll(path.Join(gwpath, "bin64", "d3d9_chainload.dll"))
	os.RemoveAll(path.Join(gwpath, "bin64", "d912pxy.dll"))

	for _, addon := range addons_list {
		var filename string
		switch addon.Name {
		case "d912pxy":
			if len(addons_list) == 3 { // arcdps, d912pxy, gwradial
				filename = "d912pxy.dll"
				break
			}
			if len(addons_list) == 2 && Find(addons_list_names, "arcdps") { // arcdps, d912pxy
				filename = "d3d9_chainload.dll"
				break
			}
			if len(addons_list) == 2 { // d912pxy, gwradial
				filename = "d912pxy.dll"
				break
			}
			filename = "d3d9.dll" // d912pxy
			break
		case "gwradial":
			if len(addons_list) == 3 {
				filename = "d3d9_chainload.dll" // arcdps, d912pxy, gwradial
				break
			}
			if len(addons_list) == 2 && Find(addons_list_names, "arcdps") { // arcdps, gwradial
				filename = "d3d9_chainload.dll"
				break
			}
			filename = "d3d9.dll" // gwradial
			break
		default:
			filename = "d3d9.dll" // arcdps will always be d3d9.dll
			break
		}
		copyAddon(addon, filename, gwpath)
	}
	return true
}

func copyAddon(addon Addon, filename string, gwpath string) {
	var err error
	switch addon.Name {
	case "arcdps":
		err = copy.Copy(path.Join(addon.Tmpdir, "arcdps.dll"), path.Join(gwpath, "bin64", filename))
		break
	case "d912pxy":
		err = copy.Copy(path.Join(addon.Tmpdir, addon.Name), path.Join(gwpath, addon.Name))
		err = copy.Copy(path.Join(gwpath, addon.Name, "dll", "release", "d3d9.dll"), path.Join(gwpath, "bin64", filename))
		break
	case "gwradial":
		err = copy.Copy(path.Join(addon.Tmpdir, "gw2addon_gw2radial.dll"), path.Join(gwpath, "bin64", filename))
		break
	}
	if err != nil {
		println(err)
		println("Error while copying: ", addon.Name)
	}
	println("Copied: ", addon.Name)
}
