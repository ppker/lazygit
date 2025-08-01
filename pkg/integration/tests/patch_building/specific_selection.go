package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SpecificSelection = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Build a custom patch with a specific selection of lines, adding individual lines, as well as a range and hunk, and adding a file directly",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("hunk-file", "1a\n1b\n1c\n1d\n1e\n1f\n1g\n1h\n1i\n1j\n1k\n1l\n1m\n1n\n1o\n1p\n1q\n1r\n1s\n1t\n1u\n1v\n1w\n1x\n1y\n1z\n")
		shell.Commit("first commit")

		// making changes in two separate places for the sake of having two hunks
		shell.UpdateFileAndAdd("hunk-file", "aa\n1b\ncc\n1d\n1e\n1f\n1g\n1h\n1i\n1j\n1k\n1l\n1m\n1n\n1o\n1p\n1q\n1r\n1s\ntt\nuu\nvv\n1w\n1x\n1y\n1z\n")

		shell.CreateFileAndAdd("line-file", "2a\n2b\n2c\n2d\n2e\n2f\n2g\n2h\n2i\n2j\n2k\n2l\n2m\n2n\n2o\n2p\n2q\n2r\n2s\n2t\n2u\n2v\n2w\n2x\n2y\n2z\n")
		shell.CreateFileAndAdd("direct-file", "direct file content")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Contains("direct-file"),
				Contains("hunk-file"),
				Contains("line-file"),
			).
			SelectNextItem().
			PressPrimaryAction().
			Tap(func() {
				t.Views().Information().Content(Contains("Building patch"))

				t.Views().Secondary().Content(Contains("direct file content"))
			}).
			NavigateToLine(Contains("hunk-file")).
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			SelectedLines(
				Contains("-1a"),
			).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains(`-1a`),
				Contains(`+aa`),
			).
			PressPrimaryAction().
			SelectedLines(
				Contains(`-1c`),
				Contains(`+cc`),
			).
			Tap(func() {
				t.Views().Information().Content(Contains("Building patch"))

				t.Views().Secondary().Content(
					// when we're inside the patch building panel, we only show the patch
					// in the secondary panel that relates to the selected file
					DoesNotContain("direct file content").
						Contains("@@ -1,6 +1,6 @@").
						Contains(" 1f"),
				)
			}).
			// Cancel hunk select
			PressEscape().
			// Escape the view
			PressEscape()

		t.Views().CommitFiles().
			IsFocused().
			NavigateToLine(Contains("line-file")).
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			SelectedLines(
				Contains("+2a"),
			).
			PressPrimaryAction().
			SelectedLines(
				Contains("+2b"),
			).
			NavigateToLine(Contains("+2c")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+2e")).
			PressPrimaryAction().
			SelectedLines(
				Contains("+2f"),
			).
			NavigateToLine(Contains("+2g")).
			PressPrimaryAction().
			SelectedLines(
				Contains("+2h"),
			).
			Tap(func() {
				t.Views().Information().Content(Contains("Building patch"))

				t.Views().Secondary().ContainsLines(
					Contains("+2a"),
					Contains("+2c"),
					Contains("+2d"),
					Contains("+2e"),
					Contains("+2g"),
				)
			}).
			PressEscape().
			Tap(func() {
				t.Views().Secondary().ContainsLines(
					// direct-file patch
					Contains(`diff --git a/direct-file b/direct-file`),
					Contains(`index`),
					Contains(`--- a/direct-file`),
					Contains(`+++ b/direct-file`),
					Contains(`@@ -0,0 +1 @@`),
					Contains(`+direct file content`),
					Contains(`\ No newline at end of file`),
					// hunk-file patch
					Contains(`diff --git a/hunk-file b/hunk-file`),
					Contains(`index`),
					Contains(`--- a/hunk-file`),
					Contains(`+++ b/hunk-file`),
					Contains(`@@ -1,6 +1,6 @@`),
					Contains(`-1a`),
					Contains(`+aa`),
					Contains(` 1b`),
					Contains(` 1c`),
					Contains(` 1d`),
					Contains(` 1e`),
					Contains(` 1f`),
					// line-file patch
					Contains(`diff --git a/line-file b/line-file`),
					Contains(`index`),
					Contains(`--- a/line-file`),
					Contains(`+++ b/line-file`),
					Contains(`@@ -0,0 +1,5 @@`),
					Contains(`+2a`),
					Contains(`+2c`),
					Contains(`+2d`),
					Contains(`+2e`),
					Contains(`+2g`),
				)
			})
	},
})
