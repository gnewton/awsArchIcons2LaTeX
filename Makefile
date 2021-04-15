
all:
	cd aws2tex; go build
	-rm icons_tex/*	
	aws2tex/aws2tex -AssetZipFile assets/AWS-Architecture-Assets-For-Light-and-Dark-BG_20200911.478ff05b80f909792f7853b1a28de8e28eac67f4.zip
	cd tex; make; mv awsAllIcons.pdf ..
	cd examples; make

clean:
	cd tex; make clean
	-rm icons_tex/*	
	rm *~
