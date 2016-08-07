;----------------------------------------------------------------------------
;				NES Sample1 "HELLO, WORLD!"
;					Copyright (C) 2007, Tekepen
;----------------------------------------------------------------------------
.setcpu		"6502"
.autoimport	on

; iNESヘッダ
.segment "HEADER"
	.byte	$4E, $45, $53, $1A	; "NES" Header
	.byte	$02			; PRG-BANKS
	.byte	$02			; CHR-BANKS
	.byte	$01			; Vetrical Mirror
	.byte	$00			;
	.byte	$00, $00, $00, $00	;
	.byte	$00, $00, $00, $00	;

.segment "STARTUP"
; リセット割り込み
.proc	Reset
	sei
	ldx	#$ff
	txs

; スクリーンオフ
	lda	#$00
	sta	$2000
	sta	$2001

; ネームテーブルへ転送(画面の中央付近)
	lda	#$21
	sta	$2006
	lda	#$c9
	sta	$2006
	ldx	#$00
	ldy	#$0d		; 13文字表示
copymap:
	lda	string, x
	sta	$2007
	inx
	dey
	bne	copymap

; スクロール設定
	lda	#$00
	sta	$2005
	sta	$2005

; スクリーンオン
	lda	#$08
	sta	$2000
	lda	#$1e
	sta	$2001

; 無限ループ
mainloop:
	jmp	mainloop
.endproc


; 表示文字列
string:
	.byte	"HELLO, WORLD!"

.segment "VECINFO"
	.word	$0000
	.word	Reset
	.word	$0000

; パターンテーブル
.segment "CHARS"
	.incbin	"character.chr"

; パレットテーブル
.segment "PAL"
	.byte	$0f, $00, $10, $20
	.byte	$0f, $06, $16, $26
	.byte	$0f, $08, $18, $28
	.byte	$0f, $0a, $1a, $2a

	.byte	$0f, $05, $15, $25
	.byte	$0f, $06, $16, $26
	.byte	$0f, $08, $18, $28
	.byte	$0f, $0a, $1a, $2a
