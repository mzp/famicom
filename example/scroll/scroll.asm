	; BGサンプル
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
	.byte	$00, $00, $00, $00

.segment "VECINFO"
	.word	mainLoop
	.word	Start

.segment "STARTUP"
.proc Start
	; PPUコントロールレジスタ1初期化
	lda #%00001000	; ここではVBlank割り込み禁止
	sta $2000

waitVSync:
	lda $2002			; VBlankが発生すると、$2002の7ビット目が1になる
	bpl waitVSync  ; bit7が0の間は、waitVSyncラベルの位置に飛んでループして待ち続ける


	; PPUコントロールレジスタ2初期化
	lda #%00000110	; 初期化中はスプライトとBGを表示OFFにする
	sta $2001

	; パレットをロード
	ldx #$00    ; Xレジスタクリア

	; VRAMアドレスレジスタの$2006に、パレットのロード先のアドレス$3F00を指定する。
	lda #$3F
	sta $2006
	lda #$00
	sta $2006

loadPal:			; ラベルは、「ラベル名＋:」の形式で記述
	lda tilepal, x ; Aに(ourpal + x)番地のパレットをロードする

	sta $2007 ; $2007にパレットの値を読み込む

	inx ; Xレジスタに値を1加算している

	cpx #32 ; Xを32(10進数。BGとスプライトのパレットの総数)と比較して同じかどうか比較している
	bne loadPal ;	上が等しくない場合は、loadpalラベルの位置にジャンプする
	; Xが32ならパレットロード終了

	; 属性(BGのパレット指定データ)をロード

	; $23C0の属性テーブルにロードする
	lda #$23
	sta $2006
	lda #$C0
	sta $2006

	ldx #$00    ; Xレジスタクリア
	lda #%00000000				; ４つともパレット0番
	; 0番か1番にする
loadAttrib:
	eor #%01010101				; XOR演算で一つおきのビットを交互に０か１にする
	sta $2007							; $2007に属性の値($0か$55)を読み込む
	; 64回(全キャラクター分)ループする
	inx
	cpx #64
	bne loadAttrib

	; ネームテーブル生成

	; $2000のネームテーブルに生成する
	lda #$20
	sta $2006
	lda #$00
	sta $2006

	lda #$00        ; 0番(真っ黒)
	ldy #$00    ; Yレジスタ初期化
loadNametable1:
	ldx Star_Tbl, y			; Starテーブルの値をXに読み込む
loadNametable2:
	sta $2007				; $2007に属性の値を読み込む
	dex							; X減算
	bne loadNametable2	; まだ0でないならばループして黒を出力する
	; 1番か2番のキャラをYの値から交互に取得
	tya							; Y→A
	and #1					; A AND 1
	adc #1					; Aに1加算して1か2に
	sta $2007				; $2007に属性の値を読み込む
	lda #$00        ; 0番(真っ黒)
	iny							; Y加算
	cpy #20					; 20回(星テーブルの数)ループする
	bne loadNametable1

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

	; BGスクロール
	lda $2002			; スクロール値クリア
	lda Scroll_X	; Xのスクロール値をロード
	sta $2005			; X方向スクロール（Y方向は固定)
	inc Scroll_X	; スクロール値を加算

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

	; 星テーブルデータ(20個)
Star_Tbl:    .byte 60,45,35,60,90,65,45,20,90,10,30,40,65,25,65,35,50,35,40,35

tilepal: .incbin "giko2.pal" ; パレットをincludeする

	.org $0000	 ; $0000から開始
Scroll_X:			 .byte  0   ; Xスクロール値

	.org $0300	 ; $0300から開始、スプライトDMAデータ配置
Sprite1_Y:     .byte  0   ; スプライト#1 Y座標
Sprite1_T:     .byte  0   ; スプライト#1 ナンバー
Sprite1_S:     .byte  0   ; スプライト#1 属性
Sprite1_X:     .byte  0   ; スプライト#1 X座標
Sprite2_Y:     .byte  0   ; スプライト#2 Y座標
Sprite2_T:     .byte  0   ; スプライト#2 ナンバー
Sprite2_S:     .byte  0   ; スプライト#2 属性
Sprite2_X:     .byte  0   ; スプライト#2 X座標


.segment "CHARS"
	.incbin "giko2.bkg"  ; 背景データのバイナリィファイルをincludeする
	.incbin "giko2.spr"  ; スプライトデータのバイナリィファイルをincludeする
