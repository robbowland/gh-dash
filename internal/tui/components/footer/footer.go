package footer

import (
	"fmt"
	"path"
	"strings"

	bbHelp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/dlvhdr/gh-dash/v4/internal/config"
	"github.com/dlvhdr/gh-dash/v4/internal/git"
	"github.com/dlvhdr/gh-dash/v4/internal/tui/context"
	"github.com/dlvhdr/gh-dash/v4/internal/tui/keys"
	"github.com/dlvhdr/gh-dash/v4/internal/utils"
)

type Model struct {
	ctx             *context.ProgramContext
	leftSection     *string
	rightSection    *string
	help            bbHelp.Model
	ShowAll         bool
	ShowConfirmQuit bool
}

func NewModel(ctx *context.ProgramContext) Model {
	help := bbHelp.New()
	help.ShowAll = true
	help.Styles = ctx.Styles.Help.BubbleStyles
	l := ""
	r := ""
	return Model{
		ctx:          ctx,
		help:         help,
		leftSection:  &l,
		rightSection: &r,
	}
}

func (m Model) View() string {
	var footer string

	if m.ShowConfirmQuit {
		footer = lipgloss.NewStyle().Render("Really quit? (Press y/enter to confirm, any other key to cancel)")
	} else {
		showRightSection := m.rightSection != nil && *m.rightSection != ""
		viewSwitcher := m.renderViewSwitcher(m.ctx)
		leftSection := ""
		if m.leftSection != nil {
			leftSection = *m.leftSection
		}
		rightSection := ""
		if m.rightSection != nil {
			rightSection = *m.rightSection
		}
		repoInfo := m.renderRepoInfo(m.ctx)
		repoStyle := m.repoInfoStyle(m.ctx)
		rightContent := ""
		if showRightSection {
			if repoInfo != "" {
				rightContent = lipgloss.JoinHorizontal(
					lipgloss.Top,
					repoInfo,
					repoStyle.Render(" • "),
					rightSection,
				)
			} else {
				rightContent = rightSection
			}
		} else {
			rightContent = repoInfo
		}
		spacing := lipgloss.NewStyle().
			Background(m.ctx.Styles.Common.FooterStyle.GetBackground()).
			Render(
				strings.Repeat(
					" ",
					utils.Max(0,
						m.ctx.ScreenWidth-lipgloss.Width(
							viewSwitcher,
						)-lipgloss.Width(leftSection)-
							lipgloss.Width(rightContent),
					)))

		footer = m.ctx.Styles.Common.FooterStyle.
			Render(lipgloss.JoinHorizontal(lipgloss.Top, viewSwitcher, leftSection, spacing,
				rightContent))
	}

	if m.ShowAll {
		keymap := keys.CreateKeyMapForView(m.ctx.View)
		fullHelp := m.help.View(keymap)
		return lipgloss.JoinVertical(lipgloss.Top, footer, fullHelp)
	}

	return footer
}

func (m *Model) SetShowConfirmQuit(val bool) {
	m.ShowConfirmQuit = val
}

func (m *Model) SetWidth(width int) {
	m.help.Width = width
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	m.help.Styles = ctx.Styles.Help.BubbleStyles
}

func (m *Model) renderViewButton(view config.ViewType) string {
	v := " PRs"
	if view == config.IssuesView {
		v = " Issues"
	}

	if m.ctx.View == view {
		return m.ctx.Styles.ViewSwitcher.ActiveView.Render(v)
	}
	return m.ctx.Styles.ViewSwitcher.InactiveView.Render(v)
}

func (m *Model) renderViewSwitcher(ctx *context.ProgramContext) string {
	separator := ctx.Styles.Tabs.TabSeparator.Render("|")
	spacer := lipgloss.NewStyle().Render(" ")
	items := []string{
		lipgloss.NewStyle().PaddingLeft(1).Render(m.renderViewButton(config.PRsView)),
		spacer,
		separator,
		spacer,
		m.renderViewButton(config.IssuesView),
		spacer,
	}
	view := lipgloss.JoinHorizontal(lipgloss.Top, items...)

	return ctx.Styles.ViewSwitcher.Root.Render(view)
}

func (m *Model) repoInfoStyle(ctx *context.ProgramContext) lipgloss.Style {
	style := ctx.Styles.Common.FooterStyle
	if fg := ctx.Styles.ViewSwitcher.InactiveView.GetForeground(); fg != nil {
		style = style.Foreground(fg)
	} else {
		style = style.Foreground(ctx.Theme.FaintText)
	}
	return style
}

func (m *Model) renderRepoInfo(ctx *context.ProgramContext) string {
	grayStyle := m.repoInfoStyle(ctx)

	var repo string
	if m.ctx.RepoPath != "" {
		name := path.Base(m.ctx.RepoPath)
		if m.ctx.RepoUrl != "" {
			name = git.GetRepoShortName(m.ctx.RepoUrl)
		}
		repo = grayStyle.Render(fmt.Sprintf(" %s", name))
	}

	var user string
	if ctx.User != "" {
		user = grayStyle.Render("@" + ctx.User)
	}

	if repo == "" && user == "" {
		return ""
	}

	leadingSpace := lipgloss.NewStyle().Render(" ")
	if repo == "" {
		return lipgloss.JoinHorizontal(lipgloss.Top, leadingSpace, user)
	}

	if user == "" {
		return lipgloss.JoinHorizontal(lipgloss.Top, leadingSpace, repo)
	}

	separator := grayStyle.Render(" • ")

	return lipgloss.JoinHorizontal(lipgloss.Top, leadingSpace, repo, separator, user)
}

func (m *Model) SetLeftSection(leftSection string) {
	*m.leftSection = leftSection
}

func (m *Model) SetRightSection(rightSection string) {
	*m.rightSection = rightSection
}
