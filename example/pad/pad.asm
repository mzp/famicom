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

.segment "DATA"
X_Pos:
  .byte 0
Y_Pos:
  .byte 0

.segment "STARTUP"
.proc Reset
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
	lda #$3F    ; have $2006 tell
	sta $2006   ; $2007 to start
	lda #$00    ; at $3F00 (pallete).
	sta $2006

loadPal:			; ラベルは、「ラベル名＋:」の形式で記述
	lda tilepal, x ; Aに(ourpal + x)番地のパレットをロードする

	sta $2007 ; $2007にパレットの値を読み込む

	inx ; Xレジスタに値を1加算している

	cpx #32 ; Xを32(10進数。BGとスプライトのパレットの総数)と比較して同じかどうか比較している
	bne loadPal ;	上が等しくない場合は、loadpalラベルの位置にジャンプする
	; Xが32ならパレットロード終了

	; スプライト座標初期化
	lda X_Pos_Init
	sta X_Pos
	lda Y_Pos_Init
	sta Y_Pos

	; PPUコントロールレジスタ2初期化
	lda #%00011110	; スプライトとBGの表示をONにする
	sta $2001

mainLoop:
	lda $2002  ; VBlankが発生すると、$2002の7ビット目が1になります。
	bpl mainLoop  ; bit7が0の間は、mainLoopラベルの位置に飛んでループして待ち続けます。

	; スプライト描画
	lda #$00   ; $00(スプライトRAMのアドレスは8ビット長)をAにロード
	sta $2003  ; AのスプライトRAMのアドレスをストア

	lda Y_Pos  ; Y座標の値をロード
	sta $2004	 ; Y座標をレジスタにストアする

	lda #00     ; 0(10進数)をAにロード
	sta $2004   ; 0をストアして0番のスプライトを指定する
	sta $2004   ; 反転や優先順位は操作しないので、再度$00をストアする

	lda X_Pos  ; X座標の値をロード
	sta $2004	 ; X座標をレジスタにストアする

	; パッドI/Oレジスタの準備
	lda #$01
	sta $4016
	lda #$00
	sta $4016

	; パッド入力チェック
	lda $4016  ; Aボタンをスキップ
	lda $4016  ; Bボタンをスキップ
	lda $4016  ; Selectボタンをスキップ
	lda $4016  ; Startボタンをスキップ
	lda $4016  ; 上ボタン
	and #1     ; AND #1
	bne UPKEYdown  ; 0でないならば押されてるのでUPKeydownへジャンプ

	lda $4016  ; 下ボタン
	and #1     ; AND #1
	bne DOWNKEYdown ; 0でないならば押されてるのでDOWNKeydownへジャンプ

	lda $4016  ; 左ボタン
	and #1     ; AND #1
	bne LEFTKEYdown ; 0でないならば押されてるのでLEFTKeydownへジャンプ

	lda $4016  ; 右ボタン
	and #1     ; AND #1
	bne RIGHTKEYdown ; 0でないならば押されてるのでRIGHTKeydownへジャンプ
	jmp NOTHINGdown  ; なにも押されていないならばNOTHINGdownへ

UPKEYdown:
	dec Y_Pos	; Y座標を1減算。ゼロページなので、以下のコードをこの１命令に短縮できる
;	lda Y_Pos ; Y座標をロード
;	sbc #1  ; #1減算する
;	sta Y_Pos ; Y座標をストア

	jmp NOTHINGdown

DOWNKEYdown:
	inc Y_Pos ; Y座標を1加算
	jmp NOTHINGdown

LEFTKEYdown:
	dec X_Pos	; X座標を1減算
	jmp NOTHINGdown

RIGHTKEYdown:
	inc X_Pos	; X座標を1加算
	; この後NOTHINGdownなのでジャンプする必要無し

NOTHINGdown:
	jmp mainLoop				; mainLoopの最初に戻る
.endproc

	; 初期データ
X_Pos_Init:   .word 20       ; X座標初期値
Y_Pos_Init:   .word 40       ; Y座標初期値

tilepal: .incbin "giko.pal" ; パレットをincludeする

.segment "VECINFO"
	.word	$0000
	.word	Reset
	.word	$0000

.segment "CHARS"
	.incbin "giko.bkg"  ; 背景データのバイナリィファイルをincludeする
	.incbin "giko.spr"  ; スプライトデータのバイナリィファイルをincludeする
