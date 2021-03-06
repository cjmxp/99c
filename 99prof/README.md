# Table of Contents

1. Usage
1. Installation
1. Changelog
1. Sample
1. Bogomips

# 99prof

Command 99prof profiles programs produced by the 99c compiler.

The profile is written to stderr.

### Usage

Profile a program by issuing

    99prof [-functions] [-lines] [-instructions] [-rate] a.out [arguments]

    -functions
      	profile functions
    -instructions
      	profile instructions
    -lines
      	profile lines
    -rate int
      	profile rate (default 1000)

### Installation

To install or update 99prof

     $ go get [-u] -tags virtual.profile github.com/cznic/99c/99prof

Online documentation: [godoc.org/github.com/cznic/99c/99prof](http://godoc.org/github.com/cznic/99c/99prof)

### Changelog

2017-10-09: Initial public release.

### Sample

    $ cd ../examples/prof/
    $ ls *
    bogomips.c fib.c
    $ cat fib.c 
    #include <stdlib.h>
    #include <stdio.h>
    
    int fib(int n) {
    	switch (n) {
    	case 0:
    		return 0;
    	case 1:
    		return 1;
    	default:
    		return fib(n-1)+fib(n-2);
    	}
    }
    
    int main(int argc, char **argv) {
    	if (argc != 2) {
    		return 2;
    	}
    
    	int n = atoi(argv[1]);
    	if (n<=0 || n>40) {
    		return 1;
    	}
    
    	printf("%i\n", fib(n));
    }
    $ 99c fib.c && 99prof -functions -lines -instructions a.out 31 2>log
    1346269
    $ cat log
    # [99prof -functions -lines -instructions a.out 31] 781.384628ms, 72.483 MIPS
    # functions
    fib   	     56636    100.00%    100.00%
    _start	         1      0.00%    100.00%
    # lines
    fib.c:11:	     32707     57.75%     57.75%
    fib.c:5:	      8738     15.43%     73.18%
    fib.c:4:	      4350      7.68%     80.86%
    fib.c:9:	      4002      7.07%     87.92%
    fib.c:7:	      2476      4.37%     92.29%
    fib.c:10:	      2184      3.86%     96.15%
    fib.c:8:	      1357      2.40%     98.55%
    fib.c:6:	       822      1.45%    100.00%
    /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:13:	         1      0.00%    100.00%
    # instructions
    Argument32      8796     15.53%	     15.53%
    Push32          5601      9.89%	     25.42%
    AddSP           4426      7.81%	     33.23%
    Return          4413      7.79%	     41.03%
    SubI32          4375      7.72%	     48.75%
    AP              4363      7.70%	     56.45%
    Arguments       4359      7.70%	     64.15%
    Func            4351      7.68%	     71.83%
    Call            4346      7.67%	     79.51%
    SwitchI32       4331      7.65%	     87.15%
    Store32         4283      7.56%	     94.72%
    AddI32          2187      3.86%	     98.58%
    Zero32           806      1.42%	    100.00%
    $ 

### Bogomips

Let's try to estimate the VM bogomips value on an older Intel® Xeon(R) CPU X5450 @ 3.00GHz × 4 machine.

    $ cd ../examples/prof/
    $ ls *
    bogomips.c  fib.c
    $ cat bogomips.c 
    #include <stdlib.h>
    #include <stdio.h>
    
    // src: https://en.wikipedia.org/wiki/BogoMips#Computation_of_BogoMIPS
    static void delay_loop(long loops) {
    	long d0 = loops;
    	do {
    		--d0;
    	} while (d0 >= 0);
    }
    
    int main(int argc, char **argv) {
    	if (argc != 2) {
    		return 2;
    	}
    
    	int n = atoi(argv[1]);
    	if (n<=0) {
    		return 1;
    	}
    
    	delay_loop(n);
    }
    $ 99c bogomips.c && 99prof -functions a.out 7000000
    # [99prof -functions a.out 7000000] 1.04654176s, 53.511 MIPS
    # functions
    delay_loop	     56000    100.00%    100.00%
    _start    	         1      0.00%    100.00%
    $ time ./a.out 18300000
    
    real	0m0,996s
    user	0m0,996s
    sys	0m0,000s
    $

In both cases the program executes for ~1 second. 18300000/7000000 = 2.614 and that's the slowdown coefficient when running the binary under 99prof. The bogomips value is thus ~140 MIPS on this machine.

    $ 99dump a.out 
    virtual.Binary a.out: code 0x0004d, text 0x00000, data 0x00030, bss 0x00020, pc2func 3, pc2line 23
    0x00000		call           0x2	; -
    0x00001		ffireturn      		; -
    
    # _start
    0x00002	func	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:13:1
    0x00003		arguments      			; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:14:1
    0x00004		push64         (ds)		; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:14:1
    0x00005		push64         (ds+0x10)	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:14:1
    0x00006		push64         (ds+0x20)	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:14:1
    0x00007		#register_stdfiles		; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:14:1
    0x00008		arguments      			; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x00009		sub            sp, 0x8		; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x0000a		arguments      			; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x0000b		push32         (ap-0x8)		; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x0000c		push64         (ap-0x10)	; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x0000d		call           0x16		; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    0x0000e		#exit          			; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:15:1
    
    0x0000f		builtin        		; /home/jnml/src/github.com/cznic/ccir/libc/crt0.c:16:1
    0x00010		#register_stdfiles	; __register_stdfiles:89:1
    0x00011		ffireturn      		; __register_stdfiles:89:1
    
    0x00012		add            sp, 0x8	; __builtin_exit:86:1
    0x00013		#exit          		; __builtin_exit:86:1
    
    0x00014		call           0x16	; __builtin_exit:86:1
    0x00015		ffireturn      		; __builtin_exit:86:1
    
    # main
    0x00016	func	variables      [0x8]byte	; bogomips.c:12:1
    0x00017		push           ap		; bogomips.c:12:1
    0x00018		zero32         			; bogomips.c:12:1
    0x00019		store32        			; bogomips.c:12:1
    0x0001a		add            sp, 0x8		; bogomips.c:12:1
    0x0001b		push32         (ap-0x8)		; bogomips.c:13:1
    0x0001c		push32         0x2		; bogomips.c:13:1
    0x0001d		neqi32         			; bogomips.c:13:1
    0x0001e		jz             0x23		; bogomips.c:13:1
    
    0x0001f		push           ap	; bogomips.c:14:1
    0x00020		push32         0x2	; bogomips.c:14:1
    0x00021		store32        		; bogomips.c:14:1
    0x00022		return         		; bogomips.c:14:1
    
    0x00023		push           bp-0x8		; bogomips.c:13:1
    0x00024		sub            sp, 0x8		; bogomips.c:17:1
    0x00025		arguments      			; bogomips.c:17:1
    0x00026		push64         (ap-0x10)	; bogomips.c:17:1
    0x00027		push32         0x1		; bogomips.c:17:1
    0x00028		indexi32       0x8		; bogomips.c:17:1
    0x00029		load64         0x0		; bogomips.c:17:1
    0x0002a		#atoi          			; bogomips.c:17:1
    0x0002b		store32        			; bogomips.c:17:1
    0x0002c		add            sp, 0x8		; bogomips.c:17:1
    0x0002d		push32         (bp-0x8)		; bogomips.c:18:1
    0x0002e		zero32         			; bogomips.c:18:1
    0x0002f		leqi32         			; bogomips.c:18:1
    0x00030		jz             0x35		; bogomips.c:18:1
    
    0x00031		push           ap	; bogomips.c:19:1
    0x00032		push32         0x1	; bogomips.c:19:1
    0x00033		store32        		; bogomips.c:19:1
    0x00034		return         		; bogomips.c:19:1
    
    0x00035		arguments      		; bogomips.c:18:1
    0x00036		push32         (bp-0x8)	; bogomips.c:22:1
    0x00037		convi32i64     		; bogomips.c:22:1
    0x00038		call           0x3f	; bogomips.c:22:1
    0x00039		return         		; bogomips.c:23:1
    
    0x0003a		builtin        	; atoi:69:1
    0x0003b		#atoi          	; atoi:69:1
    0x0003c		ffireturn      	; atoi:69:1
    
    0x0003d		call           0x3f	; atoi:69:1
    0x0003e		ffireturn      		; atoi:69:1
    
    # delay_loop
    0x0003f	func	variables      [0x8]byte		; bogomips.c:5:1
    0x00040		push           bp-0x8			; bogomips.c:6:1
    0x00041		push64         (ap-0x8)			; bogomips.c:6:1
    0x00042		store64        				; bogomips.c:6:1
    0x00043		add            sp, 0x8			; bogomips.c:6:1
    0x00044		push           bp-0x8			; bogomips.c:7:1
    0x00045		preinci64      0xffffffffffffffff	; bogomips.c:8:1
    0x00046		add            sp, 0x8			; bogomips.c:8:1
    0x00047		push64         (bp-0x8)			; bogomips.c:9:1
    0x00048		zero32         				; bogomips.c:9:1
    0x00049		convi32i64     				; bogomips.c:9:1
    0x0004a		geqi64         				; bogomips.c:9:1
    0x0004b		jnz            0x44			; bogomips.c:9:1
    
    0x0004c		return         	; bogomips.c:10:1
    
    Data segment
    00000000  30 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |0...............|
    00000010  38 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |8...............|
    00000020  40 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |@...............|
    
    DS relative bitvector
    00000000  01 00 01 00 01                                    |.....|
    
    Symbol table
    0x00012	function	__builtin_exit
    0x0000f	function	__register_stdfiles
    0x00000	function	_start
    0x0003a	function	atoi
    0x00014	function	main
    $

Alternatively, using 99dump, we can see that the loop consists of 8 instructions at addresses 0x00044-0x0004b. 18300000*8 = 146400000 confirming the estimated ~140MIPS value.
