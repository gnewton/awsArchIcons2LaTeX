
all:
	cd aws2tex; go build
	rm icons_tex/*	
	aws2tex/aws2tex >tex/icons.tex 
	cd tex; make; mv awsIcons.pdf ..

clean:
	cd tex; make clean
	rm *~
