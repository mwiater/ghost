# Setup Instructions

Follow these steps to set up the necessary tools and environment variables.

---

## 1. Install `jq`

### For Windows

#### Option 1: Manually Install `jq`

1. Download `jq` from the [official releases page](https://stedolan.github.io/jq/download/).
   - Choose `jq-win64.exe` if you’re on a 64-bit system.
2. Rename the downloaded file to `jq.exe`.
3. Move `jq.exe` to a directory that is in your system’s `PATH` (e.g., `C:\Windows\System32` or create a folder like `C:\Program Files\jq` and add it to the `PATH`).
4. Verify the installation by opening **Command Prompt** or **PowerShell** and running:
   ```bash
   jq --version
   ```

#### Option 2: Using Chocolatey (Package Manager)

1. If you use **Chocolatey** (Windows package manager), you can install `jq` using:
   ```bash
   choco install jq
   ```
2. Verify the installation by running:
   ```bash
   jq --version
   ```

### For Linux

#### Option 1: Using the Package Manager

Most Linux distributions include `jq` in their package managers. Use the following commands based on your Linux distribution:

- **Debian/Ubuntu**:
  ```bash
  sudo apt-get install jq
  ```

- **Fedora**:
  ```bash
  sudo dnf install jq
  ```

- **Arch Linux**:
  ```bash
  sudo pacman -S jq
  ```

#### Option 2: Manually Download and Install `jq`

1. Download the binary from the [official releases page](https://stedolan.github.io/jq/download/).
2. Move the `jq` binary to `/usr/local/bin` and set the correct permissions:
   ```bash
   sudo mv jq-linux64 /usr/local/bin/jq
   sudo chmod +x /usr/local/bin/jq
   ```
3. Verify the installation by running:
   ```bash
   jq --version
   ```

---

## 2. Set Up `GITHUB_TAGS_TOKEN`

### For Windows

To permanently add the GitHub token as an environment variable:

1. Press `Win + X` and choose **System**.
2. Click **Advanced system settings** > **Environment Variables**.
3. Under **User variables**, click **New**.
4. Set the **Variable name** to `GITHUB_TAGS_TOKEN` and the **Variable value** to your personal GitHub token.
5. Click **OK** and **Apply** to save the changes.

To verify, open a new Command Prompt or PowerShell and run:
```bash
echo %GITHUB_TAGS_TOKEN%
```

### For Linux

To add the GitHub token permanently:

1. Open your shell’s configuration file (`~/.bashrc` or `~/.zshrc`).
2. Add the following line to the file:
   ```bash
   export GITHUB_TAGS_TOKEN=your_generated_token
   ```
3. Save the file and apply the changes by running:
   ```bash
   source ~/.bashrc  # or ~/.zshrc if using zsh
   ```

To verify, open a new terminal and run:
```bash
echo $GITHUB_TAGS_TOKEN
```