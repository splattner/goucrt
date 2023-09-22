package denonavr

type DenonCursorControl string

const (
	DenonCursorControlUp     DenonCursorControl = "CUP"
	DenonCursorControlDown                      = "CDN"
	DenonCursorControlLeft                      = "CLT"
	DenonCursorControlRight                     = "CRT"
	DenonCursorControlEnter                     = "ENT"
	DenonCursorControlReturn                    = "RTN"
)

func (d *DenonAVR) CursorControl(cursorControl DenonCursorControl) int {
	status, _ := d.sendCommandToDevice(DenonCommandCursorControl, string(cursorControl))

	return status
}
