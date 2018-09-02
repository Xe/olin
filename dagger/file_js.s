#include "textflag.h"

TEXT ·openFD(SB), NOSPLIT, $0
  CallImport
  RET

TEXT ·read(SB), NOSPLIT, $0
  CallImport
  RET
  