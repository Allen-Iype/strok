# strok — Lesson Generation Evolution (design, pre-implementation)

Status: **proposed — awaiting approval, no implementation yet.**

Scope: evolve lesson *texture* — what the generated text looks like — from
random character drills toward pronounceable pseudo-words, real vocabulary,
and (later) programming-oriented lessons. Progression rules (unlock order,
advance gates), the typing engine, stats, and persistence are all out of
scope and unchanged.

---

## 1. Review of the current generator

`internal/lesson` today:

```go
type Generator interface { Next(p domain.Profile) domain.Lesson }   // generator.go
type Progressive struct { rng *rand.Rand; wordsPerLesson, minWordLen, maxWordLen int }
```

`Progressive.Next` reads `Profile.UnlockedLevel`, slices `unlockOrder` into a
keyset, and emits ~8 space-separated "words" of 2–5 uniformly random keyset
runes. It is pure (injected `*rand.Rand`, no I/O), deterministic under test,
and returns a plain `domain.Lesson{Text, Keyset}`.

What works — and must be preserved:

- **Perfect for tiny keysets.** With 2–6 keys the learner is building
  finger-to-key mappings; `fjf jjf` is exactly right, and no richer texture is
  even expressible in that alphabet.
- **Clean seam.** The UI, engine, mode and storage know nothing about how text
  is generated; `cmd/strok` wires one `Generator`. This is the extension point
  the original design promised ("plugin generators").

What breaks down: once ~8+ keys are unlocked, uniform random runes read as
noise. Typing skill past the key-location stage comes from *chunking* —
executing word-sized letter groups as single motor programs — and random
sequences structurally prevent chunking. The mid-game trains stage-1 skills
with stage-2 inventory, and the text is joyless to type.

## 2. Architectural approach: textures behind the existing interface

One concept — a **texture** — and one policy — a **curriculum**:

- Each texture (drill, pseudo-word, vocabulary, code) is its own
  `lesson.Generator` implementation. Same package, same interface, no new
  abstractions.
- A `Curriculum` generator composes them. Its `Next` inspects the profile's
  keyset and delegates to the richest texture the keyset can support. It owns
  the graduation thresholds and nothing else.
- `cmd/strok` still injects exactly one `Generator` (the curriculum). Engine,
  UI, mode, storage: zero changes. `domain.Lesson` unchanged.

```
cmd/strok ── lesson.NewCurriculum(rng)
                    │ delegates by keyset capability
        ┌───────────┼─────────────┬─────────────────┐
     Drill      PseudoWord     Vocabulary        (future: Code)
   (today's    (phonotactic   (embedded word     (identifiers,
  Progressive)  non-words)     list, filtered)    keywords, symbols)
```

Graduation is capability-driven, not level-driven: a texture declares what it
needs from the keyset (pseudo-words need ≥1 vowel; vocabulary needs "enough"
spellable words), so the curriculum stays a few lines and reordering
`unlockOrder` can never break it. Notably, with the letters-first unlock
order, the first vowel (`a`) arrives at level 6 — random drills naturally
carry the first six levels, exactly where they shine, and pseudo-words take
over the moment they become possible.

Because completing a lesson regenerates from the same interface, `tab` (new
lesson), weak-key data, and the advance gate all keep working untouched.

## 3. Stage 2 — pseudo-words (the next milestone)

Pronounceable non-words constrained to the unlocked keyset, keybr-style:
`las gals fask dallas`-shaped output instead of `slka dkfj`.

Mechanism (no data files, pure function of keyset + rng):

- Partition the keyset into vowels and consonants.
- Build words as alternating consonant/vowel clusters with small weighted
  irregularities (occasional double letter, consonant pair like `st`/`fl`
  when both halves are unlocked), 2–6 letters.
- Same lesson shape as today: ~8 words, 30–60 chars, so stats and the advance
  gate stay comparable across textures.

Testing is property-based and cheap: every rune ∈ keyset; every word contains
a vowel; length bounds hold; deterministic under a seeded rng. The lesson
package keeps needing no Bubble Tea, no I/O.

Why not skip straight to dictionary words: with 7–12 unlocked letters the
spellable-English pool is thin and repetitive; phonotactic generation gives
word-shaped chunking practice over *any* keyset, and remains useful later as
a fallback whenever the vocabulary filter comes up short.

## 4. Stages 3+ — how richer textures slot in unchanged

**Vocabulary (stage 3).** A frequency-ranked English word list (top ~10k,
lowercase-ASCII filtered) embedded with `go:embed` (~80 KB — the pure-Go,
single-binary property holds). `Next` filters to words spellable from the
keyset and samples by frequency, weighted toward words containing
`Profile.WeakKeys(n)` — the first time the persisted weak-key data shapes
lesson *content* rather than just stats. The curriculum graduates
pseudo→vocabulary when the filtered pool passes a size threshold (~200
words), falling back to pseudo-words otherwise.

**Programming lessons (future, per DESIGN.md "programming-symbols mode").**
A fourth texture emitting identifiers, keywords, and operators. It pairs with
a wider `unlockOrder` (symbols) and possibly a distinct `mode.Mode` — both
already designed-for extension points. Nothing about the curriculum shape
changes; it's one more delegate with its own capability rule.

**Adaptive lessons (explicitly-later list).** Orthogonal by construction:
adaptivity decides *what to emphasize* (which keys/words), textures decide
*what output looks like*. Building textures first gives a future adaptive
policy a better substrate to steer.

## 5. Smallest first implementation step

Milestone **"Lesson texture I"**, two commits:

1. `PseudoWord` generator in `internal/lesson` (pure, tested, unused).
2. `Curriculum` with one rule — delegate to pseudo-words when the keyset
   contains a vowel, else drill — wired in `cmd/strok/main.go` (one line).

Immediately improves every lesson from level 6 up, touches no other package,
and forces the delegation seam into existence where vocabulary (stage 3)
later drops in as a third delegate.

## Open questions (to settle at implementation time)

1. **Advance-gate calibration.** Word-shaped text types ~10–20% faster than
   noise; 20 wpm / 90% gets slightly easier to pass. Lean: keep global
   thresholds — simplicity over per-texture fairness; revisit only if
   progression feels too fast.
2. **Repetition within a lesson.** Allow repeats (frequency-faithful) or
   dedupe for variety. Lean: allow; drills repeat today and nobody minds.
3. **Texture visibility.** Should the header/status hint the texture change
   ("pseudo-words unlocked")? Lean: no new UI; the text speaks for itself.
