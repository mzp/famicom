; ���ץ饤��ɽ������ץ�
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
; �ꥻ�åȳ�����
.proc	Reset
Start:
	lda $2002  ; VBlank��ȯ������ȡ�$2002��7�ӥå��ܤ�1�ˤʤ�
	bpl Start  ; bit7��0�δ֤ϡ�Start��٥�ΰ��֤�����ǥ롼�פ����Ԥ�

	; PPU����ȥ���쥸���������
	lda #%00001000
	sta $2000
	lda #%00000110		; �������ϥ��ץ饤�Ȥ�BG��ɽ��OFF�ˤ���
	sta $2001

	ldx #$00    ; X�쥸�������ꥢ

	; VRAM���ɥ쥹�쥸������$2006�ˡ��ѥ�åȤΥ�����Υ��ɥ쥹$3F00����ꤹ�롣
	lda #$3F
	sta $2006
	lda #$00
	sta $2006

loadPal: ; ��٥�ϡ��֥�٥�̾��:�פη����ǵ���
	lda tilepal, x ; A��(ourpal + x)���ϤΥѥ�åȤ���ɤ���

	sta $2007 ; $2007�˥ѥ�åȤ��ͤ��ɤ߹���

	inx ; X�쥸�������ͤ�1�û����Ƥ���

	cpx #32 ; X��32(10�ʿ���BG�ȥ��ץ饤�ȤΥѥ�åȤ����)����Ӥ���Ʊ�����ɤ�����Ӥ��Ƥ���
	bne loadPal ;	�夬�������ʤ����ϡ�loadpal��٥�ΰ��֤˥����פ���
	; X��32�ʤ�ѥ�åȥ��ɽ�λ

	; ���ץ饤������
	lda #$00   ; $00(���ץ饤��RAM�Υ��ɥ쥹��8�ӥå�Ĺ)��A�˥���
	sta $2003  ; A�Υ��ץ饤��RAM�Υ��ɥ쥹�򥹥ȥ�

	lda #50     ; 50(10�ʿ�)��A�˥���
	sta $2004   ; Y��ɸ��쥸�����˥��ȥ�����
	lda #00     ; 0(10�ʿ�)��A�˥���
	sta $2004   ; 0�򥹥ȥ�����0�֤Υ��ץ饤�Ȥ���ꤹ��
	sta $2004   ; ȿž��ͥ���̤����ʤ��Τǡ�����$00�򥹥ȥ�����
	lda #20		;	20(10�ʿ�)��A�˥���
	sta $2004   ; X��ɸ��쥸�����˥��ȥ�����

	; PPU����ȥ���쥸����2�����
	lda #%00011110	; ���ץ饤�Ȥ�BG��ɽ����ON�ˤ���
	sta $2001

infinityLoop:
	jmp infinityLoop	; ��������褷�ƽ����ʤΤ�̵�¥롼�פ��ɤ�
.endproc

tilepal:
  .incbin "giko.pal" ; �ѥ�åȤ�include����

.segment "VECINFO"
	.word	$0000
	.word	Reset
	.word	$0000

.segment "CHARS"
	.incbin "giko.bkg"  ; �طʥǡ����ΥХ��ʥꥣ�ե������include����
	.incbin "giko.spr"  ; ���ץ饤�ȥǡ����ΥХ��ʥꥣ�ե������include����
