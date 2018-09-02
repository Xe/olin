package wasmgo

func (w *wasmGo) daggerOpenFD(sp int32) {
	fname := string(w.goLoadSlice(sp))
	flags := uint32(w.getInt32(sp + 4))
	w.setInt64(sp+8, w.child.OpenFD(fname, flags))
}
