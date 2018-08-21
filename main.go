package main

import (
	"path/filepath"
	"os"
	"fmt"
	"os/exec"
	"archive/zip"
	"io/ioutil"
	"time"
	"strconv"
)

/**

 * Created by Connor Wright on 21/08/2018 at 15:58
 * zenith_backups
 * https://github.com/cxnky/
 * Protected by the GNU GPL v3 License (https://www.gnu.org/licenses/gpl-3.0.en.html)

**/

func main() {

	fmt.Println("Closing FiveM servers...")
	fmt.Println("")
	killFiveMServers()
	fmt.Println("Closed FiveM servers...")

	fmt.Println("Deleting cache folders from all servers...")
	deleteCacheFolders()
	fmt.Println("Cache folders deleted")

	fmt.Println("Performing speed test...")
	fmt.Println("")
	performSpeedTest()

	fmt.Println("Fetching list of files...")
	fiveMFileList := getFilesInFolder("C:\\Users\\Administrator\\Desktop\\FiveM")
	fmt.Println("Fetched " + string(len(fiveMFileList)) + " files to be backed up")

	fmt.Println("Zipping all files and folders within FiveM directory")
	zipWriter()
	fmt.Println("Successfully zipped.")

	fmt.Println("Moving backup to appropriate folder...")
	moveBackup(constructProperFilePath())
	fmt.Println("Backup complete.")

	fmt.Println("Starting FiveM servers back up...")
	startFiveMServers()

	fmt.Println("FiveM servers started back up. Exiting.")

	time.Sleep(10 * time.Second)

}

func startFiveMServers() {

	startServer("1")
	startServer("2")
	startServer("3")

}

func startServer(number string) {

	cmd := exec.Command("C:\\Users\\Administrator\\Desktop\\FiveM\\server" + number + "\\" + number + "start.bat")
	err := cmd.Start()

	if err != nil {

		fmt.Println(err)

	}

}

func deleteCacheFolders() {

	os.RemoveAll("C:\\Users\\Administrator\\Desktop\\FiveM\\server1\\cache\\files")
	os.RemoveAll("C:\\Users\\Administrator\\Desktop\\FiveM\\server2\\cache\\files")
	os.RemoveAll("C:\\Users\\Administrator\\Desktop\\FiveM\\server3\\cache\\files")

}

func moveBackup(newPath string) {

	err := os.Rename("C:\\Users\\Administrator\\Documents\\zip.zip", newPath)

	if err != nil {

		fmt.Println(err)

	}

}

func formatDate(day time.Time) string {

	suffix := "th"

	switch day.Day() {

	case 1, 21, 31:
		suffix = "st"

	case 2, 22:
		suffix = "nd"

	case 3, 23:
		suffix = "rd"

	}

	return suffix

}

func constructProperFilePath() string {

	driveLetter := "X:\\"
	timeObject := time.Now()

	date := strconv.Itoa(timeObject.Day())
	month := timeObject.Month().String()

	return driveLetter + month + "\\" + date  + formatDate(timeObject) + " of " + month + " Backup.zip"


}

func zipWriter() {

	baseFolder := "C:\\Users\\Administrator\\Desktop\\FiveM\\"

	outFile, err := os.Create("C:\\Users\\Administrator\\Documents\\zip.zip")

	if err != nil {

		fmt.Println(err)

	}

	defer outFile.Close()

	w := zip.NewWriter(outFile)

	addFiles(w, baseFolder, "")

	if err != nil {

		fmt.Println(err)

	}

	err = w.Close()

	if err != nil {

		fmt.Println(err)

	}

}

func addFiles(w *zip.Writer, basePath, baseInZip string) {

	files, err := ioutil.ReadDir(basePath)

	if err != nil {

		fmt.Println(err)

	}

	for _, file := range files {

		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, file.Name() + "/")
		}

	}

}

func performSpeedTest() {

	cmd := exec.Command("cmd.exe", "/C", "speedtest --server=14911 > C:\\Users\\Administrator\\Desktop\\FiveM\\speedtest.txt")
	out, err := cmd.CombinedOutput()

	if err != nil {

		fmt.Println("cmd.Run() failed with ", err)
		return

	}

	fmt.Println("Successfully wrote speedtest results to file")
	fmt.Println(out)

}

func killFiveMServers() {

	cmd := exec.Command("taskkill", "/F", "/IM", "cmd.exe", "/T")
	out, err := cmd.CombinedOutput()

	if err != nil {

		fmt.Errorf("cmd.Run() failed with %s\n", err)

	}

	fmt.Printf("Result from command: \n%s\n", string(out))

}

func getFilesInFolder(path string) []string {

	fileList := []string{}

	filepath.Walk(path, func (path string, f os.FileInfo, err error) error {

		fileList = append(fileList, path)
		return nil

	})

	return fileList

}