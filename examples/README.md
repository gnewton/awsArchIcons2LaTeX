# Examples

## Simple
TODO

## AWS Architecture Blog Diagram
The Dec 22 2020 AWS Architecture blog has a great article on agriculture (a topic close to my heart) called [Building a Controlled Environment Agriculture Platform](https://aws.amazon.com/blogs/architecture/building-a-controlled-environment-agriculture-platform/ "Building a Controlled Environment Agriculture Platform"), by Ashu Joshi.
One of the architecture diagrams in the blog shows the cloud data pipeline and dashboard.
Here is the **ORIGINAL** diagram: 

![](https://d2908q01vomqb2.cloudfront.net/fc074d501302eb2b93e2554793fcaf50b3bf7291/2020/12/21/Data-pipeline-Grov-Technologies-1024x374.png "Building a Controlled Environment Agriculture Platform Diagram")(C) AWS

### Diagram re-created using aswicons.sty
I have re-created this architecture diagram using awsicons.sty, and [LaTeX tikz](https://ctan.org/pkg/pgf?lang=en) trying to reproduce as best I can the original, with a few exceptions.

<img src="https://github.com/gnewton/awsArchIcons2LaTeX/raw/main/examples/Data-pipeline-Grov-Technologies.png" alt="Diagram re-created using aswicons.sty" style="width:2000px;"/>

Some differences:
1. The "Amazon API Gateway" has a reddish colour in the original; in the icons from the assets package used here ([AWS-Architecture-Assets-For-Light-and-Dark-BG_20200911.478ff05b80f909792f7853b1a28de8e28eac67f4.zip](https://d1.awsstatic.com/webteam/architecture-icons/Q32020/AWS-Architecture-Assets-For-Light-and-Dark-BG_20200911.478ff05b80f909792f7853b1a28de8e28eac67f4.zip)), this icon is purplish.
1. I used an orange line for the block indicating AWS cloud instead of the black one in the original.
1. I used an orange square with aws written in it, as I could not find an (unencumbered) icon with "aws" and the aws swoosh.
1. I straightened out one of the lines (Grov Edge Greengrass to the S3 Data Bucket).
1. The arrows/lines connecting icons do not attach to each icon, instead leave a small space. This looked good and was the default tikz behaviour, so I did not change this.
1. There was no React icon that I could find in the AWS assets package, so I found something close on [stackexchange](https://tex.stackexchange.com/questions/349960/how-to-draw-a-symbol-on-a-tikz-node-that-can-be-reused) and modified it to be pretty close.



Things to improve/change:
1. I need to use styles instead of all of the, for example, `\draw[->,gray]` everywhere.
2. Maybe I should add the React icon as a macro in the LaTeX style?


