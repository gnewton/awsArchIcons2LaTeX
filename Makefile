
all:
	cd aws2tex; go build
	-rm icons_tex/*	
	aws2tex/aws2tex
	cd tex; make; mv awsAllIcons.pdf ..
	cd examples; make

clean:
	cd tex; make clean
	rm *~
