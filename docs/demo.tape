# strok README demo — recorded with vhs (https://github.com/charmbracelet/vhs).
#
# Regenerate the GIF from the repo root:
#
#     vhs docs/demo.tape        # writes docs/demo.gif
#
# The demo runs cmd/demorec, a throwaway harness that builds the real strok UI
# with a FIXED-seed generator, so the first f/j lesson is deterministic and the
# keystrokes below reproduce an all-correct (all-green) run every time.

Output "docs/demo.gif"

Set Shell "bash"
Set FontSize 16
Set Width 1320
Set Height 780
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 130ms
Set PlaybackSpeed 1.0

# Launch the deterministic demo harness.
Hide
Type "go run ./cmd/demorec"
Enter
Sleep 2.5s
Show

# Let the opening frame breathe so the keyboard + legend read clearly.
Sleep 1.5s

# Type the exact deterministic first lesson: jffj jff jjj fjf jfff jfj jjjf jfjf
# Typed in bursts so live WPM / accuracy / timer visibly climb.
Type "jffj "
Sleep 400ms
Type "jff "
Sleep 400ms
Type "jjj "
Sleep 400ms
Type "fjf "
Sleep 400ms
Type "jfff "
Sleep 400ms
Type "jfj "
Sleep 400ms
Type "jjjf "
Sleep 400ms
Type "jfjf"
Sleep 2s

# Show the "nice — new key unlocked" outcome / next lesson, then quit.
Sleep 1.5s
Ctrl+C
Sleep 800ms
