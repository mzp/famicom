	; VBlank�����ߥϥ�ɥ饵��ץ�

.setcpu		"6502"
.autoimport	on

; iNES�إå�
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
	; PPU����ȥ���쥸����1�����
	lda #%00001000	; �����Ǥ�VBlank�����߶ػ�
	sta $2000

waitVSync:
	lda $2002			; VBlank��ȯ������ȡ�$2002��7�ӥå��ܤ�1�ˤʤ�
	bpl waitVSync  ; bit7��0�δ֤ϡ�waitVSync��٥�ΰ��֤�����ǥ롼�פ����Ԥ�³����

	; PPU����ȥ���쥸����2�����
	lda #%00000110		; �������ϥ��ץ饤�Ȥ�BG��ɽ��OFF�ˤ���
	sta $2001

	ldx #$00    ; X�쥸�������ꥢ

	; VRAM���ɥ쥹�쥸������$2006�ˡ��ѥ�åȤΥ�����Υ��ɥ쥹$3F00����ꤹ�롣
	lda #$3F    ; have $2006 tell
	sta $2006   ; $2007 to start
	lda #$00    ; at $3F00 (pallete).
	sta $2006

loadPal:			; ��٥�ϡ��֥�٥�̾��:�פη����ǵ���
	lda tilepal, x ; A��(ourpal + x)���ϤΥѥ�åȤ���ɤ���

	sta $2007 ; $2007�˥ѥ�åȤ��ͤ��ɤ߹���

	inx ; X�쥸�������ͤ�1�û����Ƥ���

	cpx #32 ; X��32(10�ʿ���BG�ȥ��ץ饤�ȤΥѥ�åȤ����)����Ӥ���Ʊ�����ɤ�����Ӥ��Ƥ���
	bne loadPal ;	�夬�������ʤ����ϡ�loadpal��٥�ΰ��֤˥����פ���
	; X��32�ʤ�ѥ�åȥ��ɽ�λ

	; �����ܤΥ��ץ饤�Ⱥ�ɸ�����
	lda X_Pos_Init
	sta Sprite1_X
	lda Y_Pos_Init
	sta Sprite1_Y
	; �����ܤΥ��ץ饤�Ⱥ�ɸ�������֥롼����򥳡���
	jsr setSprite2
	; �����ܤΥ��ץ饤�Ȥ��ʿȿž
	lda #%01000000
	sta Sprite2_S

	; PPU����ȥ���쥸����2�����
	lda #%00011110	; ���ץ饤�Ȥ�BG��ɽ����ON�ˤ���
	sta $2001

	; PPU����ȥ���쥸����1�γ����ߵ��ĥե饰��Ω�Ƥ�
	lda #%10001000
	sta $2000

infinityLoop:					; VBlank������ȯ�����ԤĤ�����̵�¥롼��
	jmp infinityLoop
.endproc

.proc mainLoop
	; ���ץ饤������(DMA������)
	lda #$3  ; ���ץ饤�ȥǡ�����$0300���Ϥ���ʤΤǡ�3����ɤ��롣
	sta $4014 ; ���ץ饤��DMA�쥸������A�򥹥ȥ����ơ����ץ饤�ȥǡ�����DMAž������

	; �ѥå�I/O�쥸�����ν���
	lda #$01
	sta $4016
	lda #$00
	sta $4016

	; �ѥå����ϥ����å�
	lda $4016  ; A�ܥ���򥹥��å�
	lda $4016  ; B�ܥ���򥹥��å�
	lda $4016  ; Select�ܥ���򥹥��å�
	lda $4016  ; Start�ܥ���򥹥��å�
	lda $4016  ; ��ܥ���
	and #1     ; AND #1
	bne UPKEYdown  ; 0�Ǥʤ��ʤ�в�����Ƥ�Τ�UPKeydown�إ�����

	lda $4016  ; ���ܥ���
	and #1     ; AND #1
	bne DOWNKEYdown ; 0�Ǥʤ��ʤ�в�����Ƥ�Τ�DOWNKeydown�إ�����

	lda $4016  ; ���ܥ���
	and #1     ; AND #1
	bne LEFTKEYdown ; 0�Ǥʤ��ʤ�в�����Ƥ�Τ�LEFTKeydown�إ�����

	lda $4016  ; ���ܥ���
	and #1     ; AND #1
	bne RIGHTKEYdown ; 0�Ǥʤ��ʤ�в�����Ƥ�Τ�RIGHTKeydown�إ�����
	jmp NOTHINGdown  ; �ʤˤⲡ����Ƥ��ʤ��ʤ��NOTHINGdown��

UPKEYdown:
	dec Sprite1_Y	; Y��ɸ��1����
	jmp NOTHINGdown

DOWNKEYdown:
	inc Sprite1_Y ; Y��ɸ��1�û�
	jmp NOTHINGdown

LEFTKEYdown:
	dec Sprite1_X	; X��ɸ��1����
	jmp NOTHINGdown

RIGHTKEYdown:
	inc Sprite1_X	; X��ɸ��1�û�
	; ���θ�NOTHINGdown�ʤΤǥ����פ���ɬ��̵��

NOTHINGdown:
	; �����ܤΥ��ץ饤�Ⱥ�ɸ�������֥롼����򥳡���
	jsr setSprite2

	rti									; �����ߤ�������
.endproc

.proc setSprite2
	; �����ܤΥ��ץ饤�Ȥκ�ɸ�������֥롼����
	lda Sprite1_X
	adc #8 		; 8�ɥåȱ��ˤ��餹
	sta Sprite2_X
	lda Sprite1_Y
	sta Sprite2_Y
	rts
.endproc

	; ����ǡ���
X_Pos_Init:   .byte 20       ; X��ɸ�����
Y_Pos_Init:   .byte 40       ; Y��ɸ�����

tilepal: .incbin "giko2.pal" ; �ѥ�åȤ�include����

.org $0300	 ; $0300���鳫�ϡ����ץ饤��DMA�ǡ�������
  Sprite1_Y:     .byte  0   ; ���ץ饤��#1 Y��ɸ
  Sprite1_T:     .byte  0   ; ���ץ饤��#1 �ʥ�С�
  Sprite1_S:     .byte  0   ; ���ץ饤��#1 °��
  Sprite1_X:     .byte  0   ; ���ץ饤��#1 X��ɸ
  Sprite2_Y:     .byte  0   ; ���ץ饤��#2 Y��ɸ
  Sprite2_T:     .byte  0   ; ���ץ饤��#2 �ʥ�С�
  Sprite2_S:     .byte  0   ; ���ץ饤��#2 °��
  Sprite2_X:     .byte  0   ; ���ץ饤��#2 X��ɸ

.segment "VECINFO"
	.word	mainLoop
	.word	Reset
	.word	$0000

.segment "CHARS"
	.incbin "giko.bkg"  ; �طʥǡ����ΥХ��ʥꥣ�ե������include����
	.incbin "giko2.spr"  ; ���ץ饤�ȥǡ����ΥХ��ʥꥣ�ե������include����
