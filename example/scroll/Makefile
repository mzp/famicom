CL65	= cl65
LD65	= ld65

CSOURCES =

ASMSOURCES =	scroll.asm
TARGET = scroll

OBJECTS	=	$(CSOURCES:.c=.o) $(ASMSOURCES:.asm=.o)

LIBRARIES =
all :	$(OBJECTS) $(LIBRARIES)
	LD65 -o $(TARGET).nes --config $(TARGET).cfg --obj $(OBJECTS)

.SUFFIXES : .asm .o

.c.o :
	CL65 -t none -o $*.o -c -O $*.c

.asm.o :
	CL65 -t none -o $*.o -c $*.asm

clean :
	rm -rf *.o *.nes
