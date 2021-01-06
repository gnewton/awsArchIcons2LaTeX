
all:
	pdflatex simple
	pdflatex Data-pipeline-Grov-Technologies
	gs -dBATCH -dNOPAUSE -dQUIET -dSAFER -sDEVICE=png16m -dTextAlphaBits=4 -dGraphicsAlphaBits=4 -dDOINTERPOLATE -dCOLORSCREEN -r600 -o Data-pipeline-Grov-Technologies.png Data-pipeline-Grov-Technologies.pdf	

clean:
	rm *.aux *.idx *.ilg *.ind *.out *.toc *~ *.log
