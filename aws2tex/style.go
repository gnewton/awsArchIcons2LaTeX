package main

type texStyle struct {
	FixedZipFile    string
	InkscapeVersion string
}

const styleTemplate = `
%%
%%  Copyright (C) 2021 by Glen Newton
%%
%%  This file may be distributed and/or modified under the
%%  conditions of the LaTeX Project Public License, either
%%  version 1.3 of this license or (at your option) any later
%%  version. The latest version of this license is in:
%%
%%  http://www.latex-project.org/lppl.txt
%%
%%  and version 1.3 or later is part of all distributions of
%%  LaTeX version 2005/12/01 or later.
%%
%%
\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{awsicons}[2020/01/03 AWS Architectural Icons]
%%
%%
\newcommand{\assetZipFile}{%s}\n", fixedZipFile)
\newcommand{\assetZipFileSplit}{\\seqsplit{%s}}\n", fixedZipFile)
\newcommand{\inkscapeVersion}{%s}\n", strings.ReplaceAll(inkscapeVersion, "_", "\\_)
\definecolor{awsOrange}{RGB}{255 153 0}
%%%%%%%%%%%%%%%%
`
