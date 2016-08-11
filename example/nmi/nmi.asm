	; VBlank割り込みハンドラサンプル

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
.proc Reset
	; PPUコントロールレジスタ1初期化
	lda #%00001000	; ここではVBlank割り込み禁止
	sta $2000

waitVSync:
	lda $2002			; VBlankが発生すると、$2002の7ビット目が1になる
	bpl waitVSync  ; bit7が0の間は、waitVSyncラベルの位置に飛んでループして待ち続ける

	; PPUコントロールレジスタ2初期化
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

	; １番目のスプライト座標初期化
	lda X_Pos_Init
	sta Sprite1_X
	lda Y_Pos_Init
	sta Sprite1_Y
	; ２番目のスプライト座標更新サブルーチンをコール
	jsr setSprite2
	; ２番目のスプライトを水平反転
	lda #%01000000
	sta Sprite2_S

	; PPUコントロールレジスタ2初期化
	lda #%00011110	; スプライトとBGの表示をONにする
	sta $2001

	; PPUコントロールレジスタ1の割り込み許可フラグを立てる
	lda #%10001000
	sta $2000

infinityLoop:					; VBlank割り込み発生を待つだけの無限ループ
	jmp infinityLoop
.endproc

.proc mainLoop
	; スプライト描画(DMAを利用)
	lda #$3  ; スプライトデータは$0300番地からなので、3をロードする。
	sta $4014 ; スプライトDMAレジスタにAをストアして、スプライトデータをDMA転送する

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
	dec Sprite1_Y	; Y座標を1減算
	jmp NOTHINGdown

DOWNKEYdown:
	inc Sprite1_Y ; Y座標を1加算
	jmp NOTHINGdown

LEFTKEYdown:
	dec Sprite1_X	; X座標を1減算
	jmp NOTHINGdown

RIGHTKEYdown:
	inc Sprite1_X	; X座標を1加算
	; この後NOTHINGdownなのでジャンプする必要無し

NOTHINGdown:
	; ２番目のスプライト座標更新サブルーチンをコール
	jsr setSprite2

	rti									; 割り込みから復帰
.endproc

.proc setSprite2
	; ２番目のスプライトの座標更新サブルーチン
	lda Sprite1_X
	adc #8 		; 8ドット右にずらす
	sta Sprite2_X
	lda Sprite1_Y
	sta Sprite2_Y
	rts
.endproc

	; 初期データ
X_Pos_Init:   .byte 20       ; X座標初期値
Y_Pos_Init:   .byte 40       ; Y座標初期値

tilepal: .incbin "giko2.pal" ; パレットをincludeする

.org $0300	 ; $0300から開始、スプライトDMAデータ配置
  Sprite1_Y:     .byte  0   ; スプライト#1 Y座標
  Sprite1_T:     .byte  0   ; スプライト#1 ナンバー
  Sprite1_S:     .byte  0   ; スプライト#1 属性
  Sprite1_X:     .byte  0   ; スプライト#1 X座標
  Sprite2_Y:     .byte  0   ; スプライト#2 Y座標
  Sprite2_T:     .byte  0   ; スプライト#2 ナンバー
  Sprite2_S:     .byte  0   ; スプライト#2 属性
  Sprite2_X:     .byte  0   ; スプライト#2 X座標

.segment "VECINFO"
	.word	mainLoop
	.word	Reset
	.word	$0000

.segment "CHARS"
	.incbin "giko.bkg"  ; 背景データのバイナリィファイルをincludeする
	.incbin "giko2.spr"  ; スプライトデータのバイナリィファイルをincludeする
