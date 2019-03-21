#include "textflag.h"

// +build js,wasm go1.11

TEXT ·openFD(SB), NOSPLIT, $0
  CallImport
  RET

TEXT ·read(SB), NOSPLIT, $0
  CallImport
  RET
  