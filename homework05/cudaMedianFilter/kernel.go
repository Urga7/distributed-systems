// Code generated by cudago. Edit at your own risk.
package cudaMedianFilter

import (
    "github.com/InternatBlackhole/cudago/cuda"
	"unsafe"
)


//here just to force usage of unsafe package
var __kernel_useless_var__ unsafe.Pointer = nil

const (
	KeyKernel = "kernel"
)


type medianfilterArgs struct {
    imgIn uintptr
    imgOut uintptr
    width int32
    height int32

}

/*var (
    medianfilterArgs = medianfilterArgs{}

)*/







func MedianFilter(grid, block cuda.Dim3, imgIn uintptr, imgOut uintptr, width int32, height int32) error {
	err := autoloadLib_kernel()
	if err != nil {
		return err
	}
	kern, err := getKernel("kernel", "medianFilter")
	if err != nil {
		return err
	}
	
	params := medianfilterArgs{
	    imgIn: imgIn,
	    imgOut: imgOut,
	    width: width,
	    height: height,
	
	}
	
	return kern.Launch(grid, block, unsafe.Pointer(&params.imgIn), unsafe.Pointer(&params.imgOut), unsafe.Pointer(&params.width), unsafe.Pointer(&params.height))
}

func MedianFilterEx(grid, block cuda.Dim3, sharedMem uint64, stream *cuda.Stream, imgIn uintptr, imgOut uintptr, width int32, height int32) error {
	err := autoloadLib_kernel()
	if err != nil {
		return err
	}
	kern, err := getKernel("kernel", "medianFilter")
	if err != nil {
		return err
	}
	
	params := medianfilterArgs{
	    imgIn: imgIn,
	    imgOut: imgOut,
	    width: width,
	    height: height,
	
	}
	
	return kern.LaunchEx(grid, block, sharedMem, stream, unsafe.Pointer(&params.imgIn), unsafe.Pointer(&params.imgOut), unsafe.Pointer(&params.width), unsafe.Pointer(&params.height))
}



var loaded_kernel = false


func autoloadLib_kernel() error {
	if loaded_kernel {
		return nil
	}
	err := InitLibrary([]byte(Kernel_ptxCode), "kernel")
	if err != nil {
		return err
	}
	loaded_kernel = true
	return nil
}

const Kernel_ptxCode = `//
// Generated by NVIDIA NVVM Compiler
//
// Compiler Build ID: CL-34431801
// Cuda compilation tools, release 12.6, V12.6.20
// Based on NVVM 7.0.1
//

.version 8.5
.target sm_52
.address_size 64

	// .globl	medianFilter

.visible .entry medianFilter(
	.param .u64 medianFilter_param_0,
	.param .u64 medianFilter_param_1,
	.param .u32 medianFilter_param_2,
	.param .u32 medianFilter_param_3
)
{
	.local .align 1 .b8 	__local_depot0[9];
	.reg .b64 	%SP;
	.reg .b64 	%SPL;
	.reg .pred 	%p<40>;
	.reg .b16 	%rs<83>;
	.reg .b32 	%r<42>;
	.reg .b64 	%rd<59>;


	mov.u64 	%SPL, __local_depot0;
	ld.param.u64 	%rd18, [medianFilter_param_0];
	ld.param.u64 	%rd19, [medianFilter_param_1];
	ld.param.u32 	%r3, [medianFilter_param_2];
	ld.param.u32 	%r4, [medianFilter_param_3];
	add.u64 	%rd1, %SPL, 0;
	add.u64 	%rd2, %SPL, 1;
	add.u64 	%rd3, %SPL, 2;
	add.u64 	%rd4, %SPL, 3;
	add.u64 	%rd5, %SPL, 4;
	add.u64 	%rd6, %SPL, 5;
	add.u64 	%rd7, %SPL, 6;
	add.u64 	%rd8, %SPL, 7;
	add.u64 	%rd58, %SPL, 8;
	mov.u32 	%r5, %ntid.x;
	mov.u32 	%r6, %ctaid.x;
	mov.u32 	%r7, %tid.x;
	mad.lo.s32 	%r1, %r6, %r5, %r7;
	mov.u32 	%r8, %ntid.y;
	mov.u32 	%r9, %ctaid.y;
	mov.u32 	%r10, %tid.y;
	mad.lo.s32 	%r2, %r9, %r8, %r10;
	setp.ge.s32 	%p1, %r1, %r3;
	setp.ge.s32 	%p2, %r2, %r4;
	or.pred  	%p3, %p1, %p2;
	@%p3 bra 	$L__BB0_46;

	cvta.to.global.u64 	%rd29, %rd18;
	add.s32 	%r11, %r3, -1;
	add.s32 	%r12, %r1, -1;
	max.s32 	%r13, %r12, 0;
	min.s32 	%r14, %r13, %r11;
	add.s32 	%r15, %r2, -1;
	max.s32 	%r16, %r15, 0;
	add.s32 	%r17, %r4, -1;
	min.s32 	%r18, %r16, %r17;
	mul.lo.s32 	%r19, %r18, %r3;
	add.s32 	%r20, %r14, %r19;
	cvt.s64.s32 	%rd30, %r20;
	add.s64 	%rd31, %rd29, %rd30;
	ld.global.nc.u8 	%rs1, [%rd31];
	st.local.u8 	[%rd1], %rs1;
	max.s32 	%r21, %r1, 0;
	min.s32 	%r22, %r21, %r11;
	add.s32 	%r23, %r22, %r19;
	cvt.s64.s32 	%rd32, %r23;
	add.s64 	%rd33, %rd29, %rd32;
	ld.global.nc.u8 	%rs2, [%rd33];
	st.local.u8 	[%rd2], %rs2;
	add.s32 	%r24, %r1, 1;
	max.s32 	%r25, %r24, 0;
	min.s32 	%r26, %r25, %r11;
	add.s32 	%r27, %r26, %r19;
	cvt.s64.s32 	%rd34, %r27;
	add.s64 	%rd35, %rd29, %rd34;
	ld.global.nc.u8 	%rs3, [%rd35];
	st.local.u8 	[%rd3], %rs3;
	max.s32 	%r28, %r2, 0;
	min.s32 	%r29, %r28, %r17;
	mul.lo.s32 	%r30, %r29, %r3;
	add.s32 	%r31, %r14, %r30;
	cvt.s64.s32 	%rd36, %r31;
	add.s64 	%rd37, %rd29, %rd36;
	ld.global.nc.u8 	%rs4, [%rd37];
	st.local.u8 	[%rd4], %rs4;
	add.s32 	%r32, %r22, %r30;
	cvt.s64.s32 	%rd38, %r32;
	add.s64 	%rd39, %rd29, %rd38;
	ld.global.nc.u8 	%rs5, [%rd39];
	st.local.u8 	[%rd5], %rs5;
	add.s32 	%r33, %r26, %r30;
	cvt.s64.s32 	%rd40, %r33;
	add.s64 	%rd41, %rd29, %rd40;
	ld.global.nc.u8 	%rs6, [%rd41];
	st.local.u8 	[%rd6], %rs6;
	add.s32 	%r34, %r2, 1;
	max.s32 	%r35, %r34, 0;
	min.s32 	%r36, %r35, %r17;
	mul.lo.s32 	%r37, %r36, %r3;
	add.s32 	%r38, %r14, %r37;
	cvt.s64.s32 	%rd42, %r38;
	add.s64 	%rd43, %rd29, %rd42;
	ld.global.nc.u8 	%rs7, [%rd43];
	st.local.u8 	[%rd7], %rs7;
	add.s32 	%r39, %r22, %r37;
	cvt.s64.s32 	%rd44, %r39;
	add.s64 	%rd45, %rd29, %rd44;
	ld.global.nc.u8 	%rs8, [%rd45];
	st.local.u8 	[%rd8], %rs8;
	add.s32 	%r40, %r26, %r37;
	cvt.s64.s32 	%rd46, %r40;
	add.s64 	%rd47, %rd29, %rd46;
	ld.global.nc.u8 	%rs9, [%rd47];
	st.local.u8 	[%rd58], %rs9;
	setp.le.u16 	%p4, %rs1, %rs2;
	mov.u64 	%rd51, %rd2;
	@%p4 bra 	$L__BB0_3;

	st.local.u8 	[%rd2], %rs1;
	mov.u64 	%rd51, %rd1;

$L__BB0_3:
	st.local.u8 	[%rd51], %rs2;
	ld.local.u8 	%rs10, [%rd2];
	setp.le.u16 	%p5, %rs10, %rs3;
	mov.u64 	%rd52, %rd3;
	@%p5 bra 	$L__BB0_6;

	st.local.u8 	[%rd3], %rs10;
	ld.local.u8 	%rs11, [%rd1];
	setp.le.u16 	%p6, %rs11, %rs3;
	mov.u64 	%rd52, %rd2;
	@%p6 bra 	$L__BB0_6;

	st.local.u8 	[%rd2], %rs11;
	mov.u64 	%rd52, %rd1;

$L__BB0_6:
	st.local.u8 	[%rd52], %rs3;
	ld.local.u8 	%rs12, [%rd3];
	setp.le.u16 	%p7, %rs12, %rs4;
	mov.u64 	%rd53, %rd4;
	@%p7 bra 	$L__BB0_10;

	st.local.u8 	[%rd4], %rs12;
	ld.local.u8 	%rs13, [%rd2];
	setp.le.u16 	%p8, %rs13, %rs4;
	mov.u64 	%rd53, %rd3;
	@%p8 bra 	$L__BB0_10;

	st.local.u8 	[%rd3], %rs13;
	ld.local.u8 	%rs14, [%rd1];
	setp.le.u16 	%p9, %rs14, %rs4;
	mov.u64 	%rd53, %rd2;
	@%p9 bra 	$L__BB0_10;

	st.local.u8 	[%rd2], %rs14;
	mov.u64 	%rd53, %rd1;

$L__BB0_10:
	st.local.u8 	[%rd53], %rs4;
	ld.local.u8 	%rs15, [%rd4];
	setp.le.u16 	%p10, %rs15, %rs5;
	mov.u64 	%rd54, %rd5;
	@%p10 bra 	$L__BB0_15;

	st.local.u8 	[%rd5], %rs15;
	ld.local.u8 	%rs16, [%rd3];
	setp.le.u16 	%p11, %rs16, %rs5;
	mov.u64 	%rd54, %rd4;
	@%p11 bra 	$L__BB0_15;

	st.local.u8 	[%rd4], %rs16;
	ld.local.u8 	%rs17, [%rd2];
	setp.le.u16 	%p12, %rs17, %rs5;
	mov.u64 	%rd54, %rd3;
	@%p12 bra 	$L__BB0_15;

	st.local.u8 	[%rd3], %rs17;
	ld.local.u8 	%rs18, [%rd1];
	setp.le.u16 	%p13, %rs18, %rs5;
	mov.u64 	%rd54, %rd2;
	@%p13 bra 	$L__BB0_15;

	st.local.u8 	[%rd2], %rs18;
	mov.u64 	%rd54, %rd1;

$L__BB0_15:
	st.local.u8 	[%rd54], %rs5;
	ld.local.u8 	%rs19, [%rd5];
	setp.le.u16 	%p14, %rs19, %rs6;
	mov.u64 	%rd55, %rd6;
	@%p14 bra 	$L__BB0_21;

	st.local.u8 	[%rd6], %rs19;
	ld.local.u8 	%rs20, [%rd4];
	setp.le.u16 	%p15, %rs20, %rs6;
	mov.u64 	%rd55, %rd5;
	@%p15 bra 	$L__BB0_21;

	st.local.u8 	[%rd5], %rs20;
	ld.local.u8 	%rs21, [%rd3];
	setp.le.u16 	%p16, %rs21, %rs6;
	mov.u64 	%rd55, %rd4;
	@%p16 bra 	$L__BB0_21;

	st.local.u8 	[%rd4], %rs21;
	ld.local.u8 	%rs22, [%rd2];
	setp.le.u16 	%p17, %rs22, %rs6;
	mov.u64 	%rd55, %rd3;
	@%p17 bra 	$L__BB0_21;

	st.local.u8 	[%rd3], %rs22;
	ld.local.u8 	%rs23, [%rd1];
	setp.le.u16 	%p18, %rs23, %rs6;
	mov.u64 	%rd55, %rd2;
	@%p18 bra 	$L__BB0_21;

	st.local.u8 	[%rd2], %rs23;
	mov.u64 	%rd55, %rd1;

$L__BB0_21:
	st.local.u8 	[%rd55], %rs6;
	ld.local.u8 	%rs24, [%rd6];
	setp.le.u16 	%p19, %rs24, %rs7;
	mov.u64 	%rd56, %rd7;
	@%p19 bra 	$L__BB0_28;

	st.local.u8 	[%rd7], %rs24;
	ld.local.u8 	%rs25, [%rd5];
	setp.le.u16 	%p20, %rs25, %rs7;
	mov.u64 	%rd56, %rd6;
	@%p20 bra 	$L__BB0_28;

	st.local.u8 	[%rd6], %rs25;
	ld.local.u8 	%rs26, [%rd4];
	setp.le.u16 	%p21, %rs26, %rs7;
	mov.u64 	%rd56, %rd5;
	@%p21 bra 	$L__BB0_28;

	st.local.u8 	[%rd5], %rs26;
	ld.local.u8 	%rs27, [%rd3];
	setp.le.u16 	%p22, %rs27, %rs7;
	mov.u64 	%rd56, %rd4;
	@%p22 bra 	$L__BB0_28;

	st.local.u8 	[%rd4], %rs27;
	ld.local.u8 	%rs28, [%rd2];
	setp.le.u16 	%p23, %rs28, %rs7;
	mov.u64 	%rd56, %rd3;
	@%p23 bra 	$L__BB0_28;

	st.local.u8 	[%rd3], %rs28;
	ld.local.u8 	%rs29, [%rd1];
	setp.le.u16 	%p24, %rs29, %rs7;
	mov.u64 	%rd56, %rd2;
	@%p24 bra 	$L__BB0_28;

	st.local.u8 	[%rd2], %rs29;
	mov.u64 	%rd56, %rd1;

$L__BB0_28:
	st.local.u8 	[%rd56], %rs7;
	ld.local.u8 	%rs30, [%rd7];
	setp.le.u16 	%p25, %rs30, %rs8;
	mov.u64 	%rd57, %rd8;
	@%p25 bra 	$L__BB0_36;

	st.local.u8 	[%rd8], %rs30;
	ld.local.u8 	%rs31, [%rd6];
	setp.le.u16 	%p26, %rs31, %rs8;
	mov.u64 	%rd57, %rd7;
	@%p26 bra 	$L__BB0_36;

	st.local.u8 	[%rd7], %rs31;
	ld.local.u8 	%rs32, [%rd5];
	setp.le.u16 	%p27, %rs32, %rs8;
	mov.u64 	%rd57, %rd6;
	@%p27 bra 	$L__BB0_36;

	st.local.u8 	[%rd6], %rs32;
	ld.local.u8 	%rs33, [%rd4];
	setp.le.u16 	%p28, %rs33, %rs8;
	mov.u64 	%rd57, %rd5;
	@%p28 bra 	$L__BB0_36;

	st.local.u8 	[%rd5], %rs33;
	ld.local.u8 	%rs34, [%rd3];
	setp.le.u16 	%p29, %rs34, %rs8;
	mov.u64 	%rd57, %rd4;
	@%p29 bra 	$L__BB0_36;

	st.local.u8 	[%rd4], %rs34;
	ld.local.u8 	%rs35, [%rd2];
	setp.le.u16 	%p30, %rs35, %rs8;
	mov.u64 	%rd57, %rd3;
	@%p30 bra 	$L__BB0_36;

	st.local.u8 	[%rd3], %rs35;
	ld.local.u8 	%rs36, [%rd1];
	setp.le.u16 	%p31, %rs36, %rs8;
	mov.u64 	%rd57, %rd2;
	@%p31 bra 	$L__BB0_36;

	st.local.u8 	[%rd2], %rs36;
	mov.u64 	%rd57, %rd1;

$L__BB0_36:
	st.local.u8 	[%rd57], %rs8;
	ld.local.u8 	%rs37, [%rd8];
	setp.le.u16 	%p32, %rs37, %rs9;
	@%p32 bra 	$L__BB0_45;

	st.local.u8 	[%rd58], %rs37;
	ld.local.u8 	%rs38, [%rd7];
	setp.le.u16 	%p33, %rs38, %rs9;
	mov.u64 	%rd58, %rd8;
	@%p33 bra 	$L__BB0_45;

	st.local.u8 	[%rd8], %rs38;
	ld.local.u8 	%rs39, [%rd6];
	setp.le.u16 	%p34, %rs39, %rs9;
	mov.u64 	%rd58, %rd7;
	@%p34 bra 	$L__BB0_45;

	st.local.u8 	[%rd7], %rs39;
	ld.local.u8 	%rs40, [%rd5];
	setp.le.u16 	%p35, %rs40, %rs9;
	mov.u64 	%rd58, %rd6;
	@%p35 bra 	$L__BB0_45;

	st.local.u8 	[%rd6], %rs40;
	ld.local.u8 	%rs41, [%rd4];
	setp.le.u16 	%p36, %rs41, %rs9;
	mov.u64 	%rd58, %rd5;
	@%p36 bra 	$L__BB0_45;

	st.local.u8 	[%rd5], %rs41;
	ld.local.u8 	%rs42, [%rd3];
	setp.le.u16 	%p37, %rs42, %rs9;
	mov.u64 	%rd58, %rd4;
	@%p37 bra 	$L__BB0_45;

	st.local.u8 	[%rd4], %rs42;
	ld.local.u8 	%rs43, [%rd2];
	setp.le.u16 	%p38, %rs43, %rs9;
	mov.u64 	%rd58, %rd3;
	@%p38 bra 	$L__BB0_45;

	st.local.u8 	[%rd3], %rs43;
	ld.local.u8 	%rs44, [%rd1];
	setp.le.u16 	%p39, %rs44, %rs9;
	mov.u64 	%rd58, %rd2;
	@%p39 bra 	$L__BB0_45;

	st.local.u8 	[%rd2], %rs44;
	mov.u64 	%rd58, %rd1;

$L__BB0_45:
	st.local.u8 	[%rd58], %rs9;
	ld.local.u8 	%rs82, [%rd5];
	mad.lo.s32 	%r41, %r2, %r3, %r1;
	cvt.s64.s32 	%rd48, %r41;
	cvta.to.global.u64 	%rd49, %rd19;
	add.s64 	%rd50, %rd49, %rd48;
	st.global.u8 	[%rd50], %rs82;

$L__BB0_46:
	ret;

}

`