package search

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dlvhdr/gh-dash/v4/internal/tui/context"
)

type Model struct {
	ctx          *context.ProgramContext
	initialValue string
	textInput    textinput.Model
}

type SearchOptions struct {
	Prefix       string
	InitialValue string
	Placeholder  string
}

func NewModel(ctx *context.ProgramContext, opts SearchOptions) Model {
	prompt := fmt.Sprintf(" %s ", opts.Prefix)
	ti := textinput.New()
	ti.Placeholder = opts.Placeholder
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(ctx.Theme.FaintText)
	ti.Width = ctx.MainContentWidth - lipgloss.Width(prompt) - 6
	ti.Prompt = prompt
	ti.SetValue(opts.InitialValue)
	ti.CursorStart()
	ti.Blur()

	m := Model{
		ctx:          ctx,
		textInput:    ti,
		initialValue: opts.InitialValue,
	}

	m.setBlurredStyles()

	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.textInput.Width = m.getInputWidth(m.ctx)
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View(ctx *context.ProgramContext) string {
	return lipgloss.NewStyle().
		Width(ctx.MainContentWidth-4).
		Margin(1, 0).
		Render(m.textInput.View())
}

func (m *Model) Focus() {
	m.textInput.Focus()
	m.textInput.CursorEnd()
	m.setFocusedStyles()
}

func (m *Model) Blur() {
	m.textInput.Blur()
	m.textInput.CursorStart()
	m.setBlurredStyles()
}

func (m *Model) SetValue(val string) {
	m.textInput.SetValue(val)
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	oldWidth := m.textInput.Width
	m.textInput.Width = m.getInputWidth(ctx)
	if m.textInput.Width != oldWidth {
		m.textInput.CursorEnd()
	}
	if m.textInput.Focused() {
		m.setFocusedStyles()
	} else {
		m.setBlurredStyles()
	}
}

func (m *Model) getInputWidth(ctx *context.ProgramContext) int {
	// leave space for at least 2 characters - one character of the input and 1 for the cursor
	// - deduce 4 - 2 for the padding, 2 for the borders
	// - deduce 1 for the cursor
	// - deduce 1 for the spacing between the prompt and text
	return max(2, ctx.MainContentWidth-lipgloss.Width(m.textInput.Prompt)-4-1-1) // borders + cursor
}

func (m Model) Value() string {
	return m.textInput.Value()
}

func (m *Model) inactiveStyle() lipgloss.Style {
	if fg := m.ctx.Styles.ViewSwitcher.InactiveView.GetForeground(); fg != nil {
		return lipgloss.NewStyle().Foreground(fg)
	}
	return lipgloss.NewStyle().Foreground(m.ctx.Theme.FaintText)
}

func (m *Model) activeStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(m.ctx.Theme.PrimaryText)
}

func (m *Model) activePromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(m.ctx.Theme.SecondaryText)
}

func (m *Model) setBlurredStyles() {
	style := m.inactiveStyle()
	m.textInput.TextStyle = style
	m.textInput.PromptStyle = style
	m.textInput.Cursor.Style = style
	m.textInput.Cursor.TextStyle = style
}

func (m *Model) setFocusedStyles() {
	textStyle := m.activeStyle()
	promptStyle := m.activePromptStyle()
	m.textInput.TextStyle = textStyle
	m.textInput.PromptStyle = promptStyle
	m.textInput.Cursor.Style = textStyle
	m.textInput.Cursor.TextStyle = textStyle
}
