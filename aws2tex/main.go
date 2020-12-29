package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Entry struct {
	sortName  string
	filename  string
	value     string
	fullPath  string
	modified  time.Time
	chosen    bool
	macroName string
}

const MACOS_PREFIX = "__MACOSX"

const StartArch = "\\archStart"
const EndArch = "\\archEnd"
const StartRes = "\\resStart"
const EndRes = "\\resEnd"

var argAssetZipFile string = "AWS-Architecture-Assets-For-Light-and-Dark-BG_20200911.478ff05b80f909792f7853b1a28de8e28eac67f4.zip"
var argConvertSvgWithInkscape = true
var argIconsFile string = "tex/icons.tex"
var argInkscapeBinPath string = "/usr/bin/inkscape"
var argMacrosFile string = "tex/macros.tex"
var argPdfTexOutDir string = "icons_tex"
var argShowAllArch bool = false
var argShowAllRes bool = false

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create output dir, if needed
	err := os.MkdirAll(argPdfTexOutDir, 0700)
	if err != nil {
		log.Fatal(err)
	}

	arch_entries, res_entries, err := Unzip(argAssetZipFile, argPdfTexOutDir)
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(ByName(res_entries))
	sort.Sort(ByName(arch_entries))

	macrosFile, err := os.Create(argMacrosFile)
	if err != nil {
		log.Panic(err)
	}
	wmacros := bufio.NewWriter(macrosFile)
	defer func() {
		if err = wmacros.Flush(); err != nil {
			log.Panic(err)
		}
		if err := macrosFile.Close(); err != nil {
			log.Panic(err)
		}
	}()

	printMacros(wmacros, res_entries)
	printMacros(wmacros, arch_entries)

	iconsFile, err := os.Create(argIconsFile)
	if err != nil {
		log.Panic(err)
	}

	wicons := bufio.NewWriter(iconsFile)
	defer func() {
		if err = wicons.Flush(); err != nil {
			log.Panic(err)
		}
		if err := iconsFile.Close(); err != nil {
			log.Panic(err)
		}
	}()

	fmt.Fprintln(wicons, StartRes)
	printEntries(wicons, res_entries)
	fmt.Fprintln(wicons, EndRes)
	fmt.Fprintln(wicons, StartArch)
	printEntries(wicons, arch_entries)
	fmt.Fprintln(wicons, EndArch)
}

func Unzip(src string, dest string) ([]*Entry, []*Entry, error) {
	arch_entries := make([]*Entry, 0)
	res_entries := make([]*Entry, 0)
	entryMap := make(map[string]*Entry)
	macroMap := make(map[string]*Entry)

	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasPrefix(f.Name, MACOS_PREFIX) || !strings.HasSuffix(f.Name, ".svg") {
			continue
		}

		// Only take the Arch 64, not the 16,32,48
		if !argShowAllArch && strings.HasPrefix(filepath.Base(f.Name), "Arch") && !strings.HasSuffix(filepath.Base(f.Name), "_64.svg") {
			continue
		}

		//Only take the Res Light not the Dark
		if !argShowAllRes && strings.HasPrefix(filepath.Base(f.Name), "Res") && strings.HasSuffix(filepath.Base(f.Name), "_Dark.svg") {
			continue
		}

		// Store filename/path for returning and using later on

		fpath := strings.ReplaceAll(strings.ReplaceAll(filepath.Join(dest, filepath.Base(f.Name)), " ", "_"), "&", "And")

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return nil, nil, fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			continue
		}

		// Is a file:
		var rc io.ReadCloser

		// Open file (in zip stream)
		rc, err = f.Open()
		if err != nil {
			log.Fatal(err)
		}

		// Change svg to pdf file
		pdfFile := strings.TrimSuffix(fpath, "svg") + "pdf"
		justF := strings.TrimSuffix(filepath.Base(fpath), "svg") + "pdf"
		cleanName := cleanAll(justF)
		macroName := string(justF[0]) + makeMacroName(cleanName)

		entryString := fmt.Sprintf("\\gxs{%s %s} {\\includegraphics[width=\\iconsize\\textwidth]{%s}} {%s} {{\\textbackslash}%s}",
			makeSearchLink(cleanName), index(cleanName), justF, strings.ReplaceAll(justF, "_", "\\_"), macroName)

		entry := Entry{
			sortName:  cleanName,
			value:     entryString,
			filename:  justF,
			fullPath:  f.Name,
			modified:  f.FileHeader.Modified,
			chosen:    true,
			macroName: macroName,
		}
		if filenameCollision, ok := entryMap[entry.filename]; ok {
			log.Println("****************************************Collision!!!  " + entry.fullPath + " " + entry.value + "|||" + filenameCollision.fullPath + " " + filenameCollision.value)
			log.Println(entry.modified.Format("2006-01-02 3:4:5 pm") + " " + filenameCollision.modified.Format("2006-01-02 3:4:5 pm"))
			log.Println(entry.filename + "  |  " + filenameCollision.filename)
			// Decided to take the most recent file, as per zipfile timestamp
			//   Somewhat arbitrary. Have looked at the files and seen only pixel level differences
			if entry.modified.After(filenameCollision.modified) {
				entry.chosen = true
				filenameCollision.chosen = false
			} else {
				entry.chosen = false
				filenameCollision.chosen = true
			}
			//continue
			log.Println("----entry.chosen", entry.chosen, "  collision.chosen", filenameCollision.chosen)

		}

		// To deal with when macros names collide (same files):
		// Examples:
		//    Res_AWS-Amazon-Simple-Storage_S3-Replication-Time-Control_48_Light.pdf
		//    Res_Amazon-Simple-Storage_S3-Replication-Time-Control_48_Light.pdf
		//
		//    Res_Amazon-Simple-Storage_S3-Replication_48_Light.pdf
		//    Res_AWS-Amazon-Simple-Storage_S3-Replication_48_Light.pdf
		//  Same icons.
		if macroCollision, ok := macroMap[entry.macroName]; ok {
			log.Println("****************************************Macro Collision!!!  " + entry.fullPath + " " + entry.value + "|||" + macroCollision.fullPath + " " + macroCollision.value)
			if entry.modified.After(macroCollision.modified) {
				entry.chosen = true
				macroCollision.chosen = false
			} else {
				entry.chosen = false
				macroCollision.chosen = true
			}
			log.Println("---- macro entry.chosen", entry.chosen, "  collision.chosen", macroCollision.chosen)
		}
		if !entry.chosen {
			continue
		}

		entryMap[entry.filename] = &entry
		macroMap[entry.macroName] = &entry

		// Is Arch or Res entry?
		if strings.HasPrefix(justF, "Arch_") {
			arch_entries = append(arch_entries, &entry)
		} else {
			res_entries = append(res_entries, &entry)
		}

		if argConvertSvgWithInkscape {
			// Setup running inkscape to do conversion, writing out to pdfFile
			cmd := exec.Command(argInkscapeBinPath, "--file=-", "-D", "-z", "--export-latex", "--export-pdf="+pdfFile)
			stdin, err := cmd.StdinPipe()
			if err != nil {
				log.Fatal(err)
			}

			// Read the zip file and pipe in to inkscape
			go func() {
				defer func() {
					if err := stdin.Close(); err != nil {
						log.Fatal(err)
					}
				}()
				_, err = io.Copy(stdin, rc)
				if err != nil {
					log.Fatal(err)
				}
			}()

			// If there is a problem with running inkscape...
			stdoutStderr, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("%s\n", stdoutStderr)
				log.Fatal(err)
			}
		}
		// Close the file without defer to close before next iteration of loop
		err = rc.Close()

		if err != nil {
			log.Fatal(err)
		}

	}
	return arch_entries, res_entries, nil
}

var unwanted = []string{
	".pdf",
	"AWS-",
	"Amazon-",
	"Arch_",
	"Res_",
	// Do this due to "_Lightbulb"
	"_16_Dark",
	"_32_Dark",
	"_48_Dark",
	"_64_Dark",
	"_16_Light",
	"_32_Light",
	"_48_Light",
	"_64_Light",
	"_16",
	"_32",
	"_48",
	"_64",
}

var reEnd = regexp.MustCompile(`(_16|_32|_48|_64)(_Dark|_Light)*.pdf$`)
var reFront = regexp.MustCompile(`(^Arch_|^Res_)(Amazon-|AWS-)*`)

// Fixing inconsistent naming:
// Most use "Simple-Storage" (29) and others use "S3"; "S3" is more compact
var reS3 = regexp.MustCompile(`(^Amazon-Simple-Storage|^Simple-Storage)`)

func cleanAll(s string) string {

	// Fixing inconsistent naming:
	// Res_AWS-loT-Device-Defender_IoT-Device-Jobs_48_Light.svg
	// "-loT-" should be "-IoT"
	// N: 4 instances
	s = strings.ReplaceAll(s, "loT-", "IoT-")

	// Fixing inconsistent naming:
	// Res-Amazon-Simple-Storage_S3-Replication_48_Dark.svg
	// "Res-" should be "Res_"
	// N: 2 instances
	s = strings.ReplaceAll(s, "Res-", "Res_")

	// Fixing inconsistent naming:
	// Res_Amazon-Aurora_Amazon-RDS-Instance-Aternate_48_Light.svg
	// "Aternate" should be "Alternate"
	// N: 2 instances
	s = strings.ReplaceAll(s, "Aternate_", "Alternate_")

	// Fixing inconsistent naming:
	// Arch_AWS-Console-Mobile-Application__48.svg
	// Double "__" is incorrect
	// N: 3 instances
	s = strings.ReplaceAll(s, "__", "_")

	// Fixing inconsistent naming:
	// Res_Amazon-Aurora_Amazon-Aurora-Instance-alternate 48 Dark.pdf
	// versus
	// Res_Amazon-Aurora-SQL-Server-Instance_48_Light.pdf
	// "-Aurora-" should be "-Aurora_"
	// Latter uncorrect
	// N: Several instances
	if strings.HasPrefix(s, "Res_Amazon-Aurora-") {
		s = strings.Replace(s, "-Aurora-", "-Aurora_", 1)
	}

	// Fixing inconsistent naming:
	// Res_AWS-OpsWorks-Stack2_48_Dark.pdf
	// versus
	// Res_AWS-OpsWorks_Deployments_48_Dark.pdf
	// N: 4 instances
	if strings.HasPrefix(s, "Res_AWS-OpsWorks-") {
		s = strings.Replace(s, "OpsWorks-", "OpsWorks_", 1)
	}

	s = string(reEnd.ReplaceAll([]byte(s), []byte("")))
	s = string(reFront.ReplaceAll([]byte(s), []byte("")))
	s = string(reS3.ReplaceAll([]byte(s), []byte("S3")))

	// ":" used to separate major sections
	s = strings.ReplaceAll(s, "_", ": ")
	s = strings.ReplaceAll(s, "-", " ")

	return s
}

func makeSearchLink(s string) string {
	q := strings.ReplaceAll(s, "Simple Storage:", "S3:")
	q = strings.ReplaceAll(s, ":", "")
	q = strings.ReplaceAll(q, " ", "+")
	return "\\href{" + "https://www.google.com/search?q=AWS+" + q + "}{" + s + "}"
}

func index(s string) (r string) {
	s = strings.ReplaceAll(s, "Simple Storage:", "S3:")

	r = "\\index{" + strings.ReplaceAll(s, ": ", "!") + "}"
	sections := strings.Split(s, ": ")
	if len(sections) > 1 {
		r = r + "\\index{" + sections[len(sections)-1] + " (" + strings.Join(sections[0:len(sections)-1], " ") + ")}"
	}
	return r
}

type ByName []*Entry

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return strings.Compare(a[i].sortName, a[j].sortName) == -1 }

func printEntries(w io.Writer, entries []*Entry) {
	previous := ""
	previousChosen := false
	for _, entry := range entries {
		if entry.chosen {
			if entry.filename == previous {
				log.Println("---------------", previous, previousChosen, entry.chosen)
			} else {
				previous = entry.filename
				previousChosen = entry.chosen
			}

			fmt.Fprintln(w, entry.value)
			fmt.Fprintln(w, "")
		}
	}
}

func printMacros(w io.Writer, entries []*Entry) {
	fmt.Fprintln(w, "\n\n%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
	for _, entry := range entries {
		if entry.macroName == "RSThreeSThreeReplication" {
			log.Printf("*******************ZZZ  %+v", *entry)
		}
		if entry.chosen {
			fmt.Fprintf(w, "\\newcommand{\\%s}[1]{\\includegraphics[width=#1]{%s}}\n",
				entry.macroName, entry.filename)
		}
	}
}

func makeMacroName(s string) string {
	return replaceNumberWithStrings(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), ":", "_"), "_", ""))
}

// LaTeX does not allow numbers in macro names!!!!
// http://www.texfaq.org/FAQ-linmacnames
var numMap = map[string]string{"0": "Zero", "1": "One", "2": "Two", "3": "Three", "4": "Four", "5": "Five", "6": "Six", "7": "Seven", "8": "Eight", "9": "Nine"}

func replaceNumberWithStrings(s string) string {
	for k, v := range numMap {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}
