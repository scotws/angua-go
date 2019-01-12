
// Unused mode constants

const (
        // Modes
        ABSOLUTE          = 1  // Absolute                  lda $1000
        ABSOLUTE_X        = 2  // Absolute X indexed        lda.x $1000
        ABSOLUTE_Y        = 3  // Absolute Y indexed        lda.y $1000
        ABSOLUTE_IND      = 3  // Absolute indirect         jmp.i $1000
        ABSOLUTE_IND_LONG = 4  // Absolute indirect long    jmp.il $1000
        ABSOLUTE_LONG     = 5  // Absolute long             jmp.l $101000
        ABSOLUTE_LONG_X   = 6  // Absolute long X indexed   jmp.lx $101000
        ACCUMULATOR       = 7  // Accumulator               inc.a
        BLOCK_MOVE        = 8  // Block move                mvp
        DP                = 9  // Direct page (DP)          lda.d $10
        DP_IND            = 10 // Direct page indirect      lda.di $10
        DP_IND_X          = 11 // DP indirect X indexed     lda.dxi $10
        DP_IND_Y          = 12 // DP indirect Y indexed     lda.diy $10
        DP_IND_LONG       = 13 // DP indirect long          lda.dil $10
        DP_IND_LONG_Y     = 14 // DP indirect long Y index  lda.dily $10
        DP_X              = 15 // Direct page X indexed     lda.dx $10
        DP_Y              = 16 // Direct page Y indexed     ldx.dy $10
        IMMEDIATE         = 17 // Immediate                 lda.# $00
        IMPLIED           = 18 // Implied                   dex
        INDEX_IND         = 19 // Indexed indirect          jmp.xi $1000
        RELATIVE          = 20 // PC Relative               bra <LABEL>
        RELATIVE_LONG     = 21 // PC Relative long          bra.l <LABEL>
        STACK             = 22 // Stack                     pha
        STACK_REL_IND_Y   = 23 // Stack rel ind Y indexed   lda.siy 3
        STACK_REL         = 24 // Stack relative            lda.s 3
)

