package wasmgo

import "log"

func (w *wasmGo) daggerOpenFD(sp int32) {
	fname := string(w.goLoadSlice(sp + 8))
	flags := uint32(w.getInt32(sp + 12))

	log.Printf("dagger.OpenFD(%q, %x)", fname, flags)

	w.setInt64(sp+8, w.child.OpenFD(fname, flags))
}
