// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "go_tls.h"
#include "funcdata.h"
#include "textflag.h"

TEXT runtime·rt0_go(SB), NOSPLIT, $0
              // save m->g0 = g0
                I32Const runtime·m0(SB)
                I64Const runtime·g0(SB)
              I64Store $m_g0

              // save m0 to g0->m
                I32Const runtime·g0(SB)
                I64Const runtime·m0(SB)
              I64Store $g_m

              // set g to g0
                I64Const runtime·g0(SB)
              Set g

              CALL runtime·check(SB)
            End
            CALL runtime·args(SB)
          End
          CALL runtime·osinit(SB)
        End
        CALL runtime·schedinit(SB)
      End
        Get SP
        I64Const $0
      I64Store $0
        Get SP
        I64Const runtime·mainPC(SB)
      I64Store $8
      CALL runtime·newproc(SB)
    End
    CALL runtime·mstart(SB)
  End
  Unreachable

DATA  runtime·mainPC+0(SB)/8,$runtime·main(SB)
GLOBL runtime·mainPC(SB),RODATA,$8

// func checkASM() bool
TEXT ·checkASM(SB), NOSPLIT, $0-1
    Get SP
    I32Const $1
  I32Store8 $8
  RET

TEXT runtime·gogo(SB), NOSPLIT, $0-0
    I64Load buf+0(FP)
  Set I0
      I64Load gobuf_sp(I0)
    I32WrapI64
  Set SP
        I64Load gobuf_pc(I0)
      I32WrapI64
      I32Const $16
    I32ShrU
  Set PC_F
        I64Load gobuf_pc(I0)
        I64Const $0xFFFF
      I64And
    I32WrapI64
  Set PC_B
    I64Load gobuf_g(I0)
  Set g
    I64Load gobuf_ctxt(I0)
  Set CTX
    I64Load gobuf_ret(I0)
  // clear to help garbage collector
      Get I0
    I32WrapI64
    I64Const $0
  I64Store $gobuf_sp
      Get I0
    I32WrapI64
    I64Const $0
  I64Store $gobuf_ret
      Get I0
    I32WrapI64
    I64Const $0
  I64Store $gobuf_ctxt
  Set RET0
  I32Const $1

// func mcall(fn func(*g))
// Switch to m->g0's stack, call fn(g).
// Fn must never return. It should gogo(&g->sched)
// to keep running g.
TEXT runtime·mcall(SB), NOSPLIT, $0-8
    // CTX = fn
      I64Load fn+0(FP)
    Set CTX

    // I1 = g.m
      I64Load g_m(g)
    Set I1

    // I2 = g0
      I64Load m_g0(I1)
    Set I2

    // save state in g->sched
        Get g
      I32WrapI64
        Get SP
      I64Load $0 // caller's PC
    I64Store $g_sched+gobuf_pc

        Get g
      I32WrapI64
          Get SP
        I64ExtendUI32
        I64Const $fn+0(FP)
      I64Add
    I64Store $g_sched+gobuf_sp

        Get g
      I32WrapI64
      Get g
    I64Store $g_sched+gobuf_g

    //I64Store (g_sched+gobuf_bp)(g)

    // if g == g0 call badmcall
        Get g
        Get I2
      I64Eq
    If
      JMP runtime·badmcall(SB)
    End

    // switch to g0's stack
          I64Load (g_sched+gobuf_sp)(I2)
          I64Const $8
        I64Sub
      I32WrapI64
    Set SP

    // set arg to current g
      Get SP
      Get g
    I64Store $0

    // switch to g0
      Get I2
    Set g

    // call fn
          Get CTX
        I32WrapI64
      I64Load $0
    CALL $0
  End

      Get SP
      I32Const $8
    I32Add
  Set SP

  JMP runtime·badmcall2(SB)

      Set I0

      // I1 = g.m
        I64Load g_m(g)
      Set I1

      // I2 = g0
        I64Load m_g0(I1)
      Set I2

      // if g == g0
          Get g
          Get I2
        I64Eq
      If
        // no switch:
          Get I0
        Set CTX

              Get CTX
            I32WrapI64
          I64Load $0
        JMP $0
      End

      // if g != m.curg
          Get g
          I64Load m_curg(I1)
        I64Ne
      If $0
        CALL runtime·badsystemstack(SB)
      End

      // switch:

      // save state in g->sched. Pretend to
      // be systemstack_switch if the G stack is scanned.
          Get g
        I32WrapI64
        I64Const $runtime·systemstack_switch(SB)
      I64Store $g_sched+gobuf_pc

          Get g
        I32WrapI64
          Get SP
        I64ExtendUI32
      I64Store $g_sched+gobuf_sp

          Get g
        I32WrapI64
        Get g
      I64Store $g_sched+gobuf_g

      //I64Store (g_sched+gobuf_bp)(g)

      // switch to g0
        Get I2
      Set g

          I64Load (g_sched+gobuf_sp)(I2)
          // make it look like mstart called systemstack on g0, to stop traceback
          I64Const $8
        I64Sub
      Set I3

          Get I3
        I32WrapI64
        I64Const $runtime·mstart(SB)
      I64Store $0

          Get I3
        I32WrapI64
      Set SP

      // call fn
        Get I0
      Set CTX

            Get CTX
          I32WrapI64
        I64Load $0
      CALL $0
    End

    // switch back to g
      I64Load g_m(g)
    Set I1

      I64Load m_curg(I1)
    Set I2

      Get I2
    Set g

        I64Load (g_sched+gobuf_sp)(I2)
      I32WrapI64
    Set SP

        Get I2
      I32WrapI64
      I64Const $0
    I64Store $g_sched+gobuf_sp
  End

TEXT runtime·systemstack_switch(SB), NOSPLIT, $0-0
  RET

TEXT runtime·return0(SB), NOSPLIT, $0-0
    I64Const $0
  Set RET0
  RET

TEXT runtime·jmpdefer(SB), NOSPLIT, $0-16
    I64Load fn+0(FP)
  Set CTX

      Get CTX
    I64Eqz
  If
    CALL runtime·sigpanic(SB)
  End

  // caller sp after CALL
        I64Load argp+8(FP)
        I64Const $8
      I64Sub
    I32WrapI64
  Set SP

  // decrease PC_B by 1 to CALL again
    Get SP
      I32Load16U (SP)
      I32Const $1
    I32Sub
  I32Store16 $0

  // but first run the deferred function
        Get CTX
      I32WrapI64
    I64Load $0
  JMP $0

TEXT runtime·asminit(SB), NOSPLIT, $0-0
  // No per-thread init.
  RET  

TEXT ·publicationBarrier(SB), NOSPLIT, $0-0
  RET

TEXT runtime·procyield(SB), NOSPLIT, $0-0 // FIXME
  RET

TEXT runtime·breakpoint(SB), NOSPLIT, $0-0
  Unreachable

// Called during function prolog when more stack is needed.
//
// The traceback routines see morestack on a g0 as being
// the top of a stack (for example, morestack calling newstack
// calling the scheduler calling newm calling gc), so we must
// record an argument size. For that purpose, it has no arguments.
TEXT runtime·morestack(SB), NOSPLIT, $0-0
      // I1 = g.m
        I64Load g_m(g)
      Set I1

      // I2 = g0
        I64Load m_g0(I1)
      Set I2

      // Cannot grow scheduler stack (m->g0).
          Get g
          Get I1
        I64Eq
      If
        CALL runtime·badmorestackg0(SB)
      End

      // Cannot grow signal stack (m->gsignal).
          Get g
          I64Load m_gsignal(I1)
        I64Eq
      If
        CALL runtime·badmorestackgsignal(SB)
      End

      // Called from f.
      // Set m->morebuf to f's caller.

          Get I1
        I32WrapI64
        I64Load 8(SP) // f's caller's PC
      I64Store $m_morebuf+gobuf_pc

          Get I1
        I32WrapI64
            Get SP
          I64ExtendUI32
          I64Const $16
        I64Add // f's caller's SP
      I64Store $m_morebuf+gobuf_sp

          Get I1
        I32WrapI64
        Get g
      I64Store $m_morebuf+gobuf_g

      // Set g->sched to context in f.
          Get g
        I32WrapI64
        I64Load 0(SP) // f's PC
      I64Store $g_sched+gobuf_pc

          Get g
        I32WrapI64
        Get g
      I64Store $g_sched+gobuf_g

          Get g
        I32WrapI64
            Get SP
          I64ExtendUI32
          I64Const $8
        I64Add // f's SP
      I64Store $g_sched+gobuf_sp

          Get g
        I32WrapI64
        Get CTX
      I64Store $g_sched+gobuf_ctxt

      // Call newstack on m->g0's stack.
        Get I2
      Set g

          I64Load (g_sched+gobuf_sp)(I2)
        I32WrapI64
      Set SP

      CALL runtime·newstack(SB)
    End
    Unreachable // crash if newstack returns
  End

// morestack but not preserving ctxt.
TEXT runtime·morestack_noctxt(SB),NOSPLIT,$0
    I64Const $0
  Set CTX
  JMP runtime·morestack(SB)

TEXT ·asmcgocall(SB), NOSPLIT, $0-0
  Unreachable

TEXT ·cgocallback_gofunc(SB), NOSPLIT, $16-32
  Unreachable

#define DISPATCH(NAME, MAXSIZE) \
      Get I0; \
      I64Const $MAXSIZE; \
    I64LeU; \
  If; \
    JMP NAME(SB); \
  End

TEXT reflect·call(SB), NOSPLIT, $0-0
  JMP ·reflectcall(SB)

TEXT ·reflectcall(SB), NOSPLIT, $0-32
      I64Load f+8(FP)
    I64Eqz
  If
    CALL runtime·sigpanic(SB)
  End

    I64Load32U argsize+24(FP)
  Set I0

  DISPATCH(runtime·call32, 32)
  DISPATCH(runtime·call64, 64)
  DISPATCH(runtime·call128, 128)
  DISPATCH(runtime·call256, 256)
  DISPATCH(runtime·call512, 512)
  DISPATCH(runtime·call1024, 1024)
  DISPATCH(runtime·call2048, 2048)
  DISPATCH(runtime·call4096, 4096)
  DISPATCH(runtime·call8192, 8192)
  DISPATCH(runtime·call16384, 16384)
  DISPATCH(runtime·call32768, 32768)
  DISPATCH(runtime·call65536, 65536)
  DISPATCH(runtime·call131072, 131072)
  DISPATCH(runtime·call262144, 262144)
  DISPATCH(runtime·call524288, 524288)
  DISPATCH(runtime·call1048576, 1048576)
  DISPATCH(runtime·call2097152, 2097152)
  DISPATCH(runtime·call4194304, 4194304)
  DISPATCH(runtime·call8388608, 8388608)
  DISPATCH(runtime·call16777216, 16777216)
  DISPATCH(runtime·call33554432, 33554432)
  DISPATCH(runtime·call67108864, 67108864)
  DISPATCH(runtime·call134217728, 134217728)
  DISPATCH(runtime·call268435456, 268435456)
  DISPATCH(runtime·call536870912, 536870912)
  DISPATCH(runtime·call1073741824, 1073741824)
  JMP runtime·badreflectcall(SB)

#define CALLFN(NAME, MAXSIZE) \
TEXT NAME(SB), WRAPPER, $MAXSIZE-32; \
  NO_LOCAL_POINTERS; \
          I64Load32U argsize+24(FP); \
        Set I0; \
        \
              Get I0; \
            I64Eqz; \
          I32Eqz; \
        If $0; \
            Get SP; \
              I64Load argptr+16(FP); \
            I32WrapI64; \
                I64Load argsize+24(FP); \
                I64Const $3; \
              I64ShrU; \
            I32WrapI64; \
          Call runtime·wasmmove(FP); \
        End; \
        \
          I64Load f+8(FP); \
        Set CTX; \
              Get CTX; \
            I32WrapI64; \
          I64Load $0; \
        CALL $0; \
      End; \
      \
        I64Load32U retoffset+28(FP); \
      Set I0; \
        I64Load argtype+0(FP); \
      Set RET0; \
          I64Load argptr+16(FP); \
          Get I0; \
        I64Add; \
      Set RET1; \
            Get SP; \
          I64ExtendUI32; \
          Get I0; \
        I64Add; \
      Set RET2; \
          I64Load32U argsize+24(FP); \
          Get I0; \
        I64Sub; \
      Set RET3; \
      CALL callRet<>(SB); \
    End; \
  End

// callRet copies return values back at the end of call*. This is a
// separate function so it can allocate stack space for the arguments
// to reflectcallmove. It does not follow the Go ABI; it expects its
// arguments in registers.
TEXT callRet<>(SB), NOSPLIT, $32-0
      NO_LOCAL_POINTERS
        Get SP
        Get RET0
      I64Store $0
        Get SP
        Get RET1
      I64Store $8
        Get SP
        Get RET2
      I64Store $16
        Get SP
        Get RET3
      I64Store $24
      CALL runtime·reflectcallmove(SB)
    End
  End

CALLFN(·call32, 32)
CALLFN(·call64, 64)
CALLFN(·call128, 128)
CALLFN(·call256, 256)
CALLFN(·call512, 512)
CALLFN(·call1024, 1024)
CALLFN(·call2048, 2048)
CALLFN(·call4096, 4096)
CALLFN(·call8192, 8192)
CALLFN(·call16384, 16384)
CALLFN(·call32768, 32768)
CALLFN(·call65536, 65536)
CALLFN(·call131072, 131072)
CALLFN(·call262144, 262144)
CALLFN(·call524288, 524288)
CALLFN(·call1048576, 1048576)
CALLFN(·call2097152, 2097152)
CALLFN(·call4194304, 4194304)
CALLFN(·call8388608, 8388608)
CALLFN(·call16777216, 16777216)
CALLFN(·call33554432, 33554432)
CALLFN(·call67108864, 67108864)
CALLFN(·call134217728, 134217728)
CALLFN(·call268435456, 268435456)
CALLFN(·call536870912, 536870912)
CALLFN(·call1073741824, 1073741824)

TEXT runtime·goexit(SB), NOSPLIT, $0-0
    CALL runtime·goexit1(SB) // does not return
  End
  Unreachable

TEXT runtime·cgocallback(SB), NOSPLIT, $32-32
  Unreachable

// gcWriteBarrier performs a heap pointer write and informs the GC.
//
// gcWriteBarrier does NOT follow the Go ABI. It takes two wasm arguments:
// $0: the destination of the write
// $1: the value being written
TEXT runtime·gcWriteBarrier(SB), NOSPLIT, $16
      // I3 = g.m
        I64Load g_m(g)
      Set I3

      // I4 = p
        I64Load m_p(I3)
      Set I4

      // I5 = wbBuf.next
        I64Load (p_wbBuf+wbBuf_next)(I4)
      Set I5

      // Record value
          Get I5
        I32WrapI64
        Get $1
      I64Store $0

      // Record *slot
          Get I5
        I32WrapI64
          Get $0
        I64ExtendUI32
      I64Store $8

      // Increment wbBuf.next
          Get I5
          I64Const $16
        I64Add
      Set I5
          Get I4
        I32WrapI64
        Get I5
      I64Store $p_wbBuf+wbBuf_next

          Get I5
          I64Load (p_wbBuf+wbBuf_end)(I4)
        I64Ne
      BrIf $0

      // Flush
        Get SP
          Get $0
        I64ExtendUI32
      I64Store $0
        Get SP
        Get $1
      I64Store $8
      CALL runtime·wbBufFlush(SB)
    End

    // Do the write
      Get $0
      Get $1
    I64Store $0
  End
