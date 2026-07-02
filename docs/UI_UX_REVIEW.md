# strok — UI/UX Review (Interaction & Experience Polish milestone)

Scope: refine the existing experience only. No new features, no architecture
changes, no adaptive lessons. Every proposal below touches rendering, copy,
styling or small amounts of UI state — nothing changes the engine, generator,
persistence or progression rules.

Method: full read of the `ui`, `engine`, `mode`, `lesson`, `stats` packages,
plus headless frame dumps of real render output (fresh lesson, mid-lesson,
after a wrong keystroke, after lesson completion) to verify each observation
against what the app actually draws.

Current frame for reference:

```
╭─────────────────────────────────────────────────────────────╮
│ ⌨  strok · QWERTY · keys: f j                               │  header (left)
│ WPM 0   ACC 100%   ERR 0   CHARS 0   TIME 00:00             │  stats (left)
│                jjjj·jf·ffj·jf·jjjf·fffj·jff·jf              │  lesson (center)
│  ╭───╮╭───╮ ... full keyboard, 15 lines ...                 │  keyboard
│            ● L-pinky   ● L-ring   ● L-mid   ● L-index       │  legend (center, 2 rows)
│       ● R-index   ● R-mid   ● R-ring   ● R-pinky   ● thumb  │
│ esc/ctrl+c quit · backspace correct · tab restart lesson    │  footer (left)
╰─────────────────────────────────────────────────────────────╯
```

---

## High Impact

### H1. Completion note causes a visible layout jump (bug-level polish issue)

- **Current behaviour:** On lesson completion, the outcome message is appended
  inline after the lesson text: `text + note` in `view.go`. The lesson text is
  already width-padded and centered to keyboard width, so the note lands
  *after* the padding — the note line becomes ~40 columns wider than every
  other line, the outer box stretches to fit it, and the whole frame visibly
  jumps wider. On the next keystroke the note clears and the frame snaps back.
  Verified in frame dumps: box grows from 95 → 138 columns.
- **Why it hurts:** This fires at the end of *every* lesson — the app's most
  frequent transition. A frame that lurches sideways at the moment of
  completion reads as broken and destroys the "calm, intentional" feel.
- **Proposed improvement:** Give the outcome message a dedicated, always-reserved
  status line below the lesson text (blank when there is nothing to say),
  centered to the same width as the text. Frame dimensions never change.
- **Expected user benefit:** Rock-stable frame across the type → complete →
  next-lesson loop; completion feedback appears where the eye already is.
- **Implementation complexity:** Low (view.go only, ~10 lines).

### H2. A wrong keystroke gives no feedback in the lesson text

- **Current behaviour:** The engine records the wrong rune at the cursor
  (`Entry{Typed, Status: Incorrect}`) but `renderText` checks `i == cursor`
  *before* checking `Status == Incorrect`, and the cursor never advances past
  an error — so the red "incorrect" style is effectively dead code. After a
  wrong press, the text area is pixel-identical to before (verified in dumps).
  The only signals are a 130 ms keyboard flash and the ERR counter.
- **Why it hurts:** The lesson line is where the user is looking. DESIGN.md
  itself specifies "wrong = red" in the typing area. Today a typo produces no
  change at the point of attention; if you blink through the keyboard flash
  you can sit confused about why the cursor won't advance.
- **Proposed improvement:** When the entry at the cursor is in the `Incorrect`
  state, render the cursor cell in the error style (red, underlined — showing
  the wrongly typed character) until it is corrected or backspaced. This is a
  case-order fix in `renderText`, using state the engine already exposes.
- **Expected user benefit:** Instant, unmissable "that was wrong, fix it here"
  feedback exactly where the eye is focused. Closes the biggest gap in the
  keystroke feedback loop.
- **Implementation complexity:** Low (reorder/merge two cases in render_text.go).

### H3. Error feedback on the keyboard is too transient

- **Current behaviour:** After a wrong press, the expected key (yellow) and the
  wrongly pressed key (red) light up for only 130 ms — the same duration as
  the correct-key flash. DESIGN.md specifies incorrect feedback should persist
  "until next keystroke".
- **Why it hurts:** 130 ms is tuned for confirming success, not diagnosing
  failure. By the time the user's eyes drop to the keyboard to see *which* key
  they hit and which they should have hit, the highlight is gone.
- **Proposed improvement:** Keep the correct-key flash at 130 ms, but hold the
  incorrect pair (expected = yellow, pressed = red) until the next keystroke,
  matching the original design. The engine's `Feedback` already persists until
  the next press; only the `flashing` gate in the renderer needs to distinguish
  correct from incorrect.
- **Expected user benefit:** Errors become teachable moments — the user can
  glance down, see exactly what happened, and re-orient. Correct typing stays
  snappy and unobtrusive.
- **Implementation complexity:** Low (render_keyboard.go / view.go flash gate).

### H4. Lesson result vanishes at the moment it matters

- **Current behaviour:** Completing a lesson immediately swaps in the next one;
  the stats bar resets to `WPM 0 · ACC 100% · CHARS 0`. The user never sees the
  WPM/accuracy of the lesson they just finished. The failure message says
  "need 20 wpm & 90% to advance" without saying what they actually scored, so
  the gate feels arbitrary.
- **Why it hurts:** The end of a lesson is the natural feedback moment — the
  one time the user looks up from the text. Wiping the result before it can be
  read removes the payoff of finishing and makes the progression gate opaque
  ("was I close?").
- **Proposed improvement:** Include the just-finished result in the status line
  (H1's line), e.g. `18 wpm · 94% — need 20 wpm & 90% to advance` on a miss,
  `24 wpm · 96% — new key unlocked: d` on a pass. The `Session` is already in
  hand in `completeLesson`; this is copy + formatting, no new state of
  substance.
- **Expected user benefit:** Every lesson ends with a readable result; the gate
  becomes transparent ("2 wpm short — go again"). Big motivational win for
  zero feature scope.
- **Implementation complexity:** Low–Medium (extend `mode.Outcome` message
  composition or format in the UI from the session it already has).

### H5. Stats bar jitters and says nothing about the goal

- **Current behaviour:** Stat values are rendered at natural width, so when WPM
  goes 9 → 10 → 100 or ACC 100% → 83%, every cell to the right shifts
  horizontally — a small lurch on nearly every keystroke. The values are also
  neutral white; the 20 wpm / 90% thresholds that gate progression are
  invisible until you fail a lesson.
- **Why it hurts:** Constant micro-movement in the peripheral vision is exactly
  the kind of noise a typing app must avoid, and the user has no live sense of
  whether they're on pace to advance.
- **Proposed improvement:** (a) Fixed-width value slots per stat cell so the
  bar never shifts. (b) Color the WPM and ACC values green when they meet the
  advance thresholds (`domain.AdvanceWPM` / `AdvanceAccuracy`, already
  exported), neutral otherwise.
- **Expected user benefit:** A perfectly still HUD, plus an at-a-glance "both
  green = I'll advance" signal that makes the progression system legible while
  typing — without adding any UI.
- **Implementation complexity:** Low (render_stats.go formatting + one style
  choice).

---

## Medium Impact

### M1. Keyboard doesn't distinguish unlocked keys from locked ones

- **Current behaviour:** All character keys render in full finger color
  regardless of progression; modifier keys (ctrl, alt, ⇪, ⇥, ⇧, ⏎, ⌫) —
  which the MVP never asks you to press — render just as brightly. Only the
  header's `keys: f j` text says what's in play.
- **Why it hurts:** The keyboard is the centerpiece, but it presents 60+ keys
  with equal visual weight when only 2–10 matter. The user's eye has no
  guidance toward the active keyset, and unlocking a new key changes nothing
  on the board itself.
- **Proposed improvement:** Render locked character keys and unused modifiers
  dimmed (gray fg + gray border); unlocked keys keep their full finger color.
  The renderer already receives everything needed (`state.Keyset()` is on the
  model). Optionally, when a lesson starts right after an unlock, the newest
  key keeps a subtle accent for that lesson.
- **Expected user benefit:** The active keyset pops out of the board; the
  moment of unlocking a key becomes visible *on the keyboard* — a new key
  literally lights up. Progression becomes something you can see.
- **Implementation complexity:** Medium (pass keyset into renderKeyboard,
  one new dim style; the optional new-key accent needs the outcome's key,
  which the mode/generator already know).

### M2. Mixed alignment makes the frame feel unbalanced

- **Current behaviour:** Header, stats bar and footer are left-aligned; lesson
  text, keyboard and legend are centered. The right half of the top and bottom
  of the frame is empty.
- **Why it hurts:** The eye is pulled to two different axes. The frame reads
  as two different apps stacked: a left-aligned status tool and a centered
  typing surface.
- **Proposed improvement:** Commit to the centered axis for everything in the
  play loop: center the stats bar and the footer to the keyboard width. Keep
  the header as a balanced top bar — title on the left, keyset/layout on the
  right edge of the box (fill the empty space, frame the content).
- **Expected user benefit:** One visual axis; the frame reads as a single
  composed instrument. Calm, intentional, symmetric.
- **Implementation complexity:** Low–Medium (render_stats.go, view.go width
  plumbing that already exists for text/legend).

### M3. Text cursor is weaker than the keyboard highlight

- **Current behaviour:** The current character is underlined + bold in its
  finger color; the current key on the keyboard is a bright inverse block.
  On a rainbow-colored line (untyped chars are finger-colored), an underline
  on one colored glyph is easy to lose.
- **Why it hurts:** The two "what's next" indicators have very different
  strengths — the weaker one is at the primary point of attention. Losing the
  cursor mid-line costs a scan-and-refind pause, breaking flow.
- **Proposed improvement:** Render the cursor cell as an inverse block
  (finger-colored or white background, dark glyph), visually matching the
  keyboard's current-key treatment. The two indicators then read as one linked
  pair: same key, same style, two places.
- **Expected user benefit:** The eye never hunts for the cursor; text and
  keyboard highlights reinforce each other and teach the finger mapping
  faster.
- **Implementation complexity:** Low (one style in render_text.go/theme.go).

### M4. Footer copy is inaccurate and cryptic

- **Current behaviour:** `esc/ctrl+c quit · backspace correct · tab restart lesson`.
  "backspace correct" is not a sentence a new user parses; "tab restart
  lesson" is wrong — tab generates a *fresh* lesson (new random text), it does
  not restart the current one (`restartLesson` calls `Generator.Next`).
- **Why it hurts:** The footer is the only in-app documentation. One hint is
  confusing and another promises behaviour the app doesn't have.
- **Proposed improvement:** `esc quit · ⌫ fix · tab new lesson` (exact wording
  to taste). Keep it dim and short.
- **Expected user benefit:** Hints match reality; less first-session confusion.
- **Implementation complexity:** Trivial (one string).

### M5. Vertical rhythm gives the lesson line no room of its own

- **Current behaviour:** Every section is separated by exactly one blank line;
  the lesson text — the primary interaction — is a single thin line wedged
  between the stats bar and a 15-line keyboard, with no extra emphasis.
- **Why it hurts:** Uniform spacing means uniform importance. The element the
  user must focus on has the smallest visual footprint in the frame.
- **Proposed improvement:** Add one extra blank line above and below the
  lesson/status block so it sits in its own pocket of whitespace (whitespace
  is the only "font size" a terminal has). Requires bumping `minHeight` by 2
  and re-checking the resize guard.
- **Expected user benefit:** The eye lands on the lesson line naturally;
  clearer hierarchy: HUD → *lesson* → keyboard → legend.
- **Implementation complexity:** Low (view.go + minHeight constant).

---

## Nice to Have

### N1. Legend is bulky for what it says

- **Current behaviour:** Two centered rows, nine `● label` items
  (`L-pinky … thumb`), consuming 2 lines + separators every frame — repeating
  information the key colors already encode.
- **Proposed improvement:** Compress to one centered line mirroring hand order,
  e.g. `● pinky ● ring ● mid ● index · ● index ● mid ● ring ● pinky · ● thumb`
  with a left/right grouping — or keep two rows but drop redundant `L-`/`R-`
  prefixes and tighten gaps. Frees a line for M5's whitespace at no height
  cost.
- **Expected user benefit:** Less visual bulk below the keyboard; keyboard and
  lesson gain relative weight.
- **Implementation complexity:** Low (render_legend.go + ShortName usage).

### N2. Header keyset will overflow at higher levels

- **Current behaviour:** `keys: f j` grows one letter per level; by level 20+
  it's a ~60-char spaced list (`f j d k s l a ; g h e i r u w o q p t y c …`)
  crowding the header.
- **Proposed improvement:** Past ~8 keys, summarize: show count plus the newest
  key, e.g. `12 keys · new: e`. Below that, keep the explicit list (it's
  genuinely useful early on).
- **Expected user benefit:** Header stays calm at every level; the newest key —
  the one thing worth reading there — is called out.
- **Implementation complexity:** Low (render_stats.go header formatting).

### N3. Too-small screen doesn't say how far off you are

- **Current behaviour:** "Terminal too small. Resize to at least 97×30 to
  play." — doesn't show the current size, so the user can't tell if they're
  one row short or twenty.
- **Proposed improvement:** Append the current size: `now 92×24`. Dim styling,
  same box.
- **Expected user benefit:** Resize becomes a guided action instead of trial
  and error.
- **Implementation complexity:** Trivial (tooSmallView already has m.width/height).

---

## Explicitly considered and rejected

- **Per-keystroke sounds, animations, spinners** — against the calm/minimal
  goal and the "subtle over flashy" instruction.
- **Completion screen / modal between lessons** — would interrupt flow; the
  status line (H1/H4) delivers the same feedback without a stop.
- **Showing best-WPM / lifetime stats in the HUD** — surfacing more profile
  data starts to smell like feature scope; the exit summary already covers it.
- **Reworking the finger-colored lesson text to monochrome** — the rainbow
  line is a deliberate, shipped teaching feature (commit `abc78ba`); M3
  strengthens the cursor within it instead of replacing it.

## Suggested implementation order (post-approval)

Each step is independently shippable and leaves the app compiling/tested:

1. H1 + H4 (status line with real results — fixes the layout jump and the
   vanishing result together)
2. H2 + H3 (error feedback in text and persistent on keyboard)
3. H5 (stable, threshold-aware stats bar)
4. M1 (locked-key dimming), M2 + M5 (alignment & rhythm), M3 (cursor), M4 (copy)
5. N1–N3 as time allows

## Note on required reading

`docs/PROJECT_CONTEXT.md` referenced in the milestone brief does not exist in
the repo — only `docs/DESIGN.md`. This review treats DESIGN.md, the README and
the code as the source of truth.
