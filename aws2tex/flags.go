package main

import (
	"flag"
)

func init() {
	flag.BoolVar(&argConvertSvgWithInkscape, "ConvertSvgWithInkscape", true, "Convert svg using Inkscape")
	flag.StringVar(&argAssetZipFile, "AssetZipFile", "AWS-Architecture-Assets-For-Light-and-Dark-BG_20200911.478ff05b80f909792f7853b1a28de8e28eac67f4.zip", "Zip archive of AWS SVG files")

	flag.StringVar(&argIconsFile, "argIconsFile", "tex/icons.tex", "Location of icons file")
	flag.StringVar(&argStyleFile, "StyleFile", "sty/awsicons.sty", "Location of style file")
	flag.StringVar(&argInkscapeBinPath, "InkscapeBinPath", "/usr/bin/inkscape", "Location of inkscape binary")
	//flag.StringVar(&argInkscapeHelpArgs []string = []string{"--version"}
	flag.StringVar(&argPdfTexOutDir, "PdfTexOutDir", "icons_tex", "Location of LaTeX output dir")

	flag.BoolVar(&argShowAllArch, "ShowAllArch", false, "List all of the architecture icons (with repeats)")
	flag.BoolVar(&argShowAllRes, "ShowAllRes", false, "List all of the resource icons (with repeats)")
}
