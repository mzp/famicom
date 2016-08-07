; スプライト表示サンプル
.setcpu		"6502"
.autoimport	on

; iNESヘッダ
.segment "HEADER"
	.byte	$4E, $45, $53, $1A	; "NES" Header
	.byte	$02			; PRG-BANKS
	.byte	$01			; CHR-BANKS
	.byte	$00			; Vetrical Mirror
	.byte	$00			;
	.byte	$00, $00, $00, $00	;
	.byte	$00, $00, $00, $00	;

.segment "STARTUP"
; リセット割り込み
.proc	Reset
Start:
	lda $2002  ; VBlankが発生すると、$2002の7ビット目が1になる
	bpl Start  ; bit7が0の間は、Startラベルの位置に飛んでループして待つ

	; PPUコントロールレジスタ初期化
	lda #%00001000
	sta $2000
	lda #%00000110		; 初期化中はスプライトとBGを表示OFFにする
	sta $2001

	ldx #$00    ; Xレジスタクリア

	; VRAMアドレスレジスタの$2006に、パレットのロード先のアドレス$3F00を指定する。
	lda #$3F
	sta $2006
	lda #$00
	sta $2006

loadPal: ; ラベルは、「ラベル名＋:」の形式で記述
	lda tilepal, x ; Aに(ourpal + x)番地のパレットをロードする

	sta $2007 ; $2007にパレットの値を読み込む

	inx ; Xレジスタに値を1加算している

	cpx #32 ; Xを32(10進数。BGとスプライトのパレットの総数)と比較して同じかどうか比較している
	bne loadPal ;	上が等しくない場合は、loadpalラベルの位置にジャンプする
	; Xが32ならパレットロード終了

	; スプライト描画
	lda #$00   ; $00(スプライトRAMのアドレスは8ビット長)をAにロード
	sta $2003  ; AのスプライトRAMのアドレスをストア

	lda #50     ; 50(10進数)をAにロード
	sta $2004   ; Y座標をレジスタにストアする
	lda #00     ; 0(10進数)をAにロード
	sta $2004   ; 0をストアして0番のスプライトを指定する
	sta $2004   ; 反転や優先順位は操作しないので、再度$00をストアする
	lda #20		;	20(10進数)をAにロード
	sta $2004   ; X座標をレジスタにストアする

	; PPUコントロールレジスタ2初期化
	lda #%00011110	; スプライトとBGの表示をONにする
	sta $2001

infinityLoop:
	jmp infinityLoop	; 今回は描画して終わりなので無限ループで良い
.endproc

tilepal:
  .incbin "giko.pal" ; パレットをincludeする

.segment "VECINFO"
	.word	$0000
	.word	Reset
	.word	$0000

.segment "CHARS"
	.incbin "giko.bkg"  ; 背景データのバイナリィファイルをincludeする
	.incbin "giko.spr"  ; スプライトデータのバイナリィファイルをincludeする
