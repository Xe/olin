package wasmgo

import "log"

func (w *wasmGo) daggerOpenFD(sp int32) {
	fname := string(w.goLoadSlice(sp))
	flags := uint32(w.getInt32(sp + 4))

	log.Printf("dagger.OpenFD(%q, %x)", fname, flags)

	w.setInt64(sp+8, w.child.OpenFD(fname, flags))
}
