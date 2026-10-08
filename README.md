# 🐳 DockAgent

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

**DockAgent** is an intelligent, agentic CLI tool written in Go that doesn't just *generate* Dockerfiles—it actively verifies and **self-heals** them. 

Tired of copying a Dockerfile from an LLM, running `docker build`, getting a cryptic dependency error, and pasting it back into the chat? DockAgent automates that entire trial-and-error loop right in your terminal.

## ✨ Features
*   **Context-Aware Generation:** Scans your project directory (`package.json`, `go.mod`, `requirements.txt`) to understand your stack.
*   **Self-Healing Build Loop:** Automatically executes `docker build`. If it fails, it captures the raw `stderr`, feeds it back to the LLM, and patches the `Dockerfile` iteratively until the build succeeds.
*   **Lightweight & Fast:** A single compiled Go binary. No heavy agent frameworks or browser sandboxes required.
*   **Best Practices:** Prompts are tuned to generate secure, multi-stage, non-root Docker images.

## 🚀 How it Works
The core architecture follows a simple Observe-Diagnose-Act loop:
1. **Analyze:** Reads the local directory structure and configuration files.
2. **Generate:** Uses Gemini to draft an initial `Dockerfile`.
3. **Build:** Hooks into the local Docker daemon to attempt a build.
4. **Fix:** If the build fails, the error logs are sent to the LLM for a targeted fix. The loop repeats (up to 3 times) until a successful build is achieved.

## 🛠️ Installation

```bash
# Clone the repository
git clone https://github.com/Abh-igyan/dock-agent.git
cd dock-agent

# Build the binary
go build -o dock-agent main.go
```

## 💻 Usage

DockAgent requires a Gemini API Key to function. 

1. Set your API key:
   ```bash
   export GEMINI_API_KEY="your-api-key-here"
   ```
2. Run the agent in any project directory:
   ```bash
   ./dock-agent -dir /path/to/your/project
   ```

### Command Line Flags
- `-dir`: The directory you want to dockerize (default: `.`)
- `-retries`: The maximum number of times the agent should attempt to self-heal a broken build (default: `3`)

## 🤝 Contributing
Contributions are welcome! Whether it's adding support for more LLM providers (OpenAI, Anthropic, local Ollama models), integrating security scanners (like Trivy), or improving the prompt engineering, feel free to open a Pull Request.

## 📝 License
This project is licensed under the MIT License.
