package mr

import (
	"bytes"
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "mr_cmd_test")
	cmdtest.InitTest(m, "mr_cmd_autofill")
}

func TestMrCmd_noARgs(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	assert.Nil(t, NewCmdMR(&cmdutils.Factory{}).Execute())

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	assert.Contains(t, out, "Use \"mr [command] --help\" for more information about a command.\n")

}

func TestMrCmd(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	repo := cmdtest.CopyTestRepo(t, "mr_cmd_test")
	var mrID string
	t.Run("create", func(t *testing.T) {
		git := exec.Command("git", "checkout", "test-cli")
		git.Dir = repo
		b, err := git.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		cmd := exec.Command(cmdtest.GlabBinaryPath, "mr", "create",
			"-t", fmt.Sprintf("MR Title %v", rand.Int()),
			"-d", "This MR is created as a test",
			"-l", "test,bug",
			"--assignee", "profclems",
			"--milestone", "1",
		)
		cmd.Dir = repo

		b, _ = cmd.CombinedOutput()
		out := string(b)
		t.Log(out)
		out = stripansi.Strip(out)
		assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests")
		r := regexp.MustCompile(`!\S+`)
		t.Log(out)
		//i := strings.Index(out, "/diffs\n")
		//mrID = strings.TrimPrefix(out[:i], "https://gitlab.com/glab-cli/test/-/merge_requests/")
		mrID = strings.TrimPrefix(r.FindStringSubmatch(out)[0], "!")
		t.Log(mrID)
	})
	t.Run("show", func(t *testing.T) {
		if mrID == "" {
			t.Skip("mrID is empty, create likely failed")
		}
		cmd := exec.Command(cmdtest.GlabBinaryPath, "mr", "show", mrID)
		cmd.Dir = repo

		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		out := string(b)
		outStripped := stripansi.Strip(out) // To remove ansi chars added by glamour
		require.Contains(t, outStripped, "This MR is created as a test")
		assert.Contains(t, out, fmt.Sprintf("https://gitlab.com/glab-cli/test/-/merge_requests/%s", mrID))
	})
	t.Run("delete", func(t *testing.T) {
		if mrID == "" {
			t.Skip("mrID is empty, create likely failed")
		}
		cmd := exec.Command(cmdtest.GlabBinaryPath, "mr", "delete", mrID)
		cmd.Dir = repo

		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		out := stripansi.Strip(string(b))
		require.Contains(t, out, fmt.Sprintf("- Deleting Merge Request !%s\n", mrID))
		require.Contains(t, out, fmt.Sprintf("Merge request !%s deleted", mrID))
	})
}

func Test_mrCmd_autofill(t *testing.T) {
	repo := cmdtest.CopyTestRepo(t, "mr_cmd_autofill")
	var mrID string
	t.Run("create", func(t *testing.T) {
		git := exec.Command("git", "checkout", "test-cli")
		git.Dir = repo
		b, err := git.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		cmd := exec.Command(cmdtest.GlabBinaryPath, "mr", "create", "-f")
		cmd.Dir = repo

		b, _ = cmd.CombinedOutput()
		out := string(b)
		t.Log(out)
		out = stripansi.Strip(out)
		require.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests")
		r := regexp.MustCompile(`!\S+`)
		mrID = strings.TrimPrefix(r.FindStringSubmatch(out)[0], "!")
		t.Log(mrID)

	})
	t.Run("delete", func(t *testing.T) {
		if mrID == "" {
			t.Skip("mrID is empty, create -F likely failed")
		}
		cmd := exec.Command(cmdtest.GlabBinaryPath, "mr", "delete", mrID)
		cmd.Dir = repo

		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}
		out := stripansi.Strip(string(b))
		require.Contains(t, out, fmt.Sprintf("Deleting Merge Request !%s\n", mrID))
		require.Contains(t, out, fmt.Sprintf("Merge request !%s deleted", mrID))
	})

}
