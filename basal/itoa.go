package basal

// ItoA w>数字宽度时补0, w <=数字宽度时不补
func ItoA(dst *[]byte, i int, w int)

// ItoAW 强制取数字宽度的w,足够时截断,不足时补0
func ItoAW(dst *[]byte, i int, w int)
