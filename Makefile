
all:
	cd aws2tex; go build
	-rm icons_tex/*	
	aws2tex/aws2tex
	cd tex; make; mv awsIcons.pdf ..
	cd examples; make

clean:
	cd tex; make clean
	rm *~
