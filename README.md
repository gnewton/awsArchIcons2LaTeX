# What is this?}
LaTeX style awsicons.sty which allows you to use the [AWS architectural icons](https://aws.amazon.com/architecture/icons/) in your LaTeX  documents.

The AWS SVG icons are converted to LaTeX-compatible PDF form using [Inkscape](https://inkscape.org/) using the purpose-written Go program [aws2tex] (https://github.com/gnewton/awsArchIcons2LaTeX/aws2tex). 

This Go program, aws2tex, creates:
1. `awsicons.sty` file
1. LaTeX-compatible PDF versions of icons (and written to the `icons_tex` directory, using Inkscape).
1. [awsAll.pdf](https://github.com/gnewton/awsArchIcons2LaTeX/raw/main/awsAllIcons.pdf), which lists all of the LaTeX macros. It is also to illustrate as well as index these icons and their corresponding LaTeX macros. 
  Not all icons are converted:
    *  For Resources, only the `_Light` versions. The `_Dark` are omitted.
    *  For Services, only the `_64`} versions. The lower resolution `\_16,\_32,\_48` are omitted.
The [awsAll.pdf](https://github.com/gnewton/awsArchIcons2LaTeX/raw/main/awsAllIcons.pdf) document also includes hyperlinks to Google searches of the titles of each icon.


# How to use it?}
There are two parts:
1. \texttt{awsicons.sty} file
1. The directory with all of the icons rendered into LaTeX compatible form, [icons_tex](https://github.com/gnewton/awsArchIcons2LaTeX/tree/main/icons_tex)

All you need to do is use the package to your LaTeX document using `\usepackage` (after having added the `awsicons.sty` file to where LaTeX can see it, or put it in the same directory as your document) **AND** add the directory containing all the generated icons to the `\graphicspath{...}`.
