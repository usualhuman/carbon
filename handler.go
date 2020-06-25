package carbon

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/fogleman/gg"
)

type (
	// Handler is a basic element for mouse action handling as press, release or hover.
	// It also can be styled for different states using HandlerStyle. Most common are stored in styles package.
	Handler struct {
		Location
		Style    HandlerStyle
		Final    Drawing
		Disabled bool

		OnPress, OnRelease, OnHover func(this *Handler)

		background, foreground     *Style
		hovered, pressed, selected bool
	}

	// HandlerStyle is a set of styles for the Handler.
	HandlerStyle struct {
		Idle, Hover, Active, Focus, Disabled *Style

		Common Drawing
	}
)

func (handler *Handler) Handle(event Event, _, _ float64) {
	if handler.Disabled {
		return
	}
	switch {

	case event == Press.The(pixelgl.MouseButtonLeft):
		if handler.OnPress != nil {
			handler.OnPress(handler)
		}
		handler.pressed = true
		handler.selected = true
		focused = handler
		handler.Update()

	case event == Release.The(pixelgl.MouseButtonLeft):
		if handler.pressed {
			if handler.OnRelease != nil {
				handler.OnRelease(handler)
			}
			handler.pressed = false
			handler.Update()
		}
		fallthrough

	case event.Action == Move:
		if !handler.hovered {
			if handler.OnHover != nil {
				handler.OnHover(handler)
			}
			hovered = append(hovered, handler)
			handler.hovered = true
			handler.Update()
		}
	}
}

func (handler *Handler) Rasterize() {
	handler.Style.finish(handler.Final)
	handler.Style.Rasterize(handler.Size())
	handler.Update()
}

func (handler *Handler) Draw(win *Window) {
	x, y := handler.Center()
	handler.background.Draw(win, x, y)
	handler.foreground.Draw(win, x, y)
}

func (handler *Handler) Land() {
	handler.pressed = false
	handler.hovered = false
	handler.Update()
}

func (handler *Handler) Defocus() {
	if handler == nil {
		return
	}
	handler.selected = false
	handler.Update()
}

func (handler *Handler) Update() {
	switch {
	case handler.pressed:
		handler.background = handler.Style.Active
	case handler.hovered:
		handler.background = handler.Style.Hover
	default:
		handler.background = handler.Style.Idle
	}
	switch {
	case handler.selected:
		handler.foreground = handler.Style.Focus
	case handler.Disabled:
		handler.foreground = handler.Style.Disabled
	default:
		handler.foreground = nil
	}
}

func (style HandlerStyle) finish(final Drawing) {
	for _, state := range []*Style{style.Idle, style.Hover, style.Active, style.Disabled} {
		if state == nil || state.sprite != nil {
			continue
		}
		local := state.Drawing
		state.Drawing = func(ctx *gg.Context) {
			if local != nil {
				local(ctx)
			}
			if style.Common != nil {
				style.Common(ctx)
			}
			if final != nil {
				final(ctx)
			}
		}
	}
}

func (style HandlerStyle) Rasterize(w, h float64) {
	style.Idle.Rasterize(w, h)
	style.Hover.Rasterize(w, h)
	style.Active.Rasterize(w, h)
	style.Focus.Rasterize(w, h)
	style.Disabled.Rasterize(w, h)
}

var hovered []*Handler

func HandleHovered(event Event, x, y float64) {
	if event.Action == Move && event.Button != NilButton {
		return
	}
	for i, hover := range hovered {
		if !hover.Contains(x, y) {
			hover.Land()
			last := len(hovered) - 1
			hovered[i], hovered[last] = hovered[last], nil
			hovered = hovered[:last]
		}
	}
}

var focused *Handler
