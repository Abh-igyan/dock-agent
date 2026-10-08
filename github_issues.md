# Open these as Issues on your GitHub Repository!

Add the label `good first issue` or `help wanted` to these.

---

### Issue 1: Add Support for Local Models (Ollama / Llama.cpp)
**Title:** Support Local LLMs (Ollama) for offline Dockerfile generation
**Body:**
Currently, `dock-agent` relies on the Gemini API to generate and patch Dockerfiles. We want to support developers who want to run this completely offline or without API costs.

**Goal:**
- Add an `ollama` client implementation in `pkg/llm/ollama.go` alongside the existing `gemini.go`.
- Allow users to select the provider via a CLI flag (e.g., `-provider ollama -model llama3`).

---

### Issue 2: Support Multi-Provider APIs (OpenAI / Anthropic)
**Title:** Add multi-provider LLM support (OpenAI / Anthropic)
**Body:**
To make `dock-agent` flexible, we should support OpenAI and Anthropic Claude APIs.

**Goal:**
- Implement `pkg/llm/openai.go` and `pkg/llm/anthropic.go` that satisfy the LLM interface.
- Allow users to use `OPENAI_API_KEY` or `ANTHROPIC_API_KEY` based on a `-provider` flag.

---

### Issue 3: Migrate to the Official Docker Go SDK
**Title:** Refactor `builder.go` to use the official Docker Go SDK
**Body:**
Currently, `dock-agent` triggers builds by using `os/exec` to run `docker build .` in the shell (in `pkg/docker/builder.go`). 

**Goal:**
- Replace the shell execution with the official [Docker Engine API / Go SDK](https://docs.docker.com/engine/api/sdk/).
- This will make capturing build logs, handling authentication, and detecting successful builds much more robust across different operating systems.

---

### Issue 4: Add a "Run & Probe" validation step
**Title:** Add a container run & health probe validation step
**Body:**
Right now, `dock-agent` considers a self-heal successful if `docker build` exits with code 0. However, a container might build successfully but fail immediately upon starting (e.g., missing runtime dependencies).

**Goal:**
- After a successful build, add an optional step to actually `docker run` the generated image.
- Probe the container (e.g., ensure it doesn't immediately exit with a fatal crash, or check if it responds on a port).
- If it crashes, capture the runtime logs and feed them back to the LLM for a runtime patch!

---

### Issue 5: Interactive/Dry-Run Mode
**Title:** Add interactive confirmation before overwriting existing Dockerfiles
**Body:**
By default, `dock-agent` overwrites the `Dockerfile` in the directory. 

**Goal:**
- Add a `-dry-run` or `-interactive` flag.
- When enabled, `dock-agent` should output the proposed Dockerfile (or a diff) to the terminal and prompt the user `[Y/n]` before writing it to disk and starting the build loop.
