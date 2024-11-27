package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/inoth/toybox/cmd/toybox/internal/base"
)

type Project struct {
	Name string
	Path string
}

func (p *Project) New(ctx context.Context, dir, layout, branch string) error {
	to := filepath.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)

		var model base.SelectModel
		p := tea.NewProgram(&model)
		if _, err := p.Run(); err != nil {
			return err
		}
		if model.Choice == "no" {
			return err
		}
		os.RemoveAll(to)
	}
	fmt.Printf("🚀 Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)
	repo := base.NewRepo(layout, branch)
	if err := repo.CopyTo(ctx, to, p.Name, []string{".git", ".github"}); err != nil {
		return err
	}
	// e := os.Rename(
	// 	filepath.Join(to, "cmd", "server"),
	// 	filepath.Join(to, "cmd", p.Name),
	// )
	// if e != nil {
	// 	return e
	// }
	base.Tree(to, dir)

	fmt.Printf("\n🍺 Project creation succeeded %s\n", color.GreenString(p.Name))
	fmt.Print("💻 Use the following command to start the project 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	return nil
}
